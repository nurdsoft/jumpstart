package pkg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
)

func SynthesizeProject(ctx context.Context, tid string, dm *DerivedMetadata, noRemote bool) error {
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

	err = SetupGithubActions(dm.OutDir, dm.Name)
	if err != nil {
		return err
	}

	// Setup version control.  This must be the last step in order to capture all the
	// changes made to the project directory
	repo, err := VersionControl(dm.OutDir, dm.Name)
	if err != nil {
		return err
	}

	// No remote setup requested
	if noRemote {
		return nil
	}

	// Determine namespace. Github requires delineation between user and org
	var namespace string
	if dm.IsGithubOrg() {
		namespace = dm.Namespace
	}

	return SetupGithubRemote(ctx, namespace, dm.Name, repo, true)
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

// SynthesizePipelineConfigurationFile generates a Dockerfile from a pipeline configuration
func SynthesizePipelineConfigurationFile(pipeline Pipeline, outDir string) (string, error) {
	tmpfile, err := os.CreateTemp("", "Dockerfile")
	defer os.Remove(tmpfile.Name())
	if err != nil {
		return "", err
	}

	dockerfile := pipeline.Dockerfile()
	if dockerfile == "" {
		return "", nil
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
