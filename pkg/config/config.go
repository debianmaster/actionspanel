package config

// Config includes all configuration values used by actionspanel
type Config struct {
	// The port to run the web server on
	ServerPort int `required:"true" default:"8080" envconfig:"server_port"`
	// The port to run healthchecks on
	HealthServerPort int `required:"true" default:"8081" envconfig:"health_server_port"`
}
