package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/scv-go-tools/v3/api/utils"
)

type healthHandler struct {
	cfg config.Config
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(cfg config.Config) healthHandler {
	return healthHandler{
		cfg: cfg,
	}
}

// SetHealthRoutes creates health routes
func SetHealthRoutes(router *mux.Router, h healthHandler) {
	router.HandleFunc("/health", h.healthCheck).Methods(http.MethodGet)
}

// @Summary Health Check
// @Description Runs a Health Check
// @Tags Health
// @Success 200 "OK"
// @Failure 500 {object} object
// @Failure 503 {object} object
// @Router /health [get]
func (h *healthHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Version", h.cfg.Version)
	r.Header.Add("Environment", h.cfg.Environment)
	r.Header.Add("Port", strconv.Itoa(h.cfg.Port))
	r.Header.Add("Database", h.cfg.Database)

	if h.cfg.Environment == "local" {
		r.Header.Add("DSN", h.cfg.DSN)
	}

	utils.ResponseJSON(w, r, nil, http.StatusOK, nil)
}
