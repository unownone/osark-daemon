package osquery

import (
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
