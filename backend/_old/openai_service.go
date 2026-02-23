package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Choice 表示响应中的选择项
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// ChatCompletionRequest 聊天完成请求参数
type ChatCompletionRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

// ChatCompletionResponse 聊天完成响应
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

// Usage 表示 token 使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamChunk 表示流式响应的数据块
type StreamChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice 表示流式响应中的选择项
type StreamChoice struct {
	Index        int          `json:"index"`
	Delta        MessageDelta `json:"delta"`
	FinishReason *string      `json:"finish_reason"`
}

// MessageDelta 表示流式响应中的消息增量
type MessageDelta struct {
	Role    *string `json:"role"`
	Content *string `json:"content"`
}

// Logger 日志接口
type Logger interface {
	Info(msg string, fields ...map[string]interface{})
	Error(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
}

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	logger *log.Logger
}

// NewDefaultLogger 创建默认日志实例
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "[OpenAI] ", log.LstdFlags),
	}
}

// Info 记录信息日志
func (l *DefaultLogger) Info(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("%s %+v", msg, fields[0])
	} else {
		l.logger.Println(msg)
	}
}

// Error 记录错误日志
func (l *DefaultLogger) Error(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("ERROR: %s %+v", msg, fields[0])
	} else {
		l.logger.Printf("ERROR: %s", msg)
	}
}

// Warn 记录警告日志
func (l *DefaultLogger) Warn(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("WARN: %s %+v", msg, fields[0])
	} else {
		l.logger.Printf("WARN: %s", msg)
	}
}

// OpenAIService OpenAI 服务封装
type OpenAIService struct {
	apiKey  string
	baseURL string
	client  *http.Client
	logger  Logger
}

// NewOpenAIService 创建新的 OpenAI 服务实例
// apiKey: OpenAI API 密钥
// logger: 日志记录器，如果为 nil 则使用默认日志
func NewOpenAIService(apiKey string, logger Logger) *OpenAIService {
	baseURL := DeepSeekBaseURL
	if logger == nil {
		logger = NewDefaultLogger()
	}

	if apiKey == "" {
		logger.Warn("DEEPSEEK_API_KEY is not set")
	} else {
		logger.Info("OpenAI instance created successfully with DeepSeek API", map[string]interface{}{
			"baseURL":          baseURL,
			"apiKeyConfigured": apiKey != "",
		})
	}

	return &OpenAIService{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
		logger:  logger,
	}
}

// CreateChatCompletion 创建聊天完成（非流式）
// req: 聊天完成请求参数
// 返回: 聊天完成响应
func (s *OpenAIService) CreateChatCompletion(req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	startTime := time.Now()

	s.logRequest("Creating chat completion", req)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		})
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	duration := time.Since(startTime)
	s.logger.Info("Chat completion created successfully", map[string]interface{}{
		"id":       response.ID,
		"model":    response.Model,
		"duration": duration.Milliseconds(),
	})

	return &response, nil
}

// CreateChatCompletionStream 创建流式聊天完成
// req: 聊天完成请求参数
// onChunk: 接收每个数据块的回调函数
// 返回: 完整的响应内容
func (s *OpenAIService) CreateChatCompletionStream(req ChatCompletionRequest, onChunk func(string)) (string, error) {
	startTime := time.Now()

	s.logRequest("Creating streaming chat completion", req)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		s.logger.Error("Failed to send request", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		})
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var fullResponse strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if line == "data: [DONE]" {
			break
		}

		if len(line) > 6 && line[:6] == "data: " {
			var chunk StreamChunk
			if err := json.Unmarshal([]byte(line[6:]), &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) > 0 {
				delta := chunk.Choices[0].Delta
				if delta.Content != nil {
					content := *delta.Content
					fullResponse.WriteString(content)
					onChunk(content)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		s.logger.Error("Error reading response", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("error reading response: %w", err)
	}

	duration := time.Since(startTime)
	s.logger.Info("Streaming chat completion finished", map[string]interface{}{
		"contentLength": fullResponse.Len(),
		"duration":      duration.Milliseconds(),
	})

	return fullResponse.String(), nil
}

// logRequest 记录请求日志（脱敏处理）
func (s *OpenAIService) logRequest(action string, req ChatCompletionRequest) {
	messages := make([]map[string]interface{}, len(req.Messages))
	for i, msg := range req.Messages {
		content := msg.Content
		if msg.Role == "user" && len(content) > 100 {
			content = content[:100] + "..."
		}
		messages[i] = map[string]interface{}{
			"role":    msg.Role,
			"content": content,
		}
	}

	s.logger.Info(action, map[string]interface{}{
		"model":       req.Model,
		"temperature": req.Temperature,
		"stream":      req.Stream,
		"messages":    messages,
	})
}

// GetDefaultOpenAIService 获取默认的 OpenAI 服务实例
// 从环境变量中读取 API 密钥
func GetDefaultOpenAIService() *OpenAIService {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	return NewOpenAIService(apiKey, nil)
}
