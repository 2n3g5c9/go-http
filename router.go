package go_http

import (
	"net/http"

	"github.com/2n3g5c9/go-http/middlewares/cors"
	"github.com/2n3g5c9/go-http/middlewares/logging"
	"github.com/2n3g5c9/go-http/middlewares/telemetry"
)

// Router is a custom HTTP router that supports middlewares.
type Router struct {
	*http.ServeMux
	middlewares []Middleware
}

// NewRouter creates a new Router with the specified MiddlewareOptions.
// It sets up middlewares for CORS, logging, and tracing based on the provided options.
func NewRouter(opts ...MiddlewareOption) *Router {
	var (
		r       = Router{ServeMux: http.NewServeMux(), middlewares: []Middleware{}}
		options = &middlewareOptions{}
	)

	for _, opt := range opts {
		opt(options)
	}

	// Configure and add CORS middleware if CORS options are provided.
	if options.CORS != nil {
		corsCfg := cors.NewConfig()
		corsCfg.AllowedMethods = options.CORS.AllowedMethods
		corsCfg.ValidateOrigin = cors.ValidateOriginFromList(options.CORS.AllowedOrigins)
		r.middlewares = append(r.middlewares, cors.Middleware(corsCfg))
	}

	// Configure and add logging middleware if logging options are provided.
	if options.Logging != nil {
		r.middlewares = append(r.middlewares,
			func(next http.Handler) http.Handler {
				return logging.Middleware(next, logging.WithExcludedPrefixes(options.Logging.ExcludedPrefixes))
			})
	}

	// Configure and add telemetry middleware if metrics or tracing options are provided.
	if options.Telemetry != nil {
		r.middlewares = append(r.middlewares,
			func(next http.Handler) http.Handler {
				return telemetry.Middleware(next, telemetry.WithExcludedPrefixes(options.Telemetry.ExcludedPrefixes))
			})
	}

	return &r
}

// HandlerFunc method returns a http.HandlerFunc that wraps the Router with the configured middlewares.
func (r *Router) HandlerFunc() *http.HandlerFunc {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		wrappedHandler := http.Handler(r)
		for _, middleware := range r.middlewares {
			wrappedHandler = middleware(wrappedHandler)

		}
		wrappedHandler.ServeHTTP(w, req)
	})
	return &h
}
