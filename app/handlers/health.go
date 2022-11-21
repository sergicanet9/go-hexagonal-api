package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/scv-go-tools/v3/api/middlewares"
	"github.com/sergicanet9/scv-go-tools/v3/api/utils"
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
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Environment", cfg.Environment)
		r.Header.Add("Database", cfg.Database)
		r.Header.Add("Version", cfg.Version)

		utils.ResponseJSON(w, r, nil, http.StatusOK, nil)
	})
}
