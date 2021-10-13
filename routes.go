package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"therebelsource/api/staticTypes"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	registerBlogRoutes(r)

	return r
}

func registerBlogRoutes(r *mux.Router) {
	r.HandleFunc("/api/environment-emulator/get-environments", getEnvironmentsHandler).Methods("GET")
	r.HandleFunc("/api/environment-emulator/execute/single-file", executeSingleCodeBlockHandler).Methods("POST")

	r.PathPrefix("/api/v2/static/").Handler(http.StripPrefix("/api/v2/static/", http.FileServer(http.Dir(staticTypes.ImgDir()))))
}
