package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger interface {
	Setup()
	Finish()

	Info(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type DefaultLogger struct {
	sugar zap.SugaredLogger
}

func (l *DefaultLogger) Setup() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}

	l.sugar = *logger.Sugar()
}

func (l *DefaultLogger) Finish() {
	l.sugar.Desugar().Sync()
}

func (l *DefaultLogger) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *DefaultLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

func (l *DefaultLogger) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}
