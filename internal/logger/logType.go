package logger

import (
	"log"
	"os"
)

type LogType string

const (
	LogTypeProd LogType = "prod"
	LogTypeTest LogType = "test"
)

func getCurrentLogType() LogType {
	typeAsString := os.Getenv("SHORTENER_ENVIRONMENT")
	switch typeAsString {
	case string(LogTypeProd):
		return LogTypeProd

	case string(LogTypeTest):
		return LogTypeTest
	}
	log.Println("unrecognized logger type, returning prod logger as default")
	return LogTypeProd
}
