package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	log, err := config.Build(zap.AddCallerSkip(1))
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
