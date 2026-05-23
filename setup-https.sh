#!/bin/bash

# ============================================
# 知信聊天系统 - HTTPS 自动配置脚本
# ============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印彩色消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否为root用户
check_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "请使用 sudo 运行此脚本"
        exit 1
    fi
}

# 检查域名参数
check_domain() {
    if [ -z "$1" ]; then
        echo ""
        print_info "用法: sudo $0 <your-domain.com>"
        echo ""
        print_info "示例: sudo $0 chat.example.com"
        echo ""
        exit 1
    fi
}

# 检查Docker是否运行
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker未运行，请先启动Docker"
        exit 1
    fi
}

# 安装必要软件
install_dependencies() {
    print_info "安装必要软件..."
    
    # 更新包列表
    apt-get update -qq
    
    # 安装必要软件
    apt-get install -y -qq curl wget git vim certbot python3-certbot-nginx ufw
    
    print_success "依赖软件安装完成"
}

# 配置防火墙
configure_firewall() {
    print_info "配置防火墙..."
    
    # 启用UFW
    ufw --force enable
    
    # 允许SSH
    ufw allow ssh
    
    # 允许HTTP和HTTPS
    ufw allow http
    ufw allow https
    
    # 允许WebSocket（如果使用）
    ufw allow 8080/tcp
    
    print_success "防火墙配置完成"
}

# 获取SSL证书
get_ssl_certificate() {
    DOMAIN=$1
    EMAIL=$2
    
    print_info "获取SSL证书 for $DOMAIN..."
    
    # 停止nginx服务（如果运行）
    docker stop $(docker ps -q --filter name=chat-frontend) 2>/dev/null || true
    systemctl stop nginx 2>/dev/null || true
    
    # 创建证书目录
    mkdir -p /etc/letsencrypt/live/$DOMAIN
    mkdir -p /var/www/html
    
    # 获取证书
    certbot certonly --standalone \
        -d $DOMAIN \
        -d www.$DOMAIN \
        --agree-tos \
        --email $EMAIL \
        --non-interactive \
        --http-01-port 80 \
        --keep-until-expiring
    
    if [ $? -eq 0 ]; then
        print_success "SSL证书获取成功"
    else
        print_error "SSL证书获取失败"
        exit 1
    fi
}

# 配置Nginx with HTTPS
configure_nginx_https() {
    DOMAIN=$1
    
    print_info "配置Nginx HTTPS..."
    
    cat > /etc/nginx/sites-available/chat-system-https << 'EOF'
server {
    listen 80;
    server_name YOUR_DOMAIN www.YOUR_DOMAIN;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name YOUR_DOMAIN www.YOUR_DOMAIN;

    # SSL证书配置
    ssl_certificate /etc/letsencrypt/live/YOUR_DOMAIN/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/YOUR_DOMAIN/privkey.pem;
    ssl_trusted_certificate /etc/letsencrypt/live/YOUR_DOMAIN/chain.pem;

    # SSL安全配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # 安全头部
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # 前端静态文件
    root /var/www/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api/ {
        proxy_pass http://localhost:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket代理
    location /ws {
        proxy_pass http://localhost:8080/ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_read_timeout 86400;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # 日志
    access_log /var/log/nginx/chat_access.log;
    error_log /var/log/nginx/chat_error.log;
}
EOF

    # 替换域名
    sed -i "s/YOUR_DOMAIN/$DOMAIN/g" /etc/nginx/sites-available/chat-system-https
    
    # 启用站点
    ln -sf /etc/nginx/sites-available/chat-system-https /etc/nginx/sites-enabled/
    
    # 删除默认站点
    rm -f /etc/nginx/sites-enabled/default
    
    # 测试配置
    nginx -t
    
    # 重载nginx
    systemctl reload nginx
    
    print_success "Nginx HTTPS配置完成"
}

# 配置自动续期
setup_auto_renewal() {
    print_info "配置SSL自动续期..."
    
    # 创建续期脚本
    cat > /etc/cron.d/certbot-renewal << 'EOF'
# 每天凌晨2点检查并续期证书
0 2 * * * root certbot renew --quiet --deploy-hook "systemctl reload nginx"
EOF

    chmod 644 /etc/cron.d/certbot-renewal
    
    print_success "自动续期配置完成"
}

# 启动服务
start_services() {
    print_info "启动服务..."
    
    # 重载nginx
    systemctl restart nginx
    
    print_success "服务启动完成"
}

# 显示完成信息
show_completion() {
    DOMAIN=$1
    
    echo ""
    print_success "============================================"
    print_success "  HTTPS配置完成！"
    print_success "============================================"
    echo ""
    print_info "访问地址: https://$DOMAIN"
    print_info "API地址: https://$DOMAIN/api/v1"
    print_info ""
    print_info "SSL证书会自动续期，无需手动操作"
    print_info ""
    print_warning "请确保在DNS服务商处配置好A记录:"
    print_warning "  @    A    你的服务器IP"
    print_warning "  www  A    你的服务器IP"
    echo ""
}

# 主函数
main() {
    echo ""
    print_info "=========================================="
    print_info " 知信聊天系统 - HTTPS自动配置脚本"
    print_info "=========================================="
    echo ""
    
    # 检查
    check_root
    check_domain $1
    check_docker
    
    # 变量
    DOMAIN=$1
    EMAIL="admin@$DOMAIN"
    
    # 安装依赖
    install_dependencies
    
    # 配置防火墙
    configure_firewall
    
    # 获取SSL证书
    get_ssl_certificate $DOMAIN $EMAIL
    
    # 配置Nginx
    configure_nginx_https $DOMAIN
    
    # 配置自动续期
    setup_auto_renewal
    
    # 启动服务
    start_services
    
    # 显示完成信息
    show_completion $DOMAIN
}

# 运行主函数
main $@