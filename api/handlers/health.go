package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-mongo-restapi/config"
	"github.com/sergicanet9/scv-go-framework/v2/api/utils"
	infrastructure "github.com/sergicanet9/scv-go-framework/v2/infrastructure/mongo"
)

// SetHealthRoutes creates health routes
func SetHealthRoutes(ctx context.Context, cfg config.Config, r *mux.Router) {
	r.Handle("/api/health", healthCheck(ctx, cfg)).Methods(http.MethodGet)
}

// @Summary Health Check
// @Description Runs a Health Check
// @Tags Health
// @Success 200 "OK"
// @Failure 500 {object} object
// @Failure 503 {object} object
// @Router /api/health [get]
func healthCheck(ctx context.Context, cfg config.Config) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		if err := run(r, ctx, cfg); err != nil {
			utils.ResponseError(w, r, nil, http.StatusServiceUnavailable, err.Error())
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, nil)
	})
}

func run(r *http.Request, ctx context.Context, cfg config.Config) error {
	r.Header.Add("Environment", cfg.Environment)
	r.Header.Add("Version", cfg.Version)

	db, err := infrastructure.ConnectMongoDB(ctx, cfg.DBName, cfg.DBConnectionString)
	if err != nil {
		return err
	}
	return db.Client().Ping(ctx, nil)
}
