package manager

type Contexts map[string]Context

type Context struct {
	Extends    []string `yaml:"extends,omitempty"`
	ContextID  string   `yaml:"context_id"`
	Secrets    Secrets  `yaml:"secrets"`
	SkipDeploy bool     `yaml:"skip_deploy"`
	Name       string
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
