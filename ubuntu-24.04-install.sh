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
    echo "==============================================="
    echo "  Chat System Pro - Ubuntu 24.04 一键部署"
    echo "==============================================="
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
        ufw
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
}

# 获取项目代码
get_project() {
    log_info "获取项目代码..."
    
    if [ -d "jinjuyi" ]; then
        log_warning "项目目录已存在"
        read -p "是否更新现有代码？(y/n): " choice
        if [ "$choice" = "y" ] || [ "$choice" = "Y" ]; then
            cd jinjuyi
            git pull
            cd ..
            log_success "代码更新完成"
        fi
    else
        log_info "从 GitHub 克隆项目..."
        git clone https://github.com/fjy123123/jinjuyi.git
        log_success "项目克隆完成"
    fi
    
    cd jinjuyi
    PROJECT_DIR=$(pwd)
}

# 配置环境变量
setup_env() {
    log_info "配置环境变量..."
    
    if [ ! -f ".env" ]; then
        cp .env.example .env
        log_success ".env 文件已创建"
    else
        log_warning ".env 文件已存在"
    fi
    
    # 获取本机 IP
    LOCAL_IP=$(hostname -I | awk '{print $1}')
    
    log_info "检测到的本机 IP: $LOCAL_IP"
    log_warning "请根据需要手动修改 .env 文件"
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
}

# 打印部署完成信息
print_complete() {
    LOCAL_IP=$(hostname -I | awk '{print $1}')
    
    echo ""
    echo "==============================================="
    echo "  部署完成！"
    echo "==============================================="
    echo ""
    echo "  访问地址:"
    echo "    前端: http://$LOCAL_IP"
    echo "    API:  http://$LOCAL_IP:8080"
    echo ""
    echo "  常用命令 (在 jinjuyi 目录下执行):"
    echo "    查看日志: $COMPOSE_CMD logs -f"
    echo "    停止服务: $COMPOSE_CMD down"
    echo "    重启服务: $COMPOSE_CMD restart"
    echo "    更新代码: git pull && $COMPOSE_CMD pull && $COMPOSE_CMD up -d"
    echo ""
    echo "  项目目录: $PROJECT_DIR"
    echo "  配置文件: $PROJECT_DIR/.env"
    echo ""
    echo "  重要提示: 如果 Docker 命令无权限，请重新登录"
    echo ""
}

# 主函数
main() {
    print_header
    check_root
    
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
            get_project
            setup_env
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
}

# 运行主函数
main "$@"
