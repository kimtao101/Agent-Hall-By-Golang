package openai

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

// Service OpenAI 服务封装
type Service struct {
	apiKey  string
	baseURL string
	client  *http.Client
	Logger  Logger
}

// NewService 创建新的 OpenAI 服务实例
// apiKey: OpenAI API 密钥
// baseURL: API 基础地址
// logger: 日志记录器，如果为 nil 则使用默认日志
func NewService(apiKey, baseURL string, logger Logger) *Service {
	if logger == nil {
		logger = NewDefaultLogger()
	}
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}

	if apiKey == "" {
		logger.Warn("API Key 未设置")
	} else {
		logger.Info("OpenAI 服务初始化成功", map[string]interface{}{
			"baseURL":          baseURL,
			"apiKeyConfigured": apiKey != "",
		})
	}

	return &Service{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
		Logger:  logger,
	}
}

// NewDefaultService 从环境变量创建默认服务实例
func NewDefaultService() *Service {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	return NewService(apiKey, baseURL, nil)
}

// CreateChatCompletion 创建聊天完成（非流式）
// req: 聊天完成请求参数
// 返回: 聊天完成响应
func (s *Service) CreateChatCompletion(req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	startTime := time.Now()
	s.logRequest("创建聊天完成", req)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.Logger.Error("API 请求失败", map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		})
		return nil, fmt.Errorf("API 请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	duration := time.Since(startTime)
	s.Logger.Info("聊天完成成功", map[string]interface{}{
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
func (s *Service) CreateChatCompletionStream(req ChatCompletionRequest, onChunk func(string)) (string, error) {
	startTime := time.Now()
	s.logRequest("创建流式聊天完成", req)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		s.Logger.Error("发送请求失败", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.Logger.Error("API 请求失败", map[string]interface{}{
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
		s.Logger.Error("读取响应失败", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	duration := time.Since(startTime)
	s.Logger.Info("流式聊天完成成功", map[string]interface{}{
		"contentLength": fullResponse.Len(),
		"duration":      duration.Milliseconds(),
	})

	return fullResponse.String(), nil
}

// logRequest 记录请求日志（脱敏处理）
func (s *Service) logRequest(action string, req ChatCompletionRequest) {
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
		"model":       req.Model,
		"temperature": req.Temperature,
		"stream":      req.Stream,
		"messages":    messages,
	})
}
