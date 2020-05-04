package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func applyCMD(c *cli.Context) error {
	client, err := circleCIClient(c)
	if err != nil {
		return err
	}
	managerContexts, err := loadYAML(c)
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
		fmt.Printf("Reconfiguring context %s...\n", context.Name)
		err := context.reconfigure(client, managerContexts[contextKey])
		if err != nil {
			return err
		}
	}
	return nil
}

func applyProjectCMD(c *cli.Context) error {
	client, err := circleCISDKClient(c)
	if err != nil {
		return err
	}
	managerProjects, err := loadYAMLProjects(c)
	if err != nil {
		return err
	}
	managerProjects = managerProjects.Process()
	projectsEnvVars, err := getProjects(client, managerProjects)
	if err != nil {
		return err
	}
	for projectKey, managerProject := range managerProjects {
		if managerProject.SkipDeploy == true {
			continue
		}
		fmt.Printf("Reconfiguring project %s...\n", projectName(managerProject.ProjectSlug))
		err := reconfigureProject(client, managerProject, projectsEnvVars[projectKey])
		if err != nil {
			return err
		}
	}
	return nil
}
