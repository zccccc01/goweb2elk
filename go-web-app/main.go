package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogEntry 定义日志结构
type LogEntry struct {
	Method    string        `json:"method"`
	Path      string        `json:"path"`
	IP        string        `json:"ip"`
	UserAgent string        `json:"user_agent"`
	Latency   time.Duration `json:"latency"`
	Message   string        `json:"message,omitempty"` // 用于错误日志
}

func sendToLogstash(entry LogEntry) {
	// 将结构体转换为 JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry to JSON: %v", err)
		return
	}

	// 连接到 Logstash
	conn, err := net.Dial("tcp", "172.20.0.20:5044")
	if err != nil {
		log.Printf("Failed to connect to Logstash: %v", err)
		return
	}
	defer conn.Close()

	// 发送 JSON 日志
	_, err = conn.Write(jsonData)
	if err != nil {
		log.Printf("Failed to send log to Logstash: %v", err)
		return
	}

	log.Printf("Sending log to Logstash: %s", jsonData)
}

func main() {
	// 设置日志输出到文件和控制台
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log", // 日志文件路径
		MaxSize:    10,             // 日志文件最大大小（MB）
		MaxBackups: 3,              // 保留的旧日志文件最大数量
		MaxAge:     28,             // 保留旧日志文件的最大天数
		Compress:   true,           // 是否压缩旧日志文件
	})

	// 创建 Gin 路由
	r := gin.Default()

	// 记录访问日志的中间件
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// 构造日志结构体
		logEntry := LogEntry{
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Latency:   latency,
		}

		// 打印日志到控制台和文件
		log.Println(formatLog(logEntry))

		// 发送日志到 Logstash
		sendToLogstash(logEntry)
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, ELK!",
		})
	})

	r.GET("/error", func(c *gin.Context) {
		// 构造错误日志结构体
		logEntry := LogEntry{
			Message: "Error: Something went wrong!",
		}

		// 打印错误日志到控制台和文件
		log.Println(logEntry.Message)

		// 发送错误日志到 Logstash
		sendToLogstash(logEntry)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	})

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// formatLog 将日志结构体格式化为字符串
func formatLog(entry LogEntry) string {
	return fmt.Sprintf("[%s] %s %s %s %v",
		entry.Method,
		entry.Path,
		entry.IP,
		entry.UserAgent,
		entry.Latency,
	)
}
