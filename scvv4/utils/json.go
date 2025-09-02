package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// LoadJSON opens the specified file and unmarshals its JSON content in the received struct
func LoadJSON(filePath string, target interface{}) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("ignoring config file %v: %w", filePath, err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %v: %w", filePath, err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file %v: %w", filePath, err)
	}

	err = json.Unmarshal(byteValue, target)
	if err != nil {
		return fmt.Errorf("error unmarshaling file %v: %w", filePath, err)
	}

	return nil
}

// Duration allows to unmarshal time into time.Duration
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var v interface{}
	json.Unmarshal(b, &v)

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
		return fmt.Errorf("invalid duration")
	}
	return nil
}
