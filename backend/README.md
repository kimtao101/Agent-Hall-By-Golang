# Agent Hall Backend - 智能体大厅后端

基于 Go + Gin 框架实现的智能体大厅后端服务，采用标准 Go 项目布局，支持多种 AI 智能体扩展。

## 📁 项目结构

```
backend/
├── cmd/
│   └── server/              # 应用程序入口
│       └── main.go          # 主入口文件
├── internal/                # 内部包（不可被外部导入）
│   ├── agent/               # 基础对话 Agent
│   │   ├── agent.go         # Agent 核心逻辑
│   │   └── types.go         # 类型定义
│   ├── openai/              # OpenAI/DeepSeek 服务
│   │   ├── service.go       # API 封装
│   │   └── types.go         # 类型定义
│   ├── server/              # HTTP 服务器
│   │   └── server.go        # 服务器配置和路由
│   └── xiaohongshu/         # 小红书智能体 ⭐示例Agent
│       ├── handler.go       # HTTP 处理器
│       ├── service.go       # 业务逻辑
│       ├── prompts.go       # 提示词模板
│       └── types.go         # 类型定义
├── pkg/                     # 公共包（可被外部导入）
│   └── utils/               # 工具函数
│       └── string.go        # 字符串工具
├── _old/                    # 旧代码备份
├── .ENV                     # 环境变量配置
├── go.mod                   # Go 模块配置
├── go.sum                   # 依赖锁定
└── agent-backend.exe        # 编译后的可执行文件
```

## 🚀 快速开始

### 1. 环境配置

```bash
# 复制环境变量模板
cp .ENV.example .ENV

# 编辑 .ENV 文件，填入你的 DeepSeek API Key
DEEPSEEK_API_KEY=sk-your-api-key-here
PORT=8016
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
```

## 📡 API 接口

### 基础接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/health` | 健康检查 |
| POST | `/chat` | 基础聊天（流式） |
| GET | `/history` | 获取聊天历史 |
| POST | `/clear` | 清除聊天历史 |

### 小红书智能体接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/xiaohongshu/scenes` | 获取所有文案场景 |
| POST | `/xiaohongshu/copy` | 生成小红书文案 |

**生成文案示例**:
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

## 🧩 如何添加新的 Agent

参考 `internal/xiaohongshu` 的实现，按以下步骤添加新的智能体：

### 1. 创建 Agent 目录

```bash
mkdir internal/youragent
```

### 2. 创建文件结构

```
internal/youragent/
├── types.go      # 请求/响应类型定义
├── prompts.go    # 提示词模板（如需要）
├── service.go    # 业务逻辑
└── handler.go    # HTTP 处理器
```

### 3. 实现 Service

```go
// internal/youragent/service.go
package youragent

type Service struct {
    openai *openai.Service
}

func NewService(openaiSvc *openai.Service) *Service {
    return &Service{openai: openaiSvc}
}

func (s *Service) DoSomething(req Request) (*Response, error) {
    // 实现业务逻辑
}
```

### 4. 实现 Handler

```go
// internal/youragent/handler.go
package youragent

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
    group := router.Group("/youragent")
    {
        group.POST("/action", h.DoSomething)
    }
}

func (h *Handler) DoSomething(c *gin.Context) {
    // 处理 HTTP 请求
}
```

### 5. 注册路由

在 `internal/server/server.go` 中添加：

```go
import "agent-backend/internal/youragent"

func (s *Server) registerRoutes() {
    // ... 现有路由
    
    // 你的 Agent 路由
    yourHandler := youragent.NewHandler(youragent.NewService(openaiSvc))
    yourHandler.RegisterRoutes(s.router.Group("/"))
}
```

## 🔧 开发指南

### 模块说明

- **`cmd/`** - 应用程序入口，每个子目录对应一个可执行文件
- **`internal/`** - 私有应用代码，其他项目无法导入
  - **`agent/`** - 基础聊天对话功能
  - **`openai/`** - AI API 封装，供其他模块使用
  - **`xiaohongshu/`** - 小红书智能体（示例）
  - **`server/`** - HTTP 服务器配置
- **`pkg/`** - 公共库代码，可被其他项目导入

### 添加新的场景（小红书）

在 `internal/xiaohongshu/prompts.go` 中添加新的场景提示词：

```go
func (pb *PromptBuilder) buildNewScenePrompt() string {
    // 实现提示词构建逻辑
}
```

在 `internal/xiaohongshu/prompts.go` 的 `Build` 方法中添加场景分支：

```go
func (pb *PromptBuilder) Build(scene string) string {
    switch scene {
    // ... 现有场景
    case "newscene":
        return pb.buildNewScenePrompt()
    }
}
```

## 📝 环境变量

| 变量名 | 必填 | 默认值 | 描述 |
|--------|------|--------|------|
| `DEEPSEEK_API_KEY` | 是 | - | DeepSeek API 密钥 |
| `DEEPSEEK_BASE_URL` | 否 | `https://api.deepseek.com` | API 基础地址 |
| `PORT` | 否 | `8016` | 服务端口 |

## 🎯 技术栈

- **语言**: Go 1.24+
- **Web 框架**: Gin
- **AI 服务**: DeepSeek API
- **项目布局**: 标准 Go 项目布局

## 📄 许可证

MIT
