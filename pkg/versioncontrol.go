package pkg

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git"
	"github.com/go-git/go-git/config"
	"github.com/go-git/go-git/plumbing/object"
	"github.com/go-git/go-git/plumbing/transport/http"
	"github.com/google/go-github/github"
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
		},
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// namespace should be empty if it is not an org
func SetupGithubRemote(ctx context.Context, namespace string, repoName string, localRepo *git.Repository, push bool) error {
	// Create github repo

	// tokens := &oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")}
	// ts := oauth2.StaticTokenSource(tokens)
	// tc := oauth2.NewClient(ctx, ts)

	// client := github.NewClient(tc)

	// auth := &http.BasicAuth{
	// 	Username: "Darshan072",
	// 	Password: os.Getenv("GITHUB_TOKEN"),
	// }

	// newR, _, _ := client.Repositories.Get(ctx, "Darshan072", repoName)

	// if newR != nil {
	// 	return errors.New("repository already exists")
	// }

	// repo := &github.Repository{
	// 	Name: github.String(repoName),
	// }

	// result, _, err := client.Repositories.Create(ctx, namespace, repo)

	// if err != nil {
	// 	return err
	// }

	// repository := "https://github.com/" + "Darshan072" + "/" + repoName + ".git"
	// if result != nil {
	// 	repository = *result.CloneURL
	// }
	// // localRepo.CreateRemote(&)

	// newRepo, err := locak.PlainClone(newClonePath, false, &git.CloneOptions{
	// 	URL:        repository,
	// 	RemoteName: "origin/main",
	// 	Auth:       auth,
	// })

	// if err != nil {
	// 	return err
	// }

	// err = localRepo.Push(&git.PushOptions{
	// 	RemoteName: "origin/main",
	// 	Auth:       auth,
	// })

	// if err != nil {
	// 	log.Printf("ERR: final commit newRepo push event")
	// 	return err
	// }

	/////////

	ghRepo, err := NewGithubAPI().CreateRepo(ctx, namespace, repoName)
	if err != nil {
		return err
	}

	// Add github as remote to local repo
	_, err = localRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{ghRepo.GetCloneURL()},
	})

	if err != nil {
		return err
	}

	if push {

		// Push to local repo to github
		err = localRepo.Push(&git.PushOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username: "Darshan072",
				Password: os.Getenv("GITHUB_TOKEN"),
			},
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
