# Go Tools SDK

## 簡介
這是一個用 Go 語言開發的 SDK，提供以下功能模組，方便整合到其他專案中：
- **Logger**: 提供日誌記錄功能。
- **Prometheus**: 與 Prometheus 進行整合的客戶端。
- **RabbitMQ**: 與 RabbitMQ 進行整合的客戶端。
- **Splunk HEC**: 寫入 Splunk HTTP Event Collector，支援 ACK、重試與 Cluster 容錯。


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

### 🔍 Splunk HEC Client / Cluster

支援單一 HEC endpoint 或多節點 cluster，含 gzip 壓縮、ACK 等待與 raw/batch 送出。
- `NewClient(serverURL, token)` 單點；`NewCluster([]string, token)` 支援容錯/負載分散。
- `WriteEvent`/`WriteBatch` 傳送結構化事件，可設定 host/index/source/sourcetype/time/fields。
- `WriteRaw` 送原始文字流，會依 `maxContentLength` 自動分片；可搭配 `SetCompression("gzip")`。
- `WaitForAcknowledgement` 等待 Splunk HEC ACK（需開啟 channel ACK）。

```go
    // 也可改用 splunkhec.NewClient("https://splunk.example.com:8088", "token")
    hec := splunkhec.NewCluster([]string{
        "https://splunk-1.example.com:8088",
        "https://splunk-2.example.com:8088",
    }, "your-hec-token")

    hec.SetCompression("gzip") // 可選：縮小 payload
    hec.SetMaxRetry(3)         // Cluster 會自動換下一個 endpoint

    event := splunkhec.NewEvent(map[string]any{
        "level": "info",
        "msg":   "service started",
    })
    event.SetIndex("main")
    event.SetSourceType("app:log")

    if err := hec.WriteEvent(event); err != nil {
        log.Fatal(err)
    }
    if err := hec.WaitForAcknowledgement(); err != nil {
        log.Printf("ack failed: %v", err)
    }

    // 也可以寫 raw（會依大小自動分片）
    raw := strings.NewReader("{\"msg\":\"raw log\"}\n{\"msg\":\"another\"}\n")
    _ = hec.WriteRaw(raw, &splunkhec.EventMetadata{
        Host:   splunkhec.String("app-01"),
        Source: splunkhec.String("stdout"),
    })

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

