package monitor

import "context"

type (
	Metrics struct {
		CpuUsage    float64
		DiskTotal   float64
		DiskUsed    float64
		MemoryUsed  float64
		MemoryTotal float64
	}
)

func GetMetrics(ctx context.Context, appBaseUrl string) (Metrics, error) {
	return Metrics{}, nil
}
