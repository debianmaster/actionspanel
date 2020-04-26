package config_test

import (
	"testing"

	"github.com/phunki/actionspanel/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_NewSessionManagerFactory_PanicsOnInvalidType(t *testing.T) {
	cfg := config.Config{
		SessionManagerType: "invalid",
	}
	assert.Panics(t, func() {
		config.NewSessionManagerFactory(cfg)
	})
}
