package logger

import (
	"log"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func Setup() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}

	sugar = *logger.Sugar()
}

func Finish() {
	sugar.Desugar().Sync()
}

func Info(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	sugar.Fatalw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}
