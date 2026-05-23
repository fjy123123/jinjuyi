# Render部署配置

## 🚀 部署方式

### 方式1：使用Render的Go模板

1. 访问 https://render.com
2. 点击 **"New +"** → **"Web Service"**
3. 连接你的GitHub仓库: `https://github.com/fjy123123/jinjuyi`
4. 配置构建设置：
   - **Root Directory**: `backend`
   - **Build Command**: `go mod tidy && go build -o server main.go`
   - **Start Command**: `./server`

### 方式2：使用Docker（推荐）

1. 创建 `render.yaml` 文件（见下方）
2. 在Render控制台选择 **"New +"** → **"Blueprint"**
3. 上传或连接包含 `render.yaml` 的仓库

## 📄 render.yaml

```yaml
# render.yaml
services:
  - type: web
    name: chat-backend
    env: golang
    repo: https://github.com/fjy123123/jinjuyi.git
    rootDir: backend
    buildCommand: go mod tidy && go build -o server main.go
    startCommand: ./server
    healthCheckPath: /health
    envVars:
      - key: GIN_MODE
        value: release
      - key: JWT_SECRET
        generateValue: true
      - key: MYSQL_HOST
        sync: false
      - key: MYSQL_PORT
        value: "3306"
      - key: MYSQL_USER
        sync: false
      - key: MYSQL_PASSWORD
        sync: false
      - key: MYSQL_DATABASE
        value: chat_system
      - key: MONGO_HOST
        sync: false
      - key: MONGO_PORT
        value: "27017"
      - key: MONGO_DB
        value: chat_system
      - key: REDIS_HOST
        sync: false
      - key: REDIS_PORT
        value: "6379"
      - key: REDIS_PASSWORD
        sync: false
```

## 🗄️ 数据库配置

Render提供免费的PostgreSQL和Redis：

### 创建PostgreSQL
1. **"New +"** → **"PostgreSQL"**
2. 创建后获取连接信息
3. 在环境变量中设置：
   - `MYSQL_HOST` → PostgreSQL的Internal Connection URL
   - `MYSQL_PORT` → 5432
   - `MYSQL_USER` → PostgreSQL用户名
   - `MYSQL_PASSWORD` → PostgreSQL密码
   - `MYSQL_DATABASE` → chat_system

### 创建Redis
1. **"New +"** → **"Redis"**
2. 创建后获取连接信息
3. 在环境变量中设置：
   - `REDIS_HOST` → Redis的Internal Connection URL
   - `REDIS_PORT` → 10000
   - `REDIS_PASSWORD` → Redis密码

### MongoDB（需要外置）
使用MongoDB Atlas免费套餐：
1. 注册 https://www.mongodb.com/atlas
2. 创建免费集群
3. 获取连接字符串
4. 设置环境变量

## 🌐 访问URL

部署成功后，Render会提供：
- **Production**: `https://chat-backend.onrender.com`
- **Staging**: `https://chat-backend-staging.onrender.com`

## ⚙️ 环境变量

在Render控制台的 **"Environment"** 标签中设置：

```bash
# 必需的环境变量
GIN_MODE=release
JWT_SECRET=your-very-long-secret-key-here-make-it-random

# MySQL配置
MYSQL_HOST=your-postgres-host
MYSQL_PORT=5432
MYSQL_USER=your-postgres-user
MYSQL_PASSWORD=your-postgres-password
MYSQL_DATABASE=chat_system

# MongoDB配置
MONGO_HOST=your-mongo-host
MONGO_PORT=27017
MONGO_DB=chat_system

# Redis配置
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
```

## 🔧 自定义域名（可选）

1. 在 **"Settings"** → **"Custom Domains"**
2. 添加你的域名
3. 按提示配置DNS记录

## 💰 费用

- **Free Tier**: 每月750小时，自动休眠
- **Starter**: $7/月，永不休眠
- **Starter Plus**: $14/月，Never Sleeps选项

## 🛠️ 故障排除

### 问题1：构建失败

检查Go版本，确保使用 `go.mod` 文件。

### 问题2：数据库连接失败

确保所有数据库服务已启动，并正确配置环境变量。

### 问题3：端口错误

确保后端监听端口与Render配置一致（通常是8080）。

## 📚 更多信息

- **Render文档**: https://render.com/docs
- **Go部署**: https://render.com/docs/go
- **PostgreSQL**: https://render.com/docs/postgres
- **Redis**: https://render.com/docs/redis