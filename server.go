package main

import (
	"context"
	"fmt"
	"github.com/coreos/go-systemd/daemon"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	errorHandler "therebelsource/emulator/appErrors"
	"therebelsource/emulator/logger"
	"therebelsource/emulator/slack"
	"time"
)

func InitServer(r *mux.Router) *http.Server {
	origins := []string{"https://rebelsource.dev"}

	if os.Getenv("APP_ENV") != "prod" {
		origins = []string{"https://dev.therebelsource.local:8000"}
	}

	if os.Getenv("APP_ENV") == "staging" {
		origins = []string{"https://staging-rebelsource.com"}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete},
		ExposedHeaders:   []string{"Content-Length", "Content-Range", "Content-Type", "Cookie", "Set-Cookie"},
		AllowedHeaders:   []string{"Content-Range", "Set-Cookie", "Cookie", "Range", "Content-Type", "User-Agent", "X-Requested-With", "Cache-Control", "If-Modified-Since", "Keep-Alive", "DNT", "Origin", "Authorization", "x-rebel-source-auth", "Accept"},
		// Enable Debugging for testing, consider disabling in production
		Debug: os.Getenv("APP_ENV") != "prod" && os.Getenv("APP_ENV") != "staging",
	})

	handler := c.Handler(r)

	srv := &http.Server{
		Handler:      handler,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		Addr:         os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logger.Info(fmt.Sprintf("Starting server on %s:%v...\n", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")))
		slack.SendLog("Deploy", "Server is booted and ready", "deploy_log")

		if os.Getenv("APP_ENV") == "prod" || os.Getenv("APP_ENV") == "staging" {
			daemon.SdNotify(false, daemon.SdNotifyReady)
		}

		if err := srv.ListenAndServe(); err != nil {
			logger.Error(fmt.Sprintf("Server closed: %s", err.Error()))
			slack.SendLog("Server closed", err.Error(), "deploy_log")
		}
	}()

	return srv
}

func WatchServerShutdown(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Info("Stopping emulator workers...")
	closeExecutioners()
	logger.Info("Emulator workers stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)

	if err != nil {
		logger.Info(errorHandler.ConstructError(errorHandler.ServerError, 0, err.Error()))
	}

	logger.Info("Server is terminated. App shutting down!")
	fmt.Println("")

	slack.SendLog("Shutdown", "Server has successfully shutdown", "deploy_log")

	if os.Getenv("APP_ENV") == "prod" || os.Getenv("APP_ENV") == "staging" {
		daemon.SdNotify(false, daemon.SdNotifyStopping)
	}

	os.Exit(0)
}
