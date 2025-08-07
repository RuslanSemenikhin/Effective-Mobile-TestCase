package models

import (
	"github.com/kelseyhightower/envconfig"
)

type DbConfig struct {
	DBHost     string `envconfig:"POSTGRES_HOST"`
	DBPort     int    `envconfig:"POSTGRES_PORT"`
	DBUser     string `envconfig:"POSTGRES_USER"`
	DBPassword string `envconfig:"POSTGRES_PASSWORD"`
	DBName     string `envconfig:"POSTGRES_DB_NAME"`
}

func NewDbConfig() (*DbConfig, error) {
	cnfg := &DbConfig{}
	err := envconfig.Process("", cnfg)
	if err != nil {
		return nil, err
	}
	return cnfg, nil
}
