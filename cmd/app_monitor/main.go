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

	monitorConf := &service.MonitorConfig{
		UrlsToMonitor:   viper.GetStringSlice("SPRINGBOOT_APPLICATION_BASE_URLS"),
		MonitorInterval: time.Minute * time.Duration(viper.GetInt("MONITOR_INTERVAL_MINUTES")),
		AppLogoPath:     appLogoPath,
	}

	monitorService := service.NewMonitorService(*monitorConf)
	monitorService.Start()
}
