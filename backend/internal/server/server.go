package server

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"agent-backend/internal/agent"
	"agent-backend/internal/anthropic"
	"agent-backend/internal/chat"
	"agent-backend/internal/openai"
	"agent-backend/internal/xiaohongshu"
)

// Server HTTP 服务器
type Server struct {
	router *gin.Engine
	port   string
}

// New 创建新的服务器实例
func New(port string) *Server {
	if port == "" {
		port = "8016"
	}

	aiType := strings.ToUpper(strings.TrimSpace(os.Getenv("AI_TYPE")))

	openaiSvc := openai.NewDefaultService()
	anthropicSvc := anthropic.NewDefaultService()

	agentInstance := agent.New(openaiSvc, anthropicSvc, aiType)
	xhsSvc := xiaohongshu.NewService(openaiSvc, anthropicSvc, aiType)

	r := gin.Default()
	r.Use(corsMiddleware())
	r.Use(gin.Recovery())
	r.Use(rateLimitMiddleware(100, 15*time.Minute))

	s := &Server{router: r, port: port}
	s.registerRoutes(agentInstance, xhsSvc)

	return s
}

// registerRoutes 注册所有路由
func (s *Server) registerRoutes(agentInstance *agent.Agent, xhsSvc *xiaohongshu.Service) {
	chatHandler := chat.NewHandler(agentInstance)
	chatHandler.RegisterRoutes(s.router)

	xhsHandler := xiaohongshu.NewHandler(xhsSvc)
	xhsHandler.RegisterRoutes(s.router.Group("/"))
}

// Run 启动服务器
func (s *Server) Run() error {
	return s.router.Run(":" + s.port)
}
