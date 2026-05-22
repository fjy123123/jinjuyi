# Chat System Pro - 完整商业化聊天系统

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

### 📱 多端支持
- ✅ **Web端**：React + Ant Design
- ✅ **小程序**：微信小程序（UniApp）
- ✅ **Android**：UniApp 打包
- ✅ **iOS**：UniApp 打包
- ✅ **多端消息同步**
- ✅ 设备管理

### 🎪 朋友圈
- ✅ 发布朋友圈（文字/图片/位置）
- ✅ 点赞/取消点赞
- ✅ 评论/回复
- ✅ 可见范围设置（所有人/仅好友/指定可见/隐藏）
- ✅ 时间线展示

### 数据库管理
- ✅ 清除指定日期前的消息
- ✅ 清空数据库（危险操作）
- ✅ 初始化数据库（重置系统）
- ✅ 删除指定用户及其所有数据
- ✅ 删除指定群及其所有数据
- ✅ 消息归档功能
- ✅ 数据库统计查询

### 支付系统
- ✅ Stripe 支付集成
- ✅ 微信支付（预留接口）
- ✅ 支付宝（预留接口）
- ✅ 积分系统
- ✅ 订单管理
- ✅ 积分历史记录
- ✅ 手动增扣积分（管理员）

### 安全防护
- ✅ 邀请码注册开关
- ✅ 验证码验证
- ✅ 请求限流（防止暴力请求）
- ✅ SQL 注入防护
- ✅ XSS 防护
- ✅ 安全头部设置

### 前端与UI
- ✅ 现代化 Web 界面（React + Ant Design）
- ✅ 多主题支持（Modern/Dark/Ocean/Purple）
- ✅ 主题切换功能
- ✅ 响应式设计
- ✅ 移动端适配（UniApp）

### 后台管理
- ✅ 系统配置管理
- ✅ 用户管理
- ✅ 数据统计
- ✅ 数据库操作

### 性能优化
- ✅ 消息异步处理
- ✅ 分页加载历史消息
- ✅ WebSocket 断线重连
- ✅ 消息推送与本地存储
- ✅ 支持水平扩展

### 商业化功能
- ✅ 支付系统
- ✅ 积分消费体系
- ✅ 会员系统（预留）
- ✅ 广告位（预留）

## 🏗️ 技术架构

```
┌───────────────────────────────────────────────────────────┐
│                       Client Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────┐  │
│  │   Web App    │  │  Mobile App  │  │  Mini Program │  │
│  │  (React)     │  │  (UniApp)    │  │              │  │
│  └──────────────┘  └──────────────┘  └───────────────┘  │
└───────────────────────────────────────────────────────────┘
                          ▼
┌───────────────────────────────────────────────────────────┐
│                    Nginx / CDN                             │
│              (Static + Reverse Proxy)                     │
└───────────────────────────────────────────────────────────┘
                          ▼
┌───────────────────────────────────────────────────────────┐
│                    Backend (Go + Gin)                     │
│  ┌───────────────────────────────────────────────────────┐│
│  │  Auth  │  Message  │  Payment  │  System  │  Group ││
│  │  Service  Service   Service    Service   Service  ││
│  └───────────────────────────────────────────────────────┘│
└───────────────────────────────────────────────────────────┘
         ▼                  ▼                  ▼
┌──────────────────┐┌─────────────────┐┌──────────────────┐
│      MySQL       ││    MongoDB      ││     Redis        │
│  (Relational Data)││ (Message Data) ││  (Cache/Limit)   │
└──────────────────┘└─────────────────┘└──────────────────┘
```

## 📦 快速开始

### 环境要求
- Docker & Docker Compose
- Node.js 18+ (前端开发)
- Go 1.21+ (后端开发)

### 方式一：Docker 一键部署（推荐）

```bash
# 1. 克隆项目
cd chat-system-pro

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env，填入你的配置

# 3. 启动所有服务
docker-compose up -d

# 4. 查看状态
docker-compose ps

# 5. 访问应用
# 前端：http://your-domain
# 后端API：http://your-domain:8080
```

### 方式二：本地开发

#### 后端
```bash
cd backend

# 1. 安装依赖
go mod tidy

# 2. 配置数据库
# 编辑 config.yaml

# 3. 运行
go run main.go
```

#### 前端
```bash
cd web

# 1. 安装依赖
npm install

# 2. 运行
npm run dev
```

#### 移动端
```bash
cd mobile

# 1. 安装依赖
npm install

# 2. 运行（需要 HBuilderX）
# 或通过命令行
npm run dev:mp-weixin
```

## 🔧 配置说明

### 环境变量配置 (.env)
```env
# 域名配置
DOMAIN=your-domain.com
API_BASE_URL=https://your-domain.com

# 数据库配置
MYSQL_ROOT_PASSWORD=your_password
MYSQL_DATABASE=chat_system_pro
MYSQL_USER=chat_user
MYSQL_PASSWORD=chat_password

# Redis配置
REDIS_PASSWORD=your_redis_password

# JWT配置
JWT_SECRET=your_super_secret_key
JWT_EXPIRE_HOURS=720

# 服务模式
GIN_MODE=release
```

### 系统配置（后台管理）

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| security.invite_code_enabled | 邀请码注册开关 | false |
| security.captcha_enabled | 验证码开关 | false |
| security.rate_limit_enabled | 限流开关 | true |
| payment.wechat_pay_enabled | 微信支付 | false |
| payment.alipay_enabled | 支付宝 | false |
| system.ui_default | 默认主题 | modern |

## 📡 API 接口文档

### 认证接口
| 接口 | 方法 | 说明 |
|------|------|------|
| /api/v1/auth/register | POST | 注册 |
| /api/v1/auth/login | POST | 登录 |
| /api/v1/auth/logout | POST | 登出 |

### 用户接口
| 接口 | 方法 | 说明 |
|------|------|------|
| /api/v1/users/me | GET | 获取当前用户信息 |
| /api/v1/users/profile | PUT | 更新个人资料 |
| /api/v1/users/search | GET | 搜索用户 |
| /api/v1/users/friends | GET | 获取好友列表 |
| /api/v1/users/friends | POST | 添加好友 |

### 消息接口
| 接口 | 方法 | 说明 |
|------|------|------|
| /api/v1/messages | POST | 发送消息 |
| /api/v1/messages/private/:id | GET | 获取私聊消息 |
| /api/v1/messages/group/:id | GET | 获取群消息 |
| /api/v1/messages/:id/recall | POST | 撤回消息 |
| /api/v1/messages/unread | GET | 获取未读数 |

### 群聊接口
| 接口 | 方法 | 说明 |
|------|------|------|
| /api/v1/groups | POST | 创建群聊 |
| /api/v1/groups | GET | 获取我的群聊 |
| /api/v1/groups/:id | GET | 获取群信息 |
| /api/v1/groups/:id/members | GET | 获取群成员 |

### 支付接口
| 接口 | 方法 | 说明 |
|------|------|------|
| /api/v1/payment/orders | POST | 创建订单 |
| /api/v1/payment/orders | GET | 获取订单列表 |
| /api/v1/payment/points/history | GET | 积分历史 |
| /api/v1/payment/orders/:id/pay | POST | 支付订单 |

### 数据库管理（管理员）
| 接口 | 方法 | 说明 |
|------|------|------|
| /api/v1/admin/db/clear-old-messages | POST | 清除旧消息 |
| /api/v1/admin/db/clear-all | POST | 清空所有数据 |
| /api/v1/admin/db/init | POST | 初始化数据库 |
| /api/v1/admin/db/users/:id/delete | POST | 删除用户数据 |
| /api/v1/admin/db/groups/:id/delete | POST | 删除群数据 |
| /api/v1/admin/db/archive-old | POST | 归档旧消息 |
| /api/v1/admin/db/stats | GET | 数据库统计 |

## 📱 WebSocket 协议

### 连接
```
ws://your-domain/ws?token=YOUR_JWT_TOKEN
```

### 消息格式
```json
{
  "type": "chat",
  "from_id": 1001,
  "to_id": 1002,
  "group_id": 0,
  "content": "Hello world",
  "timestamp": 1699999999999
}
```

### 消息类型
- chat: 聊天消息
- typing: 正在输入
- read: 已读回执
- recall: 撤回消息
- notify: 系统通知

## 🎨 UI 主题配置

### 内置主题
| 主题名称 | 说明 |
|----------|------|
| modern | 现代化蓝白主题（默认） |
| dark | 暗黑主题 |
| ocean | 海洋蓝绿主题 |
| purple | 紫色主题 |

### 自定义主题
在 `Chat.tsx` 中修改 `THEMES` 配置即可。

## 🔐 安全说明

### 生产环境建议
1. 配置 HTTPS（必须）
2. 修改 JWT_SECRET 为强密码
3. 开启 RateLimit
4. 定期备份数据库
5. 使用 WAF 防护
6. 定期更新依赖包

### 防攻击措施
- ✅ 请求限流
- ✅ SQL 注入防护
- ✅ XSS 防护
- ✅ CSRF 防护
- ✅ 密码哈希存储
- ✅ 敏感操作确认

## 📈 性能优化建议

### 单机配置（5万用户）
- CPU: 4核+
- 内存: 8GB+
- 磁盘: SSD
- MySQL: 8GB内存池
- Redis: 2GB+

### 集群配置（50万+用户）
- 3台以上应用服务器
- MySQL 主从复制
- Redis 集群
- MongoDB 分片集群
- CDN加速静态资源

## 📚 二次开发指南

### 项目结构
```
chat-system-pro/
├── backend/
│   ├── config/         # 配置
│   ├── handlers/       # API 处理器
│   ├── middleware/     # 中间件
│   ├── models/         # 数据模型
│   ├── services/       # 业务逻辑
│   ├── websocket/      # WebSocket
│   └── main.go
├── web/
│   ├── src/
│   │   ├── components/ # 组件
│   │   ├── pages/      # 页面
│   │   ├── services/   # API
│   │   └── store/      # 状态
├── mobile/
│   ├── pages/
│   ├── store/
│   └── utils/
├── docker/
└── docker-compose.yml
```

### 添加新功能
1. 在 `models/` 定义数据结构
2. 在 `services/` 实现业务逻辑
3. 在 `handlers/` 添加 API
4. 在前端 `pages/` 实现界面

## 🐛 故障排查

### 常见问题
**端口被占用**
```bash
# 修改 docker-compose.yml 的 ports
```

**数据库连接失败**
```bash
# 检查 MySQL 是否启动
docker-compose logs mysql

# 检查网络连接
docker-compose exec mysql ping -c 3 127.0.0.1
```

**WebSocket 连接断开**
```bash
# 检查 Nginx 配置
# 确认 PingInterval 设置
```

## 📞 技术支持

- 问题反馈：提交 Issue
- 技术交流：待开放讨论区
- 定制开发：联系作者

## 📄 License

MIT License - 可商用
