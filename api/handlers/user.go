package handlers

import (
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
func SetUserRoutes(cfg config.Config, r *mux.Router, s user.UserService) {
	r.Handle("/api/users/login", loginUser(s)).Methods(http.MethodPost)
	r.Handle("/api/users", createUser(s)).Methods(http.MethodPost)
	r.Handle("/api/users", utils.JWTMiddleware(getAllUsers(s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/api/users/email/{email}", utils.JWTMiddleware(getUserByEmail(s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/api/users/{id}", utils.JWTMiddleware(getUserByID(s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/api/users/{id}", utils.JWTMiddleware(updateUser(s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodPatch)
	r.Handle("/api/users/{id}", utils.JWTMiddleware(deleteUser(s), cfg.JWTSecret, jwt.MapClaims{"admin": true})).Methods(http.MethodDelete)
	r.Handle("/api/claims", utils.JWTMiddleware(getUserClaims(s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/api/users/atomic", utils.JWTMiddleware(atomicTransactionProof(s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodPost)
}

// @Summary Login user
// @Description Logs in an user
// @Tags Users
// @Param login body requests.LoginUser true "Login request"
// @Success 200 {object} responses.LoginUser "OK"
// @Router /api/users/login [post]
func loginUser(s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var credentials requests.LoginUser
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		response, err := s.Login(credentials)
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
func createUser(s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var user requests.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		result, err := s.Create(user)
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
func getAllUsers(s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		users, err := s.GetAll()
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
func getUserByEmail(s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		user, err := s.GetByEmail(params["email"])
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
func getUserByID(s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		user, err := s.GetByID(params["id"])
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
func updateUser(s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		var user requests.UpdateUser
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		err = s.Update(params["id"], user)
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
func deleteUser(s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		err := s.Delete(params["id"])
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
func getUserClaims(s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		claims := s.GetClaims()
		utils.ResponseJSON(w, r, http.StatusOK, claims)
	})
}

// @Summary Atomic transaction proof
// @Description Inserts two users atomically
// @Tags Users
// @Security Bearer
// @Success 200 "OK"
// @Router /api/users/atomic [post]
func atomicTransactionProof(s user.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		err := s.AtomicTransationProof()
		if err != nil {
			utils.ResponseError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		utils.ResponseJSON(w, r, http.StatusOK, nil)
	})
}
