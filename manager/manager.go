package manager

type Contexts map[string]Context

type Context struct {
	Extends    []string `yaml:"extends,omitempty"`
	ContextID  string   `yaml:"context_id,omitempty"`
	Secrets    Secrets  `yaml:"secrets"`
	SkipDeploy bool     `yaml:"skip_deploy,omitempty"`
	Name       string
}

type Projects map[string]Project

type Project struct {
	Extends     []string `yaml:"extends,omitempty"`
	ProjectSlug string   `yaml:"project_slug,omitempty"`
	Secrets     Secrets  `yaml:"secrets"`
	SkipDeploy  bool     `yaml:"skip_deploy,omitempty"`
}

type Secrets map[string]string

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
