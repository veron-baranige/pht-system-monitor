package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/martinlindhe/notify"
	"github.com/spf13/viper"
	"github.com/veron-baranige/pht-system-monitor/internal/config"
)

func main() {
	log.Println("loading configurations")
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	log.Println("loaded configurations successfully")

	logoAbsPath, err := filepath.Abs(config.LogoPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, baseUrl := range viper.GetStringSlice("SPRINGBOOT_APPLICATION_BASE_URLS") {
		go func(baseUrl string) {
			notify.Alert(baseUrl, baseUrl, time.Now().Format("15:04") + " High Resource Usage [RAM: 2GB, CPU: ]", logoAbsPath)
		} (baseUrl)	
	}

	ticker := time.NewTicker(time.Minute * time.Duration(viper.GetInt("MONITOR_INTERVAL_MINUTES")))

	for range ticker.C {
		for _, baseUrl := range viper.GetStringSlice("SPRINGBOOT_APPLICATION_BASE_URLS") {
			go func(baseUrl string) {
				notify.Alert(baseUrl, baseUrl, time.Now().Format("15:04") + " High Resource Usage [RAM: 2GB, CPU: ]", logoAbsPath)
			} (baseUrl)	
		}
	}

	// /actuator/metrics/system.cpu.usage

}

// {
//     "status": "UP",
//     "components": {
//         "diskSpace": {
//             "status": "UP",
//             "details": {
//                 "total": 63278391296,
//                 "free": 15634952192,
//                 "threshold": 10485760,
//                 "path": "/pht/.",
//                 "exists": true
//             }
//         },
//         "mongo": {
//             "status": "UP",
//             "details": {
//                 "maxWireVersion": 17
//             }
//         },
//         "ping": {
//             "status": "UP"
//         }
//     }
// }