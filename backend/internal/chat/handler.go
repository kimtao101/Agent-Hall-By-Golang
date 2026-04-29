package chat

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"agent-backend/internal/agent"
)

type Handler struct {
	agent *agent.Agent
}

func NewHandler(a *agent.Agent) *Handler {
	return &Handler{agent: a}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.healthHandler)
	router.POST("/chat", h.chatHandler)
	router.GET("/history", h.historyHandler)
	router.POST("/clear", h.clearHandler)
}

type chatRequest struct {
	Message string `json:"message" binding:"required"`
}

func (h *Handler) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (h *Handler) chatHandler(c *gin.Context) {
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

	err := h.agent.GenerateResponse(req.Message, func(chunk string) {
		c.Writer.WriteString(chunk)
		flusher.Flush()
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成响应失败"})
		return
	}
}

func (h *Handler) historyHandler(c *gin.Context) {
	history := h.agent.GetHistory()
	c.JSON(http.StatusOK, history)
}

func (h *Handler) clearHandler(c *gin.Context) {
	h.agent.ClearHistory()
	c.JSON(http.StatusOK, gin.H{"success": true})
}
