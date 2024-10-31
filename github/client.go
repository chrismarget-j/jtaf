package github

import (
	"net/http"
	"os"

	"github.com/google/go-github/v66/github"
)

const envGithubToken = "GITHUB_PUB_API_TOKEN"

func GithubClient(httpClient *http.Client) *github.Client {
	ghc := github.NewClient(httpClient)

	if token, ok := os.LookupEnv(envGithubToken); ok {
		ghc = ghc.WithAuthToken(token)
	}

	return ghc
}
