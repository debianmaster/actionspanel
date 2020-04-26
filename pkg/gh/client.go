package gh

//go:generate mockgen -package mock -destination=../../mock/mock_gh_client.go github.com/phunki/actionspanel/pkg/gh AppsService

import (
	"context"

	"github.com/google/go-github/v30/github"
)

// AppsService wraps a GitHub client AppService
type AppsService interface {
	ListUserInstallations(ctx context.Context, opt *github.ListOptions) ([]*github.Installation, *github.Response, error)
}
