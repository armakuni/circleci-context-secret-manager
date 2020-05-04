package main

import (
	"fmt"
	"strings"

	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/client"
	"github.com/armakuni/circleci-context-secret-manager/manager"
	"github.com/armakuni/circleci-workflow-dashboard/circleci"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type Data struct {
	Context Context `json:"context"`
}

type Context struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Resources []api.Resource `json:"resources"`
}

type SecretDiffs struct {
	New         []string
	ToBeDeleted []string
	ToUpdate    []string
}

func (c *Context) getChanges(managerContext manager.Context) SecretDiffs {
	var diffs SecretDiffs
	matchingDiffs := make(map[string]interface{})

	for _, resource := range c.Resources {
		if _, ok := managerContext.Secrets[resource.Variable]; ok {
			matchingDiffs[resource.Variable] = nil
			diffs.ToUpdate = append(diffs.ToUpdate, resource.Variable)
		} else {
			diffs.ToBeDeleted = append(diffs.ToBeDeleted, resource.Variable)
		}
	}
	for secretKey, _ := range managerContext.Secrets {
		if _, ok := matchingDiffs[secretKey]; !ok {
			diffs.New = append(diffs.New, secretKey)
		}
	}
	return diffs
}

func getProjectChanges(managerProject manager.Project, projectEnvVars circleci.ProjectEnvVars) SecretDiffs {
	var diffs SecretDiffs
	matchingDiffs := make(map[string]interface{})

	for _, envVar := range projectEnvVars {
		if _, ok := managerProject.Secrets[envVar.Name]; ok {
			matchingDiffs[envVar.Name] = nil
			diffs.ToUpdate = append(diffs.ToUpdate, envVar.Name)
		} else {
			diffs.ToBeDeleted = append(diffs.ToBeDeleted, envVar.Name)
		}
	}
	for secretKey, _ := range managerProject.Secrets {
		if _, ok := matchingDiffs[secretKey]; !ok {
			diffs.New = append(diffs.New, secretKey)
		}
	}
	return diffs
}

func (c *Context) reconfigure(cl *client.Client, managerContext manager.Context) error {
	diffs := c.getChanges(managerContext)
	for _, new := range diffs.New {
		if err := api.StoreEnvironmentVariable(cl, c.ID, new, managerContext.Secrets[new]); err != nil {
			return fmt.Errorf("Error storing secret %s on context %s: %v", new, c.Name, err)
		}
	}
	for _, update := range diffs.ToUpdate {
		if err := updateEnvironmentVariable(cl, c, update, managerContext.Secrets[update]); err != nil {
			return err
		}
	}
	for _, delete := range diffs.ToBeDeleted {
		if err := api.DeleteEnvironmentVariable(cl, c.ID, delete); err != nil {
			return fmt.Errorf("Error removing secret %s on context %s: %v", delete, c.Name, err)
		}
	}
	return nil
}

func reconfigureProject(cl *circleci.Client, managerProject manager.Project, projectEnvVars circleci.ProjectEnvVars) error {
	diffs := getProjectChanges(managerProject, projectEnvVars)
	for _, new := range diffs.New {
		if _, err := cl.CreateProjectEnvVar(managerProject.ProjectSlug, new, managerProject.Secrets[new]); err != nil {
			return fmt.Errorf("Error storing secret %s on project %s: %v", new, managerProject.ProjectSlug, err)
		}
	}
	for _, update := range diffs.ToUpdate {
		if err := cl.DeleteProjectEnvVar(managerProject.ProjectSlug, update); err != nil {
			return fmt.Errorf("Error updating secret %s on project %s: %v", update, managerProject.ProjectSlug, err)
		}
		if _, err := cl.CreateProjectEnvVar(managerProject.ProjectSlug, update, managerProject.Secrets[update]); err != nil {
			return fmt.Errorf("Error updating secret %s on project %s: %v", update, managerProject.ProjectSlug, err)
		}
	}
	for _, delete := range diffs.ToBeDeleted {
		if err := cl.DeleteProjectEnvVar(managerProject.ProjectSlug, delete); err != nil {
			return fmt.Errorf("Error removing secret %s on project %s: %v", delete, managerProject.ProjectSlug, err)
		}
	}
	return nil
}

func updateEnvironmentVariable(cl *client.Client, context *Context, key, value string) error {
	if err := api.DeleteEnvironmentVariable(cl, context.ID, key); err != nil {
		return fmt.Errorf("Error removing secret %s for update on context %s: %v", key, context.Name, err)
	}
	if err := api.StoreEnvironmentVariable(cl, context.ID, key, value); err != nil {
		return fmt.Errorf("Error storing secret %s for update on context %s: %v", key, context.Name, err)
	}
	return nil
}

func circleCIClient(c *cli.Context) (*client.Client, error) {
	apiToken := c.String("api-token")
	if apiToken == "" {
		return nil, fmt.Errorf("Must set api-token")
	}
	return client.NewClient(c.String("circleci-url"), "graphql-unstable", apiToken, false), nil
}

func circleCISDKClient(c *cli.Context) (*circleci.Client, error) {
	apiToken := c.String("api-token")
	if apiToken == "" {
		return nil, fmt.Errorf("Must set api-token")
	}
	return circleci.NewClient(&circleci.Config{APIToken: apiToken, APIURL: c.String("cirlceci-url")})
}

func getContexts(cl *client.Client, managerContexts manager.Contexts) (map[string]*Context, error) {
	contexts := make(map[string]*Context)
	for contextKey, managerContext := range managerContexts {
		if managerContext.SkipDeploy {
			continue
		}
		context, err := getContext(cl, managerContext.ContextID)
		if err != nil {
			return nil, fmt.Errorf("Could not get context with ID: '%s', do you have the correct ID and permissions?", managerContext.ContextID)
		}
		contexts[contextKey] = context
	}
	return contexts, nil
}

func getProjects(cl *circleci.Client, managerProjects manager.Projects) (map[string]circleci.ProjectEnvVars, error) {
	projectEnvVars := make(map[string]circleci.ProjectEnvVars)
	for contextKey, managerProject := range managerProjects {
		envVars, err := cl.GetProjectEnvVars(managerProject.ProjectSlug)
		if err != nil {
			return nil, err
		}
		projectEnvVars[contextKey] = envVars
	}
	return projectEnvVars, nil
}

func validateContexts(contexts map[string]*Context, managerContexts manager.Contexts) error {
	for contextKey, managerContext := range managerContexts {
		if managerContext.SkipDeploy {
			continue
		}
		err := validateContextID(contexts[contextKey], managerContext)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateContextID(context *Context, managerContext manager.Context) error {
	if managerContext.ContextID == "" {
		return fmt.Errorf("Context file '%s', does not contain a context ID", managerContext.Name)
	}
	if context.Name != managerContext.Name {
		return fmt.Errorf("The context name for '%s' was '%s', but the context file was called '%s', these need to match, please check your context ID and update accordingly or rename your context file", context.ID, context.Name, managerContext.Name)
	}
	return nil
}

func getContext(cl *client.Client, contextID string) (*Context, error) {
	query := fmt.Sprintf(`
	{
		context(id: "%s") {
			...Context
		}
	}
	fragment Context on Context {
		id
		name
		createdAt
		groups {
			edges {
				node {
					...SecurityGroups
				}
			}
		}
		resources {
			...EnvVars
		}
	}
	fragment EnvVars on EnvironmentVariable {
		variable
		createdAt
		truncatedValue
	}
	fragment SecurityGroups on Group {
		id
		name
	}
	`, contextID)

	request := client.NewRequest(query)
	request.SetToken(cl.Token)

	var response Data
	err := cl.Run(request, &response)
	return &response.Context, errors.Wrapf(err, "failed to load context list")
}

func projectName(projectSlug string) string {
	splitSlug := strings.Split(projectSlug, "/")
	return splitSlug[len(splitSlug)-1]
}
