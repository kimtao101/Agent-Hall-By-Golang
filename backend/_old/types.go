package main

// Message 表示聊天消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
