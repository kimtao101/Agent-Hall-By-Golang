# Agent Hall Backend - 智能体大厅后端

基于 Go + Gin 框架实现的智能体大厅后端服务，采用标准 Go 项目布局，支持多种 AI 智能体扩展，兼容 Anthropic Claude 和 DeepSeek 两种 AI 服务提供商。

## 📋 目录

- [项目概述](#项目概述)
- [项目结构](#项目结构)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [环境配置](#环境配置)
- [API 接口文档](#api-接口文档)
- [架构设计](#架构设计)
- [后端模块分析](#后端模块分析)
- [开发指南](#开发指南)
- [部署指南](#部署指南)

## 项目概述

Agent Hall Backend 是一个功能完整的智能体大厅后端服务，提供：

- **多 AI 提供商支持**: 兼容 Anthropic Claude 和 DeepSeek 两种 AI 服务
- **智能体扩展架构**: 模块化设计，支持快速添加新的 AI 智能体
- **流式响应**: 支持实时流式聊天响应
- **RESTful API**: 标准化的 REST API 接口
- **速率限制**: 内置 IP 级别速率限制保护
- **CORS 支持**: 跨域资源共享配置

## 项目结构

```
backend/
├── cmd/
│   └── server/              # 应用程序入口
│       └── main.go          # 主入口文件
├── constant/                # 全局常量定义
│   └── constant.go          # AI_TYPE、模型名、BaseURL 等常量
├── internal/                # 内部包（不可被外部导入）
│   ├── agent/               # 聊天 Agent 核心逻辑
│   │   ├── agent.go         # 会话历史、模型选择、流式响应
│   │   └── types.go         # Message / History 类型定义
│   ├── anthropic/           # Anthropic Claude 适配层（官方 SDK）
│   │   ├── service.go       # 基于 anthropic-sdk-go 的流式/非流式调用
│   │   └── types.go         # Anthropic 请求/响应结构
│   ├── chat/                # 聊天 HTTP 接口
│   │   └── handler.go       # /chat /history /clear /health
│   ├── openai/              # DeepSeek 适配层（OpenAI 兼容协议）
│   │   ├── service.go       # DeepSeek API 封装（/v1/chat/completions）
│   │   └── types.go         # OpenAI 兼容请求/响应结构
│   ├── server/              # HTTP 服务器装配层
│   │   ├── middleware.go    # CORS、限流中间件
│   │   └── server.go        # 服务初始化和路由注册
│   └── xiaohongshu/         # 小红书文案生成模块
│       ├── handler.go       # HTTP 处理器
│       ├── prompts.go       # Prompt 构建逻辑
│       ├── service.go       # 文案生成业务逻辑
│       └── types.go         # 请求/响应与场景定义
├── pkg/                     # 公共包
│   ├── logger/              # 日志抽象
│   │   └── logger.go        # Info / Warn / Error
│   └── utils/               # 工具函数
│       └── string.go        # 字符串截断和 map 取值工具
├── _old/                    # 历史实现备份（当前入口不再使用）
├── .ENV.example             # 环境变量模板
├── .ENV                     # 环境变量配置（需创建）
├── go.mod                   # Go 模块配置
├── go.sum                   # 依赖锁定
├── README.md                # 项目文档
└── agent-backend.exe        # 编译后的可执行文件
```

## 技术栈

- **语言**: Go 1.24+
- **Web 框架**: Gin v1.11.0
- **AI SDK**:
  - Anthropic Claude API via `anthropic-sdk-go` (claude-opus-4-7-20250415)
  - DeepSeek API via OpenAI 兼容协议 (deepseek-v4-flash)
- **项目布局**: 标准 Go 项目布局
- **配置管理**: godotenv
- **HTTP 客户端**: 标准库 net/http

## 快速开始

### 前置要求

- Go 1.24 或更高版本
- 有效的 Anthropic API Key 或 DeepSeek API Key
- Windows/Linux/macOS 操作系统

### 1. 克隆项目

```bash
git clone <repository-url>
cd backend
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境变量

```bash
# 复制环境变量模板
cp .ENV.example .ENV

# 编辑 .ENV 文件，配置 AI 服务
```

### 4. 运行服务

```bash
# 开发模式
go run ./cmd/server/

# 编译后运行
go build -o agent-backend.exe ./cmd/server/
./agent-backend.exe
```

### 5. 验证服务

```bash
curl http://localhost:8016/health
```

预期响应：

```json
{
  "status": "ok",
  "timestamp": "2026-03-09T19:25:23+08:00"
}
```

## 环境配置

### 环境变量说明

| 变量名               | 必填     | 默认值                                        | 描述                                              |
| -------------------- | -------- | --------------------------------------------- | ------------------------------------------------- |
| `AI_TYPE`            | 否       | `DEEPSEEK`                                    | AI 服务类型：`ANTHROPIC` 或 `DEEPSEEK`            |
| `ANTHROPIC_API_KEY`  | 条件必填 | -                                             | Anthropic API 密钥（当 AI_TYPE=ANTHROPIC 时必填） |
| `ANTHROPIC_BASE_URL` | 否       | `https://api.aicodemirror.com/api/claudecode` | Anthropic API 基础地址                            |
| `DEEPSEEK_API_KEY`   | 条件必填 | -                                             | DeepSeek API 密钥（当 AI_TYPE=DEEPSEEK 时必填）   |
| `DEEPSEEK_BASE_URL`  | 否       | `https://api.deepseek.com`                    | DeepSeek API 基础地址                             |
| `PORT`               | 否       | `8016`                                        | 服务监听端口                                      |

### 配置示例

#### 使用 Anthropic Claude

```env
AI_TYPE=ANTHROPIC
ANTHROPIC_API_KEY=sk-ant-api03-xxxxxxxxxxxxxxxxxxxxxxxx
ANTHROPIC_BASE_URL=https://api.aicodemirror.com/api/claudecode
PORT=8016
```

#### 使用 DeepSeek

```env
AI_TYPE=DEEPSEEK
DEEPSEEK_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxx
DEEPSEEK_BASE_URL=https://api.deepseek.com
PORT=8016
```

> 注意：`AI_TYPE` 决定运行时使用哪个服务，不存在自动回退机制。两个服务都会初始化，但只有 `AI_TYPE` 指定的那个会被实际调用。

## API 接口文档

### 基础接口

#### 1. 健康检查

**接口**: `GET /health`

**描述**: 检查服务运行状态

**响应示例**:

```json
{
  "status": "ok",
  "timestamp": "2026-03-09T19:25:23+08:00"
}
```

#### 2. 聊天接口（流式）

**接口**: `POST /chat`

**描述**: 发送消息并获取流式响应

**请求体**:

```json
{
  "message": "你好，请介绍一下你自己"
}
```

**响应格式**: `text/plain` 流式响应

**请求示例**:

```bash
curl -X POST http://localhost:8016/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请介绍一下你自己"}'
```

#### 3. 获取聊天历史

**接口**: `GET /history`

**描述**: 获取当前会话的聊天历史记录

**响应示例**:

```json
[
  {
    "role": "system",
    "content": "You are a helpful research assistant..."
  },
  {
    "role": "user",
    "content": "你好，请介绍一下你自己"
  },
  {
    "role": "assistant",
    "content": "你好！我是一个智能助手..."
  }
]
```

#### 4. 清除聊天历史

**接口**: `POST /clear`

**描述**: 清除当前会话的聊天历史记录

**响应示例**:

```json
{
  "success": true
}
```

### 小红书智能体接口

#### 5. 获取所有文案场景

**接口**: `GET /xiaohongshu/scenes`

**描述**: 获取小红书文案生成的所有可用场景

**响应示例**:

```json
[
  {
    "id": "beauty",
    "name": "美妆护肤",
    "description": "美妆、护肤、彩妆等产品文案"
  },
  {
    "id": "food",
    "name": "美食探店",
    "description": "餐厅、美食、饮品推荐文案"
  },
  {
    "id": "travel",
    "name": "旅行攻略",
    "description": "旅游景点、旅行体验分享文案"
  }
]
```

#### 6. 生成小红书文案

**接口**: `POST /xiaohongshu/copy`

**描述**: 根据配置生成小红书风格文案

**请求体**:

```json
{
  "scene": "beauty",
  "config": {
    "productName": "小棕瓶精华",
    "brand": "雅诗兰黛",
    "price": "999",
    "usageFeel": "吸收很快，不油腻",
    "effect": "肤色提亮明显",
    "recommendation": "适合熬夜肌"
  }
}
```

**响应示例**:

```json
{
  "copy": "✨熬夜党的救星来啦！🌟\n\n今天要给大家安利这款雅诗兰黛小棕瓶精华～💕\n\n💰价格：999元\n\n🌟使用感受：\n吸收超级快！完全不油腻，用完皮肤水润润的～\n\n✨效果：\n肤色提亮真的太明显了！坚持用下来，素颜都自信满满！\n\n👉特别推荐给经常熬夜的姐妹们，真的绝绝子！\n\n#雅诗兰黛 #小棕瓶 #护肤 #熬夜肌 #美妆分享"
}
```

**请求示例**:

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

## 架构设计

### 系统架构图

```
┌─────────────────────────────────────────────────────────┐
│                     HTTP Client                          │
│                  (Frontend / API Client)                 │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   Gin HTTP Server                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │  CORS Middleware │  │ Rate Limit  │  │ Recovery     │   │
│  └──────────────┘  └──────────────┘  └──────────────┘   │
└────────────────────┬────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ Chat Handler │ │ XHS Handler  │ │ Health Check │
└──────┬───────┘ └──────┬───────┘ └──────────────┘
       │                │
       ▼                ▼
┌──────────────┐ ┌──────────────┐
│   Agent      │ │ XHS Service  │
└──────┬───────┘ └──────┬───────┘
       │                │
       └────────┬───────┘
                ▼
       ┌────────────────┐
       │ AI Provider    │
       │ ┌────────────┐ │
       │ │ Anthropic  │ │
       │ └────────────┘ │
       │ ┌────────────┐ │
       │ │ DeepSeek   │ │
       │ └────────────┘ │
       └────────────────┘
```

### 核心组件

#### 1. Server 模块 (`internal/server/`)

- **职责**: HTTP 服务器装配、路由管理、中间件注册
- **功能**:
  - 初始化 Gin 引擎
  - 注册 CORS、Recovery、IP 级别速率限制（100 请求/15 分钟）
  - 读取 `AI_TYPE` 并装配 `agent` 与 `xiaohongshu` 业务模块

#### 2. Chat 模块 (`internal/chat/`)

- **职责**: 对外暴露聊天相关 HTTP 接口
- **功能**:
  - `GET /health`
  - `POST /chat`
  - `GET /history`
  - `POST /clear`

#### 3. Agent 模块 (`internal/agent/`)

- **职责**: 基础聊天对话能力编排
- **功能**:
  - 消息历史管理（最多保留 20 条）
  - 默认系统提示词注入
  - 根据 `AI_TYPE` 选择 Anthropic 或 DeepSeek
  - 流式响应处理

#### 4. AI 服务适配层

- **Anthropic** (`internal/anthropic/`): 基于官方 `anthropic-sdk-go` SDK，调用 Claude Messages API
- **OpenAI/DeepSeek** (`internal/openai/`): 基于 HTTP 客户端，专用于 DeepSeek（OpenAI 兼容协议 `/v1/chat/completions`）
- **功能**:
  - 流式与非流式调用
  - SSE 事件解析
  - 请求日志与错误处理

#### 5. 小红书文案模块 (`internal/xiaohongshu/`)

- **职责**: 按场景生成小红书文案
- **功能**:
  - 多场景 Prompt 构建
  - 文案生成服务
  - 场景列表查询
  - HTTP 接口暴露

### 数据流

#### 聊天请求流程

```
Client Request
    ↓
Gin Router
    ↓
Chat Handler
    ↓
Agent (Add Message)
    ↓
AI Provider Selection
    ↓
Anthropic/DeepSeek API
    ↓
Stream Response
    ↓
Client
```

#### 文案生成流程

```
Client Request (Scene + Config)
    ↓
XHS Handler
    ↓
XHS Service (Build Prompt)
    ↓
Prompt Builder (Scene-based)
    ↓
AI Provider
    ↓
Formatted Copy
    ↓
Client
```

## 后端模块分析

### 启动链路

当前后端启动入口为 [backend/cmd/server/main.go](backend/cmd/server/main.go)。

启动流程如下：

1. 加载 `.ENV` 环境变量
2. 读取 `PORT`，默认值为 `8016`
3. 调用 `server.New(port)` 创建服务实例
4. 在 [backend/internal/server/server.go](backend/internal/server/server.go) 中完成组件装配：
   - 初始化 `openai.NewDefaultService()`
   - 初始化 `anthropic.NewDefaultService()`
   - 初始化 `agent.New(...)`
   - 初始化 `xiaohongshu.NewService(...)`
   - 注册聊天与小红书路由
5. 调用 `Run()` 启动 Gin HTTP 服务

### 分层说明

当前 backend 可分为四层：

1. **入口层**
   - [backend/cmd/server/main.go](backend/cmd/server/main.go)
   - 负责启动应用与加载配置

2. **HTTP 层**
   - [backend/internal/server/server.go](backend/internal/server/server.go)
   - [backend/internal/server/middleware.go](backend/internal/server/middleware.go)
   - [backend/internal/chat/handler.go](backend/internal/chat/handler.go)
   - [backend/internal/xiaohongshu/handler.go](backend/internal/xiaohongshu/handler.go)
   - 负责路由、中间件、请求解析、响应输出

3. **业务层**
   - [backend/internal/agent/agent.go](backend/internal/agent/agent.go)
   - [backend/internal/xiaohongshu/service.go](backend/internal/xiaohongshu/service.go)
   - [backend/internal/xiaohongshu/prompts.go](backend/internal/xiaohongshu/prompts.go)
   - 负责会话管理、Prompt 构建、业务编排

4. **模型适配层**
   - [backend/internal/openai/service.go](backend/internal/openai/service.go)
   - [backend/internal/anthropic/service.go](backend/internal/anthropic/service.go)
   - 负责对接外部模型 API，屏蔽不同供应商协议差异

### 依赖关系

核心依赖关系如下：

```text
cmd/server
  -> internal/server

internal/server
  -> internal/chat
  -> internal/agent
  -> internal/xiaohongshu
  -> internal/openai
  -> internal/anthropic

internal/chat
  -> internal/agent

internal/agent
  -> internal/openai
  -> internal/anthropic
  -> pkg/logger

internal/xiaohongshu/service
  -> internal/openai
  -> internal/anthropic
  -> pkg/logger
  -> pkg/utils

internal/openai
  -> pkg/logger
  -> pkg/utils

internal/anthropic
  -> pkg/logger
  -> pkg/utils
```

整体保持单向依赖：`HTTP 层 -> 业务层 -> 模型适配层 -> 公共工具层`，没有明显循环依赖。

### 外部依赖

[backend/go.mod](backend/go.mod) 中当前直接依赖较少，主要包括：

- `github.com/gin-gonic/gin`：HTTP 服务框架
- `github.com/joho/godotenv`：加载 `.ENV` 配置
- `github.com/anthropics/anthropic-sdk-go`：Anthropic 官方 Go SDK

说明当前后端实现较轻量，尚未引入数据库、缓存、消息队列或鉴权框架。

### 当前核心能力

当前 backend 主要承载两类能力：

1. **多轮聊天能力**
   - 保存内存中的会话历史
   - 支持流式输出
   - 支持清空历史
   - 支持 Anthropic / DeepSeek 两种提供商切换

2. **小红书文案生成能力**
   - 根据不同场景生成 Prompt
   - 调用大模型生成结构化文案
   - 提供场景列表查询接口

### `_old` 目录说明

[backend/\_old/](backend/_old/) 中保留了旧版本实现，主要是早期单包结构代码与历史 TypeScript/Go 原型。

当前正式启动链路不再依赖该目录，实际运行以 `cmd/server` 和 `internal/*` 下的新结构为准。

## 开发指南

### 添加新的 AI 智能体

参考 `internal/xiaohongshu` 的实现，按以下步骤添加新的智能体：

#### 1. 创建 Agent 目录

```bash
mkdir internal/youragent
```

#### 2. 创建文件结构

```
internal/youragent/
├── types.go      # 请求/响应类型定义
├── prompts.go    # 提示词模板（如需要）
├── service.go    # 业务逻辑
└── handler.go    # HTTP 处理器
```

#### 3. 实现 Service

```go
// internal/youragent/service.go
package youragent

import (
    anthropicpkg "agent-backend/internal/anthropic"
    "agent-backend/internal/openai"
    "agent-backend/pkg/logger"
)

type Service struct {
    openaiSvc    *openai.Service
    anthropicSvc *anthropicpkg.Service
    aiType       string
    logger       logger.Logger
}

func NewService(openaiSvc *openai.Service, anthropicSvc *anthropicpkg.Service, aiType string) *Service {
    return &Service{
        openaiSvc:    openaiSvc,
        anthropicSvc: anthropicSvc,
        aiType:       aiType,
        logger:       logger.NewDefaultLogger("[YourAgent]"),
    }
}

func (s *Service) DoSomething(req Request) (*Response, error) {
    // 实现业务逻辑
    return nil, nil
}
```

#### 4. 实现 Handler

```go
// internal/youragent/handler.go
package youragent

import "github.com/gin-gonic/gin"

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

#### 5. 注册路由

在 `internal/server/server.go` 中添加：

```go
import "agent-backend/internal/youragent"

func New(port string) *Server {
    aiType := strings.ToUpper(strings.TrimSpace(os.Getenv("AI_TYPE")))

    openaiSvc := openai.NewDefaultService()
    anthropicSvc := anthropic.NewDefaultService()
    yourSvc := youragent.NewService(openaiSvc, anthropicSvc, aiType)

    // ... 省略其他初始化

    yourHandler := youragent.NewHandler(yourSvc)
    yourHandler.RegisterRoutes(s.router.Group("/"))
    return s
}
```

### 添加新的场景（小红书）

在 `internal/xiaohongshu/prompts.go` 中添加新的场景提示词：

```go
func (pb *PromptBuilder) buildNewScenePrompt() string {
    c := pb.config

    prompt := fmt.Sprintf(`
请根据以下信息生成小红书风格的文案：

场景：新场景
产品名称：%s
品牌：%s
价格：%s
使用感受：%s
效果：%s
推荐人群：%s

要求：
1. 使用适当的表情符号
2. 包含相关话题标签
3. 语言活泼有感染力
4. 字数控制在200-300字之间
    `,
        utils.GetStringValue(c, "productName", ""),
        utils.GetStringValue(c, "brand", ""),
        utils.GetStringValue(c, "price", "未提供"),
        utils.GetStringValue(c, "usageFeel", ""),
        utils.GetStringValue(c, "effect", ""),
        utils.GetStringValue(c, "recommendation", ""),
    )

    return prompt
}
```

在 `internal/xiaohongshu/prompts.go` 的 `Build` 方法中添加场景分支：

```go
func (pb *PromptBuilder) Build(scene string) string {
    switch scene {
    case "beauty":
        return pb.buildBeautyPrompt()
    case "food":
        return pb.buildFoodPrompt()
    case "travel":
        return pb.buildTravelPrompt()
    case "newscene":
        return pb.buildNewScenePrompt()
    default:
        return pb.buildDefaultPrompt()
    }
}
```

### 测试

#### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/agent/

# 运行测试并显示覆盖率
go test -cover ./...
```

#### API 测试

使用提供的测试脚本：

```bash
# PowerShell
.\test_api.ps1

# 或使用 curl
curl http://localhost:8016/health
```

## 部署指南

### 开发环境部署

1. **配置环境变量**

```bash
cp .ENV.example .ENV
# 编辑 .ENV 文件
```

2. **启动服务**

```bash
go run ./cmd/server/
```

### 生产环境部署

#### 1. 编译可执行文件

```bash
# Windows
go build -o agent-backend.exe ./cmd/server/

# Linux
GOOS=linux GOARCH=amd64 go build -o agent-backend ./cmd/server/

# macOS
GOOS=darwin GOARCH=amd64 go build -o agent-backend ./cmd/server/
```

#### 2. 配置生产环境变量

```env
AI_TYPE=ANTHROPIC
ANTHROPIC_API_KEY=your-production-api-key
PORT=8016
```

#### 3. 使用进程管理器（推荐）

**使用 systemd (Linux)**

创建服务文件 `/etc/systemd/system/agent-hall.service`:

```ini
[Unit]
Description=Agent Hall Backend Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/agent-hall/backend
ExecStart=/opt/agent-hall/backend/agent-backend
Restart=always
RestartSec=10
Environment=AI_TYPE=ANTHROPIC
Environment=ANTHROPIC_API_KEY=your-api-key
Environment=PORT=8016

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl start agent-hall
sudo systemctl enable agent-hall
```

**使用 PM2 (跨平台)**

```bash
npm install -g pm2
pm2 start agent-backend --name agent-hall-backend
pm2 save
pm2 startup
```

#### 4. 使用 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8016;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

#### 5. Docker 部署

创建 `Dockerfile`:

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o agent-backend ./cmd/server/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/agent-backend .
COPY .ENV .

EXPOSE 8016
CMD ["./agent-backend"]
```

构建和运行：

```bash
docker build -t agent-hall-backend .
docker run -p 8016:8016 --env-file .ENV agent-hall-backend
```

### 监控和日志

#### 日志管理

服务使用标准输出进行日志记录，建议：

```bash
# 保存日志到文件
./agent-backend >> /var/log/agent-hall/backend.log 2>&1

# 使用 logrotate 管理日志文件
```

#### 健康检查

定期调用健康检查接口：

```bash
# 每 30 秒检查一次
watch -n 30 'curl -f http://localhost:8016/health || echo "Service down!"'
```

## 故障排除

### 常见问题

#### 1. 端口被占用

**错误**: `bind: Only one usage of each socket address`

**解决**:

```bash
# Windows
netstat -ano | findstr :8016
taskkill /PID <PID> /F

# Linux/macOS
lsof -i :8016
kill -9 <PID>
```

#### 2. API 密钥无效

**错误**: `API 请求失败，状态码 401`

**解决**: 检查 `.ENV` 文件中的 API 密钥是否正确

#### 3. 依赖安装失败

**错误**: `go: module xxx: not found`

**解决**:

```bash
go mod tidy
go mod download
```

#### 4. CORS 错误

**错误**: 前端无法访问后端 API

**解决**: 检查 `server.go` 中的 CORS 中间件配置

## 性能优化

### 1. 连接池配置

```go
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
```

### 2. 速率限制调整

在 `server.go` 中调整速率限制参数：

```go
r.Use(rateLimitMiddleware(200, 10*time.Minute)) // 200 请求/10 分钟
```

### 3. 缓存策略

对于频繁请求的场景，考虑添加缓存层：

```go
type CacheService struct {
    cache map[string]CacheEntry
    mu    sync.RWMutex
}

func (c *CacheService) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    entry, ok := c.cache[key]
    if !ok {
        return nil, false
    }
    if time.Since(entry.Timestamp) > 5*time.Minute {
        return nil, false
    }
    return entry.Value, true
}
```

## 安全建议

1. **API 密钥管理**: 永远不要将 API 密钥提交到版本控制系统
2. **环境变量**: 使用环境变量存储敏感信息
3. **HTTPS**: 生产环境务必使用 HTTPS
4. **速率限制**: 保持适当的速率限制防止滥用
5. **输入验证**: 严格验证所有用户输入
6. **日志脱敏**: 确保日志中不包含敏感信息

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

MIT License - 详见 LICENSE 文件

## 联系方式

- 项目主页: [GitHub Repository]
- 问题反馈: [GitHub Issues]
- 文档: [项目 Wiki]

## 更新日志

### v1.0.0 (2026-03-09)

- 初始版本发布
- 支持 Anthropic Claude 和 DeepSeek AI
- 基础聊天功能
- 小红书文案生成智能体
- 流式响应支持
- 速率限制和 CORS 支持
