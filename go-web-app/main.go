package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func sendToLogstash(message string) {
	conn, err := net.Dial("tcp", "172.20.0.20:5044") // 连接到 Logstash
	if err != nil {
		log.Printf("Failed to connect to Logstash: %v", err)
		return
	}
	defer conn.Close()
	log.Printf("Sending log to Logstash: %s", message)
	_, err = conn.Write([]byte(message + "\n")) // 发送日志
	if err != nil {
		log.Printf("Failed to send log to Logstash: %v", err)
	}
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
		logMessage := formatLog(c, latency)
		log.Println(logMessage)
		sendToLogstash(logMessage) // 发送日志到 Logstash
	})

	// 示例路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, ELK!",
		})
	})

	r.GET("/error", func(c *gin.Context) {
		logMessage := "Error: Something went wrong!"
		log.Println(logMessage)
		sendToLogstash(logMessage) // 发送日志到 Logstash
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	})

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func formatLog(c *gin.Context, latency time.Duration) string {
	return fmt.Sprintf("[%s] %s %s %s %v",
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
		c.Request.UserAgent(),
		latency,
	)
}
