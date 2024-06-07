package monitor

import "context"

type (
	HealthStatus string
)

const (
	Up           HealthStatus = "UP"
	Down         HealthStatus = "DOWN"
	OutOfService HealthStatus = "OUT_OF_SERVICE"
	Unkown       HealthStatus = "UNKNOWN"
)

func GetHealthStatus(ctx context.Context, appBaseUrl string) (HealthStatus, error) {
	return Up, nil
}