// agent 提供基础对话 Agent 功能
package agent

// Message 聊天消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// History 聊天历史
type History struct {
	Messages []Message `json:"messages"`
}
