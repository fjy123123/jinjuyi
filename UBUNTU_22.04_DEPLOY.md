# 知信 (Zhixin) - Ubuntu 22.04 一键部署指南

## 📋 系统要求

### 最低配置
| 组件 | 要求 |
|------|------|
| CPU | 2 核 |
| 内存 | 4 GB |
| 存储 | 50 GB |
| 系统 | Ubuntu 22.04 LTS 或 Ubuntu 24.04 LTS |

### 推荐配置
| 组件 | 要求 |
|------|------|
| CPU | 4 核+ |
| 内存 | 8 GB+ |
| 存储 | 100 GB SSD |
| 网络 | 公网 IP |

---

## 🚀 一键部署

### 快速开始

```bash
# 1. 下载部署脚本
wget -O ubuntu-22.04-install.sh https://raw.githubusercontent.com/fjy123123/jinjuyi/main/ubuntu-22.04-install.sh

# 2. 添加执行权限
chmod +x ubuntu-22.04-install.sh

# 3. 运行部署脚本
./ubuntu-22.04-install.sh
```

部署脚本会自动完成以下操作：
- ✅ 检查系统版本和依赖
- ✅ 安装 Docker 和 Docker Compose
- ✅ 克隆项目代码
- ✅ 生成安全密钥
- ✅ 创建数据目录
- ✅ 构建并启动所有服务
- ✅ 创建管理脚本

---

## 🛠️ 管理系统

部署完成后，可以使用 `zhixin` 命令管理系统：

```bash
# 查看服务状态
zhixin status

# 查看实时日志
zhixin logs

# 查看后端日志
zhixin logs-backend

# 重启服务
zhixin restart

# 更新系统
zhixin update

# 备份数据
zhixin backup

# 清理无用资源
zhixin clean

# 查看资源使用
zhixin stats
```

---

## 📂 目录结构

```
/opt/zhixin/                    # 项目目录
├── backend/                    # 后端代码
├── web/                        # 前端代码
├── docker-compose.yml          # Docker 编排文件
├── .env                        # 环境变量
├── zhixin-manage.sh            # 管理脚本
└── credentials.txt             # 数据库凭证（请妥善保管！）

/opt/zhixin-data/               # 数据目录
├── mysql/                      # MySQL 数据
├── mongo/                      # MongoDB 数据
├── redis/                      # Redis 数据
├── uploads/                    # 上传文件
├── logs/                       # 日志
└── backups/                    # 备份
```

---

## 🔧 手动部署

如果需要手动部署，可以按照以下步骤：

### 1. 安装 Docker

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Docker
curl -fsSL https://get.docker.com | sh -

# 添加用户到 docker 组
sudo usermod -aG docker $USER
newgrp docker
```

### 2. 克隆项目

```bash
git clone https://github.com/fjy123123/jinjuyi.git /opt/zhixin
cd /opt/zhixin
```

### 3. 配置环境变量

```bash
cp .env.example .env

# 编辑 .env 文件，修改以下配置：
# - JWT_SECRET (使用 openssl rand -base64 32 生成)
# - MYSQL_ROOT_PASSWORD
# - MYSQL_PASSWORD
# - MONGO_PASSWORD
# - REDIS_PASSWORD
```

### 4. 启动服务

```bash
# 创建数据目录
sudo mkdir -p /opt/zhixin-data/{mysql,mongo,redis,uploads,logs,backups}
sudo chown -R $USER:$USER /opt/zhixin-data

# 启动服务
docker compose up -d

# 查看日志
docker compose logs -f
```

---

## 🌐 访问系统

部署完成后，通过浏览器访问：

- **前端界面**: `http://<服务器IP>`
- **API 接口**: `http://<服务器IP>:8080`
- **WebSocket**: `ws://<服务器IP>:8080/ws`

---

## 🔒 安全建议

### 1. 修改默认密码

首次登录后，请立即修改以下密码：
- 管理员账号密码
- 数据库密码（记录在 `/opt/zhixin/credentials.txt`）

### 2. 配置 HTTPS

```bash
# 安装 Certbot
sudo apt install -y certbot python3-certbot-nginx

# 获取 SSL 证书
sudo certbot --nginx -d your-domain.com
```

### 3. 配置防火墙

```bash
# 启用 UFW 防火墙
sudo ufw enable

# 允许必要端口
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS

# 查看状态
sudo ufw status
```

### 4. 定期备份

```bash
# 手动备份
zhixin backup

# 自动备份（每天凌晨 2 点）
(crontab -l 2>/dev/null; echo "0 2 * * * /opt/zhixin/zhixin-manage.sh backup") | crontab -
```

---

## 🔍 故障排查

### 服务无法启动

```bash
# 查看服务状态
docker compose ps

# 查看详细日志
docker compose logs --tail=100

# 检查端口占用
sudo ss -tuln | grep -E '(80|443|8080|3306|27017|6379)'
```

### 数据库连接失败

```bash
# 测试 MySQL 连接
docker exec zhixin-mysql mysqladmin ping -h localhost

# 测试 MongoDB 连接
docker exec zhixin-mongo mongosh --eval "db.adminCommand('ping')"

# 测试 Redis 连接
docker exec zhixin-redis redis-cli -a <redis_password> ping
```

### 重新部署

```bash
# 停止并删除所有容器
docker compose down

# 清理缓存
docker system prune -f

# 重新构建并启动
docker compose build
docker compose up -d
```

---

## 📞 获取帮助

- **项目地址**: https://github.com/fjy123123/jinjuyi
- **问题反馈**: https://github.com/fjy123123/jinjuyi/issues
- **API 文档**: [API.md](./API.md)

---

## 📝 更新日志

### v2.0.0 (2026-05-23)
- ✅ 支持 Ubuntu 22.04 和 24.04
- ✅ 完整的一键部署脚本
- ✅ 自动化管理脚本
- ✅ Docker Compose 健康检查
- ✅ 数据目录持久化
- ✅ 安全密钥自动生成
- ✅ 自动备份功能
