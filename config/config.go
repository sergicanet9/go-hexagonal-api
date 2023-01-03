package config

import (
	"fmt"
	"path"

	"github.com/sergicanet9/scv-go-tools/v3/api/utils"
)

type Async struct {
	Run      bool
	Interval utils.Duration
}

type Config struct {
	// set in flags
	Version     string
	Environment string
	Port        int
	Database    string
	DSN         string

	// set in json config files
	config
}

type config struct {
	Address               string
	PostgresMigrationsDir string
	JWTSecret             string
	Timeout               utils.Duration
	Async                 Async
}

// ReadConfig from the project´s JSON config files.
// Default values are specified in the default configuration file, config/config.json
// and can be overrided with values specified in the environment configuration files, config/config.{env}.json.
func ReadConfig(version, env string, port int, database, dsn string) (Config, error) {
	var c Config
	c.Version = version
	c.Environment = env
	c.Port = port
	c.Database = database
	c.DSN = dsn

	var cfg config
	configPath := "config"

	if err := utils.LoadJSON(path.Join(configPath, "config.json"), &cfg); err != nil {
		return c, fmt.Errorf("error parsing configuration, %s", err)
	}

	if err := utils.LoadJSON(path.Join(configPath, "config."+env+".json"), &cfg); err != nil {
		return c, fmt.Errorf("error parsing environment configuration, %s", err)
	}

	c.config = cfg

	return c, nil
}
