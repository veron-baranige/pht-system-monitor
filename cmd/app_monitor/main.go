package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"github.com/veron-baranige/springboot-app-monitor/internal/config"
	"github.com/veron-baranige/springboot-app-monitor/internal/service"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	log.Println("loaded configurations")

	log.Println("setting up application")
	config.SetHttpClientConfig()

	appLogoPath, err := filepath.Abs(config.LogoPath)
	if err != nil {
		log.Fatal(err)
	}

	monitorConf := &service.MonitorConfig{
		UrlsToMonitor:       viper.GetStringSlice("SPRINGBOOT_APPLICATION_BASE_URLS"),
		MonitorInterval:     time.Minute * viper.GetDuration("MONITOR_INTERVAL_MINUTES"),
		AppLogoPath:         appLogoPath,
		TestConnectivityUrl: config.ConnectivityTestUrl,
	}

	monitorService := service.NewMonitorService(*monitorConf)
	monitorService.Start()
}
