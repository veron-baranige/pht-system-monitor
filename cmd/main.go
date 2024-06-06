package main

import (
	"log"

	"github.com/veron-baranige/pht-system-monitor/internal/config"
)

func main() {
	log.Println("loading configurations")
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	log.Println("loaded configurations successfully")
}
