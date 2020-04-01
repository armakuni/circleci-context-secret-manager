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
