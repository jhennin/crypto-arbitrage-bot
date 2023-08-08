package internal

import (
	"net/http"
	"time"
)

var (
	SharedHttpClient *http.Client
)

/*
Returns the shared http client.
*/
func GetClient() *http.Client {
	return SharedHttpClient
}

/*
Intializes the shared http client.
*/
func init() {
	SharedHttpClient = &http.Client{Timeout: time.Second * 3}
}
