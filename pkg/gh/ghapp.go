package gh

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v30/github"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/phunki/actionspanel/pkg/config"
	"github.com/phunki/actionspanel/pkg/log"
	"github.com/pkg/errors"
)

// NewGitHubClientCreator creates a GitHub client creator
func NewGitHubClientCreator(cfg config.Config) githubapp.ClientCreator {
	githubConfig := NewGitHubConfig(cfg)
	githubClientCreator, _ := githubapp.NewDefaultCachingClientCreator(githubConfig,
		githubapp.WithClientUserAgent("actionspanel/0.0.0"),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
	)

	return githubClientCreator
}

// InstallationHandler is a webhook handler for dealing with installation
// events from the GitHub v3 API
// More info:
// https://developer.github.com/v3/apps/installations/
type InstallationHandler struct {
	githubapp.ClientCreator
}

// NewInstallationHandler creates a new installation handler
func NewInstallationHandler(clientCreator githubapp.ClientCreator) InstallationHandler {
	return InstallationHandler{
		ClientCreator: clientCreator,
	}
}

// Handles is an array of event types that this handler can respond to
func (h InstallationHandler) Handles() []string {
	return []string{"installation"}
}

// Handle is the handler that reacts to an installation webhook event
func (h InstallationHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.InstallationEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse installation event")
	}

	log.Info("Received an installation event")

	return nil
}
