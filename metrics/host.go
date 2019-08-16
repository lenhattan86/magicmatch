package metrics

// metrics on host
type Host struct {
	Name string

	CpuUsage     float64
	MemUsage  float64
	DiskUsage float64
}
