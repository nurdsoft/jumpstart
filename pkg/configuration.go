package pkg

import (
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

type Pipeline struct {
	Setup   []string
	Build   []string
	Install []string
	Run     string
}

func (pipeline *Pipeline) Dockerfile() string {
	setup := strings.Join(pipeline.Setup, "\n")
	build := strings.Join(pipeline.Build, "\n")
	install := strings.Join(pipeline.Install, "\n")

	run := "ENTRYPOINT " + pipeline.Run

	return setup + "\n" + build + "\n\n" + install + "\n\n" + run + "\n"
}

type Configuration struct {
	Name     string   `yaml:"name"`
	Binary   string   `yaml:"binary"`
	Commands []string `yaml:"commands"`
	Pipeline Pipeline `yaml:"pipeline"`
}

func loadConfigBytes(b []byte) (*Configuration, error) {
	var cfg Configuration
	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}

	// Make sure binary is in PATH
	_, err = exec.LookPath(cfg.Binary)

	return &cfg, err
}
