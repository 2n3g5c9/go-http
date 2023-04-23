package logging

import (
	"testing"

	"golang.org/x/exp/slog"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
	}{
		{
			name:     "Debug level",
			logLevel: "debug",
		},
		{
			name:     "Info level",
			logLevel: "info",
		},
		{
			name:     "Warn level",
			logLevel: "warn",
		},
		{
			name:     "Error level",
			logLevel: "error",
		},
		{
			name:     "Invalid level",
			logLevel: "invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Init(test.logLevel)
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  string
		wantLevel slog.Level
	}{
		{
			name:      "Debug",
			logLevel:  "debug",
			wantLevel: slog.LevelDebug,
		},
		{
			name:      "Info",
			logLevel:  "info",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "Warn",
			logLevel:  "warn",
			wantLevel: slog.LevelWarn,
		},
		{
			name:      "Error",
			logLevel:  "error",
			wantLevel: slog.LevelError,
		},
		{
			name:      "Invalid",
			logLevel:  "invalid",
			wantLevel: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := parseLevel(tt.logLevel)
			if level != tt.wantLevel {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.logLevel, level, tt.wantLevel)
			}
		})
	}
}
