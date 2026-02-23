# Go 语言包 (package) 说明

## Go 语言的包规则

1. **同一目录 = 同一包**：同一个目录下的所有 `.go` 文件必须声明相同的包名
2. **package main**：表示这是一个可执行程序的入口包
3. **共享作用域**：同一个包内的所有文件可以互相访问彼此的导出类型和函数

## 项目文件结构

```
backend/
├── main.go           (package main) - 程序入口
├── types.go          (package main) - 类型定义
├── agent.go          (package main) - Agent 实现
└── openai_service.go (package main) - OpenAI 服务
```

## 文件功能对比

| 文件 | 作用 | 是否包含 main 函数 |
|------|------|------------------|
| main.go | 程序入口，启动 HTTP 服务器 | ✅ 是 |
| types.go | 定义共享的数据类型 | ❌ 否 |
| agent.go | 实现 Agent 核心逻辑 | ❌ 否 |
| openai_service.go | 实现 OpenAI API 服务 | ❌ 否 |

## 工作原理

### 包的组织方式

- 所有文件属于同一个 `main` 包
- `main.go` 中的 `func main()` 是程序启动点
- 其他文件中定义的类型和函数可以在 `main.go` 中直接使用
- 编译时会将所有文件打包成一个可执行程序

### 代码示例

在 `main.go` 中可以直接使用其他文件定义的类型和函数：

```go
// 使用 agent.go 中定义的 Agent 结构体和 NewAgent 函数
agent := NewAgent(nil)

// 使用 GenerateResponse 方法
agent.GenerateResponse("你好", func(chunk string) {
    fmt.Println(chunk)
})

// 使用 types.go 中定义的 Message 类型
msg := Message{
    Role:    "user",
    Content: "测试消息",
}
```

## 为什么都使用 package main

这是 Go 语言组织代码的标准方式：

1. **代码组织**：将大型程序拆分成多个文件，便于维护
2. **共享访问**：同一包内的代码可以直接访问，无需导入
3. **编译单元**：所有文件编译成一个可执行程序
4. **标准做法**：Go 社区的常见实践

## 注意事项

- 同一目录下不能有不同包名的文件
- `package main` 必须包含且只能包含一个 `main()` 函数
- 首字母大写的类型和函数是导出的，可以被其他包访问
- 首字母小写的类型和函数是私有的，只能在当前包内访问
