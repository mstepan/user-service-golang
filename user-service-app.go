package main

import (
	"context"
	"flag"
	"github.com/mstepan/user-service-golang/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const address = "0.0.0.0:7070"

func main() {

	// configure routing
	routing := api.NewRouting()
	http.Handle("/", routing)

	server := &http.Server{
		Addr: address,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      routing, // Pass our instance of gorilla/mux routing.
	}

	// Run our server in a goroutine so that it doesn't block.
	go startServer(server)

	mainChannel := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(mainChannel, os.Interrupt)

	// Block until we receive our signal.
	<-mainChannel

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), waitDuration())
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

func startServer(server *http.Server) {
	log.Printf("Server will be started at %s\n", address)
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func waitDuration() time.Duration {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	return wait
}
