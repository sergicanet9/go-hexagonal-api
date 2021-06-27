package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scanet9/go-mongo-restapi/business"
	"github.com/scanet9/go-mongo-restapi/config"
	"github.com/scanet9/go-mongo-restapi/models/entities"
	"github.com/scanet9/go-mongo-restapi/models/requests"
	"github.com/scanet9/scv-go-framework/api/utils"
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

func loginUser(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var credentials requests.Login
		_ = json.NewDecoder(r.Body).Decode(&credentials)
		token := s.Login(credentials)
		utils.ResponseJSON(w, http.StatusCreated, token)
	})
}

func createUser(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var user entities.User
		_ = json.NewDecoder(r.Body).Decode(&user)
		insertedID := s.Create(user)
		utils.ResponseJSON(w, http.StatusCreated, insertedID)
	})
}

func getAllUsers(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		users := s.GetAll()
		utils.ResponseJSON(w, http.StatusOK, users)
	})
}

func getUserByEmail(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		user := s.GetByEmail(params["email"])
		utils.ResponseJSON(w, http.StatusOK, user)
	})
}

func getUserByID(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		users := s.GetByID(params["id"])
		utils.ResponseJSON(w, http.StatusOK, users)
	})
}

func updateUser(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		var user entities.User
		_ = json.NewDecoder(r.Body).Decode(&user)
		s.Update(params["id"], user)
		utils.ResponseJSON(w, http.StatusOK, nil)
	})
}

func deleteUser(s business.UserService) http.Handler {
	return utils.HandlerFuncErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		s.Delete(params["id"])
		utils.ResponseJSON(w, http.StatusOK, nil)
	})
}
