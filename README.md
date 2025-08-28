# Go Tools SDK

## 簡介
這是一個用 Go 語言開發的 SDK，提供以下功能模組，方便整合到其他專案中：
- **Logger**: 提供日誌記錄功能。
- **Prometheus**: 與 Prometheus 進行整合的客戶端。
- **RabbitMQ**: 與 RabbitMQ 進行整合的客戶端。


## 📦 安裝

```bash
go get github.com/b85bagent/tools@v0.2.0
```

---

## 📌 模組介紹

### 🐰 RabbitMQ Client

快速建立連線與發送訊息。

```go
client, err := tools.NewClient("amqp://guest:guest@localhost:5672/")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

err = client.Publish("my-exchange", "my-key", []byte(`{"msg":"hello"}`))
if err != nil {
    log.Printf("publish failed: %v", err)
}
```

---

### 📈 Prometheus Client (with TLS/mTLS)

支援 TLS / mTLS，包含：
- 自訂 CA 憑證
- 客戶端憑證與私鑰
- 可選擇跳過 TLS 驗證

```go
client, err := tools.NewClient(
    "https://prom.example.com",
    "username",
    "password",
    "./client.crt",
    "./client.key",
    "./ca.crt",
    false, // insecureTLS: true 表示跳過驗證
)
if err != nil {
    log.Fatal(err)
}

// 假設你有封裝 client.DoQuery(query string)
resp, err := client.DoQuery("up")
```

---

### 📋 Logrus Logger 初始化工具

快速建立一個可自訂等級與輸出的 logrus logger。

```go
log := tools.NewLogrusLogger("debug", os.Stdout)
log.Info("Logger is ready")
```

支援等級：`trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`  
輸出位置為任意符合 `io.Writer` 的實例（如檔案或 `os.Stdout`）

---

