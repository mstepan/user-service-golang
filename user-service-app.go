package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/mstepan/user-service-golang/api"
	"github.com/mstepan/user-service-golang/domain_service"
	"github.com/mstepan/user-service-golang/utils/http_utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const routPrefix = "/api/v1"

var userHolder = domain_service.NewUserHolder()

func main() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// configure all routing here
	routing := mux.NewRouter()

	routing.HandleFunc(routPrefix+"/users", addNewUser).
		Methods("POST").
		Schemes("http")

	routing.HandleFunc(routPrefix+"/users", getAllUsers).
		Methods("GET").
		Schemes("http")

	routing.HandleFunc(routPrefix+"/users/{username:[a-zA-Z][\\w-]{1,31}}", getUserByUsername).
		Methods("GET").
		Schemes("http")

	http.Handle("/", routing)

	server := &http.Server{
		Addr: "0.0.0.0:7070",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      routing, // Pass our instance of gorilla/mux routing.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Can't properly shutdown server: %v\n", err)
		os.Exit(1)
	}

	// Optionally, you could run server.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Server termination completed successfully")
	os.Exit(0)

}

func addNewUser(respWriter http.ResponseWriter, req *http.Request) {

	userReq := &api.CreateUserRequest{}

	err := json.NewDecoder(req.Body).Decode(&userReq)

	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	wasAdded := userHolder.AddUser(userReq)

	if wasAdded {
		respWriter.WriteHeader(http.StatusCreated)
	} else {
		respWriter.WriteHeader(http.StatusConflict)
	}

}

func getAllUsers(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Header().Set(http_utils.ApplicationJson())
	respWriter.WriteHeader(http.StatusOK)
}

func getUserByUsername(respWriter http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	userProfile := userHolder.GetUserByUsername(vars["username"])

	if userProfile == nil {
		respWriter.WriteHeader(http.StatusNotFound)
		return
	}

	userProfileData, marshallErr := json.Marshal(userProfile)
	if marshallErr != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, writeErr := respWriter.Write(userProfileData)
	if writeErr != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	respWriter.WriteHeader(http.StatusOK)

}
