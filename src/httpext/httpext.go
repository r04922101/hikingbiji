package httpext

import (
	"net/http"
	"time"
)

// NewHTTPClient news an HTTP client with customized transport and timeout
func NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			// connection pool size
			MaxIdleConns: 20,
			// the number of connection can be allowed to open per host basic
			MaxIdleConnsPerHost: 20,
			// connection remains open until idle connection timeout
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: time.Second * 10,
	}
}
