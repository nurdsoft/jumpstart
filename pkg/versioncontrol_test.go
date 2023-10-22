package pkg

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
)

func Test_GithubAPI_GetAuthUser(t *testing.T) {
	godotenv.Load("../.env")

	user, err := NewGithubAPI().GetAuthUser(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", user)
}
