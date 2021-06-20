package logger

import (
	"time"

	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	"github.com/sirupsen/logrus"
)

const appNameKey = "app_name"

type Config struct {
	AppName string
}

func NewLogger(config *Config) log.MainLogger {
	impl := logrus.New()
	impl.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap:        fieldMap,
	})

	return &logger{
		FieldLogger: impl.WithField(appNameKey, config.AppName),
	}
}

type logger struct {
	logrus.FieldLogger
}

var fieldMap = logrus.FieldMap{
	logrus.FieldKeyTime: "@timestamp",
	logrus.FieldKeyMsg:  "message",
}

func (l *logger) WithField(key string, value interface{}) log.Logger {
	return &logger{l.FieldLogger.WithField(key, value)}
}

func (l *logger) WithFields(fields log.Fields) log.Logger {
	return &logger{l.FieldLogger.WithFields(logrus.Fields(fields))}
}

func (l *logger) Error(err error, args ...interface{}) {
	l.FieldLogger.WithError(err).Error(args...)
}

func (l *logger) FatalError(err error, args ...interface{}) {
	l.FieldLogger.WithError(err).Fatal(args...)
}
