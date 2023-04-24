package telemetry

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWriteHeader tests the WriteHeader function with various status codes.
func TestWriteHeader(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		expectedCode int
	}{
		{
			name:         "OK status",
			statusCode:   http.StatusOK,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Not Found status",
			statusCode:   http.StatusNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Internal Server Error status",
			statusCode:   http.StatusInternalServerError,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				recorder = httptest.NewRecorder()
				sw       = &statusWriter{ResponseWriter: recorder}
			)

			sw.WriteHeader(tt.statusCode)
			assert.Equal(t, tt.expectedCode, sw.statusCode)
			assert.Equal(t, tt.expectedCode, recorder.Code)
		})
	}
}
