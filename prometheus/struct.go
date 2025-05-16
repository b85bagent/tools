package tools

import "net/http"

// Client 負責與 Prometheus API 進行通信
type Client struct {
	endpoint        string
	username        string
	password        string
	clientCert      string
	clientKey       string
	caCert          string
	enableTLSVerify bool
	httpClient      *http.Client
}
