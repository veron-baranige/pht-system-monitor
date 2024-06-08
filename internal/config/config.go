package config

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultMonitorIntervalMins = 5
	defaultWarnThreshold       = 80
	requestTimeoutSeconds      = 25
	LogoPath                   = "assets/logo.png"
	ConnectivityTestUrl        = "https://google.com"
)

func Load() error {
	viper.SetConfigType("env")
	viper.AddConfigPath("config")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := validateConfig(); err != nil {
		return err
	}

	parseConfig()

	return nil
}

func validateConfig() error {
	if viper.GetString("SPRINGBOOT_APPLICATION_BASE_URLS") == "" {
		return errors.New("invalid or missing springboot application base URLs to monitor")
	}

	for _, rawBaseUrl := range strings.Split(viper.GetString("SPRINGBOOT_APPLICATION_BASE_URLS"), ",") {
		if u, err := url.Parse(rawBaseUrl); err != nil {
			return errors.New("failed to parse springboot application base URL: " + err.Error())
		} else if u.Scheme == "" && u.Host == "" {
			return errors.New("invalid springboot application base URL: " + rawBaseUrl)
		}
	}

	if viper.GetInt("MONITOR_INTERVAL_MINUTES") <= 0 {
		log.Printf("Invalid or missing monitor interval configuration. "+
			"Setting interval to default value: %v mins", defaultMonitorIntervalMins) // warn
		viper.Set("MONITOR_INTERVAL_MINUTES", defaultMonitorIntervalMins)
	}

	if viper.GetInt("CPU_USAGE_WARN_THRESHOLD") <= 0 {
		log.Printf("invalid or missing cpu warn threshold configuration. "+
			"setting threshold to default value: %v%%", defaultWarnThreshold) // warn
		viper.Set("CPU_USAGE_WARN_THRESHOLD", defaultWarnThreshold)
	}

	if viper.GetInt("JVM_MEMORY_USAGE_WARN_THRESHOLD") <= 0 {
		log.Printf("invalid or missing jvm memory usage warn threshold configuration. "+
			"setting threshold to default value: %v%%", defaultWarnThreshold) // warn
		viper.Set("JVM_MEMORY_USAGE_WARN_THRESHOLD", defaultWarnThreshold)
	}

	if !viper.GetBool("ENABLE_DESKTOP_ALERTS") && !viper.GetBool("ENABLE_EMAIL_ALERTS") {
		return errors.New("atleast one type of alerts should be enabled")
	}

	if viper.GetBool("ENABLE_EMAIL_ALERTS") && viper.GetString("EMAIL_ALERT_RECIPIENTS") == "" {
		return errors.New("no email recipients have been configured for email alerts")
	}

	if viper.GetBool("ENABLE_EMAIL_ALERTS") && !isValidSmtpConfig() {
		return errors.New("invalid or missing SMTP configurations for email alerts")
	}

	return nil
}

func parseConfig() {
	rawApplicationBaseUrls := viper.GetString("SPRINGBOOT_APPLICATION_BASE_URLS")
	if rawApplicationBaseUrls != "" {
		parsedUrls := []string{}
		for _, rawBaseUrl := range strings.Split(rawApplicationBaseUrls, ",") {
			parsedUrl, _ := url.Parse(rawBaseUrl)
			parsedUrls = append(parsedUrls, strings.TrimRight(parsedUrl.String(), "/"))
		}
		viper.Set("SPRINGBOOT_APPLICATION_BASE_URLS", parsedUrls)
	}

	rawEmailrecipients := viper.GetString("EMAIL_ALERT_RECIPIENTS")
	if rawEmailrecipients != "" {
		viper.Set("EMAIL_ALERT_RECIPIENTS", strings.Split(rawEmailrecipients, ","))
	}
}

func isValidSmtpConfig() bool {
	return viper.GetString("SMTP_HDESKTOPT") != "" && viper.GetInt("SMTP_PORT") != 0 &&
		viper.GetString("SMTP_USER") != "" && viper.GetString("SMTP_PASSWORD") != ""
}

func SetHttpClientConfig() {
	http.DefaultClient.Timeout = time.Second * requestTimeoutSeconds
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}
