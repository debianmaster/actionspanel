// Package log is where we abstract away log functionality
package log

import (
	"github.com/rs/zerolog/log"
)

// Info prints a message at the info level
func Info(format string) {
	log.Info().Msgf(format)
}

// Infof prints a message at the info level with a formatted message
func Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

// Err prints a message at the error level
func Err(err error, format string) {
	log.Err(err).Msgf(format)
}

// Errf prints a message at the error level with a formatted message
func Errf(err error, format string, v ...interface{}) {
	log.Err(err).Msgf(format, v...)
}
