package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mstepan/user-service-golang/domain/service"
	"github.com/mstepan/user-service-golang/utils/http_utils"
	"log"
	"net/http"
	"os"
)

const (
	contextPath = "/api/v1"
	get         = "GET"
	post        = "POST"
	delete      = "DELETE"
	httpScheme  = "http"
)

var userHolder = service.NewUserHolder()

func NewRouting() *mux.Router {
	// configure all routing here
	routing := mux.NewRouter()

	// Middlewares
	if isLoggingEnabled() {
		routing.Use(loggingMiddleware)
	}

	// REST apis
	routing.HandleFunc(contextPath+"/users", addNewUser).
		Methods(post).
		Schemes(httpScheme)

	routing.HandleFunc(contextPath+"/users", getAllUsers).
		Methods(get).
		Schemes(httpScheme)

	routing.HandleFunc(contextPath+"/users/count", getUsersCount).
		Methods(get).
		Schemes(httpScheme)

	routing.HandleFunc(contextPath+"/users/{username:[a-zA-Z][\\w-]{1,31}}", getUserByUsername).
		Methods(get).
		Schemes(httpScheme)

	routing.HandleFunc(contextPath+"/users/{username:[a-zA-Z][\\w-]{1,31}}", deleteUserByUsername).
		Methods(delete).
		Schemes(httpScheme)

	return routing
}

func isLoggingEnabled() bool {
	_, exist := os.LookupEnv("LOGGING")
	return exist
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(resp, req)
	})
}

type CreateUserRequest struct {
	Username string
}

func addNewUser(resp http.ResponseWriter, req *http.Request) {

	userReq := &CreateUserRequest{}

	err := json.NewDecoder(req.Body).Decode(&userReq)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	userProfile := userHolder.AddUser(userReq.Username)

	if userProfile == nil {
		resp.WriteHeader(http.StatusConflict)
	} else {
		http_utils.WriteJsonBody(resp, http.StatusCreated, userProfile)
	}
}

func getAllUsers(resp http.ResponseWriter, req *http.Request) {

	allUsers := userHolder.GetAllUsers()

	http_utils.WriteJsonBody(resp, http.StatusOK, allUsers)
}

type counterResponse struct {
	Count int `json:"count"`
}

func getUsersCount(resp http.ResponseWriter, req *http.Request) {
	usersCount := userHolder.GetUsersCount()
	http_utils.WriteJsonBody(resp, http.StatusOK, &counterResponse{Count: usersCount})
}

func getUserByUsername(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	userProfile := userHolder.GetUserByUsername(vars["username"])

	if userProfile == nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	http_utils.WriteJsonBody(resp, http.StatusOK, userProfile)
}

func deleteUserByUsername(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	wasDeleted := userHolder.DeleteUserByUsername(vars["username"])

	if wasDeleted {
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	resp.WriteHeader(http.StatusNotFound)
}
