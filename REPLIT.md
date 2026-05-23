# 知信聊天系统 - Replit 部署指南

## 🎯 快速开始

### 方法1：从GitHub导入（推荐）

1. 访问 https://replit.com
2. 点击 **"Create Repl"**
3. 选择 **"Import from GitHub"**
4. 粘贴仓库地址: `https://github.com/fjy123123/jinjuyi`
5. 选择语言: **Go**
6. 点击 **"Import from GitHub"**

### 方法2：手动创建

1. 创建一个新的 Go Repl
2. 克隆仓库：
```bash
git clone https://github.com/fjy123123/jinjuyi.git
cd jinjuyi
```
3. 修改 `backend/main.go` 中的数据库连接为内嵌SQLite

## ⚙️ Replit配置

### 必需的文件

项目已包含以下Replit配置文件：
- `replit.nix` - Nix环境配置
- `replit.toml` - Replit运行配置
- `run.sh` - 启动脚本

### 环境变量

在Replit的Secrets中设置：

```bash
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=replit_password
MYSQL_DATABASE=chat_system

JWT_SECRET=your_super_secret_key_change_this
GIN_MODE=release

REDIS_HOST=localhost
REDIS_PORT=6379

MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB=chat_system
```

## 🔧 数据库配置

### 选项1：使用Replit内置PostgreSQL（推荐）

1. 在Replit面板中点击 **"+ Add Database"**
2. 选择 **"PostgreSQL"**
3. 获取连接信息并更新环境变量

### 选项2：使用外置数据库服务

推荐使用免费的云数据库：
- **MongoDB Atlas**: https://www.mongodb.com/atlas
- **Redis Cloud**: https://redis.com/cloud
- **Neon PostgreSQL**: https://neon.tech

## 🚀 运行项目

### 自动运行

Replit会自动执行 `run.sh` 脚本启动服务。

### 手动运行

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 启动服务
go run main.go
```

## 🌐 访问服务

启动后，Replit会提供两个访问地址：

1. **主URL**: 直接访问应用（端口8080）
2. **WebSocket**: 用于实时聊天

## 📱 API端点

基础URL: `https://your-replit-url.repl.co`

| 功能 | 端点 | 方法 |
|------|------|------|
| 注册 | `/api/v1/auth/register` | POST |
| 登录 | `/api/v1/auth/login` | POST |
| 获取用户信息 | `/api/v1/auth/me` | GET |
| 发送消息 | `/api/v1/message/send` | POST |
| 获取会话列表 | `/api/v1/conversation/list` | GET |

## 🛠️ 故障排除

### 问题1：端口被占用

修改 `backend/config/config.go` 中的端口配置。

### 问题2：数据库连接失败

检查环境变量是否正确设置。

### 问题3：依赖下载失败

```bash
cd backend
go clean -modcache
go mod download
```

## 💡 提示

- Replit免费版有CPU时间限制，不用时暂停Repl
- 建议使用外置数据库服务以获得更好的性能
- 定期保存数据，因为免费实例可能会被重置

## 📚 更多信息

- **GitHub仓库**: https://github.com/fjy123123/jinjuyi
- **后端文档**: 访问 `/api/v1` 路径查看Swagger文档