package logging

import (
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/2n3g5c9/go-http/middlewares/common"
)

var gitCommit string

type MiddlewareOption func(*middlewareOptions)

type middlewareOptions struct {
	excludedPrefixes []string
}

// WithExcludedPrefixes sets the excluded paths for the middleware.
func WithExcludedPrefixes(prefixes []string) MiddlewareOption {
	return func(options *middlewareOptions) {
		options.excludedPrefixes = prefixes
	}
}

// Middleware is a middleware that provides basic HTTP access logging.
func Middleware(next http.Handler, opts ...MiddlewareOption) http.Handler {
	options := &middlewareOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if common.ShouldSkip(r.URL.Path, options.excludedPrefixes) {
			next.ServeHTTP(w, r)
			return
		}

		slog.Info("request received",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("userAgent", r.UserAgent()),
		)
		next.ServeHTTP(w, r)
	})
}
