package models

import "github.com/kelseyhightower/envconfig"

type Config struct {
	PortEndpoint   int    `envconfig:"ENDPOINT_PORT"`
	HostController string `envconfig:"CONTROLLER_HOST"`
	PortController int    `envconfig:"CONTROLLER_PORT"`
}

func NewConfig() (*Config, error) {
	cnfg := &Config{}
	err := envconfig.Process("", cnfg)
	if err != nil {
		return nil, err
	}
	return cnfg, nil
}
