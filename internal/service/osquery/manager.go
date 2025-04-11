package osquery

import (
	"runtime"
	"time"

	"github.com/osquery/osquery-go"
	"github.com/pkg/errors"
	"github.com/unownone/osark-daemon/internal/utils"
	"github.com/unownone/osark-daemon/models"
)

// Manager is the interface for the tracking manager
// It is responsible for managing the tracking of apps and other events in the system
type Manager interface {
	GetSystemInfo() (*models.SystemInfo, error) // GetSystemInfo returns the system information
	GetApps() ([]*models.AppInfo, error) // GetApps returns all the apps in the system
	StartLoggerProcess() error // StartLoggerProcess starts the logger process
}

type manager struct {
	osArch   models.OSArch
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
	osArch, err := models.NewOSArch(runtime.GOOS)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start manager")
	}
	return &manager{
		osArch:   osArch,
		osClient: osQueryClient,
	}, nil
}

// StartLoggerProcess starts the logger process
func (m *manager) StartLoggerProcess() error {
	return nil
}

