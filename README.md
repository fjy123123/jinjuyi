# 知信 - 完整商业化聊天系统

一个功能完整、可用于生产环境的即时聊天系统。支持前后端分离、多端（Web/Android/iOS/小程序）部署、混合数据库架构、支付系统、安全防护等功能。

## 🚀 特性列表

### 核心功能
- ✅ 用户注册/登录认证（JWT）
- ✅ 好友管理（添加、删除、备注）
- ✅ 私聊消息（文字、图片、文件）
- ✅ 群聊功能（创建、管理、成员）
- ✅ WebSocket 实时通信
- ✅ 消息撤回
- ✅ 未读消息计数
- ✅ 消息分页查询
- ✅ 在线状态显示

### 🔐 安全与加密
- ✅ **消息加密**：AES-256 + RSA 混合加密
- ✅ 端到端加密（E2E）支持
- ✅ 邀请码注册
- ✅ 验证码防暴力破解
- ✅ 请求限流
- ✅ SQL注入/XSS防护

### 混合存储架构
- ✅ MySQL 存储关系数据（用户、好友、群、订单、会话）
- ✅ MongoDB 存储消息（高性能、可扩展）
- ✅ Redis 缓存、限流、会话

### 💾 文件存储
- ✅ **多存储支持**：本地存储、阿里云OSS、AWS S3
- ✅ 图片/文件上传
- ✅ 文件URL生成与管理
- ✅ CDN加速支持

### 🔔 推送服务
- ✅ **极光推送**（JPush）
- ✅ **个推推送**（Getui）
- ✅ 多端消息推送
- ✅ 离线消息推送

### 💰 支付与积分
- ✅ Stripe 支付
- ✅ **微信支付**（WeChat Pay）
- ✅ **支付宝支付**（Alipay）
- ✅ 积分系统
- ✅ 充值/消费记录
- ✅ 订单管理

### 🎁 红包功能
- ✅ 普通红包
- ✅ 拼手气红包
- ✅ 积分支付红包
- ✅ 微信/支付宝支付红包
- ✅ 红包记录管理

### 🎨 系统配置
- ✅ 系统名称配置（后台可修改）
- ✅ 系统图标配置（后台可修改）
- ✅ Logo上传
- ✅ 主题色配置
- ✅ 多种UI模板支持

## 🔧 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: MySQL + MongoDB + Redis
- **实时通信**: WebSocket
- **加密**: AES-256 + RSA
- **部署**: Docker + Docker Compose

### 前端
- **框架**: React 18 + TypeScript
- **样式**: TailwindCSS 3
- **状态管理**: Redux
- **构建工具**: Vite

### 移动端
- **框架**: UniApp
- **支持**: Android / iOS / 微信小程序

## 📦 快速开始

### 环境要求
- Docker 20+
- Docker Compose 2+

### 一键部署

```bash
# 克隆项目
git clone https://github.com/fjy123123/jinjuyi.git
cd jinjuyi

# 一键部署
chmod +x deploy.sh
./deploy.sh
```

### 访问地址
```
前端: http://localhost
API:  http://localhost:8080
```

## 📋 目录结构

```
├── backend/              # Go后端源码
│   ├── main.go          # 主程序入口
│   ├── handlers/        # API处理器
│   ├── services/        # 业务逻辑
│   ├── models/          # 数据模型
│   ├── middleware/      # 中间件
│   ├── config/          # 配置管理
│   └── utils/           # 工具函数
├── web/                  # Web前端
├── mobile/               # 移动端（UniApp）
├── docker-compose.yml    # Docker配置
├── deploy.sh            # 部署脚本
└── README.md            # 项目说明
```

## 🔌 API接口

### 认证接口
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户退出

### 消息接口
- `POST /api/v1/messages` - 发送消息
- `GET /api/v1/messages/private/:id` - 获取私聊消息
- `GET /api/v1/messages/group/:id` - 获取群消息
- `POST /api/v1/messages/:id/recall` - 撤回消息

### 红包接口
- `POST /api/v1/redpackets` - 发送红包
- `POST /api/v1/redpackets/:id/grab` - 抢红包
- `GET /api/v1/redpackets/:id` - 获取红包详情

### 系统配置接口
- `GET /api/v1/system/config` - 获取系统配置
- `PUT /api/v1/system/config` - 更新系统配置
- `POST /api/v1/system/logo` - 上传Logo

## 📚 文档

| 文档 | 说明 |
|------|------|
| [INSTALL.md](INSTALL.md) | 详细安装文档 |
| [QUICKSTART.md](QUICKSTART.md) | 快速入门指南 |
| [API.md](API.md) | 完整API文档 |
| [DEPLOY.md](DEPLOY.md) | 部署指南 |
| [mobile/README.md](mobile/README.md) | 移动端打包指南 |

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License

---

**知信** - 让沟通更简单 🚀
