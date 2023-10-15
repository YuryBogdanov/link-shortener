package config

import (
	"flag"
	"os"
)

const (
	serverAddressKey = "SERVER_ADDRESS"
	baseURLKey       = "BASE_URL"
)

func SetupFlags() {
	configFlags := flag.NewFlagSet("Configuration flag set", flag.ContinueOnError)

	serverPath := obtainServerPath(configFlags)
	baseURL := obtainBaseURL(configFlags)

	configFlags.Parse(os.Args[1:])

	if configFlags.Parsed() {
		BaseConfig.ServerPath.Value = *serverPath
		BaseConfig.ShoretnedBaseURL.Value = *baseURL
	}
}

func obtainServerPath(flags *flag.FlagSet) *string {
	var serverPath *string
	if path := os.Getenv(serverAddressKey); path != "" {
		serverPath = &path
	} else {
		serverPath = flags.String(BaseConfig.ServerPath.CliParameterName, BaseConfig.ServerPath.Value, "Server path")
	}
	return serverPath
}

func obtainBaseURL(flags *flag.FlagSet) *string {
	var baseURL *string
	if path := os.Getenv(baseURLKey); path != "" {
		baseURL = &path
	} else {
		baseURL = flags.String(BaseConfig.ShoretnedBaseURL.CliParameterName, BaseConfig.ShoretnedBaseURL.Value, "Base url for shortened links")
	}
	return baseURL
}
