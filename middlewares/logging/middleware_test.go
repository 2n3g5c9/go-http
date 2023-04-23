package logging

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/exp/slog"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		excludedPrefixes []string
		shouldLog        bool
	}{
		{
			name:             "Request not excluded",
			path:             "/api/v1/test",
			excludedPrefixes: []string{},
			shouldLog:        true,
		},
		{
			name:             "Request excluded",
			path:             "/api/v1/test",
			excludedPrefixes: []string{"/api"},
			shouldLog:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init("info")
			slog.SetDefault(slog.New(slog.NewJSONHandler(&buf).WithAttrs([]slog.Attr{})))

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			middlewareHandler := Middleware(nextHandler, WithExcludedPrefixes(tt.excludedPrefixes))

			req := httptest.NewRequest("GET", tt.path, nil)
			resp := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(resp, req)

			logged := buf.String()

			if tt.shouldLog && !strings.Contains(logged, tt.path) {
				t.Errorf("expected log to contain path %s, got: %s", tt.path, logged)
			}

			if !tt.shouldLog && strings.Contains(logged, tt.path) {
				t.Errorf("expected log not to contain path %s, got: %s", tt.path, logged)
			}
		})
	}
}
