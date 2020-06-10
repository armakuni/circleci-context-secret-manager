package main

import (
	"fmt"

	"github.com/armakuni/circleci-context-secret-manager/manager"
	"github.com/urfave/cli/v2"
)

func dryRunCMD(c *cli.Context) error {
	client, err := circleCIClient(c)
	if err != nil {
		return err
	}
	managerContexts, err := loadYAML(c)
	if err != nil {
		return err
	}
	manager := &manager.Manager{}
	managerContexts, err = manager.ProcessContexts(managerContexts)
	if err != nil {
		return err
	}
	managerContexts = managerContexts.Process()
	contexts, err := getContexts(client, managerContexts)
	if err != nil {
		return err
	}
	err = validateContexts(contexts, managerContexts)
	if err != nil {
		return err
	}
	for contextKey, context := range contexts {
		fmt.Printf("\nContext '%s':\n", context.Name)
		diffs := context.getChanges(managerContexts[contextKey])
		if len(diffs.New) > 0 {
			fmt.Printf("  New secrets:\n")
			for _, new := range diffs.New {
				fmt.Printf("    %s\n", new)
			}
		}
		if len(diffs.ToUpdate) > 0 {
			fmt.Printf("  Updating secrets:\n")
			for _, update := range diffs.ToUpdate {
				fmt.Printf("    %s\n", update)
			}
		}
		if len(diffs.ToBeDeleted) > 0 {
			fmt.Printf("  Deleting secrets:\n")
			for _, delete := range diffs.ToBeDeleted {
				fmt.Printf("    %s\n", delete)
			}
		}
	}
	return nil
}

func dryRunProjectsCMD(c *cli.Context) error {
	client, err := circleCISDKClient(c)
	if err != nil {
		return err
	}
	managerProjects, err := loadYAMLProjects(c)
	if err != nil {
		return err
	}
	manager := &manager.Manager{}
	managerProjects, err = manager.ProcessProjects(managerProjects)
	if err != nil {
		return err
	}
	projectsEnvVars, err := getProjects(client, managerProjects)
	if err != nil {
		return err
	}
	for projectKey, managerProject := range managerProjects {
		if managerProject.SkipDeploy == true {
			continue
		}
		fmt.Printf("\nProject '%s':\n", projectName(managerProject.ProjectSlug))
		diffs := getProjectChanges(managerProject, projectsEnvVars[projectKey])
		if len(diffs.New) > 0 {
			fmt.Printf("  New secrets:\n")
			for _, new := range diffs.New {
				fmt.Printf("    %s\n", new)
			}
		}
		if len(diffs.ToUpdate) > 0 {
			fmt.Printf("  Updating secrets:\n")
			for _, update := range diffs.ToUpdate {
				fmt.Printf("    %s\n", update)
			}
		}
		if len(diffs.ToBeDeleted) > 0 {
			fmt.Printf("  Deleting secrets:\n")
			for _, delete := range diffs.ToBeDeleted {
				fmt.Printf("    %s\n", delete)
			}
		}
	}
	return nil
}
