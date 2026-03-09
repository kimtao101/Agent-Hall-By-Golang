# Agent Hall Backend - 智能体大厅后端

基于 Go + Gin 框架实现的智能体大厅后端服务，采用标准 Go 项目布局，支持 Anthropic Claude 和 DeepSeek 双 AI 引擎，以及多种 AI 智能体扩展。

## 📁 项目结构

```
backend/
├── cmd/
│   └── server/              # 应用程序入口
│       └── main.go          # 主入口文件
├── internal/                # 内部包（不可被外部导入）
│   ├── agent/               # 基础对话 Agent
│   │   ├── agent.go         # Agent 核心逻辑（支持双 AI 引擎切换）
│   │   └── types.go         # 类型定义
│   ├── anthropic/           # Anthropic Claude 服务
│   │   ├── service.go       # 原生 SSE 流式 API 封装
│   │   └── types.go         # 类型定义
│   ├── openai/              # OpenAI/DeepSeek 服务
│   │   ├── service.go       # API 封装
│   │   └── types.go         # 类型定义
│   ├── server/              # HTTP 服务器
│   │   └── server.go        # 服务器配置、路由、速率限制
│   └── xiaohongshu/         # 小红书智能体 ⭐示例Agent
│       ├── handler.go       # HTTP 处理器
│       ├── service.go       # 业务逻辑（支持双 AI 引擎）
│       ├── prompts.go       # 提示词模板
│       └── types.go         # 类型定义
├── pkg/                     # 公共包（可被外部导入）
│   └── utils/               # 工具函数
│       └── string.go        # 字符串工具
├── _old/                    # 旧代码备份
├── .ENV                     # 环境变量配置
├── .ENV.example             # 环境变量模板
└── go.mod                   # Go 模块配置
```

## 🚀 快速开始

### 1. 环境配置

```bash
# 复制环境变量模板
cp .ENV.example .ENV
```

编辑 `.ENV` 文件，选择一种 AI 引擎：

```env
# 使用 Anthropic Claude
AI_TYPE=ANTHROPIC
ANTHROPIC_API_KEY=sk-ant-your-api-key-here
ANTHROPIC_BASE_URL="https://api.aicodemirror.com/api/claudecode"

# 或使用 DeepSeek
# AI_TYPE=DEEPSEEK
# DEEPSEEK_API_KEY=sk-your-api-key-here
```

### 2. 运行服务

```bash
# 开发模式
go run ./cmd/server/

# 编译运行
go build -o agent-backend.exe ./cmd/server/
./agent-backend.exe
```

### 3. 验证服务

```bash
curl http://localhost:8016/health
# {"status":"ok","timestamp":"2026-03-09T10:00:00+08:00"}
```

## 📡 API 接口

### 基础接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/health` | 健康检查（不受速率限制） |
| POST | `/chat` | 基础聊天（流式响应） |
| GET | `/history` | 获取聊天历史 |
| POST | `/clear` | 清除聊天历史 |

**聊天示例**：
```bash
curl -X POST http://localhost:8016/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，介绍一下自己"}'
```

**获取历史**：
```bash
curl http://localhost:8016/history
# [{"role":"system","content":"..."},{"role":"user","content":"..."},...]
```

### 小红书智能体接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/xiaohongshu/scenes` | 获取所有支持的文案场景 |
| POST | `/xiaohongshu/copy` | 生成小红书文案 |

**支持的场景（8 种）**：

| 场景 ID | 名称 | 场景 ID | 名称 |
|---------|------|---------|------|
| `beauty` | 💄 美妆护肤 | `fashion` | 👗 穿搭时尚 |
| `travel` | ✈️ 旅行打卡 | `food` | 🍔 美食探店 |
| `home` | 🏠 家居好物 | `fitness` | 🏋️ 健身运动 |
| `parenting` | 👶 育儿亲子 | `tech` | 📱 数码科技 |

**生成文案示例**：
```bash
curl -X POST http://localhost:8016/xiaohongshu/copy \
  -H "Content-Type: application/json" \
  -d '{
    "scene": "beauty",
    "config": {
      "productName": "小棕瓶精华",
      "brand": "雅诗兰黛",
      "price": "999",
      "usageFeel": "吸收很快，不油腻",
      "effect": "肤色提亮明显",
      "recommendation": "适合熬夜肌"
    }
  }'
```

## 📝 环境变量

| 变量名 | 必填 | 默认值 | 描述 |
|--------|------|--------|------|
| `AI_TYPE` | 是 | `DEEPSEEK` | AI 引擎类型：`ANTHROPIC` 或 `DEEPSEEK` |
| `ANTHROPIC_API_KEY` | AI_TYPE=ANTHROPIC 时必填 | - | Anthropic API 密钥 |
| `ANTHROPIC_BASE_URL` | 否 | `https://api.aicodemirror.com/api/claudecode` | Anthropic API 地址 |
| `DEEPSEEK_API_KEY` | AI_TYPE=DEEPSEEK 时必填 | - | DeepSeek API 密钥 |
| `DEEPSEEK_BASE_URL` | 否 | `https://api.deepseek.com` | DeepSeek API 地址 |
| `PORT` | 否 | `8016` | 服务监听端口 |

## 🔒 速率限制

服务内置 IP 级别速率限制：
- 每个 IP 每 **15 分钟**最多 **100 次**请求
- `/health` 端点不受限制
- 超限返回 `429 Too Many Requests`

## 🧩 如何添加新的 Agent

参考 `internal/xiaohongshu` 的实现，按以下步骤添加新的智能体：

### 1. 创建 Agent 目录

```
internal/youragent/
├── types.go      # 请求/响应类型定义
├── prompts.go    # 提示词模板（如需要）
├── service.go    # 业务逻辑
└── handler.go    # HTTP 处理器
```

### 2. 实现 Service

```go
// internal/youragent/service.go
package youragent

import (
    anthropicpkg "agent-backend/internal/anthropic"
    "agent-backend/internal/openai"
)

type Service struct {
    openaiSvc    *openai.Service
    anthropicSvc *anthropicpkg.Service
    aiType       string
}

func NewService(openaiSvc *openai.Service, anthropicSvc *anthropicpkg.Service) *Service {
    return &Service{
        openaiSvc:    openaiSvc,
        anthropicSvc: anthropicSvc,
        aiType:       os.Getenv("AI_TYPE"),
    }
}
```

### 3. 实现 Handler

```go
// internal/youragent/handler.go
package youragent

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
    group := router.Group("/youragent")
    group.POST("/action", h.DoSomething)
}
```

### 4. 注册到 Server

在 `internal/server/server.go` 中：

```go
import "agent-backend/internal/youragent"

// New() 函数中：
yourSvc := youragent.NewService(openaiSvc, anthropicSvc)

// registerRoutes() 中：
yourHandler := youragent.NewHandler(yourSvc)
yourHandler.RegisterRoutes(s.router.Group("/"))
```

## 🎯 技术栈

- **语言**: Go 1.24+
- **Web 框架**: Gin
- **AI 服务**: Anthropic Claude（`claude-sonnet-4-6`）/ DeepSeek（`deepseek-chat`）
- **流式响应**: 原生 HTTP chunked 传输 + Anthropic SSE（`content_block_delta`）
- **项目布局**: 标准 Go 项目布局

## 📄 许可证

MIT
