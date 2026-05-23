#!/bin/bash

# Replit启动脚本

echo "🚀 启动知信聊天系统..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装，正在安装..."
    apt-get update && apt-get install -y golang-go
fi

# 进入后端目录
cd backend

# 初始化Go模块
if [ ! -f "go.mod" ]; then
    echo "📦 初始化Go模块..."
    go mod init chat-system
    go mod tidy
fi

# 创建必要目录
mkdir -p uploads logs

# 启动后端服务
echo "🌐 启动后端服务..."
go run main.go &

# 等待服务启动
sleep 5

# 启动前端（如果npx可用）
if command -v npx &> /dev/null; then
    echo "🎨 启动前端开发服务器..."
    cd ../web
    npx vite --port 3000 &
fi

echo "✅ 服务启动完成！"
echo "📍 后端API: http://localhost:8080"
echo "📍 前端: http://localhost:3000"