#!/bin/bash

# ===============================================
# 知信 (Zhixin) - 一键部署脚本 (Ubuntu 24.04)
# 版本: v1.0.0
# 作者: Zhixin Team
# 描述: 自动安装部署知信聊天系统
# ===============================================

set -e  # 遇到错误立即退出

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

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否为 root 用户
check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_error "请不要以 root 用户运行此脚本！"
        exit 1
    fi
}

# 检查系统版本
check_system() {
    if ! command -v lsb_release &> /dev/null; then
        log_error "lsb_release 未安装，请确认系统为 Ubuntu 24.04"
        exit 1
    fi

    local dist=$(lsb_release -si)
    local ver=$(lsb_release -sr)

    if [[ "$dist" != "Ubuntu" ]] || [[ "$ver" != "24.04" ]]; then
        log_error "此脚本仅支持 Ubuntu 24.04，当前系统: $dist $ver"
        exit 1
    fi

    log_info "系统版本检查通过: Ubuntu $ver"
}

# 检查必要工具
check_prerequisites() {
    local missing_tools=()

    if ! command -v curl &> /dev/null; then
        missing_tools+=("curl")
    fi

    if ! command -v git &> /dev/null; then
        missing_tools+=("git")
    fi

    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi

    if ! command -v docker-compose &> /dev/null; then
        missing_tools+=("docker-compose")
    fi

    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_info "正在安装缺少的工具: ${missing_tools[*]}"
        sudo apt update
        sudo apt install -y "${missing_tools[@]}"
    else
        log_info "所有必要工具均已安装"
    fi
}

# 检查端口占用
check_ports() {
    local ports=("80" "443" "8080" "3306" "27017" "6379")
    local occupied_ports=()

    for port in "${ports[@]}"; do
        if netstat -tuln | grep -q ":$port "; then
            occupied_ports+=("$port")
        fi
    done

    if [ ${#occupied_ports[@]} -ne 0 ]; then
        log_error "以下端口已被占用: ${occupied_ports[*]}"
        log_error "请停止占用这些端口的服务后再运行此脚本"
        exit 1
    fi

    log_info "端口检查通过: ${ports[*]}"
}

# 创建项目目录
setup_project_dir() {
    PROJECT_DIR="/opt/zhixin"
    
    if [ -d "$PROJECT_DIR" ]; then
        log_warn "项目目录 $PROJECT_DIR 已存在"
        read -p "是否覆盖现有安装? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "操作已取消"
            exit 0
        fi
        sudo rm -rf "$PROJECT_DIR"
    fi

    sudo mkdir -p "$PROJECT_DIR"
    sudo chown $USER:$USER "$PROJECT_DIR"
    cd "$PROJECT_DIR"
    
    log_info "项目目录已创建: $PROJECT_DIR"
}

# 克隆项目代码
clone_project() {
    log_info "正在克隆项目代码..."
    
    git clone https://github.com/your-org/zhixin.git .
    
    if [ $? -ne 0 ]; then
        log_error "克隆项目失败"
        exit 1
    fi
    
    log_success "项目代码克隆成功"
}

# 生成强随机密钥
generate_secret_keys() {
    log_info "生成安全密钥..."

    # 生成 JWT 密钥 (32 字节 Base64 编码)
    JWT_SECRET=$(openssl rand -base64 32)
    MYSQL_ROOT_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    MYSQL_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    MONGO_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    REDIS_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    ADMIN_SECRET=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)

    # 更新 .env 文件
    sed -i "s|JWT_SECRET=.*|JWT_SECRET=$JWT_SECRET|" .env
    sed -i "s|MYSQL_ROOT_PASSWORD=.*|MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD|" .env
    sed -i "s|MYSQL_PASSWORD=.*|MYSQL_PASSWORD=$MYSQL_PASSWORD|" .env
    sed -i "s|MONGO_PASSWORD=.*|MONGO_PASSWORD=$MONGO_PASSWORD|" .env
    sed -i "s|REDIS_PASSWORD=.*|REDIS_PASSWORD=$REDIS_PASSWORD|" .env
    sed -i "s|ADMIN_SECRET=.*|ADMIN_SECRET=$ADMIN_SECRET|" .env

    log_success "安全密钥生成完成"
}

# 创建配置文件
create_config() {
    log_info "创建配置文件..."

    # 如果没有 .env 文件，则从 .env.example 创建
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            generate_secret_keys
        else
            log_error ".env.example 文件不存在"
            exit 1
        fi
    else
        log_info ".env 文件已存在，跳过生成"
    fi

    # 创建数据目录
    sudo mkdir -p /opt/zhixin-data/{mysql,mongo,redis,uploads,logs}
    sudo chown -R $USER:$USER /opt/zhixin-data

    log_success "配置文件创建完成"
}

# 启动服务
start_services() {
    log_info "正在启动服务..."

    # 使用 Docker Compose 启动所有服务
    docker-compose up -d --remove-orphans

    if [ $? -ne 0 ]; then
        log_error "服务启动失败"
        exit 1
    fi

    log_success "服务启动完成"
}

# 等待服务就绪
wait_for_services() {
    log_info "等待服务就绪..."

    # 等待数据库服务就绪
    for i in {1..30}; do
        if docker-compose exec mysql mysqladmin ping -h localhost --silent; then
            log_info "MySQL 服务就绪"
            break
        fi
        sleep 2
    done

    for i in {1..30}; do
        if docker-compose exec mongo mongosh --eval "db.runCommand({ping:1})" &>/dev/null; then
            log_info "MongoDB 服务就绪"
            break
        fi
        sleep 2
    done

    # 等待应用服务就绪
    for i in {1..60}; do
        if curl -sf http://localhost:8080/api/v1/system/config &>/dev/null; then
            log_success "应用服务就绪"
            break
        fi
        sleep 5
    done

    log_success "所有服务已就绪"
}

# 创建管理脚本
create_manage_scripts() {
    log_info "创建管理脚本..."

    # 创建管理脚本
    cat > zhixin-manage.sh << 'EOF'
#!/bin/bash

# 知信管理系统脚本

case "$1" in
    start)
        docker-compose up -d
        echo "知信服务已启动"
        ;;
    stop)
        docker-compose down
        echo "知信服务已停止"
        ;;
    restart)
        docker-compose restart
        echo "知信服务已重启"
        ;;
    logs)
        docker-compose logs -f
        ;;
    logs-backend)
        docker-compose logs -f backend
        ;;
    logs-web)
        docker-compose logs -f web
        ;;
    status)
        docker-compose ps
        ;;
    update)
        git pull
        docker-compose build
        docker-compose up -d
        echo "知信服务已更新"
        ;;
    backup)
        TIMESTAMP=$(date +%Y%m%d_%H%M%S)
        BACKUP_DIR="/opt/zhixin-backup/$TIMESTAMP"
        mkdir -p $BACKUP_DIR
        
        # 备份数据库
        docker-compose exec mysql mysqldump -u zhixin_user -p$(grep MYSQL_PASSWORD .env | cut -d'=' -f2) zhixin_chat > $BACKUP_DIR/mysql_backup.sql
        docker-compose exec mongo mongodump --out $BACKUP_DIR/mongo_backup
        
        echo "备份已保存到: $BACKUP_DIR"
        ;;
    *)
        echo "用法: $0 {start|stop|restart|logs|logs-backend|logs-web|status|update|backup}"
        exit 1
        ;;
esac
EOF

    chmod +x zhixin-manage.sh
    sudo mv zhixin-manage.sh /usr/local/bin/
    
    log_success "管理脚本已创建 (/usr/local/bin/zhixin-manage.sh)"
}

# 显示安装完成信息
show_completion_info() {
    log_success "==========================================="
    log_success "知信聊天系统安装完成！"
    log_success "==========================================="
    echo
    log_info "访问地址:"
    log_info "  - 前端界面: http://$(hostname -I | awk '{print $1}')"
    log_info "  - API 文档: http://$(hostname -I | awk '{print $1}'):8080/swagger/index.html (如果已启用)"
    echo
    log_info "管理命令:"
    log_info "  - 启动服务: zhixin-manage.sh start"
    log_info "  - 停止服务: zhixin-manage.sh stop"
    log_info "  - 重启服务: zhixin-manage.sh restart"
    log_info "  - 查看日志: zhixin-manage.sh logs"
    log_info "  - 服务状态: zhixin-manage.sh status"
    echo
    log_info "数据目录:"
    log_info "  - MySQL: /opt/zhixin-data/mysql"
    log_info "  - MongoDB: /opt/zhixin-data/mongo"
    log_info "  - Redis: /opt/zhixin-data/redis"
    log_info "  - 上传文件: /opt/zhixin-data/uploads"
    log_info "  - 日志: /opt/zhixin-data/logs"
    echo
    log_info "配置文件:"
    log_info "  - 环境变量: $PROJECT_DIR/.env"
    echo
    log_success "==========================================="
    log_success "欢迎使用知信聊天系统！"
    log_success "==========================================="
}

# 主函数
main() {
    log_info "开始安装知信聊天系统..."

    check_root
    check_system
    check_prerequisites
    check_ports
    setup_project_dir
    clone_project
    create_config
    start_services
    wait_for_services
    create_manage_scripts
    show_completion_info

    log_success "安装完成！请根据上述信息进行后续操作。"
}

# 运行主函数
main "$@"