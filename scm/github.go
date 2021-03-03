package scm

import (
	"context"

	"github.com/google/go-github/v33/github"
)

var (
	Client *GitHubClient
)

type GitHubClient struct {
	client 		*github.Client
	password 	string
}

func NewGitHubClient(username,password string) bool {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())

	Client = &GitHubClient{
		client,
		password,
	}

	return true
}

func (c *GitHubClient) GetRepositories() ([]string, error) {
	ctx := context.Background()
	repos, _, err := c.client.Repositories.List(ctx, "", nil)
	if listPorjsErr != nil {
		return nil, fmt.Errorf("Get user's project from gitlab error:  %v", listPorjsErr)
	}

	repositories := []string{}
	for _, repo := range repos {
		repositories = append(repositories, repo.Name)
	}
	return repositories, nil
}