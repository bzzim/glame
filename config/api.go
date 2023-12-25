package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	ApiConfig struct {
		App  `yaml:"app"`
		Auth `yaml:"auth"`
	}

	App struct {
		Name    string `env-required:"true" env:"G_APP_NAME" yaml:"name"`
		Version string `env-required:"true" env:"G_APP_VERSION" yaml:"version"`
	}

	Auth struct {
		Secret   string `env-required:"true" env:"G_AUTH_PASSWORD" yaml:"secret"`
		Password string `env-required:"true" env:"G_AUTH_PASSWORD" yaml:"password"`
	}
)

func NewApiConfig(path string) (*ApiConfig, error) {
	cfg := &ApiConfig{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
