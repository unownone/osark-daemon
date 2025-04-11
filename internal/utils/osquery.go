package utils

import (
	"os"
	"runtime"
)

// FindOSQuery searches for the osquery socket in common locations
// and returns the socket path if found, or an error if not found.
func FindOSQuery() (string, error) {
	// Common socket paths by OS
	var socketPaths []string

	switch runtime.GOOS {
	case "darwin":
		socketPaths = []string{
			"/tmp/osquery.sock",
			"/var/osquery/osquery.em",
			"/var/run/osquery/osquery.em",
			"/private/var/osquery/osquery.em",
		}
	case "linux":
		socketPaths = []string{
			"/var/run/ossuary.sock",
			"/var/osquery/osquery.em",
			"/var/run/osquery.em",
			"/var/run/osquery/osquery.em",
		}
	case "windows":
		socketPaths = []string{
			`\\.\pipe\osquery.em`,
		}
	}

	// Check all potential socket paths
	for _, path := range socketPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", os.ErrNotExist
}
