package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/google/go-github/v56/github"
)

const (
	DEFAULT_CODEPROVIDER = "github.com"
)

func VersionControl(dir, repoName string) (*git.Repository, error) {
	// Initialize git repo in the given dir
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		return nil, err
	}

	// Get the worktree
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	// Add all files to the repo worktree
	_, err = wt.Add(".")
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, dir)
	}

	// Commit the files
	_, err = wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Jump Start",
			Email: "jumpstart",
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// namespace should be empty if it is not an org
func SetupGithubRemote(ctx context.Context, namespace string, repoName string, localRepo *git.Repository, push bool, auth ssh.AuthMethod) error {
	// Create github repo
	ghRepo, err := NewGithubAPI().CreateRepo(ctx, namespace, repoName)
	if err != nil {
		return err
	}

	// Add github as remote to local repo
	_, err = localRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{ghRepo.GetSSHURL()},
	})

	if err != nil {
		return err
	}

	if push {
		// Push to local repo to github
		err = localRepo.Push(&git.PushOptions{
			RemoteName: "origin",
			Auth:       auth,
		})
	}
	return err
}

type GithubAPI struct {
	client *github.Client
}

func NewGithubAPI() *GithubAPI {
	return &GithubAPI{
		client: github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN")),
	}
}

func (api *GithubAPI) GetAuthUser(ctx context.Context) (*github.User, error) {
	user, _, err := api.client.Users.Get(ctx, "")
	return user, err
}

// orgName is the name of the organization to create the repo in. If it is an empty string,
// the repo will be created in the user's account.
func (api *GithubAPI) CreateRepo(ctx context.Context, orgName, repoName string) (*github.Repository, error) {
	r := &github.Repository{
		Name: &repoName,
		// Description: description,
		Private:  github.Bool(true),
		AutoInit: github.Bool(false),
	}

	repo, _, err := api.client.Repositories.Create(ctx, orgName, r)
	return repo, err
}
