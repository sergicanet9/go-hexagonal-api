package observability

import (
	"log"
	"os"
	"sync"
)

var (
	logger = log.New(os.Stdout, "", log.Default().Flags())
	m      sync.Mutex
)

// Logger returns the singleton logger instance
func Logger() *log.Logger {
	return logger
}
