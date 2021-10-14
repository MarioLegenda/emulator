package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	errorHandler "therebelsource/api/appErrors"
	"therebelsource/emulator/runner"
	"time"
)

func InitServer(r *mux.Router) *http.Server {

	srv := &http.Server{
		Handler:      r,
		Addr:         os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("")
	fmt.Println("Starting container balancer...")
	runner.StartContainerBalancer()

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		fmt.Printf("Starting server on %s:%v...\n", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	return srv
}

func WatchServerShutdown(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)

	if err != nil {
		fmt.Println(errorHandler.ConstructError(errorHandler.ServerError, errorHandler.ServerShutdownError, err.Error()))
	}

	fmt.Println("Server is terminated.")
	fmt.Println("")
	fmt.Println("Stopping container balancer...")

	runner.StopContainerBalancer()
	fmt.Println("Container balancer stopped!")

	os.Exit(0)
}
