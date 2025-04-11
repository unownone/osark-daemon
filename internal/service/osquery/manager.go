package osquery

import (
	"time"

	"github.com/osquery/osquery-go"
	"github.com/pkg/errors"
	"github.com/unownone/osark-daemon/internal/utils"
	"github.com/unownone/osark-daemon/models"
)

// Manager is the interface for the tracking manager
// It is responsible for managing the tracking of apps and other events in the system
type Manager interface {
	GetSystemInfo() (*models.SystemInfo, error)                                   // GetSystemInfo returns the system information
	GetApps() ([]*models.AppInfo, error)                                          // GetApps returns all the apps in the system
	GetCurrentRunningProcesses(bundleIDs []string) ([]*models.ProcessInfo, error) // GetCurrentRunningProcesses returns the current running processes
	StartLoggerProcess() error                                                    // StartLoggerProcess starts the logger process
}

type manager struct {
	osClient *osquery.ExtensionManagerClient
}

// NewManager creates a new manager
func NewManager() (Manager, error) {
	osQueryPath, err := utils.FindOSQuery()
	if err != nil {
		return nil, errors.Wrap(err, "osquery not found. cannot setup manager")
	}
	osQueryClient, err := osquery.NewClient(osQueryPath, 10*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create osquery client")
	}
	return &manager{
		osClient: osQueryClient,
	}, nil
}

// StartLoggerProcess starts the logger process
func (m *manager) StartLoggerProcess() error {
	return nil
}
