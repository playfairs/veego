package config

type Config struct {
	APIKey  string            `yaml:"api_key"`
	Devices map[string]Device `yaml:"devices"`
}

type Device struct {
	ID    string `yaml:"id"`
	Model string `yaml:"model"`
}
