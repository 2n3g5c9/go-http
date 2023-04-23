package common

import "strings"

// ShouldSkip checks if the given path should be skipped based on the excluded prefixes.
func ShouldSkip(path string, excludedPrefixes []string) bool {
	for _, prefix := range excludedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
