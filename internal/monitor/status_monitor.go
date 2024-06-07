package monitor

type (
	HealthStatus string
)

const (
	Up           HealthStatus = "UP"
	Down         HealthStatus = "DOWN"
	OutOfService HealthStatus = "OUT_OF_SERVICE"
	Unkown       HealthStatus = "UNKNOWN"
)

