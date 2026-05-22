#!/bin/bash
# ===============================================
# Chat System Pro - Ubuntu 22.04 快速部署脚本
# ===============================================

set -e

echo "==============================================="
echo "  Chat System Pro - Ubuntu 22.04 快速部署"
echo "==============================================="
echo ""

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then
    echo "[INFO] 请输入sudo密码以继续..."
    sudo -v
fi

# 1. 更新系统
echo "[1/6] 更新系统..."
sudo apt update && sudo apt upgrade -y

# 2. 安装依赖
echo "[2/6] 安装依赖..."
sudo apt install -y curl git wget unzip

# 3. 安装Docker
echo "[3/6] 安装Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com | sudo sh
    sudo usermod -aG docker $USER
    echo "[INFO] Docker安装完成，需要重新登录以生效"
else
    echo "[INFO] Docker已安装"
fi

# 4. 安装Docker Compose
echo "[4/6] 安装Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
else
    echo "[INFO] Docker Compose已安装"
fi

# 5. 克隆项目
echo "[5/6] 克隆项目..."
if [ ! -d "jinjuyi" ]; then
    git clone https://github.com/fjy123123/jinjuyi.git
    cd jinjuyi
else
    echo "[INFO] 项目已存在，更新代码..."
    cd jinjuyi
    git pull
fi

# 6. 配置环境变量
echo "[6/6] 配置环境变量..."
if [ ! -f ".env" ]; then
    cp .env.example .env
    echo "[INFO] .env 文件已创建，请根据需要修改配置"
fi

# 启动服务
echo ""
echo "[INFO] 启动服务..."
sudo docker-compose up -d

# 等待服务启动
echo "[INFO] 等待服务启动..."
sleep 10

# 检查服务状态
echo ""
echo "==============================================="
echo "  部署状态检查"
echo "==============================================="
sudo docker-compose ps

echo ""
echo "==============================================="
echo "  部署完成！"
echo "==============================================="
echo ""
echo "  访问地址:"
echo "    前端: http://$(hostname -I | awk '{print $1}')"
echo "    API:  http://$(hostname -I | awk '{print $1}'):8080"
echo ""
echo "  常用命令:"
echo "    查看日志: docker-compose logs -f"
echo "    停止服务: docker-compose down"
echo "    重启服务: docker-compose restart"
echo ""
echo "  配置文件: $(pwd)/.env"
echo ""
