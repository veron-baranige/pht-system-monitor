package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/veron-baranige/pht-system-monitor/pkg/utils"
)

type (
	Metrics struct {
		CpuUsage    float64
		DiskTotal   float64
		DiskUsed    float64
		MemoryUsed  float64
		MemoryTotal float64
	}
)

const (
	cpuUsageEndpoint      = "/actuator/metrics/system.cpu.usage"
	jvmMaxMemoryEndpoint  = "/actuator/metrics/jvm.memory.max"
	jvmUsedMemoryEndpoint = "/actuator/metrics/jvm.memory.used"
	diskTotalEndpoint     = "/actuator/metrics/disk.total"
	diskFreeEndpoint      = "/actuator/metrics/disk.free"
)

func GetMetrics(ctx context.Context, appBaseUrl string) (Metrics, error) {
	cpuUsage, err := getCpuUsage(ctx, appBaseUrl)
	if err != nil {
		return Metrics{}, err
	}

	jvmMaxMemory, err := getJvmMaxMemory(ctx, appBaseUrl)
	if err != nil {
		return Metrics{}, err
	}

	jvmUsedMemory, err := getJvmUsedMemory(ctx, appBaseUrl)
	if err != nil {
		return Metrics{}, err
	}

	diskTotalSpace, err := getDiskTotalSpace(ctx, appBaseUrl)
	if err != nil {
		return Metrics{}, err
	}

	diskFreeSpace, err := getDiskFreeSpace(ctx, appBaseUrl)
	if err != nil {
		return Metrics{}, err
	}

	return Metrics{
		CpuUsage:    cpuUsage,
		MemoryTotal: jvmMaxMemory,
		MemoryUsed:  jvmUsedMemory,
		DiskTotal:   diskTotalSpace,
		DiskUsed:    diskFreeSpace,
	}, nil
}

func getCpuUsage(ctx context.Context, appBaseUrl string) (float64, error) {
	value, err := getMeasurementValue(ctx, appBaseUrl+cpuUsageEndpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to get cpu usage: %v", err)
	}

	usagePercentage := value * 100
	return usagePercentage, nil
}

func getJvmMaxMemory(ctx context.Context, appBaseUrl string) (float64, error) {
	value, err := getMeasurementValue(ctx, appBaseUrl+jvmMaxMemoryEndpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to get jvm max memory: %v", err)
	}

	maxMemoryGb, err := utils.ConvertBytes(value, utils.Gigabytes)
	if err != nil {
		return 0, fmt.Errorf("failed to convert jvm max memory to gigabytes: %v", err)
	}

	return maxMemoryGb, nil
}

func getJvmUsedMemory(ctx context.Context, appBaseUrl string) (float64, error) {
	value, err := getMeasurementValue(ctx, appBaseUrl+jvmUsedMemoryEndpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to get jvm used memory: %v", err)
	}

	usedMemoryGb, err := utils.ConvertBytes(value, utils.Gigabytes)
	if err != nil {
		return 0, fmt.Errorf("failed to convert jvm used memory to gigabytes: %v", err)
	}

	return usedMemoryGb, nil
}

func getDiskTotalSpace(ctx context.Context, appBaseUrl string) (float64, error) {
	value, err := getMeasurementValue(ctx, appBaseUrl+diskTotalEndpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to get disk total space: %v", err)
	}

	totalDiskSpaceGb, err := utils.ConvertBytes(value, utils.Gigabytes)
	if err != nil {
		return 0, fmt.Errorf("failed to convert disk total space to gigabytes: %v", err)
	}

	return totalDiskSpaceGb, nil
}

func getDiskFreeSpace(ctx context.Context, appBaseUrl string) (float64, error) {
	value, err := getMeasurementValue(ctx, appBaseUrl+diskFreeEndpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to get disk free space: %v", err)
	}

	freeDiskSpaceGb, err := utils.ConvertBytes(value, utils.Gigabytes)
	if err != nil {
		return 0, fmt.Errorf("failed to convert disk free space to gigabytes: %v", err)
	}

	return freeDiskSpaceGb, nil
}

func getMeasurementValue(ctx context.Context, metricUrl string) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, metricUrl, nil)
	if err != nil {
		return 0, fmt.Errorf("client: could not create request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("client: error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("client: error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("client: received non-200 response: %s", body)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("client: error unmarshaling response: %v", err)
	}

	measurements, ok := data["measurements"].([]interface{})
	if !ok || len(measurements) == 0 {
		return 0, fmt.Errorf("client: unexpected response format: %s", body)
	}

	measurement, ok := measurements[0].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("client: unexpected measurement format: %s", body)
	}

	value, ok := measurement["value"].(float64)
	if !ok {
		return 0, fmt.Errorf("client: unexpected value format: %s", body)
	}

	return value, nil
}
