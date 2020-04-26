// Package config holds configuration that might need to be injected into various parts of the application
package config

import (
	"github.com/kelseyhightower/envconfig"
)

// DefaultCookieSessionSecret is the default value used for an unset cookie secret
const DefaultCookieSessionSecret = "defaultcookiesessionsecret"

// Config includes all configuration values used by actionspanel
type Config struct {
	// The port to run the web server on
	ServerPort int `required:"true" default:"8080" envconfig:"server_port"`
	// The port to run healthchecks on
	HealthServerPort int `required:"true" default:"8081" envconfig:"health_server_port"`

	// GitHub config
	OauthClientID     string `required:"true" split_words:"true"`
	OauthClientSecret string `required:"true" split_words:"true"`
	IntegrationID     int64  `required:"true" split_words:"true"`
	WebhookSecret     string `split_words:"true"`
	PrivateKey        string `required:"true" split_words:"true"`

	// Session management config
	SessionManagerType  string `required:"true" default:"cookie" envconfig:"session_manager_type"`
	CookieSessionSecret string `envconfig:"session_manager_type"`
}

// NewConfig scans environment variables for values
// This function will panic if the configuration can't be read successfully
func NewConfig() Config {
	var cfg Config

	err := envconfig.Process("ap", &cfg)
	if err != nil {
		panic("couldn't load config")
	}

	return cfg
}
