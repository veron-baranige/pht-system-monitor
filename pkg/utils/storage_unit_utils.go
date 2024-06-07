package utils

import "fmt"

type (
	StorageUnit string
)

const (
	Bytes     StorageUnit = "b"
	Kilobytes StorageUnit = "kb"
	Megabytes StorageUnit = "mb"
	Gigabytes StorageUnit = "gb"
)

func ConvertBytes(bytes float64, unit StorageUnit) (float64, error) {
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
