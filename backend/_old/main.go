package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	agent             *Agent
	xiaohongshuSvc    *XiaohongshuService
)

// ChatRequest 聊天请求
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

// XiaohongshuCopyRequest 小红书文案生成请求
type XiaohongshuCopyRequest struct {
	Scene  string                 `json:"scene" binding:"required"`
	Config map[string]interface{} `json:"config" binding:"required"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error     string `json:"error"`
	Timestamp string `json:"timestamp,omitempty"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func main() {
	if err := godotenv.Load(".ENV"); err != nil {
		log.Printf("Warning: Error loading .ENV file: %v", err)
	}

	agent = NewAgent(nil)
	xiaohongshuSvc = NewXiaohongshuService(nil, nil)

	router := gin.Default()

	router.Use(corsMiddleware())

	// 健康检查
	router.GET("/health", healthHandler)
	
	// 聊天相关接口
	router.POST("/chat", chatHandler)
	router.GET("/history", historyHandler)
	router.POST("/clear", clearHandler)
	
	// 小红书文案生成接口
	router.POST("/xiaohongshu/copy", xiaohongshuCopyHandler)

	log.Printf("Server starting on port %s", PORT)
	if err := router.Run(":" + PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// healthHandler 健康检查端点
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// chatHandler 聊天端点
func chatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid message",
		})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Message cannot be empty",
		})
		return
	}

	c.Header("Content-Type", "text/plain")
	c.Header("Transfer-Encoding", "chunked")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Streaming not supported",
		})
		return
	}

	err := agent.GenerateResponse(req.Message, func(chunk string) {
		c.Writer.WriteString(chunk)
		flusher.Flush()
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}
}

// historyHandler 获取聊天历史端点
func historyHandler(c *gin.Context) {
	history := agent.GetHistory()
	c.JSON(http.StatusOK, history)
}

// clearHandler 清除聊天历史端点
func clearHandler(c *gin.Context) {
	agent.ClearHistory()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// xiaohongshuCopyHandler 小红书文案生成端点
func xiaohongshuCopyHandler(c *gin.Context) {
	var req XiaohongshuCopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request: scene and config are required",
		})
		return
	}

	if req.Scene == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Scene cannot be empty",
		})
		return
	}

	if req.Config == nil || len(req.Config) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Config cannot be empty",
		})
		return
	}

	copyReq := CopyRequest{
		Scene:  req.Scene,
		Config: req.Config,
	}

	resp, err := xiaohongshuSvc.GenerateCopy(copyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate copy: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
