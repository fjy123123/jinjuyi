# Chat System Pro - 快速入门指南

## 🚀 5分钟快速部署

### 前提条件
- Docker Desktop 已安装并运行
- 2核CPU, 4GB内存可用
- 网络连接正常

### 步骤 1: 解压项目
```bash
tar -xzf chat-system-pro-v1.0.0-20240115.tar.gz
cd chat-system-pro-v1.0.0-20240115
```

### 步骤 2: 配置环境
```bash
# 复制环境变量文件
cp .env.example .env

# 编辑配置（必填项）
vi .env
```
**必填配置项：**
```env
MYSQL_ROOT_PASSWORD=your_secure_password  # MySQL密码
JWT_SECRET=your-super-secret-key-32chars  # JWT密钥（至少32字符）
```

### 步骤 3: 启动服务
```bash
# Windows
deploy.bat

# Linux/macOS
chmod +x deploy.sh
./deploy.sh
```

### 步骤 4: 验证部署
```bash
# 检查服务状态
docker compose ps

# 测试API
curl http://localhost:8080/health

# 浏览器访问
http://localhost
```

**看到 `{"status":"ok"}` 表示成功！**

---

## 📱 使用系统

### 默认账号
- **管理员**: admin / admin123
- **测试用户**: test / test123

### 访问地址
- **前端界面**: http://localhost
- **API接口**: http://localhost:8080
- **健康检查**: http://localhost:8080/health

---

## 🔧 常用操作

### 查看日志
```bash
# 查看所有日志
docker compose logs -f

# 查看后端日志
docker compose logs -f backend

# 查看数据库日志
docker compose logs -f mysql
```

### 重启服务
```bash
# 重启所有服务
docker compose restart

# 重启特定服务
docker compose restart backend
```

### 停止服务
```bash
# 停止服务（保留数据）
docker compose stop

# 完全停止并删除容器
docker compose down
```

### 更新系统
```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker compose up -d --build

# 清理旧镜像
docker image prune -f
```

---

## 🛠️ 常见问题

### Q1: Docker启动失败
```bash
# 检查Docker状态
docker --version
docker compose version

# 启动Docker
sudo systemctl start docker  # Linux
open -a Docker               # macOS
```

### Q2: 端口被占用
```bash
# 查看端口占用
lsof -i :8080

# 修改端口（编辑.env）
SERVER_PORT=8081

# 或修改docker-compose.yml中的端口映射
```

### Q3: 数据库连接失败
```bash
# 等待MySQL初始化（约30秒）
docker compose logs mysql

# 重启后端服务
docker compose restart backend
```

### Q4: 前端无法访问
```bash
# 检查Nginx日志
docker compose logs frontend

# 重启前端
docker compose restart frontend
```

---

## 📚 更多文档

- **[详细安装文档](./INSTALL.md)** - 完整的安装部署说明
- **[API接口文档](./API.md)** - 所有API接口详细说明
- **[功能说明](./FEATURES.md)** - 系统功能详细介绍
- **[部署文档](./DEPLOY.md)** - 各种部署场景说明

---

## 🎯 功能尝鲜

### 发送消息
```bash
# 登录获取Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# 发送消息
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"receiver_id":2,"content":"Hello!","message_type":1}'
```

### 发送红包
```bash
# 发送红包
curl -X POST http://localhost:8080/api/v1/redpackets \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": 2,
    "pay_type": 1,
    "amount": 100,
    "total_count": 10,
    "receiver_id": 2,
    "greeting": "恭喜发财"
  }'
```

### WebSocket测试
```javascript
// 浏览器控制台执行
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?token=YOUR_TOKEN');

ws.onopen = () => {
  console.log('已连接');
  
  // 发送消息
  ws.send(JSON.stringify({
    type: 'message',
    data: {
      receiver_id: 2,
      content: 'WebSocket测试',
      message_type: 1
    }
  }));
};

ws.onmessage = (event) => {
  console.log('收到消息:', JSON.parse(event.data));
};
```

---

## 🌟 下一步

1. **配置HTTPS** - 参考 INSTALL.md 第十一章
2. **配置微信支付** - 参考 INSTALL.md 第十章
3. **扩展多台服务器** - 参考 INSTALL.md 第十四章
4. **自定义UI模板** - 修改 web/src 目录

---

## 📞 获取帮助

遇到问题？
- 📖 查看详细文档
- 💬 加入社区讨论
- 📧 联系技术支持

---

**5分钟快速开始，就从这里开始！**
