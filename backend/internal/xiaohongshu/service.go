package xiaohongshu

import (
	"agent-backend/internal/openai"
	"agent-backend/pkg/utils"
)

// Service 小红书文案生成服务
type Service struct {
	openai *openai.Service
	logger openai.Logger
}

// NewService 创建新的小红书文案生成服务
func NewService(openaiSvc *openai.Service, logger openai.Logger) *Service {
	if openaiSvc == nil {
		openaiSvc = openai.NewDefaultService()
	}
	if logger == nil {
		logger = openai.NewDefaultLogger()
	}
	return &Service{
		openai: openaiSvc,
		logger: logger,
	}
}

// GenerateCopy 生成小红书文案
func (s *Service) GenerateCopy(req CopyRequest) (*CopyResponse, error) {
	s.logger.Info("生成小红书文案", map[string]interface{}{
		"scene": req.Scene,
	})

	// 构建提示词
	builder := NewPromptBuilder(req.Config)
	prompt := builder.Build(req.Scene)

	s.logger.Info("提示词已生成", map[string]interface{}{
		"promptPreview": utils.TruncateString(prompt, 100) + "...",
	})

	// 调用 AI 生成文案
	messages := []openai.Message{
		{
			Role:    "system",
			Content: "你是一位专业的小红书文案撰写专家，擅长撰写各种场景的优质文案。请根据用户提供的信息，生成符合小红书平台风格的文案，包含适当的表情符号、话题标签和互动引导语。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	chatReq := openai.ChatCompletionRequest{
		Messages:    messages,
		Model:       "deepseek-chat",
		Temperature: 0.8,
		Stream:      false,
	}

	resp, err := s.openai.CreateChatCompletion(chatReq)
	if err != nil {
		s.logger.Error("生成文案失败", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	var copy string
	if len(resp.Choices) > 0 {
		copy = resp.Choices[0].Message.Content
	}

	s.logger.Info("文案生成成功", map[string]interface{}{
		"copyLength": len(copy),
	})

	return &CopyResponse{Copy: copy}, nil
}

// GetScenes 获取所有可用场景
func (s *Service) GetScenes() []Scene {
	return AvailableScenes()
}
