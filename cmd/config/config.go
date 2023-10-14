package config

type Option struct {
	Name             string
	CliParameterName string
	Value            string
}

type Config struct {
	ServerPath       Option
	ShoretnedBaseURL Option
}

var BaseConfig Config = Config{
	Option{
		"Server Path",
		"a",
		"",
	},
	Option{
		"Output link host",
		"b",
		"",
	},
}

func (o *Option) String() string {
	return o.Value
}

func (o *Option) Set(s string) error {
	o.Value = s
	return nil
}
