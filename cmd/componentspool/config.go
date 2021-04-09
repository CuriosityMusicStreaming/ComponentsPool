package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func parseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	DatabaseUser     string `envconfig:"db_user" default:"root"`
	DatabasePassword string `envconfig:"db_password" default:"1234"`
	DatabaseHost     string `envconfig:"db_host" default:"componentspool-db"`
	DatabaseName     string `envconfig:"db_name" default:"ContentService"`

	MaxDatabaseConnections int `envconfig:"max_connections" default:"10"`
}
