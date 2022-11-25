package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

type Async struct {
	Run      bool
	Interval Duration
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
	Timeout                  Duration
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

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %v: %w", filePath, err)
	}

	byteValue, _ := io.ReadAll(file)

	if err = file.Close(); err != nil {
		return fmt.Errorf("error closing file %v: %w", filePath, err)
	}

	err = json.Unmarshal(byteValue, c)
	if err != nil {
		return fmt.Errorf("error unmarshaling file %v: %w", filePath, err)
	}

	return nil
}

// Duration allows to unmarshal time into time.Duration
// https://github.com/golang/go/issues/10275
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var v interface{}
	if err = json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
	case string:
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
	return nil
}
