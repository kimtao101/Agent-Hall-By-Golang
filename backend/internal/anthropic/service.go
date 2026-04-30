package anthropic

import (
	"agent-backend/pkg/logger"
	"agent-backend/pkg/utils"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	constantpkg "agent-backend/constant"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Service struct {
	client *anthropic.Client
	Logger logger.Logger
}

func NewService(apiKey, baseURL string, lg logger.Logger) *Service {
	if lg == nil {
		lg = logger.NewDefaultLogger("[Anthropic]")
	}
	if baseURL == "" {
		baseURL = constantpkg.ANBaseURL
	}

	if apiKey == "" {
		lg.Warn("Anthropic API Key 未设置")
	} else {
		lg.Info("Anthropic 服务初始化成功", map[string]interface{}{
			"baseURL": baseURL,
		})
	}

	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	return &Service{
		client: &client,
		Logger: lg,
	}
}

func NewDefaultService() *Service {
	apiKey := strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY"))
	baseURL := strings.TrimSpace(os.Getenv("ANTHROPIC_BASE_URL"))
	if baseURL == "" {
		baseURL = constantpkg.ANBaseURL
	}
	return NewService(apiKey, baseURL, nil)
}

func (s *Service) CreateChatCompletion(req ChatRequest) (*ChatResponse, error) {
	startTime := time.Now()
	s.logRequest("创建 Anthropic 聊天完成", req)

	model := req.Model
	if model == "" {
		model = constantpkg.ANModel
	}
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = constantpkg.DefaultMaxTokens
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: int64(maxTokens),
		Messages:  buildMessages(req.Messages),
	}
	if req.System != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: req.System},
		}
	}

	message, err := s.client.Messages.New(context.Background(), params)
	if err != nil {
		s.Logger.Error("Anthropic API 请求失败", map[string]interface{}{"error": err.Error()})
		return nil, fmt.Errorf("API 请求失败: %w", err)
	}

	duration := time.Since(startTime)
	resp := convertResponse(message)
	s.Logger.Info("Anthropic 聊天完成成功", map[string]interface{}{
		"id":       resp.ID,
		"model":    resp.Model,
		"duration": duration.Milliseconds(),
	})
	return resp, nil
}

func (s *Service) CreateChatCompletionStream(req ChatRequest, onChunk func(string)) (string, error) {
	startTime := time.Now()
	s.logRequest("创建 Anthropic 流式聊天完成", req)

	model := req.Model
	if model == "" {
		model = constantpkg.ANModel
	}
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = constantpkg.DefaultMaxTokens
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: int64(maxTokens),
		Messages:  buildMessages(req.Messages),
	}
	if req.System != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: req.System},
		}
	}

	stream := s.client.Messages.NewStreaming(context.Background(), params)
	defer stream.Close()

	var fullResponse strings.Builder
	for stream.Next() {
		event := stream.Current()
		if event.Type == "content_block_delta" && event.Delta.Type == "text_delta" {
			text := event.Delta.Text
			if text != "" {
				fullResponse.WriteString(text)
				onChunk(text)
			}
		}
	}
	if err := stream.Err(); err != nil {
		s.Logger.Error("流式响应失败", map[string]interface{}{"error": err.Error()})
		return "", fmt.Errorf("流式请求失败: %w", err)
	}

	duration := time.Since(startTime)
	s.Logger.Info("Anthropic 流式聊天完成成功", map[string]interface{}{
		"contentLength": fullResponse.Len(),
		"duration":      duration.Milliseconds(),
	})
	return fullResponse.String(), nil
}

func buildMessages(msgs []Message) []anthropic.MessageParam {
	var params []anthropic.MessageParam
	for _, msg := range msgs {
		switch msg.Role {
		case "user":
			params = append(params, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)))
		case "assistant":
			params = append(params, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
		}
	}
	return params
}

func convertResponse(msg *anthropic.Message) *ChatResponse {
	var blocks []ContentBlock
	for _, c := range msg.Content {
		if c.Type == "text" {
			blocks = append(blocks, ContentBlock{Type: "text", Text: c.Text})
		}
	}
	return &ChatResponse{
		ID:      msg.ID,
		Type:    string(msg.Type),
		Role:    string(msg.Role),
		Content: blocks,
		Model:   string(msg.Model),
		Usage: Usage{
			InputTokens:  int(msg.Usage.InputTokens),
			OutputTokens: int(msg.Usage.OutputTokens),
		},
	}
}

func (s *Service) logRequest(action string, req ChatRequest) {
	messages := make([]map[string]interface{}, len(req.Messages))
	for i, msg := range req.Messages {
		content := msg.Content
		if msg.Role == "user" && len(content) > 100 {
			content = utils.TruncateString(content, 100) + "..."
		}
		messages[i] = map[string]interface{}{
			"role":    msg.Role,
			"content": content,
		}
	}
	s.Logger.Info(action, map[string]interface{}{
		"model":     req.Model,
		"maxTokens": req.MaxTokens,
		"messages":  messages,
	})
}
