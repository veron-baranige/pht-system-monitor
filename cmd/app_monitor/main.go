package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"github.com/veron-baranige/springboot-app-monitor/internal/config"
	"github.com/veron-baranige/springboot-app-monitor/internal/service"
	"gopkg.in/gomail.v2"
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

	mailDialer := gomail.NewDialer(
		viper.GetString("SMTP_HOST"),
		viper.GetInt("SMTP_PORT"),
		viper.GetString("SMTP_USER"),
		viper.GetString("SMTP_PASSWORD"),
	)

	monitorConf := &service.MonitorConfig{
		UrlsToMonitor:          viper.GetStringSlice("SPRINGBOOT_APPLICATION_BASE_URLS"),
		MonitorInterval:        time.Minute * viper.GetDuration("MONITOR_INTERVAL_MINUTES"),
		AppLogoPath:            appLogoPath,
		TestConnectivityUrl:    config.ConnectivityTestUrl,
		CpuUsageWarnThreshold:  viper.GetUint32("CPU_USAGE_WARN_THRESHOLD"),
		JvmUsageWarnThreshold:  viper.GetUint32("JVM_MEMORY_USAGE_WARN_THRESHOLD"),
		IsDesktopAlertsEnabled: viper.GetBool("ENABLE_DESKTOP_ALERTS"),
		IsEmailAlertsEnabled:   viper.GetBool("ENABLE_EMAIL_ALERTS"),
		MailDialer:             mailDialer,
		EmailReceipients:       viper.GetStringSlice("EMAIL_ALERT_RECIPIENTS"),
	}

	monitorService := service.NewMonitorService(*monitorConf)
	monitorService.Start()
}
