#!/bin/bash
# ===============================================
# Chat System Pro - Linux/MacOS 部署脚本
# 支持: Linux (x86/ARM), MacOS (Intel/M1/M2/M3)
# 支持国产芯片: 鲲鹏、飞腾、龙芯、海光
# ===============================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

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

# 检测系统
detect_system() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    else
        log_error "不支持的操作系统: $OSTYPE"
        exit 1
    fi

    # 检测架构
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        mips64el)
            ARCH="mips64el"  # 龙芯
            ;;
        *)
            log_warning "未知架构: $ARCH, 默认使用 amd64"
            ARCH="amd64"
            ;;
    esac

    log_info "检测到系统: $OS, 架构: $ARCH"
}

# 检测国产CPU
detect_china_cpu() {
    if [[ "$OS" == "linux" ]]; then
        if grep -q "鲲鹏" /proc/cpuinfo 2>/dev/null || grep -q "Kunpeng" /proc/cpuinfo 2>/dev/null; then
            log_success "检测到 鲲鹏 (Kunpeng) CPU"
            CHINA_CPU="kunpeng"
        elif grep -q "飞腾" /proc/cpuinfo 2>/dev/null || grep -q "FT-" /proc/cpuinfo 2>/dev/null; then
            log_success "检测到 飞腾 (Phytium) CPU"
            CHINA_CPU="phytium"
        elif grep -q "龙芯" /proc/cpuinfo 2>/dev/null || grep -q "Loongson" /proc/cpuinfo 2>/dev/null; then
            log_success "检测到 龙芯 (Loongson) CPU"
            CHINA_CPU="loongson"
        elif grep -q "海光" /proc/cpuinfo 2>/dev/null || grep -q "Hygon" /proc/cpuinfo 2>/dev/null; then
            log_success "检测到 海光 (Hygon) CPU"
            CHINA_CPU="hygon"
        else
            log_info "未检测到国产CPU"
            CHINA_CPU=""
        fi
    fi
}

# 安装依赖
install_dependencies() {
    log_info "检查并安装依赖..."

    if [[ "$OS" == "linux" ]]; then
        if command -v apt-get &> /dev/null; then
            # Debian/Ubuntu
            log_info "使用 apt-get 包管理器"
            sudo apt-get update
            sudo apt-get install -y curl git docker.io docker-compose
            
        elif command -v yum &> /dev/null; then
            # CentOS/RHEL
            log_info "使用 yum 包管理器"
            sudo yum install -y curl git docker docker-compose
            
        elif command -v dnf &> /dev/null; then
            # Fedora
            log_info "使用 dnf 包管理器"
            sudo dnf install -y curl git docker docker-compose
            
        else
            log_warning "无法检测包管理器，请手动安装依赖"
        fi
        
    elif [[ "$OS" == "macos" ]]; then
        # MacOS
        if ! command -v brew &> /dev/null; then
            log_info "安装 Homebrew..."
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        fi
        
        log_info "使用 Homebrew 安装依赖"
        brew install curl git docker docker-compose
    fi
}

# 配置国产芯片环境
setup_china_cpu() {
    if [[ "$CHINA_CPU" == "loongson" ]]; then
        log_info "配置龙芯环境..."
        # 使用龙芯专用镜像
        export DOCKER_DEFAULT_PLATFORM=linux/mips64el
        
    elif [[ "$CHINA_CPU" == "kunpeng" ]] || [[ "$CHINA_CPU" == "phytium" ]]; then
        log_info "配置鲲鹏/飞腾环境 (ARM64)..."
        export DOCKER_DEFAULT_PLATFORM=linux/arm64
        
    elif [[ "$CHINA_CPU" == "hygon" ]]; then
        log_info "配置海光环境 (x86_64)..."
        export DOCKER_DEFAULT_PLATFORM=linux/amd64
    fi
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    cd "$PROJECT_DIR"
    
    # 选择合适的 docker-compose 文件
    if [[ "$CHINA_CPU" != "" ]]; then
        log_info "使用国产芯片专用配置"
        COMPOSE_FILE="docker-compose.china.yml"
    else
        COMPOSE_FILE="docker-compose.yml"
    fi
    
    if [[ ! -f "$COMPOSE_FILE" ]]; then
        log_warning "配置文件不存在，使用默认配置"
        COMPOSE_FILE="docker-compose.yml"
    fi
    
    # 检查 .env 文件
    if [[ ! -f ".env" ]]; then
        log_info "复制 .env.example 为 .env"
        cp .env.example .env
    fi
    
    # 启动服务
    log_info "使用配置文件: $COMPOSE_FILE"
    docker-compose -f "$COMPOSE_FILE" up -d
    
    if [[ $? -eq 0 ]]; then
        log_success "服务启动成功！"
        log_info "访问地址:"
        log_info "  前端: http://localhost"
        log_info "  API:  http://localhost:8080"
        log_info ""
        log_info "查看日志: docker-compose logs -f"
    else
        log_error "服务启动失败！"
    fi
}

# 停止服务
stop_services() {
    log_info "停止服务..."
    cd "$PROJECT_DIR"
    docker-compose down
    log_success "服务已停止"
}

# 查看状态
show_status() {
    log_info "服务状态:"
    cd "$PROJECT_DIR"
    docker-compose ps
}

# 查看日志
show_logs() {
    cd "$PROJECT_DIR"
    docker-compose logs -f
}

# 初始化数据库
init_database() {
    log_info "初始化数据库..."
    cd "$PROJECT_DIR"
    
    docker-compose exec backend ./chat-system-pro init-db
    
    log_success "数据库初始化完成！"
}

# 主菜单
show_menu() {
    clear
    echo "============================================="
    echo "  Chat System Pro - 部署管理工具"
    echo "============================================="
    echo "  系统: $OS ($ARCH)"
    if [[ "$CHINA_CPU" != "" ]]; then
        echo "  国产芯片: $CHINA_CPU"
    fi
    echo "============================================="
    echo ""
    echo "  1. 一键部署 (首次安装)"
    echo "  2. 启动服务"
    echo "  3. 停止服务"
    echo "  4. 重启服务"
    echo "  5. 查看状态"
    echo "  6. 查看日志"
    echo "  7. 初始化数据库"
    echo "  8. 更新服务"
    echo "  0. 退出"
    echo ""
    echo -n "请选择操作 [0-8]: "
}

# 主函数
main() {
    detect_system
    detect_china_cpu
    
    while true; do
        show_menu
        read choice
        
        case $choice in
            1)
                install_dependencies
                setup_china_cpu
                start_services
                ;;
            2)
                start_services
                ;;
            3)
                stop_services
                ;;
            4)
                stop_services
                start_services
                ;;
            5)
                show_status
                ;;
            6)
                show_logs
                ;;
            7)
                init_database
                ;;
            8)
                log_info "更新服务..."
                git pull
                docker-compose pull
                stop_services
                start_services
                ;;
            0)
                log_info "再见！"
                exit 0
                ;;
            *)
                log_error "无效选择"
                ;;
        esac
        
        echo ""
        echo -n "按 Enter 键继续..."
        read
    done
}

# 命令行参数
if [[ $# -gt 0 ]]; then
    case "$1" in
        start)
            detect_system
            detect_china_cpu
            start_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            stop_services
            start_services
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs
            ;;
        *)
            echo "使用方法: $0 [start|stop|restart|status|logs]"
            exit 1
            ;;
    esac
else
    # 交互模式
    main
fi
