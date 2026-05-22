#!/bin/bash
# ===============================================
# Chat System Pro - 设置系统自启动
# ===============================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then
    log_error "请使用 sudo 运行此脚本"
    exit 1
fi

echo ""
echo "================================================"
echo "  Chat System Pro - 设置系统自启动"
echo "================================================"
echo ""

# 检测项目路径
PROJECT_DIR=""
if [ -d "/home/fjya/jinjuyi" ]; then
    PROJECT_DIR="/home/fjya/jinjuyi"
elif [ -d "$HOME/jinjuyi" ]; then
    PROJECT_DIR="$HOME/jinjuyi"
elif [ -d "$(pwd)/jinjuyi" ]; then
    PROJECT_DIR="$(pwd)/jinjuyi"
else
    log_error "未找到项目目录，请先克隆项目"
    exit 1
fi

log_info "检测到项目目录: $PROJECT_DIR"

# 检测 docker compose 命令
COMPOSE_CMD=""
if command -v docker &> /dev/null && docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
elif command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
else
    log_error "未找到 docker-compose 命令"
    exit 1
fi

log_info "使用命令: $COMPOSE_CMD"

# 创建 systemd 服务文件
SERVICE_NAME="chat-system-pro"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

log_info "创建 systemd 服务文件..."

cat > "$SERVICE_FILE" << EOF
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
User=fjya
Group=docker

[Install]
WantedBy=multi-user.target
EOF

log_success "服务文件已创建: $SERVICE_FILE"

# 重新加载 systemd
log_info "重新加载 systemd 配置..."
systemctl daemon-reload

# 启用服务（开机自启动）
log_info "启用服务自启动..."
systemctl enable ${SERVICE_NAME}.service

# 立即启动服务
log_info "立即启动服务..."
systemctl start ${SERVICE_NAME}.service

# 检查服务状态
sleep 3
if systemctl is-active --quiet ${SERVICE_NAME}.service; then
    log_success "服务启动成功"
else
    log_error "服务启动失败，请检查日志"
    journalctl -u ${SERVICE_NAME}.service --no-pager -n 20
    exit 1
fi

echo ""
echo "================================================"
echo -e "${GREEN}  设置完成！${NC}"
echo "================================================"
echo ""
echo "  服务名称: ${SERVICE_NAME}"
echo "  项目路径: ${PROJECT_DIR}"
echo ""
echo "  常用命令:"
echo "    systemctl start ${SERVICE_NAME}      # 启动服务"
echo "    systemctl stop ${SERVICE_NAME}       # 停止服务"
echo "    systemctl restart ${SERVICE_NAME}    # 重启服务"
echo "    systemctl status ${SERVICE_NAME}     # 查看状态"
echo "    systemctl enable ${SERVICE_NAME}    # 启用自启动"
echo "    journalctl -u ${SERVICE_NAME} -f     # 查看日志"
echo ""
echo "  下次重启后，服务将自动启动"
echo ""
