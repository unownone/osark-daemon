package osquery

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/unownone/osark-daemon/models"
)

// GetApps returns all the apps in the system
func (m *manager) GetApps() ([]*models.AppInfo, error) {
	res, err := m.osClient.Query(getAppsQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get apps")
	}

	if res.Status.Code != 0 {
		return nil, errors.New("failed to get apps: " + res.Status.Message)
	}

	apps := make([]*models.AppInfo, 0, len(res.Response))
	for _, app := range res.Response {
		appInfo := &models.AppInfo{
			Name: app["display_name"],
			BundleName: app["bundle_name"],
			BundleID: app["bundle_identifier"],
			BundleVersion: app["bundle_version"],
			Path: app["path"],
		}
		if time, err := time.Parse(app["last_opened_time"], time.RFC3339);err == nil {
			appInfo.LastOpenedTime = time
		}
		apps = append(apps, appInfo)
	}
	return apps, nil
}

// GetCurrentRunningProcesses returns the current running processes
func (m *manager) GetCurrentRunningProcesses(bundleIDs []string) ([]*models.ProcessInfo, error) {
	query := fmt.Sprintf(getCurrentRunningProcesses, strings.Join(bundleIDs, ","))
	res, err := m.osClient.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current running processes")
	}
	if res.Status.Code != 0 {
		return nil, errors.New("failed to get current running processes: " + res.Status.Message)
	}
	
	processes := make([]*models.ProcessInfo, 0, len(res.Response))
	for _, process := range res.Response {
		processes = append(processes, &models.ProcessInfo{
			PID: process["pid"],
			Name: process["name"],
			BundleID: process["bundle_identifier"],
			BundleVersion: process["bundle_version"],
			Path: process["path"],
		})
	}
	return processes, nil
}
