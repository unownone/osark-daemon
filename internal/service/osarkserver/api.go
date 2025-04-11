package osarkserver

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/unownone/osark-daemon/models"
)

// Authenticate authenticates the push manager
func (p *pushManager) Authenticate(info *models.SystemInfo) error {
	deviceID, err := p.generateDeviceID(info)
	if err != nil {
		return errors.Wrap(err, "failed to generate device ID")
	}
	p.deviceID = deviceID
	// TODO: send device ID to the server
	return nil
}

// generateDeviceID generates a device ID
func (p *pushManager) generateDeviceID(info *models.SystemInfo) (string, error) {
	concat := fmt.Sprintf("%s-%s-%s", info.OSName, info.OSArch, info.MacAddress)
	hash := sha256.Sum256([]byte(concat))
	return hex.EncodeToString(hash[:]), nil
}

// Push pushes the data to the server
func (p *pushManager) Push(data []*models.LogEvent) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}
	req, err := http.NewRequest("POST", p.osarkServerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Identifier", p.deviceID)
	resp, err := p.service.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send request")
	}
	return nil
}

// PushError pushes an error to the server
func (p *pushManager) PushError(err error) error {
	return p.Push([]*models.LogEvent{
		{
			Error: err.Error(),
		},
	})
}
