package http

import (
	"emulator/cmd/http/rateLimiter"
	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	registerRoutes(r)

	return r
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/api/environment-emulator/get-environments", getEnvironmentsHandler).Methods("GET")
	r.HandleFunc("/api/environment-emulator/health", getEnvironmentsHandler).Methods("GET")
	r.Handle("/api/environment-emulator/execute/single-file", rateLimiter.PackageService.LimitedMiddleware(executeSingleCodeBlockHandler)).Methods("POST")
	r.Handle("/api/environment-emulator/public/execute/snippet", rateLimiter.PackageService.LimitedMiddleware(executePublicSnippet)).Methods("POST")
	r.Handle("/api/environment-emulator/execute/snippet", rateLimiter.PackageService.LimitedMiddleware(executeSnippet)).Methods("POST")
	r.Handle("/api/environment-emulator/public/execute/single-file", rateLimiter.PackageService.LimitedMiddleware(executePublicSingleFileRunResult)).Methods("POST")
	r.Handle("/api/environment-emulator/execute/project", rateLimiter.PackageService.LimitedMiddleware(executeProjectHandler)).Methods("POST")
	r.Handle("/api/environment-emulator/execute/linked-project", rateLimiter.PackageService.LimitedMiddleware(executeLinkedProjectHandler)).Methods("POST")
	r.Handle("/api/environment-emulator/public/execute/linked-project", rateLimiter.PackageService.LimitedMiddleware(executePublicLinkedProjectHandler)).Methods("POST")
}
