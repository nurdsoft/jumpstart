package pkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/flosch/pongo2/v6"
)

func SynthesizeProject(ctx context.Context, tid string, dm *DerivedMetadata) error {
	templateCtx := map[string]interface{}{
		// Project name
		"name": dm.Name,
		// Github namespace / organization
		"namespace": dm.Namespace,
		// Authenticated user
		"user":         dm.User,
		"codeprovider": dm.CodeProvider,
	}

	tp, err := NewTemplateProvider(DEFAULT_TEMPLATES_DIR)
	if err != nil {
		return fmt.Errorf("failed to create template provider: %v", err)
	}
	templatePath := tp.GetTemplate(tid)

	cfg, err := GetTemplateConfiguration(templateCtx, templatePath)
	if err != nil {
		return fmt.Errorf("failed to get template configuration: %v", err)
	}

	err = SynthesizeProjectFromDir(templateCtx, templatePath, cfg, dm.OutDir)
	if err != nil {
		return fmt.Errorf("failed to synthesize project: %v", err)
	}

	repo, err := VersionControl(dm.OutDir, dm.Name)
	if err != nil {
		return err
	}

	push := true

	// Github requires deliniation between user and org
	var namespace string
	if dm.IsGithubOrg() {
		namespace = dm.Namespace
	}

	return SetupGithubRemote(ctx, namespace, dm.Name, repo, push)
}

func GetTemplateConfiguration(ctx pongo2.Context, srcTemplateDir string) (*Configuration, error) {
	b, err := renderConfig(ctx, filepath.Join(srcTemplateDir, GEN_CFG_FILENAME))
	if err != nil {
		return nil, err
	}
	return loadConfigBytes(b)
}

func SynthesizeProjectFromDir(ctx map[string]any, srcTemplateDir string, cfg *Configuration, outDir string) (err error) {
	os.MkdirAll(outDir, 0755)

	// Remove project directory if an error occurred
	defer func() {
		if err != nil {
			os.RemoveAll(outDir) // remove directory if an error occurred
		}
	}()

	err = renderTemplates(ctx, filepath.Join(srcTemplateDir, TEMPLATE_DIRNAME), outDir)
	if err != nil {
		return err
	}

	err = runCommands(outDir, cfg.Commands)
	if err != nil {
		return err
	}

	_, err = SynthesizePipelineConfigurationFile(cfg.Pipeline, outDir)
	if err != nil {
		return err
	}

	return err
}

func SynthesizePipelineConfigurationFile(pipeline Pipeline, outDir string) (string, error) {
	tmpfile, err := os.CreateTemp("", "Dockerfile")
	if err != nil {
		return "", err
	}

	_, err = tmpfile.WriteString(pipeline.Dockerfile())
	if err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return "", err
	}
	tmpfile.Close()

	outFile := filepath.Join(outDir, "Dockerfile")
	err = os.Rename(tmpfile.Name(), outFile)
	return outFile, err
}

func runCommands(workdir string, commands Commands) error {
	var err error

	if len(commands.Windows) > 0 {
		err = runCommandsOnWindows(workdir, commands)
		if err != nil {
			return err
		}
	}
	if len(commands.UnixLike) > 0 {
		err = runCommandsOnUnixLike(workdir, commands)
		if err != nil {
			return err
		}
	} else {
		// Run all commands if no OS specific commands are specified
		for _, command := range commands.All {
			cag := strings.Split(command, " ")
			cmd := exec.Command(cag[0], cag[1:]...)
			cmd.Dir = workdir
			if err = cmd.Run(); err != nil {
				return fmt.Errorf("error running command '%s': %w", command, err)
			}
		}
	}
	return nil
}

func runCommandsOnWindows(workdir string, commands Commands) error {
	if runtime.GOOS == "windows" {
		for _, command := range commands.Windows {
			cag := strings.Split(command, " ")
			cmd := exec.Command(cag[0], cag[1:]...)
			cmd.Dir = workdir
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error running windows command '%s': %w", command, err)
			}
		}
	}
	return nil
}

func runCommandsOnUnixLike(workdir string, commands Commands) error {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		for _, command := range commands.UnixLike {
			cag := strings.Split(command, " ")
			cmd := exec.Command(cag[0], cag[1:]...)
			cmd.Dir = workdir
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error running unix command '%s': %w", command, err)
			}
		}
	}
	return nil
}
