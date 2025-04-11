package osquery

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/unownone/osark-daemon/models"
)

// GetSystemInfo returns the system information
func (m *manager) GetSystemInfo() (*models.SystemInfo, error) {
	res, err := m.osClient.Query(getSystemInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get system info")
	}
	
	if res.Status.Code != 0 || len(res.Response) == 0 {
		return nil, errors.New("failed to get system info: " + res.Status.Message)
	}

	data := res.Response[0]

	systemInfo := &models.SystemInfo{
		OSName: data["name"],
		OSVersion: data["version"],
		OSArch: data["arch"],
	}

	if uptime, err := m.getUptime(); err != nil {
		return nil, errors.Wrap(err, "failed to get uptime")
	} else {
		systemInfo.UptimeSeconds = uptime
	}

	if macAddress, err := m.getMACAddress(); err != nil {
		return nil, errors.Wrap(err, "failed to get mac address")
	} else {
		systemInfo.MacAddress = macAddress
	}

	if osqueryVersion, err := m.getOSQueryVersion(); err != nil {
		return nil, errors.Wrap(err, "failed to get osquery version")
	} else {
		systemInfo.OSQueryVersion = osqueryVersion
	}

	return systemInfo, nil
}

// getUptime returns the uptime of the system
func (m *manager) getUptime() (time.Duration, error) {
	res, err := m.osClient.Query(getSystemUptime)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get uptime")
	}

	if res.Status.Code != 0 || len(res.Response) == 0 {
		return 0, errors.New("failed to get uptime: " + res.Status.Message)
	}

	data := res.Response[0]
	uptimeSeconds, err := strconv.ParseInt(data["total_seconds"], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse uptime seconds")
	}

	return time.Duration(uptimeSeconds) * time.Second, nil
}

// getMACAddress returns the mac address of the system
func (m *manager) getMACAddress() (string, error) {
	res, err := m.osClient.Query(getMACAddress)
	if err != nil {
		return "", errors.Wrap(err, "failed to get mac address")
	}

	if res.Status.Code != 0 || len(res.Response) == 0 {
		return "", errors.New("failed to get mac address: " + res.Status.Message)
	}

	data := res.Response[0]
	return data["address"], nil
}


func (m *manager) getOSQueryVersion() (string, error) {
	res, err := m.osClient.Query(getOSQueryVersion)
	if err != nil {
		return "", errors.Wrap(err, "failed to get osquery version")
	}

	if res.Status.Code != 0 || len(res.Response) == 0 {
		return "", errors.New("failed to get osquery version: " + res.Status.Message)
	}

	data := res.Response[0]
	return data["version"], nil
}