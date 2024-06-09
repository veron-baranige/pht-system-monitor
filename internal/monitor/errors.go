package monitor

import "errors"

var (
	ErrNotResponding     = errors.New("no response from server")
	ErrMetricsFetching   = errors.New("failed to obtain metrics")
	ErrNoActuatorSupport = errors.New("no actuator support")
)
