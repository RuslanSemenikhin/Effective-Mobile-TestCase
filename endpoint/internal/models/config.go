package models

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port int `envconfig:"ENDPOINT_GRPC_PORT"`
}

func NewConfig() (*Config, error) {
	cnfg := &Config{}
	err := envconfig.Process("", cnfg)
	if err != nil {
		return nil, err
	}
	return cnfg, nil
}
