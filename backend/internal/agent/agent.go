package agent

import (
	"os"
	"strings"
	"sync"

	anthropicpkg "agent-backend/internal/anthropic"
	"agent-backend/internal/openai"
)

const (
	aiTypeAnthropic = "ANTHROPIC"
	aiTypeDeepSeek  = "DEEPSEEK"
	modelAnthropic  = "claude-sonnet-4-6"
	modelDeepSeek   = "deepseek-chat"
)

// Agent 聊天代理结构体
type Agent struct {
	messages     []Message
	mu           sync.RWMutex
	openaiSvc    *openai.Service
	anthropicSvc *anthropicpkg.Service
	aiType       string
	logger       openai.Logger
}

// New 创建新的 Agent 实例
func New(openaiSvc *openai.Service, anthropicSvc *anthropicpkg.Service) *Agent {
	if openaiSvc == nil {
		openaiSvc = openai.NewDefaultService()
	}
	if anthropicSvc == nil {
		anthropicSvc = anthropicpkg.NewDefaultService()
	}

	aiType := strings.ToUpper(strings.TrimSpace(os.Getenv("AI_TYPE")))
	if aiType != aiTypeAnthropic && aiType != aiTypeDeepSeek {
		aiType = aiTypeDeepSeek
	}

	a := &Agent{
		messages:     []Message{},
		openaiSvc:    openaiSvc,
		anthropicSvc: anthropicSvc,
		aiType:       aiType,
		logger:       openaiSvc.Logger,
	}
	a.initializeSystemPrompt()
	return a
}

// initializeSystemPrompt 初始化系统提示词
func (a *Agent) initializeSystemPrompt() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.messages = append(a.messages, Message{
		Role:    "system",
		Content: "You are a helpful research assistant. Provide detailed and accurate responses to user queries.",
	})
}

// AddMessage 添加消息到历史记录
func (a *Agent) AddMessage(role string, content string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.messages = append(a.messages, Message{Role: role, Content: content})
	a.trimHistory()
}

// trimHistory 修剪历史记录，保持最多 20 条消息（系统提示词 + 19 条）
func (a *Agent) trimHistory() {
	if len(a.messages) > 20 {
		a.messages = append([]Message{a.messages[0]}, a.messages[len(a.messages)-19:]...)
	}
}

// GenerateResponse 生成响应
func (a *Agent) GenerateResponse(userInput string, onChunk func(string)) error {
	a.logger.Info("收到聊天请求", map[string]interface{}{
		"message": truncateString(userInput, 100) + "...",
		"aiType":  a.aiType,
	})

	a.AddMessage("user", userInput)

	var (
		fullResponse string
		err          error
	)

	if a.aiType == aiTypeAnthropic {
		fullResponse, err = a.generateAnthropicResponse(onChunk)
	} else {
		fullResponse, err = a.generateDeepSeekResponse(onChunk)
	}

	if err != nil {
		a.logger.Error("生成响应失败", map[string]interface{}{"error": err.Error()})
		errorMessage := "抱歉，我在处理您的请求时遇到了错误。请稍后重试。"
		onChunk(errorMessage)
		a.AddMessage("assistant", errorMessage)
		return err
	}

	a.AddMessage("assistant", fullResponse)
	a.logger.Info("聊天响应完成", map[string]interface{}{"responseLength": len(fullResponse)})
	return nil
}

// generateAnthropicResponse 使用 Anthropic API 生成响应
func (a *Agent) generateAnthropicResponse(onChunk func(string)) (string, error) {
	history := a.getHistoryCopy()

	// 分离 system prompt 和对话消息
	var systemContent string
	var chatMessages []anthropicpkg.Message

	for _, msg := range history {
		if msg.Role == "system" {
			systemContent = msg.Content
		} else {
			chatMessages = append(chatMessages, anthropicpkg.Message{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	req := anthropicpkg.ChatRequest{
		Messages:  chatMessages,
		System:    systemContent,
		Model:     modelAnthropic,
		MaxTokens: 2048,
		Stream:    true,
	}

	return a.anthropicSvc.CreateChatCompletionStream(req, onChunk)
}

// generateDeepSeekResponse 使用 DeepSeek API 生成响应
func (a *Agent) generateDeepSeekResponse(onChunk func(string)) (string, error) {
	messages := a.getHistoryCopy()

	req := openai.ChatCompletionRequest{
		Messages:    messages,
		Model:       modelDeepSeek,
		Temperature: 0.7,
		Stream:      true,
	}

	return a.openaiSvc.CreateChatCompletionStream(req, onChunk)
}

// GetHistory 获取消息历史记录的副本
func (a *Agent) GetHistory() []Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]Message, len(a.messages))
	copy(result, a.messages)
	return result
}

// getHistoryCopy 获取消息历史记录的副本（内部使用）
func (a *Agent) getHistoryCopy() []openai.Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]openai.Message, len(a.messages))
	for i, msg := range a.messages {
		result[i] = openai.Message{Role: msg.Role, Content: msg.Content}
	}
	return result
}

// ClearHistory 清空消息历史记录并重新初始化系统提示词
func (a *Agent) ClearHistory() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.messages = []Message{}
	a.initializeSystemPrompt()
	a.logger.Info("已清空聊天历史", nil)
}

// GetMessagesCount 获取当前消息数量
func (a *Agent) GetMessagesCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.messages)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
