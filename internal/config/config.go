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
		"localhost:8080",
	},
	Option{
		"Output link host",
		"b",
		"http://localhost:8080",
	},
}

func (o *Option) String() string {
	return o.Value
}

func (o *Option) Set(s string) error {
	o.Value = s
	return nil
}
