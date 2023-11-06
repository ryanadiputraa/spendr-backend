package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
}

type logger struct {
	log *zap.SugaredLogger
}

func NewLogger() (Logger, error) {
	log, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	s := log.Sugar()
	return &logger{log: s}, nil
}

func (l *logger) Info(v ...any) {
	l.log.Info(v...)
}

func (l *logger) Infow(msg string, keyAndValues ...any) {
	l.log.Infow(msg, keyAndValues)
}

func (l *logger) Warn(v ...any) {
	l.log.Warn(v...)
}

func (l *logger) Error(v ...any) {
	l.log.Error(v...)
}

func (l *logger) Fatal(v ...any) {
	l.log.Fatal(v...)
}
