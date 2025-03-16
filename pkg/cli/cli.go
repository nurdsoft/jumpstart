package cli

import (
	"fmt"
	"os"

	"github.com/nurdsoft/jumpstart/pkg"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	VERSION   = "VERSION"
	COMMIT    = "COMMIT"
	BUILDTIME = "BUILDTIME"
)

var (
	requiredEnvVars = []string{
		"GITHUB_TOKEN",
	}
)

func NewApp() *cli.App {
	return &cli.App{
		Name:    "jumpstart",
		Usage:   "jumpstart your project",
		Version: VERSION + "-" + COMMIT + "+" + BUILDTIME,
		Commands: []*cli.Command{
			templateCommand(),
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
			&cli.BoolFlag{
				Name:    "deploy",
				Aliases: []string{"d"},
				Usage:   "setup deployment",
			},
		},
		Action: func(c *cli.Context) error {
			// Check if default template directory exists
			_, err := os.Stat(pkg.AbsPath(pkg.DEFAULT_TEMPLATES_DIR))
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

			// Check environment variables
			if err = checkEnv(); err != nil {
				fmt.Println(err)
				os.Exit(2)
			}

			// Get context metadata
			dm, err := pkg.NewDerivedMetadata(c.Context, c.Args().First())
			if err != nil {
				return err
			}
			logrus.Infof("derived metadata: %+v", *dm)

			var addDeploy *string
			if c.Bool("deploy") {
				d := "deployment"
				addDeploy = &d
			}

			return pkg.SynthesizeProject(c.Context, tid, dm, c.Bool("no-remote"), addDeploy)
		},
	}
}

func checkEnv() error {
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("missing environment variable: %s", envVar)
		}
	}
	return nil
}
