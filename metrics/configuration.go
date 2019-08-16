package metrics

// configuration for tsdb queries
// clean up the sensitive data.
var (
	COLO      = ""
	HOST      = "*"
	TSDB_URL  = ""
	TIME_AGO = "1h-ago"
	CPU_METRIC       = ""
	MEM_USAGE_METRIC = ""
	DISK_FREE_METRIC = ""
	AGREGATOR = "avg"
	CPU_USAGE_THRESHOLD  = 0.9 // [0,1]
	MEM_USAGE_THRESHOLD  = 0.9 // [0,1]
	DISK_USAGE_THRESHOLD = 0.9 // [0,1]
)
