package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scanet9/go-mongo-restapi/business"
	"github.com/scanet9/go-mongo-restapi/config"
	"github.com/scanet9/go-mongo-restapi/models/entities"
	"github.com/scanet9/go-mongo-restapi/models/requests"
	"github.com/scanet9/scv-go-framework/v2/api/utils"
)

// SetUserRoutes creates user routes
func SetUserRoutes(r *mux.Router, s business.UserService) {
	r.Handle("/api/users/login", loginUser(s)).Methods("POST")
	r.Handle("/api/users", createUser(s)).Methods("POST")
	r.Handle("/api/users", utils.JWTMiddleware(getAllUsers(s), config.JWTSecret)).Methods("GET")
	r.Handle("/api/users/email/{email}", utils.JWTMiddleware(getUserByEmail(s), config.JWTSecret)).Methods("GET")
	r.Handle("/api/users/{id}", utils.JWTMiddleware(getUserByID(s), config.JWTSecret)).Methods("GET")
	r.Handle("/api/users/{id}", utils.JWTMiddleware(updateUser(s), config.JWTSecret)).Methods("PUT")
	r.Handle("/api/users/{id}", utils.JWTMiddleware(deleteUser(s), config.JWTSecret)).Methods("DELETE")
}

// @Summary Login user
// @Description Logs in an user
// @Tags Users
// @Param login body requests.Login true "Login request"
// @Success 200 {object} responses.Login "OK"
// @Router /api/users/login [post]
func loginUser(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var credentials requests.Login
		_ = json.NewDecoder(r.Body).Decode(&credentials)
		token := s.Login(credentials)
		utils.ResponseJSON(w, http.StatusCreated, token)
	})
}

// @Summary Create user
// @Description Creates a new user
// @Tags Users
// @Param user body entities.User true "New user to be created"
// @Success 200 {object} responses.Creation "OK"
// @Router /api/users [post]
func createUser(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var user entities.User
		_ = json.NewDecoder(r.Body).Decode(&user)
		insertedID := s.Create(user)
		utils.ResponseJSON(w, http.StatusCreated, insertedID)
	})
}

// @Summary Get all users
// @Description Gets all the users
// @Tags Users
// @Security Bearer
// @Success 200 {array} entities.User "OK"
// @Router /api/users [get]
func getAllUsers(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		users := s.GetAll()
		utils.ResponseJSON(w, http.StatusOK, users)
	})
}

// @Summary Get user by email
// @Description Get a user by email
// @Tags Users
// @Security Bearer
// @Param email path string true "Email"
// @Success 200 {object} entities.User "OK"
// @Router /api/users/email/{email} [get]
func getUserByEmail(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		user := s.GetByEmail(params["email"])
		utils.ResponseJSON(w, http.StatusOK, user)
	})
}

// @Summary Get user by ID
// @Description Get a user by ID
// @Tags Users
// @Security Bearer
// @Param ID path string true "ID"
// @Success 200 {object} entities.User "OK"
// @Router /api/users/{id} [get]
func getUserByID(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		users := s.GetByID(params["id"])
		utils.ResponseJSON(w, http.StatusOK, users)
	})
}

// @Summary Update user
// @Description Update a user
// @Tags Users
// @Security Bearer
// @Param ID path string true "ID"
// @Param User body entities.User true "User"
// @Success 200 "OK"
// @Router /api/users/{id} [put]
func updateUser(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		var user entities.User
		_ = json.NewDecoder(r.Body).Decode(&user)
		s.Update(params["id"], user)
		utils.ResponseJSON(w, http.StatusOK, nil)
	})
}

// @Summary Delete user
// @Description Delete a user
// @Tags Users
// @Security Bearer
// @Param ID path string true "ID"
// @Success 200 "OK"
// @Router /api/users/{id} [delete]
func deleteUser(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		s.Delete(params["id"])
		utils.ResponseJSON(w, http.StatusOK, nil)
	})
}
