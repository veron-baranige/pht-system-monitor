package utils

import (
	"net/http"
)

func IsConnectedToInternet(connTestUrl string) bool {
	_, err := http.Get(connTestUrl)
	return err == nil
}