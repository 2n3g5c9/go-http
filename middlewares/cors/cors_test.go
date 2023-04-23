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
		name     string
		list     interface{}
		value    interface{}
		expected bool
	}{
		{
			name:     "Int value present in list",
			list:     []int{1, 2, 3, 4, 5},
			value:    3,
			expected: true,
		},
		{
			name:     "Int value not present in list",
			list:     []int{1, 2, 3, 4, 5},
			value:    6,
			expected: false,
		},
		{
			name:     "String value present in list",
			list:     []string{"apple", "banana", "cherry"},
			value:    "banana",
			expected: true,
		},
		{
			name:     "String value not present in list",
			list:     []string{"apple", "banana", "cherry"},
			value:    "orange",
			expected: false,
		},
		{
			name:     "Empty list",
			list:     []int{},
			value:    1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch l := tt.list.(type) {
			case []int:
				result = contains(l, tt.value.(int))
			case []string:
				result = contains(l, tt.value.(string))
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateHeaders(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := validateHeaders(tc.validateHeaderFunc, tc.requestedHeaders)
			if got != tc.want {
				t.Errorf("validateHeaders() = %v, want %v", got, tc.want)
			}
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
			if got != tt.want {
				t.Errorf("ValidateOriginFromList() = %v, want %v", got, tt.want)
			}
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

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			if tt.wantHeaders != nil {
				for key, wantValue := range tt.wantHeaders {
					gotValue := rr.Header().Get(key)
					if gotValue != wantValue {
						t.Errorf("handler returned wrong header value for %s: got %v want %v", key, gotValue, wantValue)
					}
				}
			}
		})
	}
}
