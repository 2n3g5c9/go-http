package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldSkip(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		excludedPrefixes []string
		want             bool
	}{
		{
			name:             "Empty excludedPrefixes",
			path:             "path/to/file",
			excludedPrefixes: []string{},
			want:             false,
		},
		{
			name:             "Not excluded",
			path:             "path/to/file",
			excludedPrefixes: []string{"other"},
			want:             false,
		},
		{
			name:             "Single exclusion match",
			path:             "path/to/file",
			excludedPrefixes: []string{"path"},
			want:             true,
		},
		{
			name:             "Multiple exclusion match",
			path:             "path/to/file",
			excludedPrefixes: []string{"other", "path"},
			want:             true,
		},
		{
			name:             "Partial match",
			path:             "path/to/file",
			excludedPrefixes: []string{"pat"},
			want:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldSkip(tt.path, tt.excludedPrefixes)
			assert.Equal(t, tt.want, got)
		})
	}
}
