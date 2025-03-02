package pkg

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	VERSION   string
	COMMIT    string
	BUILDTIME string
)

func NewApp() *cli.App {
	return &cli.App{
		Name:    "jumpstart",
		Usage:   "jumpstart your project",
		Version: VERSION + "-" + COMMIT + "+" + BUILDTIME,
		Commands: []*cli.Command{
			{
				Name: "template",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list available templates",
						Action: func(c *cli.Context) error {
							tp, err := NewTemplateProvider(DEFAULT_TEMPLATES_DIR)
							if err != nil {
								return fmt.Errorf("failed to create template provider: %v", err)
							}
							list := tp.ListTemplates()
							for _, t := range list {
								fmt.Println(t)
							}
							return nil
						},
					},
					{
						Name:  "sync",
						Usage: "sync templates from remote",
						Action: func(c *cli.Context) error {
							err := SyncTemplates(DEFAULT_TEMPLATES_DIR)
							if err != nil {
								return fmt.Errorf("failed to sync templates: %v", err)
							}
							return nil
						},
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "template to use",
			},
			&cli.BoolFlag{
				Name:  "no-remote",
				Usage: "do not setup remote repository",
			},
		},
		Action: func(c *cli.Context) error {
			// Define the name of the required environment variable
			const requiredEnvVar = "GITHUB_TOKEN"

			// Check if default template directory exists
			_, err := os.Stat(absPath(DEFAULT_TEMPLATES_DIR))
			if os.IsNotExist(err) {
				fmt.Printf("Please run `jumpstart template sync` to seed templates\n\n")
				os.Exit(2)
			}

			tid := c.String("template")
			if tid == "" {
				fmt.Printf("-t | --template is required\n\n")
				cli.ShowAppHelp(c)
				os.Exit(2)
			}

			// Check if the required environment variable is set
			if os.Getenv(requiredEnvVar) == "" {
				fmt.Printf("The environment variable %s is missing or empty.\n", requiredEnvVar)
				fmt.Printf("Please set %s and run the program again.\n", requiredEnvVar)
				return err
			}
			logrus.Info("found github token")

			dm, err := NewDerivedMetadata(c.Context, c.Args().First())
			if err != nil {
				return err
			}
			logrus.Infof("derived metadata: %+v", *dm)

			return SynthesizeProject(c.Context, tid, dm, c.Bool("no-remote"))
		},
	}
}
