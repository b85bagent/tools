package tools

import (
	"io"

	"github.com/sirupsen/logrus"
)

// NewLogrusLogger 創建一個新的 logrus logger 實例
// level: 日誌等級 (trace, debug, info, warn, error, fatal, panic)
// output: 日誌輸出位置 (os.Stdout, os.Stderr, 或者其他 io.Writer 實現)
// 這個函數會根據傳入的 level 來設置 logger 的等級
func NewLogrusLogger(level string, output io.Writer) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006/01/02 15:04:05", // 設置時間格式
		FullTimestamp:   true,                  // 顯示完整的時間
	})

	// 設置螢幕顯示
	log.SetOutput(output)

	// 根據傳入的 level 來做設定
	switch level {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.InfoLevel) // 預設為 Info
	}

	return log
}
