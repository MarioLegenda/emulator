package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	errorHandler "therebelsource/emulator/appErrors"
	"therebelsource/emulator/runner"
	"time"
)

func InitServer(r *mux.Router) *http.Server {
	origins := []string{"https://rebelsource.dev"}

	if os.Getenv("APP_ENV") != "prod" {
		origins = []string{"https://dev.therebelsource.local:8000"}
	}
	
	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete},
		ExposedHeaders:   []string{"Content-Length", "Content-Range", "Content-Type", "Cookie", "Set-Cookie"},
		AllowedHeaders:   []string{"Content-Range", "Set-Cookie", "Cookie", "Range", "Content-Type", "User-Agent", "X-Requested-With", "Cache-Control", "If-Modified-Since", "Keep-Alive", "DNT", "Origin", "Authorization", "x-rebel-source-auth", "Accept"},
		// Enable Debugging for testing, consider disabling in production
		Debug: os.Getenv("APP_ENV") != "prod",
	})

	handler := c.Handler(http.TimeoutHandler(r, 15*time.Second, "A timeout occurred"))

	srv := &http.Server{
		Handler: handler,
		Addr:    os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"),
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
		fmt.Println(errorHandler.ConstructError(errorHandler.ServerError, 0, err.Error()))
	}

	fmt.Println("Server is terminated.")
	fmt.Println("")
	fmt.Println("Stopping container balancer...")

	runner.StopContainerBalancer()
	fmt.Println("Container balancer stopped!")

	os.Exit(0)
}
