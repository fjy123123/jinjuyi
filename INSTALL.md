# Chat System Pro - 完整安装部署文档

## 目录

1. [系统概述](#1-系统概述)
2. [环境要求](#2-环境要求)
3. [快速开始](#3-快速开始)
4. [Windows部署](#4-windows部署)
5. [Linux部署](#5-linux部署)
6. [macOS部署](#6-macos部署)
7. [Docker Compose部署](#7-docker-compose部署)
8. [手动编译部署](#8-手动编译部署)
9. [国产芯片部署](#9-国产芯片部署)
10. [配置详解](#10-配置详解)
11. [域名与HTTPS配置](#11-域名与https配置)
12. [数据管理](#12-数据管理)
13. [运维监控](#13-运维监控)
14. [常见问题](#14-常见问题)
15. [技术支持](#15-技术支持)

---

## 1. 系统概述

### 1.1 系统简介

Chat System Pro 是一款功能完善的即时通讯系统，采用前后端分离架构，支持私聊、群聊、朋友圈、红包、支付等丰富功能。系统采用Go语言开发，具备高性能、高并发、易扩展等特点，适合企业级应用和商业化运营。

### 1.2 核心特性

- **实时通讯**：基于WebSocket的实时消息推送，支持万人群聊
- **红包功能**：支持积分、微信、支付宝等多种支付方式
- **朋友圈**：支持图片、文字、位置等多种形式分享
- **多端支持**：支持Web端、PC端、移动端（iOS/Android）
- **跨平台**：支持Windows、Linux、macOS等主流操作系统
- **国产化支持**：支持鲲鹏、飞腾、龙芯、海光等国产芯片
- **高可用**：支持水平扩展，满足大规模用户使用
- **安全可靠**：支持JWT认证、数据加密、安全防护等

### 1.3 技术架构

```
┌─────────────────────────────────────────────────────────┐
│                        客户端层                          │
├──────────────┬──────────────┬──────────────┬──────────────┤
│   Web端     │   PC端       │  Android    │    iOS      │
│   React    │  Electron    │   UniApp    │   UniApp    │
└──────────────┴──────────────┴──────────────┴──────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                        网关层                           │
├─────────────────────────────────────────────────────────┤
│   Nginx (负载均衡 / SSL终结 / 静态资源)                  │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                        服务层                           │
├──────────────┬──────────────┬──────────────┬──────────────┤
│  API 服务    │  WebSocket   │   推送服务   │  定时任务    │
│  Gin框架     │  长连接管理   │  极光/个推   │  消息归档    │
└──────────────┴──────────────┴──────────────┴──────────────┘
                           │
          ┌────────────────┼────────────────┐
          ▼                ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│    MySQL     │  │   Redis      │  │   MongoDB    │
│   业务数据   │  │  缓存/会话   │  │   消息存储   │
└──────────────┘  └──────────────┘  └──────────────┘
```

### 1.4 支持规模

| 部署规模 | CPU | 内存 | 磁盘 | 并发用户 | 消息并发 |
|---------|-----|------|------|---------|---------|
| 开发测试 | 1核 | 2GB | 20GB | <100 | <1000条/秒 |
| 小型生产 | 2核 | 4GB | 50GB | <5000 | <5000条/秒 |
| 中型生产 | 4核 | 8GB | 100GB | <50000 | <20000条/秒 |
| 大型生产 | 8核+ | 16GB+ | 500GB+ | >50000 | >50000条/秒 |

---

## 2. 环境要求

### 2.1 硬件要求

#### 开发测试环境
- CPU: 1核以上
- 内存: 2GB以上
- 磁盘: 20GB以上
- 网络: 1Mbps以上

#### 小型生产环境
- CPU: 2核以上
- 内存: 4GB以上
- 磁盘: 50GB以上（推荐SSD）
- 网络: 5Mbps以上

#### 中型生产环境
- CPU: 4核以上
- 内存: 8GB以上
- 磁盘: 100GB以上（推荐SSD）
- 网络: 10Mbps以上

#### 大型生产环境
- CPU: 8核以上
- 内存: 16GB以上
- 磁盘: 500GB以上（推荐SSD）
- 网络: 100Mbps以上

### 2.2 软件要求

#### 操作系统
- **Windows**: Windows 10/11, Windows Server 2016+
- **Linux**: Ubuntu 20.04+, CentOS 7+, Debian 10+, Fedora 36+
- **macOS**: macOS 11+ (Intel / Apple Silicon)

#### 运行环境
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Go**: 1.21+ (手动编译时需要)
- **Node.js**: 18+ (前端开发时需要)

#### 数据库
- **MySQL**: 8.0+
- **MongoDB**: 6.0+
- **Redis**: 7.0+

### 2.3 支持的CPU架构

| 架构 | 说明 | 支持芯片 |
|------|------|---------|
| x86_64 (amd64) | Intel/AMD | Intel, AMD, 海光(Hygon) |
| ARM64 (aarch64) | 64位ARM | 鲲鹏(Kunpeng), 飞腾(Phytium), Apple Silicon |
| MIPS64 (mips64el) | 龙芯64位 | 龙芯(Loongson) |

---

## 3. 快速开始

### 3.1 一键部署（推荐）

#### Windows
```cmd
cd chat-system-pro
deploy.bat
```

#### Linux/macOS
```bash
chmod +x deploy.sh
./deploy.sh
```

### 3.2 Docker快速启动

```bash
# 1. 克隆项目
git clone https://github.com/yourrepo/chat-system-pro.git
cd chat-system-pro

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，修改数据库密码和JWT密钥

# 3. 启动服务
docker compose up -d

# 4. 查看服务状态
docker compose ps

# 5. 访问系统
# 前端: http://localhost
# API: http://localhost:8080
# 健康检查: http://localhost:8080/health
```

### 3.3 手动编译启动

```bash
# 1. 克隆项目
git clone https://github.com/yourrepo/chat-system-pro.git
cd chat-system-pro

# 2. 安装数据库服务（MySQL, MongoDB, Redis）

# 3. 编译后端
cd backend
go mod tidy
go build -o chat-server main.go

# 4. 配置并启动
cp config.yaml.example config.yaml
# 编辑 config.yaml
./chat-server

# 5. 编译前端
cd ../web
npm install
npm run build
```

### 3.4 验证部署

#### API健康检查
```bash
curl http://localhost:8080/health
```

响应：
```json
{
  "status": "ok"
}
```

#### WebSocket连接测试
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?token=YOUR_TOKEN');
ws.onopen = () => console.log('Connected');
ws.onmessage = (event) => console.log(event.data);
```

#### 登录测试
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'
```

---

## 4. Windows部署

### 4.1 安装Docker Desktop

#### 方法一：下载安装包
1. 访问 https://www.docker.com/products/docker-desktop/
2. 下载 Windows 版本安装包
3. 双击运行安装程序
4. 按照向导完成安装
5. 重启电脑
6. 启动 Docker Desktop

#### 方法二：使用包管理器
```powershell
# 使用 Chocolatey
choco install docker-desktop

# 使用 Winget
winget install Docker.DockerDesktop
```

#### 验证安装
```cmd
docker --version
docker compose version
```

### 4.2 配置Docker Desktop

1. 打开 Docker Desktop
2. 进入 Settings → General
3. 确保勾选 "Use Docker Desktop VMM" 和 "Use the WSL 2 based engine"
4. 进入 Resources
5. 配置 CPU: 2+, Memory: 4GB+, Disk: 50GB+
6. 点击 Apply & Restart

### 4.3 部署应用

#### 一键部署
```cmd
cd chat-system-pro
deploy.bat
```

#### 手动部署
```cmd
# 1. 进入项目目录
cd chat-system-pro

# 2. 复制环境配置
copy .env.example .env

# 3. 编辑配置
notepad .env
# 修改以下内容：
# MYSQL_ROOT_PASSWORD=your_secure_password
# JWT_SECRET=your-super-secret-key-at-least-32-chars

# 4. 启动服务
docker compose up -d

# 5. 查看状态
docker compose ps

# 6. 查看日志
docker compose logs -f backend
```

### 4.4 防火墙配置

如果需要从外部访问，需要配置防火墙：

```powershell
# 添加防火墙规则
New-NetFirewallRule -DisplayName "Chat System HTTP" -Direction Inbound -Protocol TCP -LocalPort 80,443,8080 -Action Allow

# 或使用高级防火墙
netsh advfirewall firewall add rule name="Chat System" dir=in action=allow protocol=tcp localport=80,443,8080
```

### 4.5 开机自启配置

```powershell
# 创建启动脚本 start.bat
@echo off
cd /d C:\path\to\chat-system-pro
docker compose up -d
```

使用任务计划程序设置为开机启动：
1. 打开"任务计划程序"
2. 创建基本任务
3. 触发器：计算机启动
4. 操作：启动程序，选择 start.bat

---

## 5. Linux部署

### 5.1 Ubuntu/Debian

#### 安装Docker
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装依赖
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common

# 添加Docker GPG密钥
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# 添加Docker仓库
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 安装Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 启动Docker
sudo systemctl start docker
sudo systemctl enable docker

# 添加当前用户到docker组
sudo usermod -aG docker $USER
newgrp docker
```

#### 安装Docker Compose（独立版）
```bash
# 下载Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# 添加执行权限
sudo chmod +x /usr/local/bin/docker-compose

# 创建软链接
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

# 验证安装
docker-compose --version
```

### 5.2 CentOS/RHEL

#### 安装Docker
```bash
# 安装依赖
sudo yum install -y yum-utils device-mapper-persistent-data lvm2

# 添加Docker仓库
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

# 安装Docker
sudo yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 启动Docker
sudo systemctl start docker
sudo systemctl enable docker

# 添加当前用户到docker组
sudo usermod -aG docker $USER
exit
# 重新登录
```

#### 安装Docker Compose
```bash
# 下载Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# 添加执行权限
sudo chmod +x /usr/local/bin/docker-compose

# 创建软链接
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
```

### 5.3 部署应用

#### 一键部署脚本
```bash
cd chat-system-pro
chmod +x deploy.sh
./deploy.sh
```

脚本会自动：
1. 检测系统环境
2. 安装Docker（如果未安装）
3. 配置数据库
4. 构建镜像
5. 启动服务
6. 配置防火墙

#### 手动部署
```bash
# 1. 进入项目目录
cd chat-system-pro

# 2. 复制环境配置
cp .env.example .env

# 3. 编辑配置
vi .env
# 修改以下内容：
# MYSQL_ROOT_PASSWORD=your_secure_password
# MYSQL_DATABASE=chat_system_pro
# JWT_SECRET=your-super-secret-key

# 4. 启动服务
docker compose up -d

# 5. 查看状态
docker compose ps

# 6. 查看日志
docker compose logs -f backend
```

### 5.4 配置防火墙

#### UFW (Ubuntu)
```bash
# 开放端口
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp

# 启用防火墙
sudo ufw enable

# 查看状态
sudo ufw status
```

#### firewalld (CentOS)
```bash
# 开放端口
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --permanent --add-port=8080/tcp

# 重载防火墙
sudo firewall-cmd --reload

# 查看状态
sudo firewall-cmd --list-ports
```

### 5.5 配置开机自启

```bash
# 创建systemd服务文件
sudo nano /etc/systemd/system/chat-system.service
```

```ini
[Unit]
Description=Chat System Pro
Requires=docker.service
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/path/to/chat-system-pro
ExecStart=/usr/local/bin/docker compose up -d
ExecStop=/usr/local/bin/docker compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

```bash
# 启用服务
sudo systemctl enable chat-system.service

# 管理服务
sudo systemctl start chat-system.service
sudo systemctl stop chat-system.service
sudo systemctl restart chat-system.service
sudo systemctl status chat-system.service
```

---

## 6. macOS部署

### 6.1 安装Docker Desktop

#### 方法一：Homebrew
```bash
brew install --cask docker
```

#### 方法二：下载安装包
1. 访问 https://www.docker.com/products/docker-desktop/
2. 下载 macOS 版本安装包 (.dmg)
3. 双击安装
4. 将 Docker.app 拖入应用程序文件夹
5. 启动 Docker Desktop

#### 验证安装
```bash
docker --version
docker compose version
```

### 6.2 部署应用

```bash
# 1. 进入项目目录
cd chat-system-pro

# 2. 添加执行权限
chmod +x deploy.sh

# 3. 运行部署脚本
./deploy.sh
```

或者手动部署：
```bash
# 1. 复制环境配置
cp .env.example .env

# 2. 编辑配置
nano .env

# 3. 启动服务
docker compose up -d

# 4. 查看状态
docker compose ps
```

### 6.3 Apple Silicon 注意事项

如果使用 Apple Silicon (M1/M2/M3)：
1. Docker Desktop 会自动使用 ARM64 架构
2. 大部分镜像都有 ARM64 版本，无需额外配置
3. 如遇兼容性问题，可使用 Rosetta 2 转译：
```bash
# 在docker-compose.yml中指定平台
services:
  backend:
    platform: linux/amd64
    # ...
```

---

## 7. Docker Compose部署

### 7.1 目录结构

```
chat-system-pro/
├── backend/                    # 后端代码
├── web/                       # 前端代码
├── docker/                    # Docker配置
│   ├── nginx.conf            # Nginx配置
│   └── ssl/                  # SSL证书
├── data/                      # 数据持久化
│   ├── mysql/                # MySQL数据
│   ├── mongodb/              # MongoDB数据
│   └── redis/                # Redis数据
├── logs/                      # 日志文件
│   ├── backend/               # 后端日志
│   └── nginx/                # Nginx日志
├── .env                       # 环境变量
├── docker-compose.yml         # Docker配置（标准）
├── docker-compose.china.yml   # Docker配置（国产芯片）
└── deploy.sh                  # 部署脚本
```

### 7.2 环境变量配置

创建 `.env` 文件：

```env
# ========== 基础配置 ==========
COMPOSE_PROJECT_NAME=chat-system-pro
TZ=Asia/Shanghai

# ========== 域名配置 ==========
DOMAIN=your-domain.com
SSL_EMAIL=admin@your-domain.com

# ========== MySQL配置 ==========
MYSQL_ROOT_PASSWORD=your_secure_root_password_at_least_32_chars
MYSQL_DATABASE=chat_system_pro
MYSQL_USER=chat_user
MYSQL_PASSWORD=your_mysql_password

# ========== MongoDB配置 ==========
MONGO_INITDB_DATABASE=chat_system_pro

# ========== Redis配置 ==========
REDIS_PASSWORD=your_redis_password
REDIS_DATABASE=0

# ========== JWT配置 ==========
JWT_SECRET=your-super-secret-jwt-key-at-least-32-characters-long
JWT_EXPIRE_HOURS=720

# ========== 服务模式 ==========
GIN_MODE=release
NODE_ENV=production

# ========== 平台配置（国产芯片用） ==========
PLATFORM=linux/arm64
```

### 7.3 启动服务

#### 启动所有服务
```bash
docker compose up -d
```

#### 启动特定服务
```bash
# 只启动数据库
docker compose up -d mysql mongodb redis

# 只启动后端
docker compose up -d backend
```

#### 后台启动
```bash
docker compose up -d
```

#### 前台运行（查看日志）
```bash
docker compose up
```

### 7.4 服务管理

#### 查看状态
```bash
docker compose ps
```

#### 查看日志
```bash
# 查看所有日志
docker compose logs -f

# 查看特定服务日志
docker compose logs -f backend
docker compose logs -f mysql
docker compose logs -f redis

# 查看最近100行日志
docker compose logs --tail 100 backend
```

#### 重启服务
```bash
# 重启所有服务
docker compose restart

# 重启特定服务
docker compose restart backend
```

#### 停止服务
```bash
# 停止所有服务（保留数据）
docker compose stop

# 停止并删除容器（保留数据）
docker compose down

# 停止并删除容器和数据卷（删除所有数据）
docker compose down -v
```

#### 重建服务
```bash
# 重新构建镜像
docker compose build --no-cache backend

# 重新构建并启动
docker compose up -d --build
```

### 7.5 更新服务

#### 手动更新
```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker compose up -d --build

# 清理旧镜像
docker image prune -f
```

#### 自动更新（使用Watchtower）
```bash
# 启动Watchtower
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --restart unless-stopped \
  containrrr/watchtower \
  chat-system-backend chat-system-frontend

# 查看更新日志
docker logs -f watchtower
```

---

## 8. 手动编译部署

### 8.1 安装Go环境

#### Linux/macOS
```bash
# 下载Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# 解压
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 配置环境变量
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

#### Windows
1. 下载 https://go.dev/dl/go1.21.0.windows-amd64.msi
2. 运行安装程序
3. 验证：打开CMD，输入 `go version`

### 8.2 安装Node.js

#### Linux/macOS
```bash
# 使用nvm安装
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 18
nvm use 18
node -v
npm -v
```

#### Windows
下载安装包：https://nodejs.org/

### 8.3 安装数据库

#### MySQL
```bash
# Ubuntu/Debian
sudo apt install mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# CentOS
sudo yum install mysql-server
sudo systemctl start mysqld
sudo systemctl enable mysqld
```

#### MongoDB
```bash
# Ubuntu/Debian
curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | sudo gpg --dearmor -o /usr/share/keyrings/mongodb-server-7.0.gpg
echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
sudo apt update
sudo apt install mongodb-org
sudo systemctl start mongod
sudo systemctl enable mongod
```

#### Redis
```bash
# Ubuntu/Debian
sudo apt install redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server

# CentOS
sudo yum install epel-release
sudo yum install redis
sudo systemctl start redis
sudo systemctl enable redis
```

### 8.4 编译后端

```bash
cd backend

# 安装依赖
go mod tidy

# 编译当前平台
go build -o chat-server main.go

# 交叉编译
## Linux x86_64
GOOS=linux GOARCH=amd64 go build -o chat-server-linux-amd64 main.go

## Linux ARM64
GOOS=linux GOARCH=arm64 go build -o chat-server-linux-arm64 main.go

## Windows
GOOS=windows GOARCH=amd64 go build -o chat-server.exe main.go

## macOS Intel
GOOS=darwin GOARCH=amd64 go build -o chat-server-darwin-amd64 main.go

## macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o chat-server-darwin-arm64 main.go
```

### 8.5 编译前端

```bash
cd web

# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build
# 产物在 dist/ 目录
```

### 8.6 配置Nginx

```bash
# 安装Nginx
sudo apt install nginx  # Ubuntu/Debian
sudo yum install nginx  # CentOS

# 配置
sudo nano /etc/nginx/sites-available/chat-system
```

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /var/www/chat-system/dist;
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # WebSocket代理
    location /ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }
}
```

```bash
# 启用配置
sudo ln -s /etc/nginx/sites-available/chat-system /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 8.7 启动服务

#### 后端服务
```bash
cd backend
./chat-server
# 或后台运行
nohup ./chat-server > ../logs/backend.log 2>&1 &
```

#### 使用supervisor管理进程
```bash
sudo apt install supervisor
sudo nano /etc/supervisor/conf.d/chat-system.conf
```

```ini
[program:chat-system]
command=/path/to/backend/chat-server
directory=/path/to/backend
autostart=true
autorestart=true
stderr_logfile=/var/log/chat-system.err.log
stdout_logfile=/var/log/chat-system.out.log
user=www-data
```

```bash
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start chat-system
```

---

## 9. 国产芯片部署

### 9.1 支持的国产芯片

| 芯片 | 架构 | 操作系统 | Docker镜像 |
|------|------|---------|-----------|
| 鲲鹏920 | ARM64 | Ubuntu 20.04, 麒麟V10, 统信UOS | linux/arm64 |
| 飞腾FT-2000 | ARM64 | 银河麒麟, 麒麟V10 | linux/arm64 |
| 龙芯3A5000 | MIPS64 | 统信UOS, 麒麟V10 | linux/mips64el |
| 海光3000 | x86_64 | 银河麒麟, 麒麟V10, 中科方德 | linux/amd64 |

### 9.2 国产操作系统

#### 麒麟V10
```bash
# 安装Docker
sudo apt install docker.io docker-compose

# 或使用官方安装脚本
curl -fsSL https://get.docker.com | sudo sh
```

#### 统信UOS
```bash
# 安装Docker
sudo apt update
sudo apt install docker.io docker-compose

# 启用Docker服务
sudo systemctl enable docker
sudo systemctl start docker
```

#### openEuler
```bash
# 安装Docker
sudo dnf install docker
sudo systemctl enable docker
sudo systemctl start docker
```

### 9.3 部署步骤

#### 自动部署
```bash
chmod +x deploy.sh
./deploy.sh
# 选择国产芯片部署选项
```

#### 手动部署
```bash
# 1. 创建Docker配置文件
mkdir -p /etc/docker
sudo nano /etc/docker/daemon.json
```

```json
{
  "experimental": true,
  "features": {
    "buildkit": true
  },
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com"
  ]
}
```

```bash
# 2. 重启Docker
sudo systemctl daemon-reload
sudo systemctl restart docker

# 3. 使用国产芯片配置文件
cd chat-system-pro
cp .env.example .env
vi .env
# 设置 PLATFORM=linux/arm64 或 linux/mips64el

# 4. 启动服务
docker compose -f docker-compose.china.yml up -d
```

### 9.4 性能优化

#### ARM64优化
```bash
# 启用ARM64优化
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

# 构建时指定平台
docker buildx build --platform linux/arm64 -t myapp:latest .
```

#### 内存优化
```bash
# 对于鲲鹏/飞腾，建议分配8GB以上内存给Docker
# 在Docker Desktop或dockerd配置中设置
```

---

## 10. 配置详解

### 10.1 后端配置 (config.yaml)

```yaml
# 服务配置
server:
  port: 8080                    # 监听端口
  mode: release                 # 运行模式：debug / release
  name: ChatSystemPro           # 服务名称
  read_timeout: 60              # 读超时（秒）
  write_timeout: 60             # 写超时（秒）
  max_header_bytes: 1048576     # 最大头部长度

# 数据库配置
database:
  host: localhost               # MySQL主机
  port: 3306                    # MySQL端口
  user: root                    # 用户名
  password: your_password       # 密码
  dbname: chat_system_pro       # 数据库名
  charset: utf8mb4              # 字符集
  max_open_conns: 100           # 最大连接数
  max_idle_conns: 10           # 空闲连接数
  conn_max_lifetime: 3600       # 连接最大生存时间（秒）

# MongoDB配置
mongodb:
  host: localhost               # MongoDB主机
  port: 27017                   # 端口
  user:                        # 用户名（可选）
  password:                     # 密码（可选）
  dbname: chat_system_pro       # 数据库名
  auth_source: admin           # 认证数据库
  max_pool_size: 100           # 最大连接数
  min_pool_size: 10             # 最小连接数

# Redis配置
redis:
  host: localhost               # Redis主机
  port: 6379                    # 端口
  password:                     # 密码（可选）
  db: 0                         # 数据库编号
  pool_size: 100                # 连接池大小
  min_idle_conns: 10            # 最小空闲连接

# JWT配置
jwt:
  secret: your_jwt_secret       # JWT密钥（至少32字符）
  expire_hours: 720             # 过期时间（30天）
  refresh_expire_hours: 8640    # 刷新令牌过期时间（360天）

# 文件上传配置
upload:
  path: ./uploads               # 上传目录
  max_size: 10485760            # 最大文件大小（10MB）
  allowed_exts:                 # 允许的文件扩展名
    - jpg
    - jpeg
    - png
    - gif
    - mp3
    - mp4
    - pdf
    - doc
    - docx

# WebSocket配置
websocket:
  read_buffer_size: 4096        # 读缓冲区大小
  write_buffer_size: 4096       # 写缓冲区大小
  ping_interval: 30             # 心跳间隔（秒）
  pong_timeout: 60              # 心跳超时（秒）
  max_message_size: 65535       # 最大消息大小

# 支付配置
payment:
  stripe_secret_key:             # Stripe密钥
  wechat_pay_enabled: false      # 启用微信支付
  wechat_mch_id:                # 微信商户号
  wechat_apiv3_key:             # 微信APIv3密钥
  alipay_enabled: false          # 启用支付宝
  alipay_app_id:                # 支付宝应用ID
  alipay_private_key:           # 支付宝私钥
  alipay_public_key:            # 支付宝公钥

# 安全配置
security:
  invite_code_enabled: false     # 启用邀请码
  captcha_enabled: false        # 启用验证码
  rate_limit_enabled: true      # 启用限流
  rate_limit:
    rps: 100                    # 每秒请求数
    burst: 200                  # 突发流量
  cors_enabled: true            # 启用CORS
  cors_origins:                 # CORS白名单
    - http://localhost:3000
    - https://your-domain.com

# 存储配置
storage:
  type: local                   # 存储类型：local / oss / s3
  # 阿里云OSS配置
  oss:
    endpoint: oss-cn-hangzhou.aliyuncs.com
    access_key: your_access_key
    secret_key: your_secret_key
    bucket: your-bucket
    domain: https://your-domain.com
  # AWS S3配置
  s3:
    region: us-east-1
    endpoint: s3.amazonaws.com
    access_key: your_access_key
    secret_key: your_secret_key
    bucket: your-bucket
    domain: https://your-domain.com

# 推送服务配置
push:
  provider: jpush               # 推送服务：jpush / getui
  jpush:
    app_key: your_app_key
    app_secret: your_app_secret
  getui:
    app_id: your_app_id
    app_key: your_app_key
    master_secret: your_master_secret

# 微信小程序配置
miniprogram:
  wechat:
    app_id: your_app_id
    app_secret: your_app_secret

# 系统配置
system:
  ui_default: modern            # 默认UI模板
  message_retention_days: 365   # 消息保留天数
  file_retention_days: 90       # 文件保留天数
  max_group_members: 500        # 最大群成员数
  max_friends: 5000            # 最大好友数
```

### 10.2 环境变量配置 (.env)

环境变量会覆盖config.yaml中的配置：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=chat_system_pro

# MongoDB配置
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB=chat_system_pro

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# JWT配置
JWT_SECRET=your-super-secret-key-at-least-32-chars
JWT_EXPIRE_HOURS=720

# 服务配置
GIN_MODE=release
SERVER_PORT=8080

# 存储配置
STORAGE_TYPE=local
UPLOAD_PATH=./uploads

# 支付配置
STRIPE_SECRET_KEY=sk_test_xxx
WECHAT_PAY_ENABLED=false
ALIPAY_ENABLED=false

# 安全配置
INVITE_CODE_ENABLED=false
RATE_LIMIT_ENABLED=true
```

---

## 11. 域名与HTTPS配置

### 11.1 域名解析

在DNS服务商处添加以下记录：

| 记录类型 | 主机记录 | 记录值 |
|---------|---------|--------|
| A | chat | 服务器IP |
| A | api | 服务器IP |

### 11.2 Nginx + SSL配置

```bash
# 安装certbot
sudo apt install certbot python3-certbot-nginx

# 获取SSL证书
sudo certbot --nginx -d chat.your-domain.com -d api.your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

完整Nginx配置：

```nginx
# HTTP到HTTPS重定向
server {
    listen 80;
    server_name chat.your-domain.com api.your-domain.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS配置
server {
    listen 443 ssl http2;
    server_name chat.your-domain.com;

    # SSL证书
    ssl_certificate /etc/letsencrypt/live/chat.your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/chat.your-domain.com/privkey.pem;
    ssl_trusted_certificate /etc/letsencrypt/live/chat.your-domain.com/chain.pem;

    # SSL优化
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # 前端
    location / {
        root /var/www/chat-system/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # WebSocket代理
    location /ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        # WebSocket超时
        proxy_read_timeout 86400;
        proxy_send_timeout 86400;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        root /var/www/chat-system/dist;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Gzip压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/json application/xml;
}
```

### 11.3 使用自有SSL证书

如果使用商业SSL证书：

```nginx
server {
    listen 443 ssl;
    server_name chat.your-domain.com;

    ssl_certificate /path/to/your/certificate.crt;
    ssl_certificate_key /path/to/your/private.key;
    ssl_trusted_certificate /path/to/your/ca_bundle.crt;
    
    # 其他配置同上...
}
```

---

## 12. 数据管理

### 12.1 数据库连接

#### MySQL
```bash
# Docker环境
docker compose exec mysql mysql -u root -p chat_system_pro

# 本地环境
mysql -u root -p chat_system_pro
```

#### MongoDB
```bash
# Docker环境
docker compose exec mongodb mongosh -u admin -p your_password

# 本地环境
mongosh -u admin -p your_password
```

#### Redis
```bash
# Docker环境
docker compose exec redis redis-cli -a your_password

# 本地环境
redis-cli -a your_password
```

### 12.2 数据备份

#### 完整备份脚本
```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

# 备份MySQL
docker compose exec -T mysql mysqldump -u root -p$MYSQL_ROOT_PASSWORD chat_system_pro > $BACKUP_DIR/mysql_$DATE.sql

# 备份MongoDB
docker compose exec mongodb mongodump --archive=$BACKUP_DIR/mongo_$DATE.archive --gzip

# 备份Redis
docker compose exec redis redis-cli SAVE
docker cp chat-redis-pro:/data/dump.rdb $BACKUP_DIR/redis_$DATE.rdb

# 清理30天前的备份
find $BACKUP_DIR -type f -mtime +30 -delete

echo "Backup completed: $DATE"
```

### 12.3 数据恢复

#### MySQL恢复
```bash
# 从备份文件恢复
docker compose exec -T mysql mysql -u root -p$MYSQL_ROOT_PASSWORD chat_system_pro < backup_mysql_20240115.sql
```

#### MongoDB恢复
```bash
# 从archive恢复
docker compose exec -T mongodb mongorestore --archive=$BACKUP_DIR/mongo_20240115.archive --gzip

# 或从dump目录恢复
docker compose exec -T mongodb mongorestore --db chat_system_pro $BACKUP_DIR/dump/chat_system_pro
```

### 12.4 数据迁移

#### MySQL到新服务器
```bash
# 1. 导出数据
mysqldump -u root -p chat_system_pro > chat_system_pro.sql

# 2. 传输文件
scp chat_system_pro.sql user@new-server:/path/to/

# 3. 在新服务器导入
mysql -u root -p chat_system_pro < chat_system_pro.sql
```

### 12.5 清理旧数据

#### 清理30天前的消息
```bash
curl -X POST "http://localhost:8080/api/v1/admin/db/clear-old-messages?date=2024-01-15" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 归档并清理
```bash
curl -X POST "http://localhost:8080/api/v1/admin/db/archive-old?date=2024-01-01" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 13. 运维监控

### 13.1 日志配置

#### 后端日志
```yaml
# config.yaml
logging:
  level: info                  # 日志级别：debug / info / warn / error
  format: json                  # 日志格式：json / text
  output: ./logs/app.log        # 输出文件
  max_size: 100                 # 单个文件大小（MB）
  max_backups: 30              # 保留文件数
  max_age: 30                  # 保留天数
  compress: true                # 压缩旧日志
```

#### Nginx日志
```nginx
# 访问日志
access_log /var/log/nginx/chat.access.log;

# 错误日志
error_log /var/log/nginx/chat.error.log warn;
```

### 13.2 监控脚本

```bash
#!/bin/bash
# monitor.sh

# 检查服务状态
check_service() {
    local service=$1
    if docker compose ps $service | grep -q "Up"; then
        echo "✓ $service is running"
    else
        echo "✗ $service is down"
        return 1
    fi
}

# 检查端口
check_port() {
    local port=$1
    if netstat -tuln | grep -q ":$port "; then
        echo "✓ Port $port is listening"
    else
        echo "✗ Port $port is not listening"
        return 1
    fi
}

# 检查磁盘空间
check_disk() {
    local usage=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
    if [ $usage -lt 80 ]; then
        echo "✓ Disk usage: ${usage}%"
    else
        echo "⚠ Disk usage: ${usage}% (warning)"
    fi
}

# 检查内存
check_memory() {
    local total=$(free -m | awk 'NR==2 {print $2}')
    local used=$(free -m | awk 'NR==2 {print $3}')
    local percent=$((used * 100 / total))
    echo "Memory: ${used}MB / ${total}MB (${percent}%)"
}

# 运行检查
echo "=== Service Status ==="
check_service backend
check_service mysql
check_service mongodb
check_service redis

echo ""
echo "=== Port Status ==="
check_port 80
check_port 443
check_port 8080

echo ""
echo "=== System Status ==="
check_disk
check_memory

echo ""
echo "=== Recent Errors ==="
docker compose logs --tail=50 backend | grep -i error | tail -10
```

### 13.3 性能监控

#### 使用Prometheus + Grafana
```yaml
# docker-compose.monitor.yml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    restart: unless-stopped

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3001:3000"
    volumes:
      - ./grafana:/var/lib/grafana
    restart: unless-stopped
```

### 13.4 告警配置

```bash
#!/bin/bash
# alert.sh

# 发送告警邮件
send_alert() {
    local subject=$1
    local message=$2
    echo "$message" | mail -s "$subject" admin@your-domain.com
}

# 检查服务状态
check_and_alert() {
    local service=$1
    if ! docker compose ps $service | grep -q "Up"; then
        send_alert "[ALERT] $service is down" "Service $service is not running. Please check immediately."
    fi
}

# 定期检查
while true; do
    check_and_alert backend
    check_and_alert mysql
    check_and_alert mongodb
    check_and_alert redis
    sleep 300  # 每5分钟检查一次
done
```

---

## 14. 常见问题

### Q1: Docker启动失败？
```bash
# 检查Docker状态
sudo systemctl status docker

# 查看详细日志
docker compose logs

# 清理并重启
docker system prune -a
docker compose down -v
docker compose up -d
```

### Q2: 数据库连接失败？
```bash
# 检查MySQL是否就绪
docker compose exec mysql mysqladmin ping -u root -p

# 查看MySQL日志
docker compose logs mysql

# 等待数据库初始化完成（首次启动需要时间）
docker compose restart backend
```

### Q3: 端口被占用？
```bash
# 查看端口占用
lsof -i :8080
lsof -i :80
lsof -i :3306

# 或使用netstat
netstat -tulpn | grep 8080

# 杀死占用进程
sudo kill -9 <PID>

# 或修改配置使用其他端口
```

### Q4: 前端无法访问API？
```bash
# 检查API是否运行
curl http://localhost:8080/health

# 检查Nginx代理配置
sudo nginx -t
sudo systemctl reload nginx

# 查看Nginx日志
tail -f /var/log/nginx/error.log
```

### Q5: WebSocket连接失败？
```bash
# 检查WebSocket端点
curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" \
  http://localhost:8080/api/v1/ws

# 检查Nginx WebSocket配置
# 确保有以下配置：
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";
```

### Q6: 如何重置系统？
```bash
# 方式一：通过API（需要管理员Token）
curl -X POST http://localhost:8080/api/v1/admin/db/init \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"confirm":"YES"}'

# 方式二：删除数据卷重新创建
docker compose down -v
docker compose up -d
```

### Q7: 如何扩展多台服务器？
1. **负载均衡**：使用Nginx或Keepalived做负载均衡
2. **数据库主从**：MySQL配置主从复制
3. **Redis集群**：Redis配置Cluster模式
4. **MongoDB分片**：MongoDB配置Sharded Cluster
5. **WebSocket跨节点**：使用Redis Pub/Sub同步消息

### Q8: 如何优化性能？
```bash
# 1. 启用Redis缓存
# config.yaml
redis:
  enabled: true
  cache_ttl: 3600

# 2. 启用消息压缩
# config.yaml
websocket:
  compress: true

# 3. 调整数据库连接池
# config.yaml
database:
  max_open_conns: 200
  max_idle_conns: 50

# 4. 使用CDN加速静态资源
# 上传文件使用OSS/S3
storage:
  type: oss
```

### Q9: 如何迁移到新服务器？
1. 备份数据（参考12.2）
2. 在新服务器安装Docker
3. 复制项目文件
4. 恢复数据
5. 更新DNS解析
6. 测试验证

### Q10: 如何卸载系统？
```bash
# 停止服务
docker compose down -v

# 删除数据和镜像
docker system prune -a
docker volume prune

# 删除项目文件
rm -rf /path/to/chat-system-pro

# 删除日志
rm -rf /var/log/chat-system
```

---

## 15. 技术支持

### 15.1 获取帮助

- **文档**: 查看项目目录下的文档文件
- **问题反馈**: https://github.com/yourrepo/chat-system-pro/issues
- **技术支持邮箱**: support@example.com

### 15.2 社区支持

- **官方论坛**: https://forum.example.com
- **QQ群**: 123456789
- **微信群**: 请联系客服获取

### 15.3 商业支持

如需企业级技术支持和服务，请联系：
- **商务合作**: business@example.com
- **定制开发**: dev@example.com

---

**文档版本：** v1.0
**最后更新：** 2024-01-15
**联系支持：** support@example.com
