package config

import (
	"flag"
	"os"
)

func SetupFlags() {
	configFlags := flag.NewFlagSet("Configuration flag set", flag.ContinueOnError)

	serverPath := configFlags.String(BaseConfig.ServerPath.CliParameterName, BaseConfig.ServerPath.Value, "Server path")
	baseUrl := configFlags.String(BaseConfig.ShoretnedBaseURL.CliParameterName, BaseConfig.ShoretnedBaseURL.Value, "Base URL for shortened links")

	configFlags.Parse(os.Args[1:])

	if configFlags.Parsed() {
		BaseConfig.ServerPath.Value = *serverPath
		BaseConfig.ShoretnedBaseURL.Value = *baseUrl
	}
}
