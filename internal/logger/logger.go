package logger

import (
	"log"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger
var logType LogType

func Setup() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}

	logType = getCurrentLogType()
	sugar = *logger.Sugar()
}

func Finish() {
	sugar.Desugar().Sync()
}

func Info(msg string, keysAndValues ...interface{}) {
	prefixedMessage := getPrefixedLogMessage(msg)
	sugar.Infow(prefixedMessage, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	prefixedMessage := getPrefixedLogMessage(msg)
	sugar.Fatalw(prefixedMessage, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	prefixedMessage := getPrefixedLogMessage(msg)
	sugar.Errorw(prefixedMessage, keysAndValues...)
}

func getPrefixedLogMessage(msg string) string {
	switch logType {
	case LogTypeProd:
		return "[PROD Log] " + msg

	case LogTypeTest:
		return "[TEST Log] " + msg
	}
	return "[Unknown Log Type] " + msg
}
