package cors

import (
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/exp/slog"
)

// Config is a struct that holds configuration options for the CORS middleware.
type Config struct {
	AllowCredentials bool              // Flag to allow sending of credentials (cookies, etc.).
	AllowedMethods   []string          // List of allowed HTTP methods.
	AllowedHeaders   []string          // List of allowed HTTP request headers.
	ExposedHeaders   []string          // List of exposed response headers.
	MaxAge           int               // Maximum age (in seconds) of the CORS preflight cache.
	ValidateHeader   func(string) bool // Custom function to validate the request headers.
	ValidateOrigin   func(string) bool // Custom function to validate the origin.
}

// NewConfig creates a new Config struct with default values.
func NewConfig() *Config {
	return &Config{
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           0,
		ValidateHeader:   alwaysAllow, // Allow all headers by default.
		ValidateOrigin:   alwaysAllow, // Allow all origins by default.
	}
}

// alwaysAllow is a helper function that always returns true.
func alwaysAllow(_ string) bool {
	return true
}

// Middleware is the CORS middleware function that takes a Config struct and returns the middleware.
func Middleware(config *Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// If no origin is provided, continue without CORS handling.
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate the origin using the custom validation function.
			if !config.ValidateOrigin(origin) {
				slog.Error("request from origin not allowed", slog.String("origin", origin))
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Set CORS headers.
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")

			// Set Allow-Credentials header if specified in config.
			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight requests.
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				method := r.Header.Get("Access-Control-Request-Method")
				requestedHeaders := r.Header.Get("Access-Control-Request-Headers")

				// Validate the requested method.
				if !contains(config.AllowedMethods, method) {
					slog.Error("request method not allowed", slog.String("method", method))
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				// Set preflight response headers.
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ","))

				if len(config.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ","))
				}

				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
				}

				// Validate the requested headers using the custom validation function.
				if !validateHeaders(config.ValidateHeader, requestedHeaders) {
					slog.Error("request headers not allowed", slog.String("headers", requestedHeaders))
					w.WriteHeader(http.StatusForbidden)
					return
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Pass the request to the next middleware in the chain.
			next.ServeHTTP(w, r)
		})
	}
}

// contains checks if the target value exists in the list.
func contains[T comparable](list []T, value T) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// validateHeaders checks if all requested headers are allowed using the custom validation function.
func validateHeaders(validateHeaderFunc func(string) bool, requestedHeaders string) bool {
	requestedHeaderList := strings.Split(requestedHeaders, ",")

	for _, requestedHeader := range requestedHeaderList {
		requestedHeader = strings.TrimSpace(requestedHeader)
		if !validateHeaderFunc(requestedHeader) {
			return false
		}
	}

	return true
}

// ValidateOriginFromList returns a function that validates an origin against the provided list of allowed origins.
func ValidateOriginFromList(allowedOrigins []string) func(string) bool {
	return func(origin string) bool {
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				return true
			}
		}
		return false
	}
}
