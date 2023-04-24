package cors

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		list  []interface{}
		value interface{}
		want  bool
	}{
		{"Int value present in list", []interface{}{1, 2, 3, 4, 5}, 3, true},
		{"Int value not present in list", []interface{}{1, 2, 3, 4, 5}, 6, false},
		{"String value present in list", []interface{}{"apple", "banana", "cherry"}, "banana", true},
		{"String value not present in list", []interface{}{"apple", "banana", "cherry"}, "orange", false},
		{"Empty list", []interface{}{}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, contains(tt.list, tt.value))
		})
	}
}

func TestValidateHeaders(t *testing.T) {
	tests := []struct {
		name               string
		validateHeaderFunc func(string) bool
		requestedHeaders   string
		want               bool
	}{
		{"no requested headers", alwaysAllow, "", true},
		{"all headers allowed", alwaysAllow, "Content-Type, Accept", true},
		{"some headers not allowed", func(s string) bool { return s != "Content-Type" }, "Content-Type, Accept", false},
		{"no headers allowed", func(_ string) bool { return false }, "Content-Type, Accept", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateHeaders(tt.validateHeaderFunc, tt.requestedHeaders)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateOriginFromList(t *testing.T) {
	allowedOrigins := []string{"https://example.com", "https://example.org"}
	validateOrigin := ValidateOriginFromList(allowedOrigins)

	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{"allowed origin", "https://example.com", true},
		{"disallowed origin", "https://disallowed.example", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateOrigin(tt.origin)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMiddleware(t *testing.T) {
	allowedOrigins := []string{"https://example.com", "https://example.org"}
	validateOrigin := ValidateOriginFromList(allowedOrigins)

	config := NewConfig()
	config.AllowedMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowedHeaders = []string{"Content-Type", "Accept"}
	config.ExposedHeaders = []string{"X-Exposed-Header"}
	config.MaxAge = 600
	config.AllowCredentials = true
	config.ValidateOrigin = validateOrigin

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Test response")
	})

	middleware := Middleware(config)

	tests := []struct {
		name           string
		requestOrigin  string
		requestMethod  string
		requestHeaders map[string]string
		wantStatus     int
		wantHeaders    map[string]string
	}{
		{
			"allowed CORS request",
			"https://example.com",
			"GET",
			nil,
			http.StatusOK,
			map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Allow-Credentials": "true",
				"Vary":                             "Origin",
			},
		},
		{
			"disallowed CORS request",
			"https://disallowed.example",
			"GET",
			nil,
			http.StatusForbidden,
			nil,
		},
		{
			"valid preflight request",
			"https://example.com",
			"OPTIONS",
			map[string]string{
				"Access-Control-Request-Method":  "POST",
				"Access-Control-Request-Headers": "Content-Type, Accept",
			},
			http.StatusNoContent,
			map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Allow-Credentials": "true",
				"Vary":                             "Origin",
				"Access-Control-Allow-Methods":     "GET,POST,OPTIONS",
				"Access-Control-Allow-Headers":     "Content-Type,Accept",
				"Access-Control-Expose-Headers":    "X-Exposed-Header",
				"Access-Control-Max-Age":           "600",
			},
		},
		{
			"invalid preflight request",
			"https://example.com",
			"OPTIONS",
			map[string]string{
				"Access-Control-Request-Method":  "PUT",
				"Access-Control-Request-Headers": "Content-Type, Accept",
			},
			http.StatusMethodNotAllowed,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.requestMethod, "http://example.com/test", nil)
			req.Header.Set("Origin", tt.requestOrigin)

			if tt.requestHeaders != nil {
				for key, value := range tt.requestHeaders {
					req.Header.Set(key, value)
				}
			}

			rr := httptest.NewRecorder()

			middleware(testHandler).ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)

			if tt.wantHeaders != nil {
				for key, wantValue := range tt.wantHeaders {
					gotValue := rr.Header().Get(key)
					assert.Equal(t, wantValue, gotValue, "handler returned wrong header value for %s", key)
				}
			}
		})
	}
}
