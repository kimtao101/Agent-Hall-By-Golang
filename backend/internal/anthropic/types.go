package anthropic

// Message 表示 Anthropic 聊天消息（只支持 user/assistant 角色）
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest Anthropic 聊天请求参数
type ChatRequest struct {
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream"`
}

// StreamEvent SSE 事件结构
type StreamEvent struct {
	Type  string     `json:"type"`
	Index int        `json:"index,omitempty"`
	Delta *TextDelta `json:"delta,omitempty"`
}

// TextDelta content_block_delta 事件中的文本增量
type TextDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ContentBlock 非流式响应内容块
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Usage token 使用情况
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ChatResponse Anthropic 非流式响应
type ChatResponse struct {
	ID      string         `json:"id"`
	Type    string         `json:"type"`
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
	Model   string         `json:"model"`
	Usage   Usage          `json:"usage"`
}
