package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-mongo-restapi/config"
	"github.com/sergicanet9/scv-go-framework/v2/api/utils"
	infrastructure "github.com/sergicanet9/scv-go-framework/v2/infrastructure/mongo"
)

// SetHealthRoutes creates health routes
func SetHealthRoutes(cfg config.Config, r *mux.Router) {
	r.Handle("/api/health", healthCheck(cfg)).Methods(http.MethodGet)
}

// @Summary Health Check
// @Description Runs a Health Check
// @Tags Health
// @Success 200 "OK"
// @Router /api/health [get]
func healthCheck(cfg config.Config) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		if err := run(r, cfg); err != nil {
			utils.ResponseError(w, r, http.StatusServiceUnavailable, err.Error())
		}
		utils.ResponseJSON(w, r, http.StatusOK, nil)
	})
}

func run(r *http.Request, cfg config.Config) error {
	r.Header.Add("Environment", cfg.Environment)
	r.Header.Add("Version", cfg.Version)

	if _, err := infrastructure.ConnectMongoDB(cfg.DBName, cfg.DBConnectionString); err != nil {
		return err
	}
	return nil
}
