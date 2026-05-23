# 知信 (Zhixin) - 部署文档

## 📋 目录

- [系统要求](#系统要求)
- [快速部署](#快速部署)
- [手动部署](#手动部署)
- [集群部署](#集群部署)
- [配置说明](#配置说明)
- [API 文档](#api-文档)
- [运维管理](#运维管理)
- [故障排查](#故障排查)

---

## 🖥️ 系统要求

### 最低配置
| 组件 | 要求 |
|------|------|
| CPU | 2 核 |
| 内存 | 4 GB |
| 存储 | 50 GB |
| 系统 | Ubuntu 24.04 LTS |

### 推荐配置
| 组件 | 要求 |
|------|------|
| CPU | 4 核+ |
| 内存 | 8 GB+ |
| 存储 | 100 GB SSD |
| 系统 | Ubuntu 24.04 LTS |

### 软件依赖
- Docker 20.10+
- Docker Compose 2.0+
- Nginx (可选，用于反向代理)

---

## 🚀 快速部署

### 方式一：一键安装脚本（推荐）

```bash
# 下载并运行一键安装脚本
wget -O zhixin-install.sh https://raw.githubusercontent.com/your-org/zhixin/main/deploy/ubuntu-24.04-install.sh
chmod +x zhixin-install.sh
./zhixin-install.sh
```

脚本会自动完成以下操作：
1. 检查系统版本和依赖
2. 安装 Docker 和 Docker Compose
3. 克隆项目代码
4. 生成安全密钥
5. 启动所有服务
6. 创建管理脚本

### 方式二：Docker Compose

```bash
# 1. 克隆项目
git clone https://github.com/your-org/zhixin.git
cd zhixin

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，修改数据库密码等配置

# 3. 启动服务
docker-compose up -d

# 4. 查看服务状态
docker-compose ps

# 5. 查看日志
docker-compose logs -f
```

---

## 🔧 手动部署

### 1. 安装基础软件

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Docker
curl -fsSL https://get.docker.com | sh -
sudo usermod -aG docker $USER

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 重启使 Docker 权限生效
newgrp docker
```

### 2. 编译后端

```bash
cd backend

# 下载依赖
go mod download

# 编译
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o zhixin-backend .

# 复制二进制文件
cp zhixin-backend ../docker/backend/
```

### 3. 构建前端

```bash
cd web

# 安装依赖
npm install

# 构建生产版本
npm run build
```

### 4. 启动服务

```bash
cd deploy
docker-compose up -d
```

---

## 🌐 集群部署

集群部署支持水平扩展，适用于高并发场景。

### 架构说明

```
                         ┌─────────────────┐
                         │   用户浏览器     │
                         └────────┬────────┘
                                  │
                                  ▼
                    ┌─────────────────────────┐
                    │    Nginx 负载均衡器      │
                    │   (ip_hash 会话保持)     │
                    └────────────┬────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
         ▼                       ▼                       ▼
   ┌───────────┐          ┌───────────┐          ┌───────────┐
   │ Backend 1 │          │ Backend 2 │          │ Backend 3 │
   │  (Node A) │          │  (Node B) │          │  (Node C) │
   └─────┬─────┘          └─────┬─────┘          └─────┬─────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
                    ▼                         ▼
            ┌───────────────┐         ┌───────────────┐
            │    MongoDB    │         │     MySQL     │
            │  (消息存储)    │         │  (用户/订单)  │
            └───────────────┘         └───────────────┘
                    │
                    ▼
            ┌───────────────┐
            │     Redis     │
            │ (缓存+Pub/Sub) │
            └───────────────┘
```

### 启动集群

```bash
cd deploy

# 使用集群配置启动
docker-compose -f docker-compose-cluster.yml up -d

# 查看服务状态
docker-compose -f docker-compose-cluster.yml ps

# 扩展到 5 个后端实例
docker-compose -f docker-compose-cluster.yml up -d --scale backend=5
```

详细集群部署请参考: [集群部署指南](./CLUSTER_DEPLOYMENT.md)

---

## ⚙️ 配置说明

### 环境变量 (.env)

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码 | 自动生成 |
| `MYSQL_DATABASE` | MySQL 数据库名 | zhixin_chat |
| `MYSQL_USER` | MySQL 应用用户名 | zhixin_user |
| `MYSQL_PASSWORD` | MySQL 应用用户密码 | 自动生成 |
| `MONGO_USER` | MongoDB 用户名 | mongo_admin |
| `MONGO_PASSWORD` | MongoDB 密码 | 自动生成 |
| `MONGO_DB` | MongoDB 数据库名 | zhixin_chat |
| `REDIS_PASSWORD` | Redis 密码 | 自动生成 |
| `JWT_SECRET` | JWT 密钥 (必须修改) | 自动生成 |
| `GIN_MODE` | 运行模式 | release |
| `SERVER_PORT` | 服务端口 | 8080 |
| `NODE_ID` | 集群节点 ID (可选) | - |

### 安全建议

1. **必须修改**: 部署后立即修改 `.env` 中的所有密码和密钥
2. **JWT 密钥**: 使用 `openssl rand -base64 32` 生成强随机密钥
3. **HTTPS**: 生产环境必须启用 HTTPS
4. **防火墙**: 仅开放必要端口 (80, 443)

---

## 📡 API 文档

### 基础信息

| 项目 | 值 |
|------|-----|
| 基础 URL | `http://localhost:8080/api/v1` |
| 认证方式 | Bearer Token (JWT) |
| 数据格式 | JSON |
| WebSocket | `ws://localhost:8080/ws?token=<jwt_token>` |

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/refresh` | 刷新 Token |
| POST | `/api/v1/auth/logout` | 退出登录 |

### 用户接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/users/me` | 获取个人信息 | 已登录 |
| PUT | `/api/v1/users/profile` | 更新个人信息 | 已登录 |
| GET | `/api/v1/users/settings` | 获取用户设置 | 已登录 |
| PUT | `/api/v1/users/settings` | 更新用户设置 | 已登录 |
| GET | `/api/v1/users/search` | 搜索用户 | 已登录 |

### 消息接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/messages/private/:id` | 获取私聊消息 | 已登录 |
| GET | `/api/v1/messages/group/:id` | 获取群聊消息 | 已登录 |
| POST | `/api/v1/messages` | 发送消息 | 已登录 |
| POST | `/api/v1/messages/:id/recall` | 撤回消息 | 已登录 |
| POST | `/api/v1/messages/read` | 标记已读 | 已登录 |
| GET | `/api/v1/messages/export` | 导出消息 | 已登录 |

### 群组接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/groups` | 获取群列表 | 已登录 |
| POST | `/api/v1/groups` | 创建群组 | 已登录 |
| GET | `/api/v1/groups/:id` | 获取群信息 | 已登录 |
| PUT | `/api/v1/groups/:id` | 更新群信息 | 群主/管理员 |
| GET | `/api/v1/groups/:id/members` | 获取群成员 | 已登录 |
| POST | `/api/v1/groups/:id/invite` | 邀请成员 | 已登录 |
| DELETE | `/api/v1/groups/:id/members/:uid` | 移除成员 | 群主/管理员 |

### 红包接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/v1/redpackets` | 发送红包 | 已登录 |
| POST | `/api/v1/redpackets/:id/grab` | 抢红包 | 已登录 |
| GET | `/api/v1/redpackets/:id` | 获取红包详情 | 已登录 |
| GET | `/api/v1/redpackets/sent` | 获取发送记录 | 已登录 |
| GET | `/api/v1/redpackets/received` | 获取领取记录 | 已登录 |

### 管理员接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/admin/db/stats` | 数据库统计 | 管理员 |
| POST | `/api/v1/admin/db/clear-old-messages` | 清理旧消息 | 管理员 |
| DELETE | `/api/v1/admin/db/users/:id` | 删除用户 | 管理员 |
| DELETE | `/api/v1/admin/db/groups/:id` | 删除群组 | 管理员 |
| POST | `/api/v1/admin/users/points` | 调整积分 | 管理员 |
| GET | `/api/v1/admin/recharge` | 审核充值列表 | 管理员 |
| PUT | `/api/v1/admin/recharge/:id/approve` | 通过充值 | 管理员 |
| PUT | `/api/v1/admin/recharge/:id/reject` | 拒绝充值 | 管理员 |
| GET | `/api/v1/admin/withdraw` | 审核提现列表 | 管理员 |
| PUT | `/api/v1/admin/withdraw/:id/approve` | 通过提现 | 管理员 |
| PUT | `/api/v1/admin/withdraw/:id/reject` | 拒绝提现 | 管理员 |
| PUT | `/api/v1/admin/system/config` | 更新系统配置 | 超级管理员 |
| POST | `/api/v1/admin/system/maintenance` | 维护模式开关 | 超级管理员 |

完整 API 文档请查看: [API.md](./API.md)

---

## 🛠️ 运维管理

### 管理脚本

一键安装后会创建 `/usr/local/bin/zhixin-manage.sh` 管理脚本：

```bash
# 启动服务
zhixin-manage.sh start

# 停止服务
zhixin-manage.sh stop

# 重启服务
zhixin-manage.sh restart

# 查看日志
zhixin-manage.sh logs

# 查看后端日志
zhixin-manage.sh logs-backend

# 查看服务状态
zhixin-manage.sh status

# 更新系统
zhixin-manage.sh update

# 备份数据
zhixin-manage.sh backup
```

### 常用命令

```bash
# 查看服务状态
docker-compose ps

# 查看实时日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f web
docker-compose logs -f mysql
docker-compose logs -f mongo

# 重启特定服务
docker-compose restart backend
docker-compose restart web

# 进入容器
docker-compose exec backend sh
docker-compose exec mysql mysql -u zhixin_user -p

# 清理无用镜像
docker system prune -a

# 备份数据库
docker-compose exec mysql mysqldump -u zhixin_user -p zhixin_chat > backup.sql
```

### 日志管理

```bash
# 日志位置
/opt/zhixin-data/logs/

# 清理日志（保留最近 7 天）
find /opt/zhixin-data/logs/ -name "*.log" -mtime +7 -delete
```

---

## 🔍 故障排查

### 服务无法启动

```bash
# 检查端口占用
sudo netstat -tuln | grep -E '(80|443|8080|3306|27017|6379)'

# 检查 Docker 状态
docker ps -a
docker-compose ps

# 查看详细日志
docker-compose logs --tail=100
```

### 数据库连接失败

```bash
# 测试 MySQL 连接
docker-compose exec mysql mysqladmin ping -h localhost

# 测试 MongoDB 连接
docker-compose exec mongo mongosh --eval "db.runCommand({ping:1})"

# 测试 Redis 连接
docker-compose exec redis redis-cli ping
```

### 前端无法访问

```bash
# 检查 Nginx 配置
docker-compose exec nginx nginx -t

# 重载 Nginx 配置
docker-compose exec nginx nginx -s reload

# 检查前端构建
docker-compose exec web ls -la /usr/share/nginx/html/
```

### WebSocket 连接失败

```bash
# 测试 WebSocket 连接
wscat -c "ws://localhost:8080/ws?token=<your_token>"

# 检查 Nginx WebSocket 配置
grep -A5 "proxy_set_header Upgrade" docker/nginx.conf
```

### 常见问题

| 问题 | 解决方案 |
|------|----------|
| 端口冲突 | 修改 `docker-compose.yml` 中的端口映射 |
| 权限不足 | 确保 `/opt/zhixin-data` 目录所有者正确 |
| 内存不足 | 检查容器内存使用: `docker stats` |
| 磁盘空间不足 | 清理 Docker 无用镜像和卷 |

---

## 📞 获取帮助

- **文档**: [README.md](./README.md)
- **API 文档**: [API.md](./API.md)
- **集群部署**: [CLUSTER_DEPLOYMENT.md](./CLUSTER_DEPLOYMENT.md)
- **GitHub Issues**: https://github.com/your-org/zhixin/issues

---

## 📝 更新日志

### v1.0.0 (2026-05-23)
- ✅ 完整的一键部署脚本
- ✅ Docker Compose 编排
- ✅ 集群部署支持
- ✅ Nginx 负载均衡配置
- ✅ 健康检查
- ✅ 安全加固
- ✅ 完整的 API 实现
- ✅ 自动化管理脚本
