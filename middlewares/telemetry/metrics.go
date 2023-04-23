package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

type Metrics struct {
	requestCounter  instrument.Int64Counter
	requestDuration instrument.Int64Histogram
}

// NewMetrics returns a new Metrics instance.
func NewMetrics(meter *metric.Meter) *Metrics {
	requestCounter, _ := (*meter).Int64Counter(
		"http_requests_total",
		instrument.WithDescription("Total number of HTTP requests."),
	)

	requestDuration, _ := (*meter).Int64Histogram(
		"http_request_duration_ms",
		instrument.WithDescription("HTTP request duration in milliseconds."),
	)

	return &Metrics{
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

// IncreaseRequestCounter increases the request counter by 1.
func (m *Metrics) IncreaseRequestCounter(ctx context.Context, method string) {
	m.requestCounter.Add(ctx, 1, semconv.HTTPMethodKey.String(method))
}

// RecordRequestDuration records the request duration in milliseconds.
func (m *Metrics) RecordRequestDuration(ctx context.Context, method string, duration time.Duration) {
	m.requestDuration.Record(ctx, duration.Milliseconds(), semconv.HTTPMethodKey.String(method))
}
