// Agent Hall Backend Server
// 智能体大厅后端服务入口

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"agent-backend/internal/server"
)

// @title Agent Hall API
// @version 1.0
// @description 智能体大厅后端服务 - 支持多种AI智能体
// @host localhost:8016
// @BasePath /
func main() {
	// 加载环境变量
	if err := godotenv.Load(".ENV"); err != nil {
		log.Printf("[WARN] 未找到 .ENV 文件，使用系统环境变量: %v", err)
	}

	// 获取端口配置
	port := os.Getenv("PORT")
	if port == "" {
		port = "8016"
	}

	// 创建并启动服务器
	srv := server.New(port)
	
	log.Printf("[INFO] 服务器启动中，监听端口: %s", port)
	if err := srv.Run(); err != nil {
		log.Fatalf("[FATAL] 服务器启动失败: %v", err)
	}
}
