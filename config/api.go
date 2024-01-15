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
		Name    string `env:"NAME" yaml:"name" env-default:"Glame"`
		Version string `env:"VERSION" yaml:"version" env-default:"0.1"`
		Port    int    `env:"PORT" yaml:"port" env-default:"5006"`
		Debug   bool   `env:"DEBUG" yml:"debug" env-default:"false"`
	}

	Auth struct {
		Secret   string `env:"SECRET" yaml:"secret" env-default:"L3AEyZt56nGKbJuHnTeTwfV8Gy02j8v9zJnhFCRkD8qaHVb9khXfTBDWgUr2Bqc1"`
		Password string `env:"PASSWORD" yaml:"password" env-default:"glame_password"`
	}
)

func NewApiConfig(path string) (*ApiConfig, error) {
	cfg := &ApiConfig{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		if err = cleanenv.ReadEnv(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
