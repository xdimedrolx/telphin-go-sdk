package telphin

import (
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logrus.FieldLogger
}

func WrapLogrus(logger logrus.FieldLogger) *LogrusLogger {
	return &LogrusLogger{logger}
}

// WithField returns a new Logger with the field added
func (l LogrusLogger) WithField(s string, i interface{}) FieldLogger {
	return LogrusLogger{l.FieldLogger.WithField(s, i)}
}

// WithFields returns a new Logger with the fields added
func (l LogrusLogger) WithFields(m map[string]interface{}) FieldLogger {
	return LogrusLogger{l.FieldLogger.WithFields(m)}
}
