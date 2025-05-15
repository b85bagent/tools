# Go Tools SDK

## 簡介
這是一個用 Go 語言開發的 SDK，提供以下功能模組，方便整合到其他專案中：
- **Logger**: 提供日誌記錄功能。
- **Prometheus**: 與 Prometheus 進行整合的客戶端。
- **RabbitMQ**: 與 RabbitMQ 進行整合的客戶端。

## 安裝
在您的專案中，使用以下命令安裝此 SDK：
```bash
go get -u <repository-url>
```

## 使用方式

### Logger
用於記錄日誌：
```go
import "path/to/logger"

func main() {
    logger := logger.NewLogger()
    logger.Info("這是一條日誌訊息")
}
```

### Prometheus
用於監控數據：
```go
import "path/to/prometheus"

func main() {
    client := prometheus.NewClient()
    client.Monitor()
}
```

### RabbitMQ
用於消息隊列的操作：
```go
import "path/to/rabbitmq"

func main() {
    client := rabbitmq.NewClient()
    client.Publish("queue_name", "message")
}
```
