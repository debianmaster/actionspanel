package gh

//go:generate mockgen -package mock -destination=../../mock/mock_gh_client.go . AppsService
//go:generate mockgen -package mock -destination=../../mock/mock_gh_client_creator.go . ClientCreator

import (
	"context"

	"github.com/google/go-github/v30/github"
)

// AppsService wraps a GitHub client AppService
type AppsService interface {
	ListUserInstallations(ctx context.Context, opt *github.ListOptions) ([]*github.Installation, *github.Response, error)
}

// GithubClient wraps a github client for testing
type GithubClient struct {
	Apps AppsService
}

// ClientCreator wraps go-githubapp's ClientCreator
type ClientCreator interface {
	NewTokenClient(token string) (*github.Client, error)
}
