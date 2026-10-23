package utils

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// CustomFormatter is a custom logrus formatter that applies colors to log messages.
type CustomFormatter struct {
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = "\033[32m" // Green
	case logrus.ErrorLevel:
		levelColor = "\033[31m" // Red
	case logrus.WarnLevel:
		levelColor = "\033[33m" // Yellow
	default:
		levelColor = "\033[0m" // Default
	}

	timestamp := entry.Time.Format(time.RFC3339)
	message := fmt.Sprintf("%s[%s] %s\033[0m\n", levelColor, timestamp, entry.Message)

	// Check if Data is not empty and add it to the message
	if len(entry.Data) > 0 {
		message += fmt.Sprintf(" %v", entry.Data)
	}

	return []byte(message), nil
}
