package logger

import (
	"github.com/sirupsen/logrus"
)

// New creates a new logger instance
func New(level string) *logrus.Logger {
	logger := logrus.New()

	// Parse log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	logger.SetLevel(logLevel)
	return logger
}
