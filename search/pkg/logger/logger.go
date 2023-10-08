package logger

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()

	log.SetReportCaller(true)
}

func WithFields(logger *logrus.Logger, fields logrus.Fields) *logrus.Logger {
	return logger.WithFields(fields).Logger
}

func GetLogger() *logrus.Logger {
	return log
}
