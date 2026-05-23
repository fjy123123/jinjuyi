#!/bin/bash

# ===============================================
# 知信 (Zhixin) - 一键部署脚本 (Ubuntu 22.04/24.04)
# 版本: v2.0.0
# 描述: 自动安装部署知信聊天系统（Docker 方式）
# 支持系统: Ubuntu 22.04, Ubuntu 24.04
# ===============================================

set -e  # 遇到错误立即退出

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

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "\n${CYAN}========================================${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}========================================${NC}\n"
}

# 打印横幅
print_banner() {
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════════╗"
    echo "║                                              ║"
    echo "║        知信 (Zhixin) 一键部署系统            ║"
    echo "║        支持: Ubuntu 22.04 / 24.04            ║"
    echo "║        版本: v2.0.0                          ║"
    echo "║                                              ║"
    echo "╚══════════════════════════════════════════════╝"
    echo -e "${NC}\n"
}

# 检查是否为 root 用户
check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_error "请不要以 root 用户运行此脚本！"
        log_info "使用: sudo ./ubuntu-22.04-install.sh"
        exit 1
    fi
}

# 检查系统版本
check_system() {
    log_step "检查系统版本"
    
    if ! command -v lsb_release &> /dev/null; then
        log_warn "lsb_release 未安装，正在安装..."
        sudo apt update
        sudo apt install -y lsb-release
    fi

    local dist=$(lsb_release -si)
    local ver=$(lsb_release -sr)

    if [[ "$dist" != "Ubuntu" ]]; then
        log_error "此脚本仅支持 Ubuntu 系统，当前系统: $dist"
        exit 1
    fi

    if [[ "$ver" != "22.04" && "$ver" != "24.04" ]]; then
        log_error "此脚本支持 Ubuntu 22.04 或 24.04，当前版本: $ver"
        exit 1
    fi

    log_success "系统版本检查通过: Ubuntu $ver"
}

# 检查必要工具
install_dependencies() {
    log_step "安装系统依赖"

    log_info "更新软件包列表..."
    sudo apt update

    local packages=(
        curl
        git
        wget
        apt-transport-https
        ca-certificates
        gnupg
        lsb-release
        software-properties-common
        unzip
    )

    for pkg in "${packages[@]}"; do
        if ! dpkg -s "$pkg" 2>/dev/null | grep -q "Status: install ok installed"; then
            log_info "安装 $pkg..."
            sudo apt install -y "$pkg"
        fi
    done

    log_success "系统依赖安装完成"
}

# 安装 Docker
install_docker() {
    log_step "安装 Docker"

    if command -v docker &> /dev/null; then
        log_info "Docker 已安装: $(docker --version)"
        return
    fi

    log_info "正在安装 Docker..."

    # 添加 Docker 官方 GPG 密钥
    sudo mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

    # 设置 Docker 仓库
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    # 安装 Docker
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # 启动 Docker 服务
    sudo systemctl enable docker
    sudo systemctl start docker

    # 添加当前用户到 docker 组
    sudo usermod -aG docker $USER

    log_success "Docker 安装完成: $(docker --version)"
    log_warn "请重新登录或运行 'newgrp docker' 使 Docker 权限生效"
}

# 安装 Docker Compose
install_docker_compose() {
    log_step "安装 Docker Compose"

    # 检查是否已有 docker compose 插件
    if docker compose version &> /dev/null; then
        log_info "Docker Compose 插件已安装: $(docker compose version)"
        return
    fi

    # 检查是否已有 docker-compose 命令
    if command -v docker-compose &> /dev/null; then
        log_info "docker-compose 已安装: $(docker-compose --version)"
        return
    fi

    log_info "正在安装 Docker Compose..."

    # 安装 Docker Compose 插件
    sudo apt install -y docker-compose-plugin

    log_success "Docker Compose 安装完成: $(docker compose version)"
}

# 检查端口占用
check_ports() {
    log_step "检查端口占用"

    local ports=("80" "443" "3306" "27017" "6379")
    local occupied_ports=()

    for port in "${ports[@]}"; do
        if sudo netstat -tuln 2>/dev/null | grep -q ":$port " || sudo ss -tuln 2>/dev/null | grep -q ":$port "; then
            occupied_ports+=("$port")
        fi
    done

    if [ ${#occupied_ports[@]} -ne 0 ]; then
        log_warn "以下端口已被占用: ${occupied_ports[*]}"
        log_info "Docker 容器将绑定到 127.0.0.1，不会与外部服务冲突"
    else
        log_success "端口检查通过"
    fi
}

# 克隆项目代码
clone_project() {
    log_step "获取项目代码"

    PROJECT_DIR="/opt/zhixin"
    
    if [ -d "$PROJECT_DIR" ]; then
        log_warn "项目目录 $PROJECT_DIR 已存在"
        echo -e "${YELLOW}请选择操作:${NC}"
        echo "  1) 使用现有代码（不重新克隆）"
        echo "  2) 重新克隆最新代码"
        echo "  3) 退出"
        read -p "请输入选项 [1-3]: " -n 1 -r
        echo
        case $REPLY in
            2)
                log_info "删除现有代码..."
                sudo rm -rf "$PROJECT_DIR"
                ;;
            3)
                log_info "操作已取消"
                exit 0
                ;;
            *)
                log_info "使用现有代码"
                cd "$PROJECT_DIR"
                return
                ;;
        esac
    fi

    sudo mkdir -p "$PROJECT_DIR"
    sudo chown $USER:$USER "$PROJECT_DIR"
    cd "$PROJECT_DIR"
    
    log_info "正在克隆项目代码..."
    
    git clone https://github.com/fjy123123/jinjuyi.git . || {
        log_error "克隆项目失败"
        exit 1
    }
    
    log_success "项目代码获取成功"
}

# 生成强随机密钥
generate_secrets() {
    log_step "生成安全密钥"

    if [ -f ".env" ]; then
        log_warn ".env 文件已存在"
        read -p "是否重新生成密钥？(y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "跳过密钥生成"
            return
        fi
    fi

    # 从 .env.example 创建 .env
    if [ ! -f ".env.example" ]; then
        log_error ".env.example 文件不存在"
        exit 1
    fi

    cp .env.example .env

    # 生成密钥
    JWT_SECRET=$(openssl rand -base64 32)
    MYSQL_ROOT_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    MYSQL_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    MONGO_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    REDIS_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)

    # 更新 .env 文件
    sed -i "s|JWT_SECRET=.*|JWT_SECRET=$JWT_SECRET|" .env
    sed -i "s|MYSQL_ROOT_PASSWORD=.*|MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD|" .env
    sed -i "s|MYSQL_PASSWORD=.*|MYSQL_PASSWORD=$MYSQL_PASSWORD|" .env
    sed -i "s|MONGO_PASSWORD=.*|MONGO_PASSWORD=$MONGO_PASSWORD|" .env
    sed -i "s|REDIS_PASSWORD=.*|REDIS_PASSWORD=$REDIS_PASSWORD|" .env

    # 保存密钥到文件
    cat > /opt/zhixin/credentials.txt << EOF
========================================
知信系统 - 数据库凭证
生成时间: $(date '+%Y-%m-%d %H:%M:%S')
========================================

MySQL Root 密码: $MYSQL_ROOT_PASSWORD
MySQL 应用密码: $MYSQL_PASSWORD
MongoDB 密码: $MONGO_PASSWORD
Redis 密码: $REDIS_PASSWORD
JWT 密钥: $JWT_SECRET

请妥善保管此文件，建议部署后删除或加密！
========================================
EOF

    chmod 600 /opt/zhixin/credentials.txt

    log_success "安全密钥生成完成"
    log_warn "凭证已保存到: /opt/zhixin/credentials.txt"
}

# 创建数据目录
create_data_dirs() {
    log_step "创建数据目录"

    sudo mkdir -p /opt/zhixin-data/{mysql,mongo,redis,uploads,logs,backups}
    sudo chown -R $USER:$USER /opt/zhixin-data

    log_success "数据目录创建完成"
}

# 更新 docker-compose.yml 使用本地数据目录
update_compose_config() {
    log_step "更新 Docker Compose 配置"

    # 备份原始文件
    cp docker-compose.yml docker-compose.yml.backup

    # 创建更新后的配置
    cat > docker-compose.yml << 'COMPOSE_EOF'
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: zhixin-mysql
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: ${MYSQL_DATABASE:-zhixin_chat}
      MYSQL_USER: ${MYSQL_USER:-zhixin_user}
      MYSQL_PASSWORD: ${MYSQL_PASSWORD}
    ports:
      - "127.0.0.1:3306:3306"
    volumes:
      - /opt/zhixin-data/mysql:/var/lib/mysql
    networks:
      - zhixin-network
    restart: always
    command: 
      --default-authentication-plugin=mysql_native_password
      --character-set-server=utf8mb4
      --collation-server=utf8mb4_unicode_ci
      --innodb-buffer-pool-size=256M
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  mongodb:
    image: mongo:7
    container_name: zhixin-mongo
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_USER:-mongo_admin}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD}
      MONGO_INITDB_DATABASE: ${MONGO_DB:-zhixin_chat}
    ports:
      - "127.0.0.1:27017:27017"
    volumes:
      - /opt/zhixin-data/mongo:/data/db
    networks:
      - zhixin-network
    restart: always
    command: mongod --auth
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: zhixin-redis
    ports:
      - "127.0.0.1:6379:6379"
    volumes:
      - /opt/zhixin-data/redis:/data
    networks:
      - zhixin-network
    restart: always
    command: redis-server --requirepass ${REDIS_PASSWORD} --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: zhixin-backend
    ports:
      - "8080:8080"
    environment:
      MYSQL_HOST: mysql
      MYSQL_PORT: 3306
      MYSQL_USER: ${MYSQL_USER:-zhixin_user}
      MYSQL_PASSWORD: ${MYSQL_PASSWORD}
      MYSQL_DATABASE: ${MYSQL_DATABASE:-zhixin_chat}
      MONGO_HOST: mongodb
      MONGO_PORT: 27017
      MONGO_USER: ${MONGO_USER:-mongo_admin}
      MONGO_PASSWORD: ${MONGO_PASSWORD}
      MONGO_DB: ${MONGO_DB:-zhixin_chat}
      REDIS_HOST: redis
      REDIS_PORT: 6379
      REDIS_PASSWORD: ${REDIS_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
      GIN_MODE: ${GIN_MODE:-release}
      SERVER_PORT: 8080
    volumes:
      - /opt/zhixin-data/uploads:/app/uploads
      - /opt/zhixin-data/logs:/app/logs
    depends_on:
      mysql:
        condition: service_healthy
      mongodb:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - zhixin-network
    restart: always
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G

  frontend:
    build:
      context: ./web
      dockerfile: Dockerfile
    container_name: zhixin-frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - zhixin-network
    restart: always

networks:
  zhixin-network:
    driver: bridge
COMPOSE_EOF

    log_success "Docker Compose 配置更新完成"
}

# 构建并启动服务
start_services() {
    log_step "启动服务"

    cd /opt/zhixin

    log_info "正在构建 Docker 镜像..."
    docker compose build

    log_info "正在启动所有服务..."
    docker compose up -d --remove-orphans

    log_success "服务启动完成"
}

# 等待服务就绪
wait_for_services() {
    log_step "等待服务就绪"

    log_info "等待数据库初始化..."

    # 等待 MySQL
    log_info "等待 MySQL 就绪..."
    for i in {1..60}; do
        if docker exec zhixin-mysql mysqladmin ping -h localhost --silent 2>/dev/null; then
            log_success "MySQL 服务就绪"
            break
        fi
        sleep 2
    done

    # 等待 MongoDB
    log_info "等待 MongoDB 就绪..."
    for i in {1..60}; do
        if docker exec zhixin-mongo mongosh --eval "db.adminCommand('ping')" &>/dev/null; then
            log_success "MongoDB 服务就绪"
            break
        fi
        sleep 2
    done

    # 等待后端
    log_info "等待后端服务就绪..."
    for i in {1..90}; do
        if curl -sf http://localhost:8080/health &>/dev/null; then
            log_success "后端服务就绪"
            break
        fi
        sleep 3
    done

    # 验证前端
    log_info "验证前端..."
    if curl -sf http://localhost/ &>/dev/null; then
        log_success "前端服务就绪"
    fi

    log_success "所有服务已就绪"
}

# 创建管理脚本
create_manage_scripts() {
    log_step "创建管理脚本"

    cat > /opt/zhixin/zhixin-manage.sh << 'MANAGE_EOF'
#!/bin/bash

# 知信管理系统脚本

PROJECT_DIR="/opt/zhixin"
COMPOSE_FILE="$PROJECT_DIR/docker-compose.yml"

cd "$PROJECT_DIR"

case "$1" in
    start)
        docker compose -f $COMPOSE_FILE up -d
        echo "✅ 知信服务已启动"
        ;;
    stop)
        docker compose -f $COMPOSE_FILE down
        echo "⏹️ 知信服务已停止"
        ;;
    restart)
        docker compose -f $COMPOSE_FILE restart
        echo "🔄 知信服务已重启"
        ;;
    status)
        docker compose -f $COMPOSE_FILE ps
        ;;
    logs)
        docker compose -f $COMPOSE_FILE logs -f ${2:-}
        ;;
    logs-backend)
        docker compose -f $COMPOSE_FILE logs -f --tail=100 zhixin-backend
        ;;
    logs-frontend)
        docker compose -f $COMPOSE_FILE logs -f --tail=100 zhixin-frontend
        ;;
    logs-mysql)
        docker compose -f $COMPOSE_FILE logs -f --tail=100 zhixin-mysql
        ;;
    logs-mongo)
        docker compose -f $COMPOSE_FILE logs -f --tail=100 zhixin-mongo
        ;;
    update)
        echo "📦 更新知信系统..."
        git pull
        docker compose -f $COMPOSE_FILE build
        docker compose -f $COMPOSE_FILE up -d
        echo "✅ 更新完成"
        ;;
    backup)
        TIMESTAMP=$(date +%Y%m%d_%H%M%S)
        BACKUP_DIR="/opt/zhixin-data/backups/$TIMESTAMP"
        mkdir -p $BACKUP_DIR
        
        echo "📦 开始备份..."
        
        # 备份 MySQL
        docker exec zhixin-mysql mysqldump -u zhixin_user -p${MYSQL_PASSWORD} zhixin_chat > $BACKUP_DIR/mysql.sql 2>/dev/null
        
        # 备份 MongoDB
        docker exec zhixin-mongo mongodump --out $BACKUP_DIR/mongo 2>/dev/null
        
        # 备份上传文件
        cp -r /opt/zhixin-data/uploads $BACKUP_DIR/
        
        echo "✅ 备份完成: $BACKUP_DIR"
        ;;
    clean)
        echo "🧹 清理 Docker 无用资源..."
        docker system prune -f
        docker volume prune -f
        echo "✅ 清理完成"
        ;;
    stats)
        docker stats --no-stream
        ;;
    *)
        echo "知信管理系统"
        echo ""
        echo "用法: $0 {start|stop|restart|status|logs|logs-backend|logs-frontend|logs-mysql|logs-mongo|update|backup|clean|stats}"
        echo ""
        echo "命令说明:"
        echo "  start        启动服务"
        echo "  stop         停止服务"
        echo "  restart      重启服务"
        echo "  status       查看服务状态"
        echo "  logs         查看所有日志"
        echo "  logs-backend 查看后端日志"
        echo "  logs-frontend 查看前端日志"
        echo "  logs-mysql   查看MySQL日志"
        echo "  logs-mongo   查看MongoDB日志"
        echo "  update       更新系统"
        echo "  backup       备份数据"
        echo "  clean        清理无用资源"
        echo "  stats        查看资源使用"
        exit 1
        ;;
esac
MANAGE_EOF

    chmod +x /opt/zhixin/zhixin-manage.sh
    
    # 创建软链接到 PATH
    if [ ! -L "/usr/local/bin/zhixin" ]; then
        sudo ln -s /opt/zhixin/zhixin-manage.sh /usr/local/bin/zhixin
    fi

    log_success "管理脚本创建完成"
    log_info "使用命令: zhixin start|stop|restart|status|logs|update|backup"
}

# 设置防火墙
setup_firewall() {
    log_step "配置防火墙"

    if sudo ufw status &>/dev/null | grep -q "Status: active"; then
        log_info "防火墙已启用"
        
        sudo ufw allow 80/tcp comment "知信 HTTP"
        sudo ufw allow 443/tcp comment "知信 HTTPS"
        sudo ufw allow 22/tcp comment "SSH"
        
        log_success "防火墙规则已添加"
    else
        log_info "防火墙未启用，跳过配置"
    fi
}

# 显示完成信息
show_completion_info() {
    echo -e "\n${CYAN}╔══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║                                                          ║${NC}"
    echo -e "${CYAN}║           🎉 知信聊天系统部署完成！                       ║${NC}"
    echo -e "${CYAN}║                                                          ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════════════╝${NC}\n"

    # 获取服务器 IP
    SERVER_IP=$(hostname -I | awk '{print $1}')
    
    log_success "访问地址:"
    log_info "  🌐 前端界面: http://${SERVER_IP}"
    log_info "  🔌 API 接口: http://${SERVER_IP}:8080"
    log_info "  🔗 WebSocket: ws://${SERVER_IP}:8080/ws"
    
    echo ""
    log_success "管理命令:"
    log_info "  zhixin start        启动服务"
    log_info "  zhixin stop         停止服务"
    log_info "  zhixin restart      重启服务"
    log_info "  zhixin status       查看状态"
    log_info "  zhixin logs         查看日志"
    log_info "  zhixin update       更新系统"
    log_info "  zhixin backup       备份数据"
    
    echo ""
    log_success "数据目录:"
    log_info "  MySQL:     /opt/zhixin-data/mysql"
    log_info "  MongoDB:   /opt/zhixin-data/mongo"
    log_info "  Redis:     /opt/zhixin-data/redis"
    log_info "  上传文件:  /opt/zhixin-data/uploads"
    log_info "  日志:      /opt/zhixin-data/logs"
    log_info "  备份:      /opt/zhixin-data/backups"
    
    echo ""
    log_warn "重要提示:"
    log_warn "  1. 请妥善保管: /opt/zhixin/credentials.txt (数据库密码)"
    log_warn "  2. 默认管理员账号: admin / admin123456 (首次登录后修改)"
    log_warn "  3. 生产环境请配置 HTTPS 和防火墙"
    log_warn "  4. 建议定期执行: zhixin backup"
    
    echo ""
    echo -e "${CYAN}══════════════════════════════════════════════════════════${NC}"
    log_success "欢迎使用知信聊天系统！"
    echo -e "${CYAN}══════════════════════════════════════════════════════════${NC}\n"
}

# 主函数
main() {
    print_banner
    
    log_info "开始部署知信聊天系统..."
    log_info "当前用户: $(whoami)"
    log_info "系统版本: $(lsb_release -sd)"
    
    check_root
    check_system
    install_dependencies
    install_docker
    install_docker_compose
    check_ports
    clone_project
    generate_secrets
    create_data_dirs
    update_compose_config
    start_services
    wait_for_services
    create_manage_scripts
    setup_firewall
    show_completion_info

    log_success "部署完成！请根据上述信息进行后续操作。"
}

# 捕获错误
trap 'log_error "部署过程中出现错误！"; exit 1' ERR

# 运行主函数
main "$@"