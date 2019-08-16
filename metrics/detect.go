package metrics

var BadHostMap = make(map[string]bool)

// check if a host is bad.
func IsBadHost(host Host) bool {
	if host.CpuUsage/100 > CPU_USAGE_THRESHOLD {
		return true
	}

	if host.MemUsage/100 > MEM_USAGE_THRESHOLD {
		return true
	}

	if host.DiskUsage/100 > DISK_USAGE_THRESHOLD {
		return true
	}

	return false
}
