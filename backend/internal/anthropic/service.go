package anthropic

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

	"agent-backend/pkg/utils"
)

const (
	defaultBaseURL       = "https://api.aicodemirror.com/api/claudecode"
	defaultModel         = "claude-sonnet-4-6"
	anthropicVersion     = "2023-06-01"
	defaultMaxTokens     = 2048
)

// Logger 日志接口（复用 openai 包的接口定义）
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
		logger: log.New(os.Stdout, "[Anthropic] ", log.LstdFlags),
	}
}

func (l *DefaultLogger) Info(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("%s %+v", msg, fields[0])
	} else {
		l.logger.Println(msg)
	}
}

func (l *DefaultLogger) Error(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("ERROR: %s %+v", msg, fields[0])
	} else {
		l.logger.Printf("ERROR: %s", msg)
	}
}

func (l *DefaultLogger) Warn(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("WARN: %s %+v", msg, fields[0])
	} else {
		l.logger.Printf("WARN: %s", msg)
	}
}

// Service Anthropic API 服务封装
type Service struct {
	apiKey  string
	baseURL string
	client  *http.Client
	Logger  Logger
}

// NewService 创建新的 Anthropic 服务实例
func NewService(apiKey, baseURL string, logger Logger) *Service {
	if logger == nil {
		logger = NewDefaultLogger()
	}
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	if apiKey == "" {
		logger.Warn("Anthropic API Key 未设置")
	} else {
		logger.Info("Anthropic 服务初始化成功", map[string]interface{}{
			"baseURL": baseURL,
		})
	}

	return &Service{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
		Logger:  logger,
	}
}

// NewDefaultService 从环境变量创建默认 Anthropic 服务实例
func NewDefaultService() *Service {
	apiKey := strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY"))
	baseURL := strings.TrimSpace(os.Getenv("ANTHROPIC_BASE_URL"))
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return NewService(apiKey, baseURL, nil)
}

// CreateChatCompletion 创建聊天完成（非流式）
func (s *Service) CreateChatCompletion(req ChatRequest) (*ChatResponse, error) {
	startTime := time.Now()
	s.logRequest("创建 Anthropic 聊天完成", req)

	if req.Model == "" {
		req.Model = defaultModel
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = defaultMaxTokens
	}
	req.Stream = false

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/v1/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	s.setHeaders(httpReq)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.Logger.Error("Anthropic API 请求失败", map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		})
		return nil, fmt.Errorf("API 请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	duration := time.Since(startTime)
	s.Logger.Info("Anthropic 聊天完成成功", map[string]interface{}{
		"id":       response.ID,
		"model":    response.Model,
		"duration": duration.Milliseconds(),
	})

	return &response, nil
}

// CreateChatCompletionStream 创建流式聊天完成
// Anthropic SSE 格式: event: content_block_delta / data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"..."}}
func (s *Service) CreateChatCompletionStream(req ChatRequest, onChunk func(string)) (string, error) {
	startTime := time.Now()
	s.logRequest("创建 Anthropic 流式聊天完成", req)

	if req.Model == "" {
		req.Model = defaultModel
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = defaultMaxTokens
	}
	req.Stream = true

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/v1/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	s.setHeaders(httpReq)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		s.Logger.Error("发送请求失败", map[string]interface{}{"error": err.Error()})
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.Logger.Error("Anthropic API 请求失败", map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		})
		return "", fmt.Errorf("API 请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	var fullResponse strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Anthropic SSE 格式: "data: {...}"
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := line[6:]
		if data == "[DONE]" {
			break
		}

		var event StreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		// 只处理 content_block_delta 事件中的 text_delta
		if event.Type == "content_block_delta" && event.Delta != nil && event.Delta.Type == "text_delta" {
			text := event.Delta.Text
			if text != "" {
				fullResponse.WriteString(text)
				onChunk(text)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		s.Logger.Error("读取响应失败", map[string]interface{}{"error": err.Error()})
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	duration := time.Since(startTime)
	s.Logger.Info("Anthropic 流式聊天完成成功", map[string]interface{}{
		"contentLength": fullResponse.Len(),
		"duration":      duration.Milliseconds(),
	})

	return fullResponse.String(), nil
}

// setHeaders 设置 Anthropic API 所需的请求头
func (s *Service) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)
}

// logRequest 记录请求日志（脱敏处理）
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
		"stream":    req.Stream,
		"messages":  messages,
	})
}
