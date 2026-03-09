package xiaohongshu

import (
	"os"
	"strings"

	anthropicpkg "agent-backend/internal/anthropic"
	"agent-backend/internal/openai"
	"agent-backend/pkg/utils"
)

// Service 小红书文案生成服务
type Service struct {
	openaiSvc    *openai.Service
	anthropicSvc *anthropicpkg.Service
	aiType       string
	logger       openai.Logger
}

// NewService 创建新的小红书文案生成服务
func NewService(openaiSvc *openai.Service, anthropicSvc *anthropicpkg.Service, logger openai.Logger) *Service {
	if openaiSvc == nil {
		openaiSvc = openai.NewDefaultService()
	}
	if anthropicSvc == nil {
		anthropicSvc = anthropicpkg.NewDefaultService()
	}
	if logger == nil {
		logger = openai.NewDefaultLogger()
	}

	aiType := strings.ToUpper(strings.TrimSpace(os.Getenv("AI_TYPE")))
	if aiType != "ANTHROPIC" && aiType != "DEEPSEEK" {
		aiType = "DEEPSEEK"
	}

	return &Service{
		openaiSvc:    openaiSvc,
		anthropicSvc: anthropicSvc,
		aiType:       aiType,
		logger:       logger,
	}
}

// GenerateCopy 生成小红书文案
func (s *Service) GenerateCopy(req CopyRequest) (*CopyResponse, error) {
	s.logger.Info("生成小红书文案", map[string]interface{}{
		"scene":  req.Scene,
		"aiType": s.aiType,
	})

	builder := NewPromptBuilder(req.Config)
	prompt := builder.Build(req.Scene)

	s.logger.Info("提示词已生成", map[string]interface{}{
		"promptPreview": utils.TruncateString(prompt, 100) + "...",
	})

	systemPrompt := "你是一位专业的小红书文案撰写专家，擅长撰写各种场景的优质文案。请根据用户提供的信息，生成符合小红书平台风格的文案，包含适当的表情符号、话题标签和互动引导语。"

	var (
		copy string
		err  error
	)

	if s.aiType == "ANTHROPIC" {
		copy, err = s.generateWithAnthropic(systemPrompt, prompt)
	} else {
		copy, err = s.generateWithDeepSeek(systemPrompt, prompt)
	}

	if err != nil {
		s.logger.Error("生成文案失败", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	s.logger.Info("文案生成成功", map[string]interface{}{"copyLength": len(copy)})
	return &CopyResponse{Copy: copy}, nil
}

func (s *Service) generateWithAnthropic(systemPrompt, userPrompt string) (string, error) {
	req := anthropicpkg.ChatRequest{
		Messages: []anthropicpkg.Message{
			{Role: "user", Content: userPrompt},
		},
		System:    systemPrompt,
		Model:     "claude-sonnet-4-6",
		MaxTokens: 2048,
		Stream:    false,
	}

	resp, err := s.anthropicSvc.CreateChatCompletion(req)
	if err != nil {
		return "", err
	}

	if len(resp.Content) > 0 {
		return resp.Content[0].Text, nil
	}
	return "", nil
}

func (s *Service) generateWithDeepSeek(systemPrompt, userPrompt string) (string, error) {
	messages := []openai.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	req := openai.ChatCompletionRequest{
		Messages:    messages,
		Model:       "deepseek-chat",
		Temperature: 0.8,
		Stream:      false,
	}

	resp, err := s.openaiSvc.CreateChatCompletion(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", nil
}

// GetScenes 获取所有可用场景
func (s *Service) GetScenes() []Scene {
	return AvailableScenes()
}
