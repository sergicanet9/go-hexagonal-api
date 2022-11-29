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
	// set by flags
	Version     string
	Environment string
	Port        int
	Database    string

	// set in config files
	MongoAddress             string
	PostgresAddress          string
	MongoConnectionString    string
	MongoDBName              string
	PostgresConnectionString string
	PostgresMigrationsDir    string
	JWTSecret                string
	Timeout                  utils.Duration
	Async                    Async
}

// ReadConfig from the configPath passed as an argument. If the config is empty, will use config/config.json
// if env is passed will load configuration file using the env as follows : config/config.{env}.json.
// A default value can be specified in the configuration and override it in the environment configuration.
func ReadConfig(version, env string, port int, database string, dir string) (Config, error) {
	var c Config
	c.Version = version
	c.Environment = env
	c.Port = port
	c.Database = database
	configPath := path.Join(dir, "config")

	err := utils.LoadJSON(path.Join(configPath, "config.json"), &c)
	if err != nil {
		return c, err
	}

	if env != "" {
		if err = utils.LoadJSON(path.Join(configPath, "config."+env+".json"), &c); err != nil {
			return c, fmt.Errorf("error parsing environment configuration, %s", err)
		}
	}

	return c, nil
}
