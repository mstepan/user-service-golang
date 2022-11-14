package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const routPrefix = "/api/v1"

func main() {

	// curl http://localhost:7070/api/v1/users/maksym | jq
	// curl http://localhost:7070/api/v1/users/zorro | jq

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// configure all routing here
	routing := mux.NewRouter()
	routing.HandleFunc(routPrefix+"/users/{username:[a-zA-Z][\\w]{1,31}}", GetUserByUsername).
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

type UserProfile struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func GetUserByUsername(respWriter http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)

	data, err := json.Marshal(&UserProfile{Id: 133, Username: vars["username"]})

	if err != nil {
		log.Println("Can't properly marshal response")
	} else {
		_, err := respWriter.Write(data)
		if err != nil {
			log.Printf("Can't properly write reponse: %v\n", err)
		}
	}

}
