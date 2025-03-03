package cli

import (
	"fmt"

	"github.com/nurdsoft/jumpstart/pkg"
	"github.com/urfave/cli/v2"
)

func templateCommand() *cli.Command {
	return &cli.Command{
		Name: "template",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "list available templates",
				Action: func(c *cli.Context) error {
					tp, err := pkg.NewTemplateProvider(pkg.DEFAULT_TEMPLATES_DIR)
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
					err := pkg.SyncTemplates(pkg.DEFAULT_TEMPLATES_DIR)
					if err != nil {
						return fmt.Errorf("failed to sync templates: %v", err)
					}
					return nil
				},
			},
		},
	}
}
