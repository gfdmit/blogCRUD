package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres
}

type Postgres struct {
	Username   string `env:"POSTGRES_USER"`
	Password   string `env:"POSTGRES_PASSWORD"`
	Host       string `env:"POSTGRES_HOST"`
	Port       string `env:"POSTGRES_PORT"`
	DB         string `env:"POSTGRES_DB"`
	Migrations string `env:"POSTGRES_MIGRATIONS"`
}

func New(env string) (*Config, error) {
	conf := &Config{}

	if err := godotenv.Overload(env); err != nil {
		return nil, fmt.Errorf("godotenv.Overload: %v", err)
	}

	if err := cleanenv.ReadEnv(conf); err != nil {
		return nil, fmt.Errorf("cleanenv.Readenv: %v", err)
	}

	return conf, nil
}
