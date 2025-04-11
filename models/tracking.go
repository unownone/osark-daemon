package models

import (
	"errors"
	"time"
)

// LogEvent is the event that is logged to the server
type LogEvent struct {
	AppInfo    []*AppInfo  `json:"app_info,omitempty"`    // AppInfo is the information about an app
	Error      string      `json:"error,omitempty"`       // Error is the error message
	SystemInfo *SystemInfo `json:"system_info,omitempty"` // SystemInfo is the information about the system
	CreatedAt  time.Time   `json:"created_at"`            // CreatedAt is the time the event was created
}

// AppInfo is the information about an app
type AppInfo struct {
	ID             string    `json:"id"`               // ID of the app
	Name           string    `json:"name"`             // Name of the app
	BundleName     string    `json:"bundle_name"`      // Bundle name of the app
	BundleID       string    `json:"bundle_id"`        // Bundle ID of the app
	BundleVersion  string    `json:"bundle_version"`   // Bundle version of the app
	Path           string    `json:"path"`             // Path of the app
	LastOpenedTime time.Time `json:"last_opened_time"` // Last opened time of the app
}

// SystemInfo is the information about the system
type SystemInfo struct {
	UptimeSeconds time.Duration `json:"uptime_seconds"` // Uptime seconds of the system
	OSName        string        `json:"os_name"`        // Name of the operating system
	OSVersion     string        `json:"os_version"`     // Version of the operating system
	OSArch        string        `json:"os_arch"`        // Architecture of the operating system
	MacAddress    string        `json:"mac_address"`    // Mac address of the system
}

// OSArch is the architecture of the operating system
type OSArch string

const (
	OSArchDarwin  OSArch = "darwin"  // macOS
	OSArchLinux   OSArch = "linux"   // Linux
	OSArchWindows OSArch = "windows" // Windows
)

// NewOSArch creates a new OSArch from a string
func NewOSArch(maybeOSArch string) (OSArch, error) {
	switch maybeOSArch {
	case "darwin":
		return OSArchDarwin, nil
	case "linux":
		return OSArchLinux, nil
	case "windows":
		return OSArchWindows, nil
	default:
		return "", errors.New("invalid/unsupported operating system")
	}
}
