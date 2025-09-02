package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type targetType struct {
	TestDuration1 Duration
	TestDuration2 Duration
}

// TestLoadJSON_Ok checks that LoadJSON returns the expected response and parses Duration object properly when everything goes as expected
func TestLoadJSON_Ok(t *testing.T) {
	// Arrange
	target := targetType{}

	_, filePath, _, _ := runtime.Caller(0)
	dir, err := os.MkdirTemp(filepath.Join(filePath, "../../.."), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	bytes := []byte(`{"TestDuration1":10,"TestDuration2":"10s"}`)
	err = os.WriteFile(file.Name(), bytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	expectedDuration1 := Duration{time.Duration(10)}
	duration2, _ := time.ParseDuration("10s")
	expectedDuration2 := Duration{duration2}

	// Act
	err = LoadJSON(file.Name(), &target)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedDuration1, target.TestDuration1)
	assert.Equal(t, expectedDuration2, target.TestDuration2)
}

// TestLoadJSON_NonExistentFile checks that LoadJSON returns an error when the specified file does not exist
func TestLoadJSON_NonExistentFile(t *testing.T) {
	// Arrange
	target := targetType{}
	expectedError := "ignoring config file : stat : no such file or directory"

	// Act
	err := LoadJSON("", &target)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestLoadJSON_FileNotAccessible checks that LoadJSON returns an error when the specified file is not accessible
func TestLoadJSON_FileNotAccessible(t *testing.T) {
	// Arrange
	target := targetType{}

	_, filePath, _, _ := runtime.Caller(0)
	dir, err := os.MkdirTemp(filepath.Join(filePath, "../../.."), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	err = file.Chmod(os.ModeExclusive)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := fmt.Sprintf("error opening file %s: open %s: permission denied", file.Name(), file.Name())

	// Act
	err = LoadJSON(file.Name(), &target)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestLoadJSON_FileIsADirectory checks that LoadJSON returns an error when the specified file is a directory
func TestLoadJSON_FileIsADirectory(t *testing.T) {
	// Arrange
	target := targetType{}

	_, filePath, _, _ := runtime.Caller(0)
	dir, err := os.MkdirTemp(filepath.Join(filePath, "../../.."), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	expectedError := fmt.Sprintf("error reading file %s: read %s: is a directory", dir, dir)

	// Act
	err = LoadJSON(dir, &target)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestLoadJSON_InvalidDurationStringFormat checks that LoadJSON returns an error when the file contains a string that cannot be parsed to a Duration object
func TestLoadJSON_InvalidDurationStringFormat(t *testing.T) {
	// Arrange
	target := targetType{}

	_, filePath, _, _ := runtime.Caller(0)
	dir, err := os.MkdirTemp(filepath.Join(filePath, "../../.."), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	bytes := []byte(`{"TestDuration1":"10secs"}`)
	err = os.WriteFile(file.Name(), bytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := fmt.Sprintf("error unmarshaling file %s: time: unknown unit \"secs\" in duration \"10secs\"", file.Name())

	// Act
	err = LoadJSON(file.Name(), &target)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestLoadJSON_InvalidDurationType checks that LoadJSON returns an error when the file contains a type that cannot be parsed to a Duration object
func TestLoadJSON_InvalidDurationType(t *testing.T) {
	// Arrange
	target := targetType{}

	_, filePath, _, _ := runtime.Caller(0)
	dir, err := os.MkdirTemp(filepath.Join(filePath, "../../.."), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	bytes := []byte(`{"TestDuration1":{"Test":1}}`)
	err = os.WriteFile(file.Name(), bytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := fmt.Sprintf("error unmarshaling file %s: invalid duration", file.Name())

	// Act
	err = LoadJSON(file.Name(), &target)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
