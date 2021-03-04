package scm

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v33/github"
)

var (
	Client *GitHubClient
)

type GitHubClient struct {
	client 		*github.Client
	token 		string
}

func NewGitHubClient(token string) bool {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	Client = &GitHubClient{
		client,
		token,
	}

	return true
}

func (c *GitHubClient) GetRepositories() ([]string, error) {
	ctx := context.Background()
	repos, _, err := c.client.Repositories.List(ctx, "", nil)
	if err != nil {
		return nil, fmt.Errorf("Get user's project from github error:  %v", err)
	}

	repositories := []string{}
	for _, repo := range repos {
		repositories = append(repositories, *repo.Name)
	}
	return repositories, nil
}