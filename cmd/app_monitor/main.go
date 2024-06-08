package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/martinlindhe/notify"
	"github.com/spf13/viper"
	"github.com/veron-baranige/springboot-app-monitor/internal/config"
	"github.com/veron-baranige/springboot-app-monitor/internal/monitor"
	"github.com/veron-baranige/springboot-app-monitor/pkg/utils"
)

func main() {
	log.Println("loading configurations")
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	log.Println("loaded configurations successfully")

	log.Println("setting up application")
	config.SetHttpClientConfig()
	appLogoPath, err := filepath.Abs(config.LogoPath)
	if err != nil {
		log.Fatal(err)
	}
	monitorUrls := viper.GetStringSlice("SPRINGBOOT_APPLICATION_BASE_URLS")

	log.Println("application started to monitor health status and metrics")
	monitorHealthAndMetrics(appLogoPath, monitorUrls)

	ticker := time.NewTicker(time.Minute * time.Duration(viper.GetInt("MONITOR_INTERVAL_MINUTES")))
	defer ticker.Stop()

	for range ticker.C {
		monitorHealthAndMetrics(appLogoPath, monitorUrls)
	}
}

func monitorHealthAndMetrics(appLogoPath string, monitorUrls []string) {
	if !utils.IsConnectedToInternet(config.ConnectivityTestUrl) {
		msg := fmt.Sprintf("[%s] NO INTERNET CONNECTIVITY", time.Now().Format("15:04"))
		notify.Notify("PHT System Monitor", "PHT System Monitor", msg, appLogoPath)
		return
	}

	for _, baseUrl := range monitorUrls {
		go func(baseUrl string) {
			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Minute)
			defer cancel()

			status, err := monitor.GetHealthStatus(ctx, baseUrl)
			if err != nil {
				notify.Alert(baseUrl, baseUrl, time.Now().Format("15:04") + " " + string(status), appLogoPath)
			}

			if status != monitor.Up {
				notify.Alert(baseUrl, baseUrl, time.Now().Format("15:04") + " " + string(status), appLogoPath)
			}

			metrics, err := monitor.GetMetrics(ctx, baseUrl)
			if err != nil {
				return
			}
			
			content := fmt.Sprintf("[%v] CPU: %.2f%%, RAM: %.2f/%.2f GB, DISK: %.2f/%.2f GB",
				time.Now().Format("15:04"), metrics.CpuUsage, metrics.MemoryUsed, metrics.MemoryTotal, metrics.DiskUsed, metrics.DiskTotal)
			notify.Notify(baseUrl, baseUrl, content, appLogoPath)
		} (baseUrl)	

		// wait 5 secs before monitoring the next app to give notification read time
		time.Sleep(5 * time.Second) 
	}
}