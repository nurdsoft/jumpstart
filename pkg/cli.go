package pkg

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:  "jumpstart",
		Usage: "jumpstart your project",
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
		},
		Action: func(c *cli.Context) error {
			tid := c.String("template")
			if tid == "" {
				return fmt.Errorf("-t | --template is required")
			}

			dm, err := NewDerivedMetadata(c.Context, c.Args().First())
			if err != nil {
				return err
			}
			logrus.Infof("derived metadata: %+v", *dm)

			sshKeys := NewSSHKeys()
			{
				count := len(sshKeys.List())
				if count == 0 {
					fmt.Println("No SSH keys found, please generate one and try again!")
					fmt.Println("See https://docs.github.com/en/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent for more information")
					os.Exit(2)
				} else if count > 1 {
					fmt.Printf("Multiple SSH keys found! Please select one to use:\n\n")
					for i, k := range sshKeys.List() {
						fmt.Printf("%d. %s\n", i+1, k)
					}
					fmt.Printf("\nEnter selection: ")

					// Read selection from user
					var selection int
					fmt.Scanln(&selection)
					sshKeys.SetSelected(selection - 1)
				}
			}
			fmt.Printf("\nUsing ssh key: %s\n", sshKeys.Get())

			gitSSHKeys, err := ssh.NewPublicKeysFromFile("git", sshKeys.Get(), "")
			if err != nil {
				return err
			}

			return SynthesizeProject(c.Context, tid, dm, gitSSHKeys)
		},
	}
}
