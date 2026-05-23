#!/bin/bash
# ===============================================
# Chat System Pro - Ubuntu 24.04 一键部署脚本 v2.0
# ===============================================
# 功能:
#   - 自动安装 Docker 和 Docker Compose
#   - 自动配置防火墙
#   - 自动创建必要的配置文件
#   - 自动配置环境变量
#   - 自动启动所有服务
#   - 支持 Docker 权限自动修复
#   - 支持服务自启动配置
# ===============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# 全局变量
COMPOSE_CMD=""
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="${PROJECT_DIR}/install.log"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
    echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') $1" >> "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
    echo "[SUCCESS] $(date '+%Y-%m-%d %H:%M:%S') $1" >> "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1" >&2
    echo "[WARNING] $(date '+%Y-%m-%d %H:%M:%S') $1" >> "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1" >&2
    echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') $1" >> "$LOG_FILE"
}

# 打印标题
print_header() {
    echo ""
    echo "================================================"
    echo -e "${CYAN}  Chat System Pro - Ubuntu 24.04 一键部署${NC}"
    echo -e "${CYAN}  Version 2.0 - 优化版${NC}"
    echo "================================================"
    echo ""
    echo "开始时间: $(date)"
    echo "项目目录: $PROJECT_DIR"
    echo "日志文件: $LOG_FILE"
    echo ""
}

# 检查是否为root用户
check_root() {
    log_info "检查用户权限..."
    if [ "$EUID" -ne 0 ]; then
        log_warning "当前用户非 root，需要 sudo 权限"
        if ! sudo -v 2>/dev/null; then
            log_error "无法获取 sudo 权限，请检查密码"
            exit 1
        fi
    fi
}

# 检查系统要求
check_system() {
    log_info "检查系统要求..."
    
    # 检查是否为 Ubuntu
    if [ ! -f /etc/os-release ]; then
        log_error "无法检测操作系统"
        exit 1
    fi
    
    source /etc/os-release
    if [ "$ID" != "ubuntu" ]; then
        log_warning "非 Ubuntu 系统，某些功能可能不正常"
    fi
    
    # 检查 Ubuntu 版本
    if [ "$VERSION_ID" != "24.04" ] && [ "$VERSION_ID" != "22.04" ]; then
        log_warning "建议使用 Ubuntu 24.04 或 22.04，当前版本: $VERSION_ID"
    fi
    
    # 检查内存
    TOTAL_MEM=$(free -m | awk '/^Mem:/{print $2}')
    if [ "$TOTAL_MEM" -lt 2048 ]; then
        log_warning "建议内存 >= 2GB，当前: ${TOTAL_MEM}MB"
    fi
    
    # 检查磁盘空间
    AVAILABLE_DISK=$(df -BG / | awk 'NR==2 {print $4}' | sed 's/G//')
    if [ "$AVAILABLE_DISK" -lt 10 ]; then
        log_error "磁盘空间不足，需要至少 10GB可用空间，当前: ${AVAILABLE_DISK}GB"
        exit 1
    fi
    
    log_success "系统检查通过"
}

# 更新系统
update_system() {
    log_info "更新系统包列表..."
    if ! sudo apt-get update >> "$LOG_FILE" 2>&1; then
        log_error "系统更新失败"
        exit 1
    fi
    
    log_info "升级系统包..."
    if ! sudo apt-get upgrade -y >> "$LOG_FILE" 2>&1; then
        log_warning "系统升级遇到问题，继续安装..."
    fi
    log_success "系统更新完成"
}

# 安装基础依赖
install_dependencies() {
    log_info "安装基础依赖..."
    
    local packages=(
        curl
        git
        wget
        unzip
        ca-certificates
        gnupg
        lsb-release
        ufw
        net-tools
        python3
        python3-pip
        htop
        tree
    )
    
    for package in "${packages[@]}"; do
        if ! dpkg -l | grep -q "^ii  $package"; then
            log_info "安装 $package..."
            sudo apt-get install -y "$package" >> "$LOG_FILE" 2>&1 || log_warning "$package 安装失败"
        fi
    done
    
    log_success "基础依赖安装完成"
}

# 配置防火墙
setup_firewall() {
    log_info "配置防火墙..."
    
    # 检查 ufw 是否已安装
    if ! command -v ufw &> /dev/null; then
        log_info "安装 UFW 防火墙..."
        sudo apt-get install -y ufw >> "$LOG_FILE" 2>&1
    fi
    
    # 配置防火墙规则
    sudo ufw default deny incoming >> "$LOG_FILE" 2>&1
    sudo ufw default allow outgoing >> "$LOG_FILE" 2>&1
    
    # 开放必要端口
    sudo ufw allow 22/tcp comment 'SSH' >> "$LOG_FILE" 2>&1
    sudo ufw allow 80/tcp comment 'HTTP' >> "$LOG_FILE" 2>&1
    sudo ufw allow 443/tcp comment 'HTTPS' >> "$LOG_FILE" 2>&1
    sudo ufw allow 8080/tcp comment 'API' >> "$LOG_FILE" 2>&1
    
    # 检查 UFW 状态
    if sudo ufw status | grep -q "Status: inactive"; then
        log_info "启用防火墙..."
        echo "y" | sudo ufw enable >> "$LOG_FILE" 2>&1 || log_warning "防火墙启用失败"
    fi
    
    log_success "防火墙配置完成"
}

# 安装 Docker
install_docker() {
    log_info "检查 Docker..."
    
    if command -v docker &> /dev/null; then
        local docker_version=$(docker --version | awk '{print $3}' | sed 's/,//')
        log_success "Docker 已安装，版本: $docker_version"
        return 0
    fi
    
    log_info "安装 Docker..."
    
    # 添加 Docker 官方 GPG 密钥
    sudo install -m 0755 -d /etc/apt/keyrings
    local gpg_key="/etc/apt/keyrings/docker.asc"
    if [ ! -f "$gpg_key" ]; then
        sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o "$gpg_key"
        sudo chmod a+r "$gpg_key"
    fi
    
    # 添加 Docker 仓库
    local arch=$(dpkg --print-architecture)
    echo "deb [arch=${arch} signed-by=${gpg_key}] https://download.docker.com/linux/ubuntu ${VERSION_CODENAME} stable" | \
        sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # 安装 Docker
    sudo apt-get update >> "$LOG_FILE" 2>&1
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin >> "$LOG_FILE" 2>&1
    
    # 启动 Docker 服务
    sudo systemctl start docker >> "$LOG_FILE" 2>&1
    sudo systemctl enable docker >> "$LOG_FILE" 2>&1
    
    log_success "Docker 安装完成"
}

# 检查 Docker Compose
check_docker_compose() {
    log_info "检查 Docker Compose..."
    
    # 优先使用新版 docker compose
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
        local compose_version=$(docker compose version --short)
        log_success "Docker Compose 已安装，版本: $compose_version"
        return 0
    fi
    
    # 降级使用旧版 docker-compose
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
        local compose_version=$(docker-compose --version | awk '{print $3}' | sed 's/,//')
        log_success "Docker Compose 已安装，版本: $compose_version"
        return 0
    fi
    
    log_error "Docker Compose 未找到"
    exit 1
}

# 添加用户到 docker 组
setup_docker_permissions() {
    log_info "配置 Docker 权限..."
    
    local current_user=$(whoami)
    
    if groups "$current_user" | grep -q docker; then
        log_success "用户 $current_user 已在 docker 组"
    else
        log_info "添加用户到 docker 组..."
        sudo usermod -aG docker "$current_user"
        log_success "已将 $current_user 添加到 docker 组"
        log_warning "请重新登录或执行: newgrp docker"
    fi
    
    # 测试 Docker 权限
    if ! docker ps &> /dev/null; then
        log_warning "Docker 权限测试失败，尝试修复..."
        sudo chmod 666 /var/run/docker.sock 2>/dev/null || true
    fi
}

# 创建必要的目录和文件
create_directories() {
    log_info "创建必要的目录..."
    
    local dirs=(
        "docker/ssl"
        "logs"
        "logs/nginx"
        "backend/uploads"
    )
    
    for dir in "${dirs[@]}"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "$dir"
            log_info "创建目录: $dir"
        fi
    done
    
    log_success "目录创建完成"
}

# 创建 nginx 配置
create_nginx_config() {
    log_info "检查 nginx 配置..."
    
    if [ -f "docker/nginx.conf" ]; then
        log_success "nginx.conf 已存在"
        return 0
    fi
    
    log_info "创建 nginx.conf..."
    
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
    
    log_success "nginx.conf 创建完成"
}

# 创建 Dockerfile
create_dockerfile() {
    log_info "检查 Dockerfile..."
    
    if [ -f "backend/Dockerfile" ]; then
        log_success "Dockerfile 已存在"
        return 0
    fi
    
    log_info "创建 backend/Dockerfile..."
    
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
    
    log_success "Dockerfile 创建完成"
}

# 配置环境变量
setup_env() {
    log_info "配置环境变量..."
    
    if [ -f ".env" ]; then
        log_warning ".env 文件已存在，跳过创建"
        return 0
    fi
    
    if [ -f ".env.example" ]; then
        cp .env.example .env
        log_success ".env 文件已创建"
    else
        log_info "创建 .env 文件..."
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
    
    # 获取本机 IP
    local local_ip=$(hostname -I | awk '{print $1}')
    log_info "检测到的本机 IP: $local_ip"
    log_warning "请根据需要修改 .env 配置文件"
}

# 清理旧容器
cleanup_old_containers() {
    log_info "清理旧容器..."
    
    # 停止并删除旧容器
    ${COMPOSE_CMD} down 2>/dev/null || true
    
    # 清理未使用的 Docker 资源
    docker system prune -f >> "$LOG_FILE" 2>&1 || true
    
    log_success "清理完成"
}

# 构建 Docker 镜像
build_images() {
    log_info "构建 Docker 镜像..."
    
    # 构建后端
    log_info "构建后端镜像..."
    if ! ${COMPOSE_CMD} build backend >> "$LOG_FILE" 2>&1; then
        log_error "后端镜像构建失败"
        log_info "查看日志: tail -50 $LOG_FILE"
        exit 1
    fi
    log_success "后端镜像构建完成"
    
    # 构建前端
    if [ -d "web" ] && [ -f "web/Dockerfile" ]; then
        log_info "构建前端镜像..."
        if ! ${COMPOSE_CMD} build frontend >> "$LOG_FILE" 2>&1; then
            log_warning "前端镜像构建失败，继续..."
        else
            log_success "前端镜像构建完成"
        fi
    fi
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    # 启动所有服务
    if ! ${COMPOSE_CMD} up -d; then
        log_error "服务启动失败"
        log_info "查看日志: ${COMPOSE_CMD} logs"
        exit 1
    fi
    
    log_success "服务启动中..."
    
    # 等待服务启动
    log_info "等待服务初始化 (20s)..."
    sleep 20
    
    # 检查服务状态
    log_info "检查服务状态..."
    ${COMPOSE_CMD} ps
    
    # 检查健康状态
    sleep 5
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "后端服务运行正常"
    else
        log_warning "后端服务可能还未完全启动，请稍后查看日志"
    fi
    
    if curl -s http://localhost:80 > /dev/null 2>&1; then
        log_success "前端服务运行正常"
    else
        log_warning "前端服务可能还未完全启动"
    fi
}

# 配置自启动
setup_autostart() {
    log_info "配置系统自启动..."
    
    read -p "是否配置系统自启动？(y/n): " choice
    if [[ "$choice" =~ ^[Yy]$ ]]; then
        local service_file="/etc/systemd/system/chat-system-pro.service"
        
        cat > "$service_file" << EOF
[Unit]
Description=Chat System Pro - Docker Compose Service
Requires=docker.service
After=docker.service network-online.target
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=${PROJECT_DIR}
ExecStart=${COMPOSE_CMD} up -d
ExecStop=${COMPOSE_CMD} down
ExecReload=${COMPOSE_CMD} restart
TimeoutStartSec=0
Restart=on-failure
RestartSec=10s
User=$(whoami)
Group=docker

[Install]
WantedBy=multi-user.target
EOF
        
        sudo systemctl daemon-reload
        sudo systemctl enable chat-system-pro.service
        sudo systemctl start chat-system-pro.service
        
        log_success "自启动配置完成"
        log_info "服务已立即启动并设置为开机自启动"
    else
        log_info "跳过自启动配置"
    fi
}

# 打印部署完成信息
print_complete() {
    local local_ip=$(hostname -I | awk '{print $1}')
    
    echo ""
    echo "================================================"
    echo -e "${GREEN}  部署完成！${NC}"
    echo "================================================"
    echo ""
    echo "  完成时间: $(date)"
    echo ""
    echo -e "  ${CYAN}访问地址:${NC}"
    echo -e "    ${GREEN}前端:   http://$local_ip${NC}"
    echo -e "    ${GREEN}API:    http://$local_ip:8080${NC}"
    echo -e "    ${GREEN}健康检查: http://$local_ip:8080/health${NC}"
    echo ""
    echo -e "  ${CYAN}数据库端口:${NC}"
    echo "    MySQL:   localhost:3306"
    echo "    Redis:   localhost:6379"
    echo "    MongoDB: localhost:27017"
    echo ""
    echo -e "  ${CYAN}常用命令:${NC}"
    echo "    查看日志: ${COMPOSE_CMD} logs -f"
    echo "    停止服务: ${COMPOSE_CMD} down"
    echo "    重启服务: ${COMPOSE_CMD} restart"
    echo "    查看状态: ${COMPOSE_CMD} ps"
    echo ""
    echo -e "  ${CYAN}日志文件:${NC} $LOG_FILE"
    echo -e "  ${CYAN}项目目录:${NC} $PROJECT_DIR"
    echo -e "  ${CYAN}配置文件:${NC} $PROJECT_DIR/.env"
    echo ""
    echo -e "  ${YELLOW}重要提示:${NC}"
    echo "    1. 请修改 .env 配置文件中的密码"
    echo "    2. 如果 Docker 命令无权限，请重新登录"
    echo "    3. MySQL root 密码: your_secure_password_here"
    echo ""
}

# 打印帮助信息
print_help() {
    echo ""
    echo "使用方法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  --auto       自动部署（使用默认配置）"
    echo "  --skip-deps  跳过依赖安装"
    echo "  --no-build   跳过镜像构建"
    echo "  --help       显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0              # 交互式部署"
    echo "  $0 --auto       # 自动部署"
    echo "  $0 --skip-deps  # 跳过依赖安装"
    echo ""
}

# 主函数
main() {
    # 解析参数
    local auto_mode=false
    local skip_deps=false
    local skip_build=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --auto)
                auto_mode=true
                shift
                ;;
            --skip-deps)
                skip_deps=true
                shift
                ;;
            --no-build)
                skip_build=true
                shift
                ;;
            --help)
                print_help
                exit 0
                ;;
            *)
                log_error "未知选项: $1"
                print_help
                exit 1
                ;;
        esac
    done
    
    # 初始化
    print_header
    check_root
    
    # 交互模式
    if [ "$auto_mode" = false ]; then
        echo "请选择部署方式:"
        echo "  1. 完整部署 (推荐)"
        echo "  2. 仅启动服务"
        echo "  0. 退出"
        echo ""
        read -p "请选择 [0-2]: " choice
    else
        choice=1
    fi
    
    case $choice in
        1)
            log_info "开始完整部署..."
            
            # 系统检查
            check_system
            
            # 更新系统（可选）
            if [ "$skip_deps" = false ]; then
                update_system
                install_dependencies
                setup_firewall
            else
                log_info "跳过依赖安装"
            fi
            
            # 安装 Docker
            install_docker
            check_docker_compose
            setup_docker_permissions
            
            # 创建配置文件
            create_directories
            create_nginx_config
            create_dockerfile
            setup_env
            
            # 清理并构建
            cleanup_old_containers
            
            if [ "$skip_build" = false ]; then
                build_images
            else
                log_info "跳过镜像构建"
            fi
            
            # 启动服务
            start_services
            
            # 配置自启动（交互模式）
            if [ "$auto_mode" = false ]; then
                setup_autostart
            fi
            
            print_complete
            ;;
        2)
            log_info "启动已有服务..."
            
            if [ ! -f "docker-compose.yml" ]; then
                log_error "未找到 docker-compose.yml，请先执行完整部署"
                exit 1
            fi
            
            check_docker_compose
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
}

# 运行主函数
main "$@"
