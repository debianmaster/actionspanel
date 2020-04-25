package gh

import (
	"encoding/base64"

	"github.com/palantir/go-githubapp/githubapp"
	"github.com/phunki/actionspanel/pkg/config"
)

// NewGitHubConfig creates a GitHub config out of total config
func NewGitHubConfig(cfg config.Config) githubapp.Config {
	// Load GitHub client configuration
	base64dPrivateKey, _ := base64.StdEncoding.DecodeString(cfg.PrivateKey)
	githubConfig := githubapp.Config{
		WebURL:   "https://github.com",
		V3APIURL: "https://api.github.com",
		// We need to define inline structs here because of how they defined their data structures
		OAuth: struct {
			ClientID     string `yaml:"client_id" json:"clientId"`
			ClientSecret string `yaml:"client_secret" json:"clientSecret"`
		}{
			ClientID:     cfg.OauthClientID,
			ClientSecret: cfg.OauthClientSecret,
		},
		App: struct {
			IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
			WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
			PrivateKey    string `yaml:"private_key" json:"privateKey"`
		}{
			IntegrationID: cfg.IntegrationID,
			WebhookSecret: cfg.WebhookSecret,
			PrivateKey:    string(base64dPrivateKey),
		},
	}

	return githubConfig
}
