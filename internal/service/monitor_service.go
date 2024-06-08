package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/martinlindhe/notify"
	"github.com/veron-baranige/springboot-app-monitor/internal/monitor"
)

type MonitorConfig struct {
	TestConnectivityUrl string
	UrlsToMonitor       []string
	MonitorInterval     time.Duration
	AppLogoPath         string
}

type MonitorService struct {
	config *MonitorConfig
}

func NewMonitorService(config MonitorConfig) *MonitorService {
	return &MonitorService{
		config: &config,
	}
}

func (ms *MonitorService) Start() {
	log.Println("monitoring service started")

	ms.monitorHealthAndMetrics()

	ticker := time.NewTicker(ms.config.MonitorInterval)
	defer ticker.Stop()

	for range ticker.C {
		ms.monitorHealthAndMetrics()
	}
}

func (ms *MonitorService) monitorHealthAndMetrics() {
	if !ms.hasInternetConnection() {
		msg := fmt.Sprintf("[%s] NO INTERNET CONNECTIVITY", time.Now().Format("15:04"))
		notify.Notify("PHT System Monitor", "PHT System Monitor", msg, ms.config.AppLogoPath)
		return
	}

	for i, baseUrl := range ms.config.UrlsToMonitor {
		go func(baseUrl string) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			status, err := monitor.GetHealthStatus(ctx, baseUrl)
			if err != nil {
				notify.Alert(baseUrl, baseUrl, time.Now().Format("15:04")+" "+string(status), ms.config.AppLogoPath)
			}

			if status != monitor.Up {
				notify.Alert(baseUrl, baseUrl, time.Now().Format("15:04")+" "+string(status), ms.config.AppLogoPath)
			}

			metrics, err := monitor.GetMetrics(ctx, baseUrl)
			if err != nil {
				return
			}

			content := fmt.Sprintf("[%v] CPU: %.2f%%, JVM: %.1f/%.1f GB, DISK: %.1f/%.1f GB",
				time.Now().Format("15:04"), metrics.CpuUsage*metrics.CpuUsage, metrics.MemoryUsed, metrics.MemoryTotal,
				metrics.DiskUsed, metrics.DiskTotal)
			notify.Notify(baseUrl, baseUrl, content, ms.config.AppLogoPath)
		}(baseUrl)

		// wait before monitoring next app to provide notification read time
		if len(ms.config.UrlsToMonitor) > 1 && i == len(ms.config.UrlsToMonitor)-1 {
			time.Sleep(8 * time.Second)
		}
	}
}

func (ms *MonitorService) hasInternetConnection() bool {
	_, err := http.Get(ms.config.TestConnectivityUrl)
	return err == nil
}