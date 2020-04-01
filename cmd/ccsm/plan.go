package main

import (
	"fmt"

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
	managerContexts = managerContexts.Process()
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
