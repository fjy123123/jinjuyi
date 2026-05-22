#!/bin/bash
# ===============================================
# Chat System Pro - Ubuntu 24.04 一键部署脚本
# ===============================================
# 功能:
#   - 自动安装 Docker 和 Docker Compose
#   - 自动配置防火墙
#   - 自动配置环境变量
#   - 自动启动所有服务
#   - 支持本地部署和 GitHub 部署
# ===============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 打印标题
print_header() {
    echo ""
    echo "================================================"
    echo -e "${CYAN}  Chat System Pro - Ubuntu 24.04 一键部署${NC}"
    echo "================================================"
    echo ""
}

# 检查是否为root用户
check_root() {
    log_info "检查用户权限..."
    if [ "$EUID" -ne 0 ]; then
        log_warning "当前用户非 root，需要 sudo 权限"
        if ! sudo -v; then
            log_error "无法获取 sudo 权限"
            exit 1
        fi
    fi
}

# 更新系统
update_system() {
    log_info "更新系统..."
    sudo apt-get update && sudo apt-get upgrade -y
    log_success "系统更新完成"
}

# 安装基础依赖
install_dependencies() {
    log_info "安装基础依赖..."
    sudo apt-get install -y \
        curl \
        git \
        wget \
        unzip \
        ca-certificates \
        gnupg \
        lsb-release \
        ufw \
        net-tools
    log_success "基础依赖安装完成"
}

# 配置防火墙
setup_firewall() {
    log_info "配置防火墙..."
    
    # 检查 ufw 是否启用
    if sudo ufw status | grep -q "inactive"; then
        sudo ufw default deny incoming
        sudo ufw default allow outgoing
        sudo ufw allow 22/tcp    # SSH
        sudo ufw allow 80/tcp    # HTTP
        sudo ufw allow 443/tcp   # HTTPS
        sudo ufw allow 8080/tcp  # API
        sudo ufw --force enable
        log_success "防火墙配置完成"
    else
        # 确保必要端口开放
        sudo ufw allow 80/tcp
        sudo ufw allow 443/tcp
        sudo ufw allow 8080/tcp
        log_success "防火墙规则已更新"
    fi
}

# 安装 Docker
install_docker() {
    log_info "检查 Docker..."
    
    if command -v docker &> /dev/null; then
        log_success "Docker 已安装，版本: $(docker --version)"
    else
        log_info "正在安装 Docker..."
        
        # 添加 Docker 官方 GPG 密钥
        sudo install -m 0755 -d /etc/apt/keyrings
        sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
        sudo chmod a+r /etc/apt/keyrings/docker.asc
        
        # 添加 Docker 仓库
        echo \
            "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
            $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
            sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
        
        # 安装 Docker
        sudo apt-get update
        sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        
        # 添加当前用户到 docker 组
        sudo usermod -aG docker $USER
        
        log_success "Docker 安装完成"
        log_warning "需要重新登录以生效 docker 组权限"
    fi
}

# 检查 Docker Compose
check_docker_compose() {
    log_info "检查 Docker Compose..."
    
    if docker compose version &> /dev/null; then
        log_success "Docker Compose 已安装（新语法 docker compose）"
        COMPOSE_CMD="docker compose"
    elif docker-compose version &> /dev/null; then
        log_success "Docker Compose 已安装（旧语法 docker-compose）"
        COMPOSE_CMD="docker-compose"
    else
        log_warning "Docker Compose 未找到，正在安装..."
        
        # 尝试安装 docker-compose-plugin
        sudo apt-get install -y docker-compose-plugin
        
        if [ $? -eq 0 ]; then
            log_success "Docker Compose 安装完成"
            COMPOSE_CMD="docker compose"
        else
            # 如果 apt 安装失败，尝试直接下载
            log_warning "apt 安装失败，尝试直接下载..."
            sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
            sudo chmod +x /usr/local/bin/docker-compose
            COMPOSE_CMD="docker-compose"
            log_success "Docker Compose 下载安装完成"
        fi
    fi
    
    # 验证 Docker 服务
    if ! sudo systemctl is-active --quiet docker; then
        log_info "启动 Docker 服务..."
        sudo systemctl start docker
        sudo systemctl enable docker
        log_success "Docker 服务已启动"
    fi
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    
    # 创建 docker 目录（如果不存在）
    if [ ! -d "docker" ]; then
        mkdir -p docker
        log_success "创建 docker 目录"
    fi
    
    # 创建 nginx 配置（如果不存在）
    if [ ! -f "docker/nginx.conf" ]; then
        cat > docker/nginx.conf << 'EOF'
server {
    listen 80;
    server_name localhost;
    
    client_max_body_size 100M;
    
    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;

    location / {
        root /usr/share/nginx/html;
        index index.html index.htm;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    location /ws {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_read_timeout 86400;
    }

    location /uploads {
        alias /app/uploads;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    error_page 500 502 503 504 /50x.html;
    location = /50x.html {
        root /usr/share/nginx/html;
    }
}
EOF
        log_success "创建 nginx 配置文件"
    fi
    
    # 创建上传目录
    mkdir -p backend/uploads logs
    log_success "创建上传和日志目录"
}

# 创建 Dockerfile
create_dockerfile() {
    log_info "检查 Dockerfile..."
    
    if [ ! -f "backend/Dockerfile" ]; then
        cat > backend/Dockerfile << 'EOF'
# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 安装依赖
RUN apk add --no-cache git make

# 复制源码
COPY backend/ ./

# 下载依赖
RUN go mod download
RUN go mod verify

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chat-server .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装CA证书和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 创建目录
RUN mkdir -p /app/uploads /app/logs

# 从构建阶段复制二进制文件
COPY --from=builder /build/chat-server .
COPY --from=builder /build/config.yaml .

# 复制上传目录
COPY --from=builder /build/uploads /app/uploads

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./chat-server"]
EOF
        log_success "创建 backend/Dockerfile"
    fi
}

# 获取项目代码
get_project() {
    log_info "获取项目代码..."
    
    if [ -d "jinjuyi" ]; then
        log_warning "项目目录已存在"
        read -p "是否更新现有代码？(y/n): " choice
        if [ "$choice" = "y" ] ] || [ "$choice" = "Y" ]; then
            cd jinjuyi
            git pull
            cd ..
            log_success "代码更新完成"
        fi
    else
        log_info "从 GitHub 克隆项目..."
        git clone -b main https://github.com/fjy123123/jinjuyi.git
        log_success "项目克隆完成"
    fi
    
    cd jinjuyi
    PROJECT_DIR=$(pwd)
}

# 配置环境变量
setup_env() {
    log_info "配置环境变量..."
    
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            log_success ".env 文件已创建"
        else
            # 如果 .env.example 不存在，创建一个
            cat > .env << 'EOF'
# MySQL 配置
MYSQL_ROOT_PASSWORD=your_secure_password_here
MYSQL_DATABASE=chat_system_pro
MYSQL_USER=chat_user
MYSQL_PASSWORD=chat_password_here

# Redis 配置
REDIS_PASSWORD=redis_password_here

# JWT 配置
JWT_SECRET=your_super_secret_jwt_key_make_it_long_and_secure_here
GIN_MODE=release

# 应用配置
APP_NAME=知信
APP_VERSION=v1.0.0
EOF
            log_success ".env 文件已创建"
        fi
    else
        log_warning ".env 文件已存在"
    fi
    
    # 获取本机 IP
    LOCAL_IP=$(hostname -I | awk '{print $1}')
    
    log_info "检测到的本机 IP: $LOCAL_IP"
    log_warning "请根据需要手动修改 .env 文件"
}

# 构建 Docker 镜像
build_images() {
    log_info "构建 Docker 镜像..."
    
    # 构建后端
    log_info "构建后端镜像..."
    $COMPOSE_CMD build backend
    
    # 构建前端
    log_info "构建前端镜像..."
    $COMPOSE_CMD build frontend
    
    log_success "Docker 镜像构建完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    # 使用检测到的 compose 命令
    $COMPOSE_CMD up -d
    
    log_success "服务启动中..."
    
    # 等待服务启动
    log_info "等待服务初始化 (15s)..."
    sleep 15
    
    # 检查服务状态
    log_info "检查服务状态..."
    $COMPOSE_CMD ps
    
    # 检查服务健康状态
    log_info "检查服务健康状态..."
    sleep 5
    
    # 尝试访问健康检查接口
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "后端服务运行正常"
    else
        log_warning "后端服务可能还未完全启动，请稍后查看日志"
    fi
}

# 打印部署完成信息
print_complete() {
    LOCAL_IP=$(hostname -I | awk '{print $1}')
    
    echo ""
    echo "================================================"
    echo -e "${GREEN}  部署完成！${NC}"
    echo "================================================"
    echo ""
    echo "  访问地址:"
    echo -e "    ${CYAN}前端: http://$LOCAL_IP${NC}"
    echo -e "    ${CYAN}API:  http://$LOCAL_IP:8080${NC}"
    echo -e "    ${CYAN}API文档: http://$LOCAL_IP:8080/api/v1${NC}"
    echo ""
    echo "  常用命令 (在 jinjuyi 目录下执行):"
    echo "    查看日志: $COMPOSE_CMD logs -f"
    echo "    停止服务: $COMPOSE_CMD down"
    echo "    重启服务: $COMPOSE_CMD restart"
    echo "    更新代码: git pull && $COMPOSE_CMD pull && $COMPOSE_CMD up -d --build"
    echo ""
    echo "  项目目录: $PROJECT_DIR"
    echo "  配置文件: $PROJECT_DIR/.env"
    echo ""
    echo "  重要提示:"
    echo "    1. 如果 Docker 命令无权限，请重新登录"
    echo "    2. 请根据需要修改 .env 配置文件"
    echo "    3. MySQL 默认端口: 3306"
    echo "    4. Redis 默认端口: 6379"
    echo "    5. MongoDB 默认端口: 27017"
    echo ""
}

# 打印帮助信息
print_help() {
    echo ""
    echo "使用方法:"
    echo "  $0              交互式部署"
    echo "  $0 --auto       自动部署（使用默认配置）"
    echo "  $0 --update     仅更新代码并重启服务"
    echo "  $0 --stop       停止所有服务"
    echo "  $0 --status     查看服务状态"
    echo "  $0 --logs       查看日志"
    echo "  $0 --help       显示帮助信息"
    echo ""
}

# 主函数
main() {
    print_header
    check_root
    
    case "${1:-}" in
        --auto)
            log_info "开始自动部署..."
            update_system
            install_dependencies
            setup_firewall
            install_docker
            check_docker_compose
            create_directories
            create_dockerfile
            setup_env
            build_images
            start_services
            print_complete
            ;;
        --update)
            log_info "更新代码并重启服务..."
            cd jinjuyi
            git pull
            $COMPOSE_CMD pull
            $COMPOSE_CMD up -d --build
            log_success "更新完成"
            ;;
        --stop)
            log_info "停止所有服务..."
            $COMPOSE_CMD down
            log_success "服务已停止"
            ;;
        --status)
            $COMPOSE_CMD ps
            ;;
        --logs)
            $COMPOSE_CMD logs -f
            ;;
        --help)
            print_help
            ;;
        *)
            echo "请选择部署方式:"
            echo "  1. 完整部署 (推荐)"
            echo "  2. 仅启动服务"
            echo "  0. 退出"
            echo ""
            read -p "请选择 [0-2]: " choice
            
            case $choice in
                1)
                    update_system
                    install_dependencies
                    setup_firewall
                    install_docker
                    check_docker_compose
                    create_directories
                    create_dockerfile
                    setup_env
                    build_images
                    start_services
                    print_complete
                    ;;
                2)
                    if [ ! -d "jinjuyi" ]; then
                        log_error "项目目录不存在，请先执行完整部署"
                        exit 1
                    fi
                    cd jinjuyi
                    PROJECT_DIR=$(pwd)
                    start_services
                    print_complete
                    ;;
                0)
                    log_info "再见！"
                    exit 0
                    ;;
                *)
                    log_error "无效选择"
                    exit 1
                    ;;
            esac
            ;;
    esac
}

# 运行主函数
main "$@"
