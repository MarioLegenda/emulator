package rateLimiter

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"net/http"
	"os"
	"therebelsource/emulator/appErrors"
	"time"
)

var PackageService Service

type Middleware func(handlerFunc http.HandlerFunc) http.Handler

type Service struct {
	limitedStore limiter.Store
}

func InitRateLimiter() {
	var tokens uint64 = 1
	interval := 500 * time.Millisecond

	if os.Getenv("APP_ENV") == "test" {
		tokens = 100000
		interval = 1 * time.Hour
	}

	limitedStore, err := memorystore.New(&memorystore.Config{
		Tokens:   tokens,
		Interval: interval,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("ratelimiter package could not create limitedStore: %s", err.Error()))
	}

	PackageService = Service{limitedStore: limitedStore}
}

func (s Service) LimitedMiddleware(handler http.HandlerFunc) http.Handler {
	middleware, err := httplimit.NewMiddleware(s.limitedStore, httplimit.IPKeyFunc())
	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("ratelimiter package could not create middleware for limitedStore: %s", err.Error()))
	}

	return middleware.Handle(handler)
}

func (s Service) Close() error {
	var lErr error

	if err := s.limitedStore.Close(context.Background()); err != nil {
		lErr = err
	}

	return lErr
}
