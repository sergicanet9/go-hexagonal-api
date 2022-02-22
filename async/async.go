package async

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sergicanet9/go-mongo-restapi/config"
	"go.mongodb.org/mongo-driver/mongo"
)

const contentType = "application/json"

type Async struct {
	config config.Config
	db     *mongo.Database
}

func NewAsync(cfg config.Config, db *mongo.Database) *Async {
	return &Async{
		config: cfg,
		db:     db,
	}
}

func (a Async) Run() {
	for {
		errs := make(chan error)
		go func() {
			errs <- healthCheck(a.config.Address, a.config.Port)
		}()
		for i := 0; i < 1; i++ {
			if err := <-errs; err != nil {
				log.Printf("async process failure, error: %s", err)
			}
		}

		<-time.After(a.config.Async.Interval.Duration)
	}
}

func healthCheck(address string, port int) error {
	start := time.Now()

	url := fmt.Sprintf("%s:%d/api/health", address, port)

	req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	elapsed := time.Since(start)
	log.Printf("Health Check complete, time elapsed: %s", elapsed)

	return err
}
