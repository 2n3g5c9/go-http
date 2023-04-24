package logging

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		excludedPrefixes []string
		shouldLog        bool
	}{
		{"Request not excluded", "/api/v1/test", []string{}, true},
		{"Request excluded", "/api/v1/test", []string{"/api"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			Init("info")
			slog.SetDefault(slog.New(slog.NewJSONHandler(buf).WithAttrs([]slog.Attr{})))

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			middlewareHandler := Middleware(nextHandler, WithExcludedPrefixes(tt.excludedPrefixes))

			req := httptest.NewRequest("GET", tt.path, nil)
			resp := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(resp, req)

			logged := buf.String()

			if tt.shouldLog {
				require.Contains(t, logged, tt.path, "expected log to contain path %s, got: %s", tt.path, logged)
			} else {
				require.NotContains(t, logged, tt.path, "expected log not to contain path %s, got: %s", tt.path, logged)
			}
		})
	}
}
