package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-hexagonal-api/config"

	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v3/api/middlewares"
	"github.com/sergicanet9/scv-go-tools/v3/api/utils"
)

type userHandler struct {
	cfg config.Config
	svc ports.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(cfg config.Config, svc ports.UserService) userHandler {
	return userHandler{
		cfg: cfg,
		svc: svc,
	}
}

// SetUserRoutes creates user routes
func SetUserRoutes(router *mux.Router, u userHandler) {
	router.Use(middlewares.Recover)

	router.HandleFunc("/v1/users/login", u.loginUser).Methods(http.MethodPost)
	router.HandleFunc("/v1/users", u.createUser).Methods(http.MethodPost)
	router.HandleFunc("/v1/users/many", u.createManyUsers).Methods(http.MethodPost)

	secureRouter := router.PathPrefix("").Subrouter()
	secureRouter.Use(middlewares.JWT(u.cfg.JWTSecret, jwt.MapClaims{}))

	secureRouter.HandleFunc("/v1/users", u.getAllUsers).Methods(http.MethodGet)
	secureRouter.HandleFunc("/v1/users/email/{email}", u.getUserByEmail).Methods(http.MethodGet)
	secureRouter.HandleFunc("/v1/users/{id}", u.getUserByID).Methods(http.MethodGet)
	secureRouter.HandleFunc("/v1/users/{id}", u.updateUser).Methods(http.MethodPatch)
	secureRouter.HandleFunc("/v1/claims", u.getUserClaims).Methods(http.MethodGet)

	adminRouter := router.PathPrefix("").Subrouter()
	adminRouter.Use(middlewares.JWT(u.cfg.JWTSecret, jwt.MapClaims{"admin": true}))

	adminRouter.HandleFunc("/v1/users/{id}", u.deleteUser).Methods(http.MethodDelete)
}

// @Summary Login user
// @Description Logs in an user
// @Tags Users
// @Param login body models.LoginUserReq true "Login request"
// @Success 200 {object} models.LoginUserResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/users/login [post]
func (u *userHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), u.cfg.Timeout.Duration)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}

	var credentials models.LoginUserReq
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}

	response, err := u.svc.Login(ctx, credentials)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}
	utils.ResponseJSON(w, r, body, http.StatusOK, response)
}

// @Summary Create user
// @Description Creates a new user
// @Tags Users
// @Param user body models.CreateUserReq true "New user to be created"
// @Success 201 {object} models.CreationResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/users [post]
func (u *userHandler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), u.cfg.Timeout.Duration)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}

	var user models.CreateUserReq
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}

	result, err := u.svc.Create(ctx, user)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}
	utils.ResponseJSON(w, r, body, http.StatusCreated, result)
}

// @Summary Create many users
// @Description Creates many users atomically
// @Tags Users
// @Param users body []models.CreateUserReq true "New users to be created"
// @Success 201 {object} models.MultiCreationResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/users/many [post]
func (u *userHandler) createManyUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), u.cfg.Timeout.Duration)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}

	var users []models.CreateUserReq
	err = json.Unmarshal(body, &users)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}

	result, err := u.svc.CreateMany(ctx, users)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}
	utils.ResponseJSON(w, r, body, http.StatusCreated, result)
}

// @Summary Get all users
// @Description Gets all the users
// @Tags Users
// @Security Bearer
// @Success 200 {array} models.UserResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/users [get]
func (u *userHandler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), u.cfg.Timeout.Duration)
	defer cancel()

	users, err := u.svc.GetAll(ctx)
	if err != nil {
		utils.ResponseError(w, r, nil, err)
		return
	}
	utils.ResponseJSON(w, r, nil, http.StatusOK, users)
}

// @Summary Get user by email
// @Description Gets a user by email
// @Tags Users
// @Security Bearer
// @Param email path string true "Email"
// @Success 200 {object} models.UserResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/users/email/{email} [get]
func (u *userHandler) getUserByEmail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), u.cfg.Timeout.Duration)
	defer cancel()

	var params = mux.Vars(r)
	user, err := u.svc.GetByEmail(ctx, params["email"])
	if err != nil {
		utils.ResponseError(w, r, nil, err)
		return
	}
	utils.ResponseJSON(w, r, nil, http.StatusOK, user)
}

// @Summary Get user by ID
// @Description Gets a user by ID
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 {object} models.UserResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/users/{id} [get]
func (u *userHandler) getUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), u.cfg.Timeout.Duration)
	defer cancel()

	var params = mux.Vars(r)
	user, err := u.svc.GetByID(ctx, params["id"])
	if err != nil {
		utils.ResponseError(w, r, nil, err)
		return
	}
	utils.ResponseJSON(w, r, nil, http.StatusOK, user)
}

// @Summary Update user
// @Description Updates a user
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Param User body models.UpdateUserReq true "User"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/users/{id} [patch]
func (u *userHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), u.cfg.Timeout.Duration)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}

	var params = mux.Vars(r)
	var user models.UpdateUserReq
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}

	err = u.svc.Update(ctx, params["id"], user)
	if err != nil {
		utils.ResponseError(w, r, body, err)
		return
	}
	utils.ResponseJSON(w, r, body, http.StatusOK, nil)
}

// @Summary Get claims
// @Description Gets all claims
// @Tags Users
// @Security Bearer
// @Success 200 {object} object "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/claims [get]
func (u *userHandler) getUserClaims(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), u.cfg.Timeout.Duration)
	defer cancel()

	claims := u.svc.GetUserClaims(ctx)
	utils.ResponseJSON(w, r, nil, http.StatusOK, claims)
}

// @Summary Delete user
// @Description Delete a user
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/users/{id} [delete]
func (u *userHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), u.cfg.Timeout.Duration)
	defer cancel()

	var params = mux.Vars(r)
	err := u.svc.Delete(ctx, params["id"])
	if err != nil {
		utils.ResponseError(w, r, nil, err)
		return
	}
	utils.ResponseJSON(w, r, nil, http.StatusOK, nil)
}
