package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nurdsoft/jumpstart/pkg/cli"
	"github.com/sirupsen/logrus"
)

func main() {
	err := cli.NewApp().Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
