// server 提供 HTTP 服务器配置和路由管理
package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"agent-backend/internal/agent"
	"agent-backend/internal/anthropic"
	"agent-backend/internal/openai"
	"agent-backend/internal/xiaohongshu"
)

// Server HTTP 服务器
type Server struct {
	router  *gin.Engine
	port    string
	agent   *agent.Agent
	xhsSvc  *xiaohongshu.Service
}

// New 创建新的服务器实例
func New(port string) *Server {
	if port == "" {
		port = "8016"
	}

	// 创建 AI 服务实例
	openaiSvc := openai.NewDefaultService()
	anthropicSvc := anthropic.NewDefaultService()

	// 创建 Agent 和 Xiaohongshu 服务
	agentInstance := agent.New(openaiSvc, anthropicSvc)
	xhsSvc := xiaohongshu.NewService(openaiSvc, anthropicSvc, nil)

	// 创建 Gin 引擎
	r := gin.Default()

	s := &Server{
		router:  r,
		port:    port,
		agent:   agentInstance,
		xhsSvc:  xhsSvc,
	}

	// 配置中间件
	r.Use(corsMiddleware())
	r.Use(gin.Recovery())
	r.Use(rateLimitMiddleware(100, 15*time.Minute))

	// 注册路由
	s.registerRoutes()

	return s
}

// registerRoutes 注册所有路由
func (s *Server) registerRoutes() {
	// 健康检查（不受速率限制）
	s.router.GET("/health", s.healthHandler)

	// 基础聊天接口
	s.router.POST("/chat", s.chatHandler)
	s.router.GET("/history", s.historyHandler)
	s.router.POST("/clear", s.clearHandler)

	// 小红书智能体接口
	xhsHandler := xiaohongshu.NewHandler(s.xhsSvc)
	xhsHandler.RegisterRoutes(s.router.Group("/"))
}

// Run 启动服务器
func (s *Server) Run() error {
	return s.router.Run(":" + s.port)
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

// ipEntry IP 速率限制记录
type ipEntry struct {
	mu        sync.Mutex
	timestamps []time.Time
}

// rateLimiter 全局速率限制器
type rateLimiter struct {
	mu     sync.Mutex
	ips    map[string]*ipEntry
	limit  int
	window time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  limit,
		window: window,
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	entry, ok := rl.ips[ip]
	if !ok {
		entry = &ipEntry{}
		rl.ips[ip] = entry
	}
	rl.mu.Unlock()

	entry.mu.Lock()
	defer entry.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// 清除过期时间戳
	valid := entry.timestamps[:0]
	for _, t := range entry.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	entry.timestamps = valid

	if len(entry.timestamps) >= rl.limit {
		return false
	}

	entry.timestamps = append(entry.timestamps, now)
	return true
}

// rateLimitMiddleware 速率限制中间件（跳过 /health）
func rateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	rl := newRateLimiter(limit, window)
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		ip := c.ClientIP()
		if !rl.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests from this IP, please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// healthHandler 健康检查
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// chatRequest 聊天请求
type chatRequest struct {
	Message string `json:"message" binding:"required"`
}

// chatHandler 聊天接口（流式）
func (s *Server) chatHandler(c *gin.Context) {
	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息"})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息不能为空"})
		return
	}

	c.Header("Content-Type", "text/plain")
	c.Header("Transfer-Encoding", "chunked")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "不支持流式响应"})
		return
	}

	err := s.agent.GenerateResponse(req.Message, func(chunk string) {
		c.Writer.WriteString(chunk)
		flusher.Flush()
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成响应失败"})
		return
	}
}

// historyHandler 获取聊天历史
func (s *Server) historyHandler(c *gin.Context) {
	history := s.agent.GetHistory()
	c.JSON(http.StatusOK, history)
}

// clearHandler 清除聊天历史
func (s *Server) clearHandler(c *gin.Context) {
	s.agent.ClearHistory()
	c.JSON(http.StatusOK, gin.H{"success": true})
}
