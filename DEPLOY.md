# Chat System Pro - 完整安装部署文档

## 目录

1. [环境要求](#1-环境要求)
2. [Windows 部署](#2-windows-部署)
3. [Linux 部署](#3-linux-部署)
4. [macOS 部署](#4-macos-部署)
5. [国产芯片部署](#5-国产芯片部署)
6. [Docker Compose 部署](#6-docker-compose-部署)
7. [手动编译部署](#7-手动编译部署)
8. [配置说明](#8-配置说明)
9. [域名绑定与HTTPS](#9-域名绑定与https)
10. [API 接口文档](#10-api-接口文档)
11. [常见问题](#11-常见问题)

---

## 1. 环境要求

### 硬件要求

| 规模 | CPU | 内存 | 磁盘 | 并发用户 |
|------|-----|------|------|---------|
| 测试 | 1核 | 2GB | 20GB | <100 |
| 小型 | 2核 | 4GB | 50GB | <5000 |
| 中型 | 4核 | 8GB | 100GB | <50000 |
| 大型 | 8核+ | 16GB+ | 500GB+ | >50000 |

### 软件要求

| 组件 | 版本要求 |
|------|---------|
| Docker | 20.10+ |
| Docker Compose | 2.0+ |
| Go (手动编译) | 1.21+ |
| Node.js (前端开发) | 18+ |
| MySQL | 8.0+ |
| MongoDB | 6.0+ |
| Redis | 7.0+ |

### 支持的操作系统

- **Windows**: 10/11, Windows Server 2016+
- **Linux**: Ubuntu 20.04+, CentOS 7+, Debian 10+, Fedora 36+
- **macOS**: 11+ (Intel / Apple Silicon M1/M2/M3)

### 支持的CPU架构

- **x86_64** (amd64): Intel, AMD, 海光(Hygon)
- **ARM64** (aarch64): 鲲鹏(Kunpeng), 飞腾(Phytium), Apple Silicon
- **MIPS64**: 龙芯(Loongson)

---

## 2. Windows 部署

### 2.1 安装 Docker Desktop

1. 下载 Docker Desktop: https://www.docker.com/products/docker-desktop/
2. 双击安装，重启电脑
3. 启动 Docker Desktop，等待状态变为 "Running"

### 2.2 一键部署

```cmd
:: 方式一：双击运行
deploy.bat

:: 方式二：命令行
cd chat-system-pro
deploy.bat
```

### 2.3 手动部署

```cmd
cd chat-system-pro

:: 复制配置文件
copy .env.example .env

:: 编辑 .env 文件，修改数据库密码和JWT密钥
notepad .env

:: 启动服务
docker-compose up -d

:: 查看状态
docker-compose ps

:: 查看日志
docker-compose logs -f
```

### 2.4 验证部署

浏览器访问 http://localhost ，应能看到聊天界面。

API健康检查: http://localhost:8080/health

---

## 3. Linux 部署

### 3.1 Ubuntu/Debian

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Docker
curl -fsSL https://get.docker.com | sudo sh

# 安装 Docker Compose
sudo apt install docker-compose-plugin -y

# 将当前用户加入 docker 组
sudo usermod -aG docker $USER

# 重新登录使权限生效
exit

# 验证安装
docker --version
docker compose version
```

### 3.2 CentOS/RHEL

```bash
# 安装 Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli docker-compose-plugin

# 启动 Docker
sudo systemctl start docker
sudo systemctl enable docker

# 将当前用户加入 docker 组
sudo usermod -aG docker $USER
```

### 3.3 一键部署

```bash
cd chat-system-pro

# 添加执行权限
chmod +x deploy.sh

# 运行部署脚本（交互式菜单）
./deploy.sh

# 或直接命令行
./deploy.sh start
```

### 3.4 手动部署

```bash
cd chat-system-pro

# 复制配置
cp .env.example .env

# 编辑配置
vi .env

# 启动
docker compose up -d

# 查看状态
docker compose ps

# 查看日志
docker compose logs -f backend
```

---

## 4. macOS 部署

### 4.1 安装 Docker Desktop

```bash
# 使用 Homebrew 安装
brew install --cask docker

# 或手动下载: https://www.docker.com/products/docker-desktop/
```

### 4.2 部署

```bash
cd chat-system-pro
chmod +x deploy.sh
./deploy.sh start
```

---

## 5. 国产芯片部署

### 5.1 支持的芯片

| 芯片 | 架构 | Docker平台 |
|------|------|-----------|
| 鲲鹏 Kunpeng | ARM64 | linux/arm64 |
| 飞腾 Phytium | ARM64 | linux/arm64 |
| 龙芯 Loongson | MIPS64 | linux/mips64el |
| 海光 Hygon | x86_64 | linux/amd64 |

### 5.2 自动检测部署

部署脚本会自动检测CPU类型：

```bash
chmod +x deploy.sh
./deploy.sh
# 脚本自动检测并选择合适的配置
```

### 5.3 手动指定平台

```bash
# 鲲鹏/飞腾 (ARM64)
export PLATFORM=linux/arm64
docker compose -f docker-compose.china.yml up -d

# 海光 (x86_64)
export PLATFORM=linux/amd64
docker compose -f docker-compose.china.yml up -d

# 龙芯 (MIPS64) - 需要龙芯专用Docker镜像
export PLATFORM=linux/mips64el
docker compose -f docker-compose.china.yml up -d
```

### 5.4 国产操作系统

支持以下国产操作系统：
- 麒麟 V10
- 统信 UOS
- openEuler
- 龙蜥 Anolis OS

---

## 6. Docker Compose 部署

### 6.1 标准部署

```bash
cd chat-system-pro
cp .env.example .env
# 编辑 .env
docker compose up -d
```

### 6.2 国产芯片部署

```bash
cd chat-system-pro
cp .env.example .env
docker compose -f docker-compose.china.yml up -d
```

### 6.3 服务管理

```bash
# 启动
docker compose up -d

# 停止
docker compose down

# 重启
docker compose restart

# 查看状态
docker compose ps

# 查看日志
docker compose logs -f

# 更新镜像
docker compose pull
docker compose up -d

# 重新构建
docker compose up -d --build
```

### 6.4 数据备份

```bash
# 备份 MySQL
docker compose exec mysql mysqldump -u root -p chat_system_pro > backup_mysql.sql

# 备份 MongoDB
docker compose exec mongodb mongodump --db chat_system_pro --out /data/backup

# 备份 Redis
docker compose exec redis redis-cli SAVE
docker cp chat-mongo-pro:/data/dump ./backup_mongo
```

### 6.5 数据恢复

```bash
# 恢复 MySQL
docker compose exec -T mysql mysql -u root -p chat_system_pro < backup_mysql.sql

# 恢复 MongoDB
docker compose exec mongodb mongorestore --db chat_system_pro /data/backup/chat_system_pro
```

---

## 7. 手动编译部署

### 7.1 编译后端

```bash
cd backend

# 安装依赖
go mod tidy

# 编译当前平台
go build -o chat-system-pro main.go

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o chat-system-pro-linux main.go
GOOS=windows GOARCH=amd64 go build -o chat-system-pro.exe main.go
GOOS=darwin GOARCH=arm64 go build -o chat-system-pro-mac main.go
```

### 7.2 编译前端

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

### 7.3 运行

```bash
# 后端
cd backend
./chat-system-pro

# 前端（生产模式，需要nginx）
# 将 dist/ 目录配置到 nginx
```

---

## 8. 配置说明

### 8.1 环境变量 (.env)

```env
# ========== 域名 ==========
DOMAIN=your-domain.com

# ========== MySQL ==========
MYSQL_ROOT_PASSWORD=your_secure_password
MYSQL_DATABASE=chat_system_pro
MYSQL_USER=chatuser
MYSQL_PASSWORD=your_mysql_password

# ========== Redis ==========
REDIS_PASSWORD=your_redis_password

# ========== JWT ==========
JWT_SECRET=your-super-secret-jwt-key-at-least-32-chars
JWT_EXPIRE_HOURS=720

# ========== 服务模式 ==========
GIN_MODE=release

# ========== 平台（国产芯片用） ==========
PLATFORM=linux/arm64
```

### 8.2 后端配置 (config.yaml)

```yaml
server:
  port: 8080
  mode: release        # debug / release

database:
  host: mysql          # Docker内使用服务名
  port: 3306
  user: chatuser
  password: your_mysql_password
  dbname: chat_system_pro

mongodb:
  host: mongodb
  port: 27017
  dbname: chat_system_pro

redis:
  host: redis
  port: 6379
  password: your_redis_password

jwt:
  secret: your-jwt-secret
  expire_hours: 720     # 30天

upload:
  path: ./uploads
  max_size: 10485760    # 10MB

storage:
  type: local           # local / oss / s3
```

---

## 9. 域名绑定与HTTPS

### 9.1 域名解析

将域名 A 记录指向服务器IP：
```
chat.your-domain.com  →  你的服务器IP
```

### 9.2 配置 Nginx + SSL

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取SSL证书
sudo certbot --nginx -d chat.your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

### 9.3 Nginx 配置示例

```nginx
server {
    listen 80;
    server_name chat.your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    server_name chat.your-domain.com;

    ssl_certificate /etc/letsencrypt/live/chat.your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/chat.your-domain.com/privkey.pem;

    # 前端
    location / {
        proxy_pass http://127.0.0.1:80;
    }

    # API
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # WebSocket
    location /ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

---

## 10. API 接口文档

### 基础信息

- **Base URL**: `http://your-domain.com/api/v1`
- **认证方式**: Bearer Token (JWT)
- **请求头**: `Authorization: Bearer <token>`
- **响应格式**: `{"code": 0, "msg": "success", "data": {...}}`

### 10.1 认证接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /auth/register | 用户注册 | 否 |
| POST | /auth/login | 用户登录 | 否 |

**注册请求体:**
```json
{
  "username": "testuser",
  "password": "123456",
  "nickname": "测试用户",
  "phone": "13800138000",
  "email": "test@example.com"
}
```

**登录请求体:**
```json
{
  "username": "testuser",
  "password": "123456"
}
```

**登录响应:**
```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {"id": 1, "username": "testuser", "nickname": "测试用户"}
  }
}
```

### 10.2 用户接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /users/me | 获取当前用户信息 |
| PUT | /users/profile | 更新个人资料 |
| GET | /users/settings | 获取用户设置 |
| PUT | /users/settings | 更新用户设置 |
| GET | /users/search?keyword=xxx | 搜索用户 |

**更新资料请求体:**
```json
{
  "nickname": "新昵称",
  "avatar": "https://example.com/avatar.jpg",
  "gender": 1,
  "region": "北京市",
  "sign": "这是我的个性签名"
}
```

**更新设置请求体:**
```json
{
  "new_msg_notify": true,
  "sound_notify": true,
  "add_friend_confirm": true,
  "show_online": true,
  "show_read_receipt": true,
  "theme": "dark",
  "language": "zh-CN"
}
```

### 10.3 好友接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /friends | 获取好友列表 |
| POST | /friends | 添加好友 |
| DELETE | /friends/:friend_id | 删除好友 |

**添加好友请求体:**
```json
{
  "friend_id": 2,
  "remark": "张三"
}
```

### 10.4 会话接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /conversations | 获取会话列表 |
| GET | /conversations/unread | 获取未读消息总数 |

### 10.5 消息接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /messages | 发送消息 |
| GET | /messages/private/:friend_id?page=1&page_size=20 | 获取私聊消息 |
| GET | /messages/group/:group_id?page=1&page_size=20 | 获取群消息 |
| POST | /messages/read | 标记已读 |
| POST | /messages/:message_id/recall | 撤回消息 |

**发送消息请求体:**
```json
{
  "receiver_id": 2,
  "group_id": 0,
  "content": "你好",
  "message_type": 1
}
```

**标记已读请求体:**
```json
{
  "target_id": 2,
  "type": 1
}
```

### 10.6 群组接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /groups | 创建群 |
| GET | /groups | 获取我的群列表 |
| GET | /groups/:group_id | 获取群信息 |
| PUT | /groups/:group_id | 更新群信息 |
| GET | /groups/:group_id/members | 获取群成员 |
| POST | /groups/:group_id/invite | 邀请入群 |
| DELETE | /groups/:group_id/members/:member_id | 踢出成员 |
| POST | /groups/:group_id/members/:member_id/mute | 禁言成员 |
| POST | /groups/:group_id/leave | 退出群 |

**创建群请求体:**
```json
{
  "name": "技术交流群",
  "avatar": "https://example.com/group.jpg",
  "description": "技术讨论",
  "member_ids": [2, 3, 4]
}
```

**更新群信息请求体:**
```json
{
  "name": "新群名",
  "announcement": "群公告内容",
  "join_mode": 1,
  "max_members": 500,
  "is_mute_all": false
}
```

**邀请入群请求体:**
```json
{
  "user_ids": [5, 6, 7]
}
```

### 10.7 朋友圈接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /moments | 发布朋友圈 |
| GET | /moments?page=1&page_size=20 | 获取朋友圈列表 |
| POST | /moments/:moment_id/like | 点赞 |
| DELETE | /moments/:moment_id/like | 取消点赞 |
| POST | /moments/:moment_id/comments | 评论 |
| DELETE | /moments/:moment_id | 删除朋友圈 |

**发布朋友圈请求体:**
```json
{
  "content": "今天天气真好",
  "images": ["https://example.com/1.jpg", "https://example.com/2.jpg"],
  "location": "北京市朝阳区",
  "view_scope": 0
}
```

**评论请求体:**
```json
{
  "reply_to_user": 0,
  "content": "写得真好！"
}
```

### 10.8 支付接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /payment/orders | 创建支付订单 |
| POST | /payment/orders/:order_id/pay | 支付订单 |
| GET | /payment/orders?page=1 | 获取订单列表 |
| GET | /payment/points/history?page=1 | 获取积分历史 |

**创建订单请求体:**
```json
{
  "amount": 100.00,
  "pay_type": 1,
  "order_type": 1
}
```

### 10.9 管理员接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/db/stats | 数据库统计 |
| POST | /admin/db/clear-old-messages?date=2024-01-01 | 清除旧消息 |
| POST | /admin/db/clear-all | 清空所有数据 |
| POST | /admin/db/init | 初始化数据库 |
| POST | /admin/db/archive-old?date=2024-01-01 | 归档旧消息 |
| DELETE | /admin/db/users/:user_id | 删除用户及数据 |
| DELETE | /admin/db/groups/:group_id | 删除群及数据 |
| POST | /admin/users/points | 手动调整积分 |
| GET | /admin/system/configs | 获取系统配置 |
| PUT | /admin/system/configs | 更新系统配置 |

**手动调整积分请求体:**
```json
{
  "user_id": 1,
  "points": 1000,
  "remark": "活动奖励"
}
```

**更新系统配置请求体:**
```json
{
  "key": "security.invite_code_enabled",
  "value": "true"
}
```

---

## 11. 常见问题

### Q: Docker 启动失败？
```bash
# 检查 Docker 状态
sudo systemctl status docker

# 重启 Docker
sudo systemctl restart docker

# 查看详细日志
docker compose logs
```

### Q: 数据库连接失败？
```bash
# 检查 MySQL 是否就绪
docker compose exec mysql mysqladmin ping -u root -p

# 等待数据库初始化
docker compose restart backend
```

### Q: 端口被占用？
```bash
# 查看端口占用
lsof -i :8080
lsof -i :80

# 修改 .env 或 docker-compose.yml 中的端口映射
```

### Q: 如何重置系统？
```bash
# 方式一：通过API
curl -X POST http://localhost:8080/api/v1/admin/db/init \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d "confirm=YES"

# 方式二：删除数据卷
docker compose down -v
docker compose up -d
```

### Q: 如何备份数据？
参见 [6.4 数据备份](#64-数据备份)

### Q: 如何扩展到多台服务器？
1. 使用 Nginx 做负载均衡
2. MySQL 配置主从复制
3. Redis 配置集群模式
4. MongoDB 配置分片集群
5. WebSocket 使用 Redis Pub/Sub 跨节点通信

---

## 附录

### A. 项目文件结构

```
chat-system-pro/
├── backend/                    # Go后端
│   ├── main.go                 # 程序入口
│   ├── config.yaml             # 配置文件
│   ├── go.mod                  # Go依赖
│   ├── config/config.go        # 配置加载
│   ├── handlers/               # API处理器
│   │   ├── auth.go             # 认证接口
│   │   ├── friend.go           # 好友接口
│   │   ├── message.go          # 消息接口
│   │   ├── group.go            # 群组接口
│   │   ├── moment.go           # 朋友圈接口
│   │   ├── conversation.go     # 会话接口
│   │   └── admin.go            # 管理员接口
│   ├── middleware/             # 中间件
│   │   ├── auth.go             # JWT认证
│   │   └── security.go         # 安全防护
│   ├── models/                 # 数据模型
│   │   ├── models.go           # 核心模型
│   │   └── moment.go           # 朋友圈/设备模型
│   ├── services/               # 业务逻辑
│   │   ├── database.go         # 数据库管理
│   │   ├── message.go          # 消息服务
│   │   ├── payment.go          # 支付服务
│   │   ├── storage.go          # 存储/推送
│   │   └── moment.go           # 朋友圈服务
│   └── utils/                  # 工具类
│       ├── response.go         # 统一响应
│       ├── helper.go           # 密码哈希
│       └── encryption.go       # 加密工具
├── web/                        # Web前端 (React)
├── mobile/                     # 移动端 (UniApp)
├── deploy.sh                   # Linux/Mac部署脚本
├── deploy.bat                  # Windows部署脚本
├── docker-compose.yml          # 标准Docker配置
├── docker-compose.china.yml    # 国产芯片Docker配置
├── .env.example                # 环境变量模板
├── README.md                   # 项目说明
└── DEPLOY.md                   # 本文档
```

### B. 端口说明

| 端口 | 服务 | 说明 |
|------|------|------|
| 80 | Nginx | 前端 + 反向代理 |
| 443 | Nginx | HTTPS |
| 8080 | Go Backend | API服务 |
| 3306 | MySQL | 数据库 |
| 27017 | MongoDB | 消息数据库 |
| 6379 | Redis | 缓存 |
