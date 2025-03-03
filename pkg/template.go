package pkg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
)

const (
	GEN_CFG_FILENAME      = "package.yml"
	TEMPLATE_DIRNAME      = "templates"
	DEFAULT_TEMPLATES_DIR = "~/.jumpstart/" + TEMPLATE_DIRNAME
	TEMPLATE_SOURCE       = "https://github.com/nurdsoft/jumpstart-templates"
)

func SyncTemplates(dest string) error {
	absDest := AbsPath(dest)

	_, err := os.Stat(absDest)
	if err == nil {
		logrus.Info("Template directory exists, pulling latest changes")
		repo, err := git.PlainOpen(absDest)
		if err != nil {
			return err
		}
		w, err := repo.Worktree()
		if err != nil {
			return err
		}

		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
		})
		if err == git.NoErrAlreadyUpToDate {
			logrus.Info("Template directory is up to date")
			return nil
		}
		return err
	}

	logrus.Info("Template directory does not exist, cloning from remote")
	os.MkdirAll(filepath.Dir(absDest), 0755)
	_, err = git.PlainClone(absDest, false, &git.CloneOptions{
		URL:      TEMPLATE_SOURCE,
		Progress: os.Stdout,
	})

	return err
}

type TemplateProvider struct {
	baseDir string
}

func NewTemplateProvider(baseDir string) (*TemplateProvider, error) {
	absDir := AbsPath(baseDir)
	return &TemplateProvider{baseDir: absDir}, nil
}

func (tp *TemplateProvider) ListTemplates() []string {
	// Read directory
	entries, _ := os.ReadDir(tp.baseDir)
	var templates []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			templates = append(templates, entry.Name())
		}
	}
	return templates
}

func (tp *TemplateProvider) GetTemplate(name string) string {
	return tp.baseDir + "/" + name
}
