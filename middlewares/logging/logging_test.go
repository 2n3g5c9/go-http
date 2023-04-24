package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
	}{
		{"Debug level", "debug"},
		{"Info level", "info"},
		{"Warn level", "warn"},
		{"Error level", "error"},
		{"Invalid level", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.logLevel)
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  string
		wantLevel slog.Level
	}{
		{"Debug", "debug", slog.LevelDebug},
		{"Info", "info", slog.LevelInfo},
		{"Warn", "warn", slog.LevelWarn},
		{"Error", "error", slog.LevelError},
		{"Invalid", "invalid", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := parseLevel(tt.logLevel)
			assert.Equal(t, tt.wantLevel, level)
		})
	}
}
