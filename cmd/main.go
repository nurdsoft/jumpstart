package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"github.com/nurdsoft/jumpstart/pkg"
)

func main() {
	err := pkg.NewApp().Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
