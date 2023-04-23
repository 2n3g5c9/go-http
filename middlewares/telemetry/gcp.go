package telemetry

import (
	"context"
	"os"
	"path/filepath"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

// NewGCPMeterProvider returns a new OpenTelemetry MeterProvider with a Google Cloud Monitoring exporter.
func NewGCPMeterProvider(ctx context.Context, projectID string) (*sdkmetric.MeterProvider, error) {
	exporter, err := mexporter.New(mexporter.WithProjectID(projectID))
	if err != nil {
		return nil, err
	}

	res, err := newGCPResource(ctx)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)

	return mp, nil
}

// NewGCPTracerProvider returns a new OpenTelemetry TracerProvider with a Google Cloud Trace exporter.
func NewGCPTracerProvider(ctx context.Context, projectID string, tracingRatio float64) (*sdktrace.TracerProvider, error) {
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		return nil, err
	}

	res, err := newGCPResource(ctx)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(tracingRatio)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	return tp, nil
}

// newGCPResource returns a new OpenTelemetry Resource for Google Cloud with the service name set to the executable name.
func newGCPResource(ctx context.Context) (*resource.Resource, error) {
	executableName := filepath.Base(os.Args[0])

	return resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(semconv.ServiceNameKey.String(executableName)),
	)
}
