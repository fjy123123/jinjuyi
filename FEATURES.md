# 完整功能文档

## 🚀 系统特性总览

### 跨平台部署支持

#### Windows
- [deploy.bat](file:///workspace/chat-system-pro/deploy.bat) - Windows 一键部署脚本
- 自动检测 Docker Desktop
- 图形化菜单界面
- 支持 Windows 10/11, Windows Server 2016+

#### Linux / MacOS
- [deploy.sh](file:///workspace/chat-system-pro/deploy.sh) - Linux/MacOS 部署脚本
- 支持 Debian/Ubuntu/CentOS/Fedora
- 支持 MacOS (Intel/M1/M2/M3)
- 自动检测系统和架构
- 国产芯片自动识别

### 国产芯片支持

#### 支持的芯片
- ✅ **鲲鹏 (Kunpeng)** - ARM64
- ✅ **飞腾 (Phytium)** - ARM64
- ✅ **龙芯 (Loongson)** - MIPS64EL
- ✅ **海光 (Hygon)** - x86_64

#### 部署配置
- [docker-compose.china.yml](file:///workspace/chat-system-pro/docker-compose.china.yml) - 国产芯片专用配置
- [backend/Dockerfile.china](file:///workspace/chat-system-pro/backend/Dockerfile.china) - 国产芯片专用镜像
- 使用国内 Go 模块代理加速构建
- 自适应架构选择

### 个人功能扩展

#### 用户信息
- 用户名/密码/昵称
- 头像/手机号/邮箱
- 性别 (男/女/未知)
- 生日/地区
- 个性签名
- VIP 等级
- 最后登录时间/IP
- 在线状态

#### 个人设置
- 新消息通知开关
- 声音提醒开关
- 加好友验证开关
- 显示在线状态
- 显示已读回执
- 主题选择
- 语言设置

### 群组功能扩展

#### 群组基础
- 群名称/头像/描述
- 群公告
- 群主/管理员/成员三级权限
- 成员数统计
- 成员上限设置

#### 群组管理
- **入群模式**：自由加入/需要验证/禁止加入
- **全员禁言**
- **成员禁言/解除禁言**
- **设置管理员**
- **群昵称设置**
- **邀请好友加入**
- **踢出成员**

#### 群组权限
- 允许/禁止邀请成员
- 允许/禁止查看成员
- 群主转让
- 解散群聊

### 消息功能

#### 消息类型
- 文本消息
- 图片消息
- 文件消息
- 语音消息
- 视频消息
- 撤回消息

#### 消息加密
- AES-256 加密
- RSA-2048/4096 密钥交换
- 混合加密方案
- 端到端加密支持

### 文件存储

#### 存储类型
- 本地存储
- 阿里云 OSS
- AWS S3
- 通用接口，易于扩展

### 推送服务

#### 推送厂商
- 极光推送 (JPush)
- 个推 (Getui)

#### 推送功能
- 单用户推送
- 全量推送
- 离线推送
- 多端同步

### 支付功能

#### 支付方式
- Stripe (国际)
- 微信支付 (国内)
- 支付宝 (国内)

#### 积分系统
- 充值得积分
- 积分消费
- 积分历史记录
- 订单管理

### 朋友圈

#### 发布
- 文字
- 图片 (多张)
- 位置
- 可见范围设置 (所有人/仅好友/指定可见/隐藏)

#### 互动
- 点赞/取消点赞
- 评论/回复评论
- 删除自己的动态

### 数据库管理

#### 管理员功能
- 清除指定日期前的消息
- 清空所有数据
- 初始化数据库
- 删除用户及其所有数据
- 删除群及其所有数据
- 归档旧消息
- 查看数据库统计

### 安全防护

#### 安全机制
- 邀请码注册
- 验证码验证
- 请求限流 (Rate Limit)
- SQL 注入防护
- XSS 防护
- JWT 认证
- 密码加密存储 (bcrypt)

## 📋 部署指南

### 快速启动 (任意系统)

#### Windows
```bash
# 双击运行 deploy.bat
# 或在命令行执行
deploy.bat
```

#### Linux / MacOS
```bash
# 添加执行权限
chmod +x deploy.sh

# 运行
./deploy.sh

# 或直接使用命令
./deploy.sh start
```

### 国产芯片部署

#### 自动检测部署
```bash
# 脚本会自动检测国产芯片
./deploy.sh

# 或直接使用国产配置
export PLATFORM=linux/arm64
docker-compose -f docker-compose.china.yml up -d
```

#### 手动指定
```bash
# 鲲鹏/飞腾 (ARM64)
export PLATFORM=linux/arm64
docker-compose -f docker-compose.china.yml up -d

# 龙芯 (MIPS64)
export PLATFORM=linux/mips64el
docker-compose -f docker-compose.china.yml up -d

# 海光 (x86_64)
export PLATFORM=linux/amd64
docker-compose -f docker-compose.china.yml up -d
```

## 📁 项目文件结构

```
chat-system-pro/
├── backend/
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   │   ├── models.go       # 用户/好友/群/消息等
│   │   └── moment.go       # 朋友圈/小程序/设备
│   ├── services/
│   │   ├── database.go     # 数据库管理
│   │   ├── message.go      # 消息服务
│   │   ├── payment.go      # 支付服务
│   │   ├── storage.go      # 文件存储/推送
│   │   └── moment.go       # 朋友圈/小程序
│   ├── utils/
│   │   └── encryption.go   # 加密工具
│   ├── config.yaml         # 配置文件
│   ├── go.mod
│   ├── main.go
│   └── Dockerfile.china    # 国产芯片构建
├── web/                    # Web前端 (React)
├── mobile/                 # 移动端 (UniApp)
├── docker/                 # Docker相关
├── deploy.sh               # Linux/Mac部署
├── deploy.bat              # Windows部署
├── docker-compose.yml      # 标准配置
├── docker-compose.china.yml # 国产芯片配置
└── README.md
```

## 🔧 系统配置

### 环境变量 (.env)
```env
# 系统
DOMAIN=your-domain.com
GIN_MODE=release

# 数据库
MYSQL_ROOT_PASSWORD=
MYSQL_DATABASE=chat_system_pro
MYSQL_USER=chatuser
MYSQL_PASSWORD=

# Redis
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key

# 平台 (可选)
PLATFORM=linux/arm64
```

### 系统配置 (config.yaml)
- 数据库
- MongoDB
- Redis
- 存储配置 (OSS/S3)
- 推送配置
- 支付配置
- 安全配置

## 📱 API 接口

### 用户相关
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/login` - 登录
- `GET /api/v1/users/me` - 获取当前用户
- `PUT /api/v1/users/profile` - 更新个人信息
- `PUT /api/v1/users/settings` - 更新设置

### 消息相关
- `POST /api/v1/messages` - 发送消息
- `GET /api/v1/messages/private/:id` - 获取私聊消息
- `GET /api/v1/messages/group/:id` - 获取群消息
- `POST /api/v1/messages/:id/recall` - 撤回消息

### 群组相关
- `POST /api/v1/groups` - 创建群
- `PUT /api/v1/groups/:id` - 更新群信息
- `POST /api/v1/groups/:id/announcement` - 设置公告
- `POST /api/v1/groups/:id/mute` - 全员禁言
- `POST /api/v1/groups/:id/members/:uid/mute` - 禁言成员
- `POST /api/v1/groups/:id/admin` - 设置管理员
- `POST /api/v1/groups/:id/invite` - 邀请成员
- `POST /api/v1/groups/:id/join` - 申请加入

### 朋友圈
- `POST /api/v1/moments` - 发布朋友圈
- `GET /api/v1/moments` - 获取朋友圈
- `POST /api/v1/moments/:id/like` - 点赞
- `DELETE /api/v1/moments/:id/like` - 取消点赞
- `POST /api/v1/moments/:id/comments` - 评论

### 支付
- `POST /api/v1/payment/orders` - 创建订单
- `POST /api/v1/payment/orders/:id/pay` - 支付
- `GET /api/v1/payment/points/history` - 积分历史

### 管理员
- `POST /api/v1/admin/db/clear-old` - 清除旧消息
- `POST /api/v1/admin/db/clear-all` - 清空数据
- `POST /api/v1/admin/db/init` - 初始化数据库
- `DELETE /api/v1/admin/users/:id` - 删除用户
- `DELETE /api/v1/admin/groups/:id` - 删除群
- `GET /api/v1/admin/db/stats` - 数据库统计

## 📚 二次开发

### 添加新的存储后端
实现 Storage 接口即可支持更多存储服务。

### 添加新的推送渠道
实现 PushProvider 接口即可支持更多推送服务。

### 添加新的支付方式
实现 PaymentProvider 接口即可支持更多支付方式。

## 🎯 下一步建议

1. 配置 HTTPS 和域名
2. 配置对象存储 (OSS/S3)
3. 配置推送服务
4. 配置支付方式
5. 配置小程序
6. 性能测试和优化
7. 安全审计和加固

---

## 💡 联系与支持

- 问题反馈：提交 Issue
- 商业合作：联系作者
- 技术支持：文档查询

---

**祝你使用愉快！** 🎉
