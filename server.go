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
	fmt.Println("Starting container watcher...")
	runner.WatchContainers()

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
	fmt.Println("Stopping container watcher...")

	remaining, current := runner.StopWatching()
	fmt.Println("Container watcher stopped!")
	
	if current != 0 {
		fmt.Printf("Remaining workers: %d; Current workers: %d. If there are still workers to be run, there is a bug in the worker mechanism that needs to be fixed!\n", remaining, current)
	}

	os.Exit(0)
}
