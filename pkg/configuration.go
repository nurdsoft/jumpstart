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
	// return empty string if setup, build and install are empty
	if len(pipeline.Setup) == 0 && len(pipeline.Build) == 0 && len(pipeline.Install) == 0 {
		return ""
	}

	setup := strings.Join(pipeline.Setup, "\n")
	build := strings.Join(pipeline.Build, "\n")
	install := strings.Join(pipeline.Install, "\n")

	run := "ENTRYPOINT " + pipeline.Run

	return setup + "\n" + build + "\n\n" + install + "\n\n" + run + "\n"
}

type Commands struct {
	UnixLike []string `yaml:"unix"`
	Windows  []string `yaml:"windows"`
	All      []string
}

func (c *Commands) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// First, try to unmarshal into the struct
	type plain Commands
	if err := unmarshal((*plain)(c)); err == nil {
		return nil
	}

	// If that fails, try to unmarshal into a slice of strings
	if err := unmarshal(&c.All); err != nil {
		return err
	}

	return nil
}

type Configuration struct {
	Name     string   `yaml:"name"`
	Binary   string   `yaml:"binary"`
	Commands Commands `yaml:"commands"`
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
