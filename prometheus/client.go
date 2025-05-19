package tools

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

// NewClient 創建一個新的 Prometheus 客戶端
// 這個函數會根據提供的配置創建一個新的 HTTP 客戶端
// 並設置 TLS 配置，包括 CA 憑證和客戶端憑證
// 如果提供了 CA 憑證，則會將其加載到 TLS 配置中
// 如果提供了客戶端憑證和私鑰，則會將其加載到 TLS 配置中
// 如果設置了 InsecureSkipVerify，則會跳過對服務器憑證的驗證
// 返回一個新的 Prometheus 客戶端實例
func NewPrometheusClient(endpoint, username, password,
	clientCert, clientKey, caCert string,
	enableTLS bool) (*Client, error) {

	// 改成是否啟用 TLS 驗證 true 表示啟用，false 表示不啟用
	// 當 enableTLSVerify 為 true，InsecureSkipVerify 會是 false（不跳過驗證，安全）
	// 當 enableTLSVerify 為 false，InsecureSkipVerify 會是 true（跳過驗證，不安全）
	tlsConfig := &tls.Config{InsecureSkipVerify: !enableTLS}

	// 如果提供 CA 憑證，則加載 CA
	if caCert != "" {
		caCertData, err := os.ReadFile(caCert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA cert: %w", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCertData)
		tlsConfig.RootCAs = caCertPool
	}

	// 如果提供 Client Cert，則加載
	if clientCert != "" && clientKey != "" {
		cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// 建立 HTTP 客戶端
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	return &Client{
		Endpoint:   endpoint,
		Username:   username,
		Password:   password,
		clientCert: clientCert,
		clientKey:  clientKey,
		caCert:     caCert,
		enableTLS:  enableTLS,
		HttpClient: client,
	}, nil
}
