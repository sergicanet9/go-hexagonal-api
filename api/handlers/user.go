package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-mongo-restapi/business/user"
	"github.com/sergicanet9/go-mongo-restapi/config"
	"github.com/sergicanet9/go-mongo-restapi/models/requests"
	"github.com/sergicanet9/scv-go-framework/v2/api/utils"
)

// SetUserRoutes creates user routes
func SetUserRoutes(ctx context.Context, cfg config.Config, r *mux.Router, s user.UserService) {
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
// @Router /api/users/login [post]
func loginUser(ctx context.Context, cfg config.Config, s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var credentials requests.LoginUser
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		response, err := s.Login(ctx, credentials)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusOK, response)
	})
}

// @Summary Create user
// @Description Creates a new user
// @Tags Users
// @Param user body requests.User true "New user to be created"
// @Success 201 {object} responses.Creation "OK"
// @Router /api/users [post]
func createUser(ctx context.Context, cfg config.Config, s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var user requests.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		result, err := s.Create(ctx, user)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusCreated, result)
	})
}

// @Summary Get all users
// @Description Gets all the users
// @Tags Users
// @Security Bearer
// @Success 200 {array} responses.User "OK"
// @Router /api/users [get]
func getAllUsers(ctx context.Context, cfg config.Config, s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		users, err := s.GetAll(ctx)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusOK, users)
	})
}

// @Summary Get user by email
// @Description Gets a user by email
// @Tags Users
// @Security Bearer
// @Param email path string true "Email"
// @Success 200 {object} responses.User "OK"
// @Router /api/users/email/{email} [get]
func getUserByEmail(ctx context.Context, cfg config.Config, s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		user, err := s.GetByEmail(ctx, params["email"])
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusOK, user)
	})
}

// @Summary Get user by ID
// @Description Gets a user by ID
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 {object} responses.User "OK"
// @Router /api/users/{id} [get]
func getUserByID(ctx context.Context, cfg config.Config, s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		user, err := s.GetByID(ctx, params["id"])
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusOK, user)
	})
}

// @Summary Update user
// @Description Updates a user
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Param User body requests.UpdateUser true "User"
// @Success 200 "OK"
// @Router /api/users/{id} [patch]
func updateUser(ctx context.Context, cfg config.Config, s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		var user requests.UpdateUser
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		err = s.Update(ctx, params["id"], user)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusOK, nil)
	})
}

// @Summary Delete user
// @Description Delete a user
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 "OK"
// @Router /api/users/{id} [delete]
func deleteUser(ctx context.Context, cfg config.Config, s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		err := s.Delete(ctx, params["id"])
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusOK, nil)
	})
}

// @Summary Get claims
// @Description Gets all claims
// @Tags Users
// @Security Bearer
// @Success 200 {object} object "OK"
// @Router /api/claims [get]
func getUserClaims(ctx context.Context, cfg config.Config, s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		claims, err := s.GetClaims(ctx)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusOK, claims)
	})
}

// @Summary Atomic transaction proof
// @Description Inserts two users atomically
// @Tags Users
// @Security Bearer
// @Success 200 "OK"
// @Router /api/users/atomic [post]
func atomicTransactionProof(ctx context.Context, cfg config.Config, s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		err := s.AtomicTransationProof(ctx)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusOK, nil)
	})
}
