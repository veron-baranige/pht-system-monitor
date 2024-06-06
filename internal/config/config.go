package config

import (
	"errors"
	"log"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultMonitorIntervalMins = 5
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

	parseEmailConfig()

	return nil
}

func validateConfig() error {
	if viper.GetInt("MONITOR_INTERVAL_MINUTES") <= 0 {
		log.Printf("Invalid or missing monitor interval configuration. "+
			"Setting interval to default value: %v mins", defaultMonitorIntervalMins) // warn
		viper.Set("MONITOR_INTERVAL_MINUTES", defaultMonitorIntervalMins)
	}

	if !viper.GetBool("ENABLE_OS_ALERTS") && !viper.GetBool("ENABLE_EMAIL_ALERTS") {
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

func parseEmailConfig() {
	rawEmailrecipients := viper.GetString("EMAIL_ALERT_RECIPIENTS")
	if rawEmailrecipients != "" {
		viper.Set("EMAIL_ALERT_RECIPIENTS", strings.Split(rawEmailrecipients, ","))
	}
}

func isValidSmtpConfig() bool {
	return viper.GetString("SMTP_HOST") != "" && viper.GetInt("SMTP_PORT") != 0 &&
		viper.GetString("SMTP_USER") != "" && viper.GetString("SMTP_PASSWORD") != ""
}
