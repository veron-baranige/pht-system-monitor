package utils

import "fmt"

type (
	MemoryUnit string
)

const (
	Bytes     MemoryUnit = "b"
	Kilobytes MemoryUnit = "kb"
	Megabytes MemoryUnit = "mb"
	Gigabytes MemoryUnit = "gb"
)

func ConvertBytes(bytes float64, unit MemoryUnit) (float64, error) {
	const (
		kilobytesInBytes = 1024
		megabytesInBytes = 1024 * kilobytesInBytes
		gigabytesInBytes = 1024 * megabytesInBytes
	)

	switch unit {
	case Kilobytes:
		return bytes / kilobytesInBytes, nil
	case Megabytes:
		return bytes / megabytesInBytes, nil
	case Gigabytes:
		return bytes / gigabytesInBytes, nil
	default:
		return 0, fmt.Errorf("invalid storage unit: %s", unit)
	}
}
