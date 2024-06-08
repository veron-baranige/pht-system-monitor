package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

type (
	HealthStatus string
)

const (
	Up           HealthStatus = "UP"
	Down         HealthStatus = "DOWN"
	OutOfService HealthStatus = "OUT_OF_SERVICE"
	Unknown      HealthStatus = "UNKNOWN"

	healthEndpoint = "/actuator/health"
)

func GetHealthStatus(ctx context.Context, appBaseUrl string) (HealthStatus, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, appBaseUrl+healthEndpoint, nil)
	if err != nil {
		return Unknown, fmt.Errorf("http client could not create request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return "", ErrNotResponding
		}
		return Unknown, fmt.Errorf("http client error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Unknown, fmt.Errorf("http client error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Unknown, fmt.Errorf("http client received non-200 response: %s", body)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return Unknown, fmt.Errorf("http client error unmarshaling response: %v", err)
	}

	status, ok := data["status"].(string)
	if !ok {
		return Unknown, fmt.Errorf("http client unexpected response format: %s", body)
	}

	return HealthStatus(status), err
}
