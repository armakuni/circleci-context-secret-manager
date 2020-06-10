package main

import (
	"fmt"

	"github.com/armakuni/circleci-context-secret-manager/manager"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func planCMD(c *cli.Context) error {
	context := c.String("context")
	if context == "" {
		return fmt.Errorf("Must set a specific context for interpolation")
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
	if _, ok := managerContexts[context]; !ok {
		return fmt.Errorf("Could not find a context file for %s", context)
	}
	d, err := yaml.Marshal(managerContexts[context])
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(d))
	return nil
}

func planProjectCMD(c *cli.Context) error {
	project := c.String("project")
	if project == "" {
		return fmt.Errorf("Must set a specific project for interpolation")
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
	managerProjects = managerProjects.Process()
	if _, ok := managerProjects[project]; !ok {
		return fmt.Errorf("Could not find a project file for %s", project)
	}
	d, err := yaml.Marshal(managerProjects[project])
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(d))
	return nil
}
