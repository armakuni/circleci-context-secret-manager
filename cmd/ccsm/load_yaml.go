package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/armakuni/circleci-context-secret-manager/manager"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func loadYAML(c *cli.Context) (manager.Contexts, error) {
	contextsDir := c.String("contexts")
	if contextsDir == "" {
		return nil, fmt.Errorf("Must set contexts")
	}
	return loadFiles(contextsDir)
}

func listFiles(contextsDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(contextsDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func loadFiles(contextsDir string) (manager.Contexts, error) {
	files, err := listFiles(contextsDir)
	if err != nil {
		return nil, err
	}
	contexts := make(manager.Contexts)
	for _, file := range files {
		fileName := filepath.Base(file)
		context, err := loadFile(file)
		if err != nil {
			return nil, err
		}
		context.Name = strings.Split(fileName, ".")[0]
		contexts[fileName] = context
	}
	return contexts, nil
}

func loadFile(file string) (manager.Context, error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return manager.Context{}, fmt.Errorf("Could not read yaml in %s, %v", file, err)
	}
	var context manager.Context
	err = yaml.Unmarshal(yamlFile, &context)
	if err != nil {
		return manager.Context{}, fmt.Errorf("Could not parse yaml in %s, %v", file, err)
	}
	return context, nil
}
