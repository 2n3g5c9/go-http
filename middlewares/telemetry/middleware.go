package telemetry

import (
	"net/http"
	"reflect"
	"time"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/2n3g5c9/go-http/middlewares/common"
)

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

// Middleware is a simple OpenTelemetry HTTP middleware.
func Middleware(next http.Handler, opts ...MiddlewareOption) http.Handler {
	var options middlewareOptions
	for _, opt := range opts {
		opt(&options)
	}

	var (
		pkgName = reflect.TypeOf(struct{}{}).PkgPath()
		meter   = otel.GetMeterProvider().Meter(pkgName)
		metrics = NewMetrics(&meter)
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if common.ShouldSkip(r.URL.Path, options.excludedPrefixes) {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		defer span.End()

		span.SetAttributes(
			semconv.HTTPMethodKey.String(r.Method),
			semconv.HTTPURLKey.String(r.URL.String()),
			semconv.HTTPUserAgentKey.String(r.UserAgent()),
		)

		startTime := time.Now()
		wrappedWriter := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))
		duration := time.Since(startTime)

		metrics.IncreaseRequestCounter(ctx, r.Method)
		metrics.RecordRequestDuration(ctx, r.Method, duration)

		// Set the default status code if it was not set by the handler.
		if wrappedWriter.statusCode == 0 {
			wrappedWriter.statusCode = http.StatusOK
		}

		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(wrappedWriter.statusCode))
	})
}

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader implements http.ResponseWriter and records the status code.
func (w *statusWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
