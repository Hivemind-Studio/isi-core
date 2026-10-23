package utils

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger() *Logger {
	logger := &logrus.Logger{
		Out:          logrus.StandardLogger().Out,
		Formatter:    &CustomFormatter{},
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     logrus.StandardLogger().ExitFunc,
		ReportCaller: false,
	}

	return &Logger{logger: logger}
}

func (l *Logger) Request(method, path string, headers, query string) {
	l.logger.Info(fmt.Sprintf("Request: %s %s at %s\nHeaders: %s\nQuery: %s", method, path, time.Now().Format(time.RFC3339), headers, query))
}

func (l *Logger) Response(method, path string, status int, duration time.Duration) {
	l.logger.Info(fmt.Sprintf("Response: %s %s %d in %s at %s", method, path, status, duration, time.Now().Format(time.RFC3339)))
}

func (l *Logger) Success(message string, args ...interface{}) {
	l.logger.Infof(message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Errorf(message, args...)
}

func (l *Logger) Warning(message string, args ...interface{}) {
	l.logger.Warnf(message, args...)
}
