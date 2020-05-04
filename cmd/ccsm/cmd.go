package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var contextFlag = &cli.StringFlag{
	Name:    "contexts, c",
	Value:   "./contexts",
	Usage:   "The directory your contexts secrets files are stored in",
	EnvVars: []string{"CONTEXTS_DIR"},
}

var projectFlag = &cli.StringFlag{
	Name:    "projects, p",
	Value:   "./projects",
	Usage:   "The directory your projects secrets files are stored in",
	EnvVars: []string{"PROJECTS_DIR"},
}

var extensionsFlag = &cli.StringFlag{
	Name:    "extensions, e",
	Value:   "./extensions",
	Usage:   "The directory your extensions for secrets files are stored in",
	EnvVars: []string{"EXTENSIONS_DIR"},
}

func main() {
	app := &cli.App{
		Name:  "ccsm",
		Usage: "CircleCI Context Secret Manager",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "circleci-url",
				Aliases: []string{"u"},
				Value:   "https://circleci.com",
				Usage:   "The URL of your CircleCI server",
				EnvVars: []string{"CIRCLECI_URL"},
			},
			&cli.StringFlag{
				Name:    "api-token",
				Aliases: []string{"t"},
				Usage:   "Your CircleCI API Token",
				EnvVars: []string{"CIRCLECI_API_TOKEN"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "apply",
				Usage: `Apply all changes to context secrets`,
				UsageText: `
  Note: This command is idempotent, your context secrets will be set to exactly what is in the config.
    If you are unsure what changes this will make you can run 'dry-run' first to print out a basic report of changes
`,

				Flags: []cli.Flag{
					contextFlag,
				},
				Action: applyCMD,
			}, {
				Name:  "apply-projects",
				Usage: `Apply all changes to project secrets`,
				UsageText: `
  Note: This command is idempotent, your context secrets will be set to exactly what is in the config.
    If you are unsure what changes this will make you can run 'dry-run-projects' first to print out a basic report of changes
`,

				Flags: []cli.Flag{
					projectFlag,
					extensionsFlag,
				},
				Action: applyProjectCMD,
			}, {
				Name:  "dry-run",
				Usage: `Print out a dry run report`,
				UsageText: `
  Check if your Context IDs match

  Check what secrets will be added/ deleted

  Note: Due to limitations in the CircleCI APIs with secret masking we are unable to dry run any modified secrets
`,

				Flags: []cli.Flag{
					contextFlag,
				},
				Action: dryRunCMD,
			}, {
				Name:  "plan",
				Usage: `Run a plan for a context file, see exactly what secrets will be set for your context on apply`,
				Flags: []cli.Flag{
					contextFlag,
					&cli.StringFlag{
						Name:  "context",
						Usage: "The context you want to show output for (as the file name)",
					},
				},
				Action: planCMD,
			}, {
				Name:  "dry-run-projects",
				Usage: `Print out a dry run report for projects`,
				UsageText: `
  Check if your Project Slugs match

  Check what secrets will be added/ deleted

  Note: Due to limitations in the CircleCI APIs with secret masking we are unable to dry run any modified secrets
`,

				Flags: []cli.Flag{
					projectFlag,
					extensionsFlag,
				},
				Action: dryRunProjectsCMD,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
