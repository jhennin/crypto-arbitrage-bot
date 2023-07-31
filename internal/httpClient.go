package internal

import (
	"net/http"
	"time"
)

var (
	SharedHttpClient *http.Client
)

func GetClient() *http.Client {
	return SharedHttpClient
}
func init() {
	SharedHttpClient = &http.Client{Timeout: time.Second * 3}
}
