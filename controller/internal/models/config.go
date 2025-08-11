package models

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost     string `envconfig:"POSTGRES_HOST"`
	DBPort     int    `envconfig:"POSTGRES_PORT"`
	DBUser     string `envconfig:"POSTGRES_USER"`
	DBPassword string `envconfig:"POSTGRES_PASSWORD"`
	DBName     string `envconfig:"POSTGRES_DB_NAME"`
	Port       int    `envconfig:"CONTROLLER_PORT"`
}

func NewConfig() (*Config, error) {
	cnfg := &Config{}
	err := envconfig.Process("", cnfg)
	if err != nil {
		return nil, err
	}
	return cnfg, nil
}
