package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DatabaseUrl string `env:"DATABASE_URL"`
}

func New() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
