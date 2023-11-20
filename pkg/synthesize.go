package pkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func SynthesizeProjectFromDir(ctx map[string]any, srcTemplateDir string, cfg *Configuration, outDir string) error {
	os.MkdirAll(outDir, 0755)

	err := renderTemplates(ctx, filepath.Join(srcTemplateDir, TEMPLATE_DIRNAME), outDir)
	if err != nil {
		return err
	}

	err = runCommands(outDir, cfg.Commands)
	if err != nil {
		return err
	}

	_, err = SynthesizePipelineConfigurationFile(cfg.Pipeline, outDir)

	return err
}

func SynthesizePipelineConfigurationFile(pipeline Pipeline, outDir string) (string, error) {
	tmpfile, err := os.CreateTemp("", "Dockerfile")
	defer os.Remove(tmpfile.Name())
	if err != nil {
		return "", err
	}

	_, err = tmpfile.WriteString(pipeline.Dockerfile())
	if err != nil {
		tmpfile.Close()
		return "", err
	}
	tmpfile.Close()

	outFile := filepath.Join(outDir, "Dockerfile")
	err = copyFile(tmpfile.Name(), outFile)
	return outFile, err
}

func runCommands(workdir string, commands []string) error {
	var err error
	for _, instr := range commands {
		cag := strings.Split(instr, " ")
		cmd := exec.Command(cag[0], cag[1:]...)
		cmd.Dir = workdir
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("error running command '%s': %w", instr, err)
		}
	}
	return nil
}
