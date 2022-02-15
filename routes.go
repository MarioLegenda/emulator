package main

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	registerBlogRoutes(r)

	return r
}

func registerBlogRoutes(r *mux.Router) {
	r.HandleFunc("/api/environment-emulator/get-environments", getEnvironmentsHandler).Methods("GET")
	r.HandleFunc("/api/environment-emulator/execute/single-file", executeSingleCodeBlockHandler).Methods("POST")
	r.HandleFunc("/api/environment-emulator/public/execute/single-file", executePublicSingleFileRunResult).Methods("POST")
	r.HandleFunc("/api/environment-emulator/execute/project", executeProjectHandler).Methods("POST")
	r.HandleFunc("/api/environment-emulator/execute/linked-project", executeLinkedProjectHandler).Methods("POST")
}
