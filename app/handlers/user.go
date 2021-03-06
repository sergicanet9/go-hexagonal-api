package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/dto/requests"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v3/api/utils"
)

// SetUserRoutes creates user routes
func SetUserRoutes(ctx context.Context, cfg config.Config, r *mux.Router, s ports.UserService) {
	r.Handle("/api/users/login", loginUser(ctx, cfg, s)).Methods(http.MethodPost)
	r.Handle("/api/users", createUser(ctx, cfg, s)).Methods(http.MethodPost)
	r.Handle("/api/users", utils.JWTMiddleware(getAllUsers(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/api/users/email/{email}", utils.JWTMiddleware(getUserByEmail(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/api/users/{id}", utils.JWTMiddleware(getUserByID(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/api/users/{id}", utils.JWTMiddleware(updateUser(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodPatch)
	r.Handle("/api/users/{id}", utils.JWTMiddleware(deleteUser(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{"admin": true})).Methods(http.MethodDelete)
	r.Handle("/api/claims", utils.JWTMiddleware(getUserClaims(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/api/users/atomic", utils.JWTMiddleware(atomicTransactionProof(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodPost)
}

// @Summary Login user
// @Description Logs in an user
// @Tags Users
// @Param login body requests.LoginUser true "Login request"
// @Success 200 {object} responses.LoginUser "OK"
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /api/users/login [post]
func loginUser(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, http.StatusBadRequest, err.Error())
			return
		}

		var credentials requests.LoginUser
		err = json.Unmarshal(body, &credentials)
		if err != nil {
			utils.ResponseError(w, r, body, http.StatusBadRequest, err.Error())
			return
		}

		response, err := s.Login(ctx, credentials)
		if err != nil {
			utils.ResponseError(w, r, body, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusOK, response)
	})
}

// @Summary Create user
// @Description Creates a new user
// @Tags Users
// @Param user body requests.User true "New user to be created"
// @Success 201 {object} responses.Creation "OK"
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /api/users [post]
func createUser(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, http.StatusBadRequest, err.Error())
			return
		}

		var user requests.User
		err = json.Unmarshal(body, &user)
		if err != nil {
			utils.ResponseError(w, r, body, http.StatusBadRequest, err.Error())
			return
		}

		result, err := s.Create(ctx, user)
		if err != nil {
			utils.ResponseError(w, r, body, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusCreated, result)
	})
}

// @Summary Get all users
// @Description Gets all the users
// @Tags Users
// @Security Bearer
// @Success 200 {array} responses.User "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/users [get]
func getAllUsers(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		users, err := s.GetAll(ctx)
		if err != nil {
			utils.ResponseError(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, users)
	})
}

// @Summary Get user by email
// @Description Gets a user by email
// @Tags Users
// @Security Bearer
// @Param email path string true "Email"
// @Success 200 {object} responses.User "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/users/email/{email} [get]
func getUserByEmail(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		user, err := s.GetByEmail(ctx, params["email"])
		if err != nil {
			utils.ResponseError(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, user)
	})
}

// @Summary Get user by ID
// @Description Gets a user by ID
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 {object} responses.User "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/users/{id} [get]
func getUserByID[T string](ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		user, err := s.GetByID(ctx, params["id"])
		if err != nil {
			utils.ResponseError(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, user)
	})
}

// @Summary Update user
// @Description Updates a user
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Param User body requests.UpdateUser true "User"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/users/{id} [patch]
func updateUser[T string](ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, http.StatusBadRequest, err.Error())
			return
		}

		var params = mux.Vars(r)
		var user requests.UpdateUser
		err = json.Unmarshal(body, &user)
		if err != nil {
			utils.ResponseError(w, r, body, http.StatusBadRequest, err.Error())
			return
		}

		err = s.Update(ctx, params["id"], user)
		if err != nil {
			utils.ResponseError(w, r, body, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusOK, nil)
	})
}

// @Summary Delete user
// @Description Delete a user
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/users/{id} [delete]
func deleteUser[T string](ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		err := s.Delete(ctx, params["id"])
		if err != nil {
			utils.ResponseError(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, nil)
	})
}

// @Summary Get claims
// @Description Gets all claims
// @Tags Users
// @Security Bearer
// @Success 200 {object} object "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/claims [get]
func getUserClaims(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		claims, err := s.GetClaims(ctx)
		if err != nil {
			utils.ResponseError(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, claims)
	})
}

// @Summary Atomic transaction proof
// @Description Creates two users atomically
// @Tags Users
// @Security Bearer
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/users/atomic [post]
func atomicTransactionProof(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		err := s.AtomicTransationProof(ctx)
		if err != nil {
			utils.ResponseError(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, nil)
	})
}
