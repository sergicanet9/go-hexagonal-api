package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Env                string
	Address            string
	Port               int
	DbConnectionString string
	DbName             string
	JWTSecret          string
}

// ReadConfig from the configPath passed as an argument. If the config is empty, will use config/config.json
// if env is passed will load configuration file using the env as follows : config/config.{env}.json.
// A default value can be specified in the configuration and override it in the environment configuration.
func ReadConfig(env string, dir string) (Config, error) {
	var c Config
	c.Env = env
	configPath := path.Join(dir, "config")

	err := c.loadJSON(path.Join(configPath, "config.json"))
	if err != nil {
		return c, err
	}

	if env != "" {
		if err = c.loadJSON(path.Join(configPath, "config."+env+".json")); err != nil {
			return c, fmt.Errorf("error parsing environment configuration, %s", err)
		}
	}

	return c, nil
}

func (c *Config) loadJSON(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("ignoring config file %v: %w", filePath, err)
	}

	file, err := os.Open(filePath) //nolint:gosec
	if err != nil {
		return fmt.Errorf("error opening file %v: %w", filePath, err)
	}

	byteValue, _ := ioutil.ReadAll(file)

	if err = file.Close(); err != nil {
		return fmt.Errorf("error closing file %v: %w", filePath, err)
	}

	err = json.Unmarshal(byteValue, c)
	if err != nil {
		return fmt.Errorf("error unmarshaling file %v: %w", filePath, err)
	}

	return nil
}
