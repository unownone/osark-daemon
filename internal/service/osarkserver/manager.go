package osarkserver

import (
	"net/http"
	"time"

	"github.com/unownone/osark-daemon/models"
)

// Manager is the interface for the push manager
type Manager interface {
	Authenticate(*models.SystemInfo) error // Authenticate authenticates the push manager
	Push(data []*models.LogEvent) error    // Push pushes the data to the server
	PushError(error) error                 // PushError pushes an error to the server
}

type pushManager struct {
	service        *http.Client
	osarkServerURL string
	deviceID       string
}

// NewPushManager creates a new push manager
func NewPushManager(osarkServerURL string, info *models.SystemInfo) (Manager, error) {
	manager := &pushManager{
		service: &http.Client{
			Timeout: 10 * time.Second,
		},
		osarkServerURL: osarkServerURL,
	}
	err := manager.Authenticate(info)
	if err != nil {
		return nil, err
	}
	return manager, nil
}
