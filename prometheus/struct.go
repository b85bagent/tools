package tools

import "net/http"

// Client 負責與 Prometheus API 進行通信
type Client struct {
	Endpoint   string
	Username   string
	Password   string
	clientCert string
	clientKey  string
	caCert     string
	enableTLS  bool
	HttpClient *http.Client
}
