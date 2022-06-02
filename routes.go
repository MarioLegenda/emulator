package main

import (
	"github.com/gorilla/mux"
	"therebelsource/emulator/rateLimiter"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	registerBlogRoutes(r)

	return r
}

func registerBlogRoutes(r *mux.Router) {
	r.HandleFunc("/api/environment-emulator/get-environments", getEnvironmentsHandler).Methods("GET")
	r.Handle("/api/environment-emulator/execute/single-file", rateLimiter.PackageService.LimitedMiddleware(executeSingleCodeBlockHandler)).Methods("POST")
	//r.Handle("/api/environment-emulator/public/execute/single-file", rateLimiter.PackageService.LimitedMiddleware(executePublicSingleFileRunResult)).Methods("POST")
	r.HandleFunc("/api/environment-emulator/public/execute/single-file", executePublicSingleFileRunResult).Methods("POST")
	r.Handle("/api/environment-emulator/execute/project", rateLimiter.PackageService.LimitedMiddleware(executeProjectHandler)).Methods("POST")
	r.Handle("/api/environment-emulator/execute/linked-project", rateLimiter.PackageService.LimitedMiddleware(executeLinkedProjectHandler)).Methods("POST")
	//r.Handle("/api/environment-emulator/public/execute/linked-project", rateLimiter.PackageService.LimitedMiddleware(executePublicLinkedProjectHandler)).Methods("POST")
	r.HandleFunc("/api/environment-emulator/public/execute/linked-project", executePublicLinkedProjectHandler).Methods("POST")
}
