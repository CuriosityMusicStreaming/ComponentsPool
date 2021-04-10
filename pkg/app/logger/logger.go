package logger

import "github.com/sirupsen/logrus"

type Fields logrus.Fields

type Logger interface {
	WithField(string, interface{}) Logger
	WithFields(Fields) Logger
	Info(...interface{})
	Error(error, ...interface{})
}

type MainLogger interface {
	Logger
	FatalError(error, ...interface{})
}
