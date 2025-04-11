package osquery

// App Data
const (
	// GetAllApps returns all the apps in the system
	getAppsQuery = `
		SELECT 
			display_name,
			bundle_name, 
			bundle_identifier, 
			bundle_version, 
			last_opened_time,
			path 
		FROM 
			apps;
		`
)

// System data
const (
	// getSystemInfo returns the system information
	getSystemInfo = `
	SELECT
		name,
		version,
		platform,
		platform_like,
		arch 
	FROM 
		os_version;`
	// getSystemUptime returns the uptime of the system
	getSystemUptime = `
	SELECT 
		total_seconds 
	FROM 
		uptime;`
	// getMACAddress returns the mac address of the system
	getMACAddress = `
	SELECT 
		mac 
	FROM 
		interface_details 
	WHERE 
		mac IS NOT NULL AND 
		mac != ''
	LIMIT 1;
	`
)