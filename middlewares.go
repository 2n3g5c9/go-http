package go_http

import (
	"net/http"
)

// Middleware is a function type that represents an HTTP middleware.
type Middleware func(http.Handler) http.Handler

type middlewareOptions struct {
	CORS      *CORSOption
	Logging   *LoggingOption
	Telemetry *TelemetryOption
}

type MiddlewareOption func(*middlewareOptions)

type CORSOption struct {
	AllowedMethods []string
	AllowedOrigins []string
}

type LoggingOption struct {
	ExcludedPrefixes []string
}

type TelemetryOption struct {
	ExcludedPrefixes []string
}

// WithCORS returns a MiddlewareOption that sets the CORS middleware options.
func WithCORS(allowedMethods, allowedOrigins []string) MiddlewareOption {
	return func(opts *middlewareOptions) {
		opts.CORS = &CORSOption{
			AllowedMethods: allowedMethods,
			AllowedOrigins: allowedOrigins,
		}
	}
}

// WithLogging returns a MiddlewareOption that sets the Logging middleware options.
func WithLogging(prefixes []string) MiddlewareOption {
	return func(opts *middlewareOptions) {
		opts.Logging = &LoggingOption{
			ExcludedPrefixes: prefixes,
		}
	}
}

// WithTelemetry returns a MiddlewareOption that sets the Telemetry (Metrics & Traces) middleware options.
func WithTelemetry(excludedPrefixes []string) MiddlewareOption {
	return func(opts *middlewareOptions) {
		opts.Telemetry = &TelemetryOption{
			ExcludedPrefixes: excludedPrefixes,
		}
	}
}
