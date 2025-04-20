package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger() *Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &Logger{logger: logger}
}

func (l *Logger) Request(method, path string, headers, query string) {
	l.logger.Info().
		Str("type", "request").
		Str("method", method).
		Str("path", path).
		Str("headers", headers).
		Str("query", query).
		Time("timestamp", time.Now()).
		Msg("Incoming HTTP request")
}

func (l *Logger) Response(method, path string, status int, duration time.Duration) {
	l.logger.Info().
		Str("type", "response").
		Str("method", method).
		Str("path", path).
		Int("status", status).
		Dur("duration", duration).
		Time("timestamp", time.Now()).
		Msg("HTTP response sent")
}

func (l *Logger) Success(message string, args ...interface{}) {
	l.logger.Info().
		Msgf(message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Error().
		Msgf(message, args...)
}

func (l *Logger) Warning(message string, args ...interface{}) {
	l.logger.Warn().
		Msgf(message, args...)
}
