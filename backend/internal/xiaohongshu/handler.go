package xiaohongshu

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler HTTP 请求处理器
type Handler struct {
	service *Service
}

// NewHandler 创建新的处理器
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	xiaohongshu := router.Group("/xiaohongshu")
	{
		xiaohongshu.POST("/copy", h.GenerateCopy)
		xiaohongshu.GET("/scenes", h.GetScenes)
	}
}

// GenerateCopy 生成文案接口
// @Summary 生成小红书文案
// @Description 根据场景和配置生成小红书风格的文案
// @Accept json
// @Produce json
// @Param request body CopyRequest true "文案生成请求"
// @Success 200 {object} CopyResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /xiaohongshu/copy [post]
func (h *Handler) GenerateCopy(c *gin.Context) {
	var req CopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数: " + err.Error(),
		})
		return
	}

	if req.Scene == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "场景不能为空",
		})
		return
	}

	if req.Config == nil || len(req.Config) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "配置不能为空",
		})
		return
	}

	resp, err := h.service.GenerateCopy(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成文案失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetScenes 获取所有可用场景
// @Summary 获取文案场景列表
// @Description 获取所有支持的小红书文案生成场景
// @Produce json
// @Success 200 {array} Scene
// @Router /xiaohongshu/scenes [get]
func (h *Handler) GetScenes(c *gin.Context) {
	scenes := h.service.GetScenes()
	c.JSON(http.StatusOK, scenes)
}
