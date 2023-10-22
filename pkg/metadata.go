package pkg

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

type DerivedMetadata struct {
	// Project name
	Name string
	// Github namespace / organization
	Namespace string
	// Authenticated user
	User string
	// Code provider
	CodeProvider string
	// Output directory
	OutDir string
}

func NewDerivedMetadata(ctx context.Context, input string) (*DerivedMetadata, error) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	ghApi := NewGithubAPI()
	user, err := ghApi.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}

	dm := &DerivedMetadata{
		User:         *user.Login,
		CodeProvider: DEFAULT_CODEPROVIDER,
	}

	if input == "" {
		dm.Name = filepath.Base(wd)
		dm.Namespace = *user.Login
		dm.OutDir = wd
	} else {
		parts := strings.Split(input, "/")
		dm.Name = parts[len(parts)-1]
		if len(parts) > 1 {
			dm.Namespace = parts[len(parts)-2]
		} else {
			dm.Namespace = *user.Login
		}
		dm.OutDir = filepath.Join(wd, dm.Name)
	}

	// logrus.Infof("user=%s namespace=%s project=%s", *user.Login, namespace, projName)
	return dm, nil
}

func (dm *DerivedMetadata) IsGithubOrg() bool {
	return dm.Namespace != dm.User
}
