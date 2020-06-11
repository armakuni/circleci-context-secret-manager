package manager

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Contexts map[string]Context

type Context struct {
	Extends           []string           `yaml:"extends,omitempty"`
	ContextID         string             `yaml:"context_id,omitempty"`
	Secrets           Secrets            `yaml:"secrets"`
	SkipDeploy        bool               `yaml:"skip_deploy,omitempty"`
	RemoteSecretStore *RemoteSecretStore `yaml:"remote_secret_store,omitempty"`
	Name              string
}

type Projects map[string]Project

type Project struct {
	Extends           []string           `yaml:"extends,omitempty"`
	ProjectSlug       string             `yaml:"project_slug,omitempty"`
	Secrets           Secrets            `yaml:"secrets"`
	SkipDeploy        bool               `yaml:"skip_deploy,omitempty"`
	RemoteSecretStore *RemoteSecretStore `yaml:"remote_secret_store,omitempty"`
}

type RemoteSecretStore struct {
	Type string `yaml:"type"`
}

type Secrets map[string]string

type Manager struct {
	AWSSecretManager *secretsmanager.SecretsManager
}

func (m *Manager) ProcessContexts(contexts Contexts) (Contexts, error) {
	processedContext := make(Contexts)
	awsSecretsEnabled, err := contexts.HasAWSRemoteSecrets()
	if err != nil {
		return contexts, err
	}
	if awsSecretsEnabled {
		m.AWSSecretManager = secretsmanager.New(session.New())
	}
	for contextName, context := range contexts {
		newContext := context
		if context.RemoteSecretStore != nil {
			if context.RemoteSecretStore.Type == "aws-secret-manager" {
				secrets, err := m.ProcessSecrets(context.Secrets)
				if err != nil {
					return contexts, err
				}
				newContext.Secrets = secrets
			}
		}
		processedContext[contextName] = newContext
	}
	return processedContext.Process(), nil
}

func (m *Manager) ProcessProjects(projects Projects) (Projects, error) {
	processedProjects := make(Projects)
	awsSecretsEnabled, err := projects.HasAWSRemoteSecrets()
	if err != nil {
		return projects, err
	}
	if awsSecretsEnabled {
		m.AWSSecretManager = secretsmanager.New(session.New())
	}
	for projectName, project := range projects {
		newProject := project
		if project.RemoteSecretStore != nil {
			if project.RemoteSecretStore.Type == "aws-secret-manager" {
				secrets, err := m.ProcessSecrets(project.Secrets)
				if err != nil {
					return projects, err
				}
				newProject.Secrets = secrets
			}
		}
		processedProjects[projectName] = newProject
	}
	return processedProjects.Process(), nil
}

func (m *Manager) ProcessSecrets(secrets Secrets) (Secrets, error) {
	lookupSecretRegex := regexp.MustCompile(`^\(\(.*\)\)$`)
	newSecrets := make(Secrets)
	for key, value := range secrets {
		if lookupSecretRegex.MatchString(value) {
			var err error
			value, err = m.GetAWSSecret(value)
			if err != nil {
				return secrets, err
			}
		}
		newSecrets[key] = value
	}
	return newSecrets, nil
}

func (m *Manager) GetAWSSecret(name string) (string, error) {
	var key = ""
	name = strings.TrimPrefix(name, "((")
	name = strings.TrimSuffix(name, "))")

	if strings.Contains(name, ":") {
		nameParts := strings.Split(name, ":")
		name = nameParts[0]
		key = nameParts[1]
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}
	result, err := m.AWSSecretManager.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("Could not get secret '%s' from aws secrets manager: %v", name, err)
	}

	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			return "", nil
		}
		secretString = string(decodedBinarySecretBytes[:len])
	}

	if key != "" {
		var keyValueSecret map[string]string
		if err := json.Unmarshal([]byte(secretString), &keyValueSecret); err != nil {
			return "", err
		}
		if _, ok := keyValueSecret[key]; !ok {
			return "", fmt.Errorf("AWS Secret Manger secret: %s did not contain the requested key: %s", name, key)
		}
		secretString = keyValueSecret[key]
	}
	return secretString, nil
}

func (contexts Contexts) HasAWSRemoteSecrets() (bool, error) {
	for _, context := range contexts {
		if context.RemoteSecretStore != nil {
			if context.RemoteSecretStore.Type == "aws-secret-manager" {
				return true, nil
			}
			return true, fmt.Errorf("Unsupported remote secret manager '%s', supported managers are", context.RemoteSecretStore.Type)
		}
	}
	return false, nil
}

func (contexts Contexts) Process() Contexts {
	processedContext := make(Contexts)
	for contextName, context := range contexts {
		newContext := context
		newContext.Secrets = make(Secrets)
		for _, extention := range context.Extends {
			if _, ok := contexts[extention]; ok {
				extendedContext := contexts[extention]
				for key, value := range extendedContext.Secrets {
					newContext.Secrets[key] = value
				}
			}
		}
		for key, value := range context.Secrets {
			newContext.Secrets[key] = value
		}
		newContext.Extends = nil
		processedContext[contextName] = newContext
	}
	return processedContext
}

func (projects Projects) HasAWSRemoteSecrets() (bool, error) {
	for _, project := range projects {
		if project.RemoteSecretStore != nil {
			if project.RemoteSecretStore.Type == "aws-secret-manager" {
				return true, nil
			}
			return true, fmt.Errorf("Unsupported remote secret manager '%s', supported managers are", project.RemoteSecretStore.Type)
		}
	}
	return false, nil
}

func (projects Projects) Process() Projects {
	processedProjects := make(Projects)
	for projectName, project := range projects {
		newProject := project
		newProject.Secrets = make(Secrets)
		for _, extention := range project.Extends {
			if _, ok := projects[extention]; ok {
				extendedproject := projects[extention]
				for key, value := range extendedproject.Secrets {
					newProject.Secrets[key] = value
				}
			}
		}
		for key, value := range project.Secrets {
			newProject.Secrets[key] = value
		}
		newProject.Extends = nil
		processedProjects[projectName] = newProject
	}
	return processedProjects
}
