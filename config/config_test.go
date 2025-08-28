package config

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReadConfig_Ok checks that ReadConfig returns the expected config struct
func TestReadConfig_Ok(t *testing.T) {
	// Arrange
	_, filePath, _, _ := runtime.Caller(0)

	expectedConfig := Config{
		Environment: "local",
	}

	// Act
	cfg, err := ReadConfig("", expectedConfig.Environment, 0, "", "", "", path.Join(path.Dir(filePath)))

	// Assert
	assert.NotEmpty(t, cfg)
	assert.Equal(t, expectedConfig.Environment, cfg.Environment)
	assert.Nil(t, err)
}

// TestReadConfig_NonExistentConfigPath checks that ReadConfig returns the expected error when the config files does not exist in the provided path
func TestReadConfig_NonExistentConfigPath(t *testing.T) {
	// Arrange
	invalidPath := "invalid/path"
	expectedError := fmt.Sprintf("error parsing configuration, ignoring config file %s/config.json: stat %s/config.json: no such file or directory", invalidPath, invalidPath)

	// Act
	_, err := ReadConfig("", "", 0, "", "", "", invalidPath)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestReadConfig_NonExistentEnvironmentFile checks that ReadConfig returns the expected error when the provided environment does not match with any config file
func TestReadConfig_NonExistentEnvironmentFile(t *testing.T) {
	// Arrange
	_, filePath, _, _ := runtime.Caller(0)

	expectedConfig := Config{
		Environment: "invalid-environment",
	}

	expectedError := fmt.Sprintf("error parsing environment configuration, ignoring config file %s/config.invalid-environment.json: stat %s/config.invalid-environment.json: no such file or directory", path.Join(path.Dir(filePath)), path.Join(path.Dir(filePath)))

	// Act
	_, err := ReadConfig("", expectedConfig.Environment, 0, "", "", "", path.Join(path.Dir(filePath)))

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
