package telemetry

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

// TestNewGCPResource tests the newGCPResource function with various inputs.
func TestNewGCPResource(t *testing.T) {
	tests := []struct {
		name string
		want *resource.Resource
	}{
		{
			name: "success",
			want: func() *resource.Resource {
				res, err := resource.New(context.Background(),
					resource.WithDetectors(gcp.NewDetector()),
					resource.WithTelemetrySDK(),
					resource.WithAttributes(semconv.ServiceNameKey.String(filepath.Base(os.Args[0]))),
				)
				assert.NoError(t, err)
				return res
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newGCPResource(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
