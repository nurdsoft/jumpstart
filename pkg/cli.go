package pkg

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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
		Name:      "jumpstart",
		Usage:     "jumpstart your project",
		UsageText: "jumpstart command [command option]",
		Version:   VERSION + "-" + COMMIT + "+" + BUILDTIME,
		Commands: []*cli.Command{
			{
				Name:        "template",
				UsageText:   "jumpstart template [sub-command]",
				Usage:       "Perform template related actions e.g. List, Sync.",
				Description: "Use this to list the templates supported and sync the supported templates with remote",
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
			{
				Name:      "create",
				Usage:     "Create a new project with the specified supporeted template",
				UsageText: "jumpstart create -t={supported template name} [project-name]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "template",
						Aliases: []string{"t"},
						Usage:   "template to use",
					},
				},
				Action: func(c *cli.Context) error {
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

					dm, err := NewDerivedMetadata(c.Context, c.Args().First())
					if err != nil {
						return err
					}
					logrus.Infof("derived metadata: %+v", *dm)

					sshKeys := NewSSHKeys()
					{
						count := len(sshKeys.List())
						if count == 0 {
							fmt.Println("No SSH keys found! Please generate one and try again!")
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
					fmt.Print("Enter ssh passphrase if you added it while ssh-key generation: ")
					var passphrase string
					fmt.Scanln(&passphrase)
					gitSSHKeys, err := ssh.NewPublicKeysFromFile("git", sshKeys.Get(), passphrase)
					if err != nil {
						return err
					}

					return SynthesizeProject(c.Context, tid, dm, gitSSHKeys)
				},
			},
		},
	}
}
