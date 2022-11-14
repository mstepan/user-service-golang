package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mstepan/user-service-golang/domain_service"
	"github.com/mstepan/user-service-golang/utils/http_utils"
	"net/http"
)

const contextPath = "/api/v1"

var userHolder = domain_service.NewUserHolder()

func NewRouting() *mux.Router {
	// configure all routing here
	routing := mux.NewRouter()

	routing.HandleFunc(contextPath+"/users", addNewUser).
		Methods("POST").
		Schemes("http")

	routing.HandleFunc(contextPath+"/users", getAllUsers).
		Methods("GET").
		Schemes("http")

	routing.HandleFunc(contextPath+"/users/{username:[a-zA-Z][\\w-]{1,31}}", getUserByUsername).
		Methods("GET").
		Schemes("http")

	return routing
}

type CreateUserRequest struct {
	Username string
}

func addNewUser(respWriter http.ResponseWriter, req *http.Request) {

	userReq := &CreateUserRequest{}

	err := json.NewDecoder(req.Body).Decode(&userReq)

	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	userProfile := userHolder.AddUser(userReq.Username)

	if userProfile == nil {
		respWriter.WriteHeader(http.StatusConflict)
	} else {
		http_utils.WriteJsonBody(respWriter, http.StatusCreated, userProfile)
	}
}

func getAllUsers(respWriter http.ResponseWriter, req *http.Request) {

	allUsers := userHolder.GetAllUsers()

	http_utils.WriteJsonBody(respWriter, http.StatusOK, allUsers)
}

func getUserByUsername(respWriter http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	userProfile := userHolder.GetUserByUsername(vars["username"])

	if userProfile == nil {
		respWriter.WriteHeader(http.StatusNotFound)
		return
	}

	http_utils.WriteJsonBody(respWriter, http.StatusOK, userProfile)
}
