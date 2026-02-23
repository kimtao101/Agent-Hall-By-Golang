package main

import (
	"sync"
)

// Agent 聊天代理结构体
type Agent struct {
	messages []Message
	mu       sync.RWMutex
	openai   *OpenAIService
	logger   Logger
}

// NewAgent 创建新的 Agent 实例
// openai: OpenAI 服务实例，如果为 nil 则使用默认配置
func NewAgent(openai *OpenAIService) *Agent {
	if openai == nil {
		openai = GetDefaultOpenAIService()
	}
	agent := &Agent{
		messages: []Message{},
		openai:   openai,
		logger:   openai.logger,
	}
	agent.initializeSystemPrompt()
	return agent
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
// role: 消息角色，可以是 "user" 或 "assistant"
// content: 消息内容
func (a *Agent) AddMessage(role string, content string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.messages = append(a.messages, Message{
		Role:    role,
		Content: content,
	})
	a.trimHistory()
}

// trimHistory 修剪历史记录，保持最多 20 条消息
func (a *Agent) trimHistory() {
	if len(a.messages) > 20 {
		a.messages = append([]Message{a.messages[0]}, a.messages[len(a.messages)-19:]...)
	}
}

// GenerateResponse 生成响应
// userInput: 用户输入
// onChunk: 接收流式响应数据块的回调函数
// 返回: 错误信息
func (a *Agent) GenerateResponse(userInput string, onChunk func(string)) error {
	a.logger.Info("Received chat request", map[string]interface{}{
		"message": truncateString(userInput, 100) + "...",
	})

	a.AddMessage("user", userInput)

	messages := a.getHistoryCopy()
	req := ChatCompletionRequest{
		Messages:    messages,
		Model:       "deepseek-chat",
		Temperature: 0.7,
		Stream:      true,
	}

	fullResponse, err := a.openai.CreateChatCompletionStream(req, onChunk)
	if err != nil {
		a.logger.Error("Error generating response", map[string]interface{}{
			"error": err.Error(),
		})
		errorMessage := "I apologize, but I encountered an error while processing your request. Please try again later."
		onChunk(errorMessage)
		a.AddMessage("assistant", errorMessage)
		return err
	}

	a.AddMessage("assistant", fullResponse)
	a.logger.Info("Chat response completed", map[string]interface{}{
		"responseLength": len(fullResponse),
	})
	return nil
}

// GetHistory 获取消息历史记录的副本
// 返回: 消息历史记录切片
func (a *Agent) GetHistory() []Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]Message, len(a.messages))
	copy(result, a.messages)
	return result
}

// getHistoryCopy 获取消息历史记录的副本（内部使用）
// 返回: 消息历史记录切片
func (a *Agent) getHistoryCopy() []Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]Message, len(a.messages))
	copy(result, a.messages)
	return result
}

// ClearHistory 清空消息历史记录并重新初始化系统提示词
func (a *Agent) ClearHistory() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.messages = []Message{}
	a.initializeSystemPrompt()
	a.logger.Info("Cleared chat history", nil)
}

// GetMessagesCount 获取当前消息数量
// 返回: 消息数量
func (a *Agent) GetMessagesCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.messages)
}

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
