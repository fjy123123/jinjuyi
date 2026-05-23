# 知信聊天系统 - 集群部署指南

## 📋 概述

本指南介绍如何将知信聊天系统部署为高可用集群架构，支持水平扩展以应对大并发场景。

## 🏗️ 架构设计

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

## 🔑 核心技术

### 1. Redis Pub/Sub 消息同步

各后端实例通过 Redis 发布/订阅机制同步消息：

- **ChannelNewMessage**: 新消息广播
- **ChannelRecall**: 消息撤回
- **ChannelReadReceipt**: 已读回执
- **ChannelSystemNotify**: 系统通知
- **ChannelCall**: 音视频通话信令

### 2. Nginx ip_hash 会话保持

WebSocket 连接使用 `ip_hash` 算法，确保同一用户会话路由到同一后端实例。

### 3. 节点 ID 标识

每个后端实例通过 `NODE_ID` 环境变量标识，避免消息循环。

## 🚀 部署方式

### 方式一：Docker Compose 集群（推荐）

```bash
# 进入部署目录
cd deploy

# 启动集群
docker-compose -f docker-compose-cluster.yml up -d

# 查看服务状态
docker-compose -f docker-compose-cluster.yml ps

# 扩展后端实例（增加到 5 个）
docker-compose -f docker-compose-cluster.yml up -d --scale backend_1=2 --scale backend_2=2 --scale backend_3=1
```

### 方式二：systemd 多实例部署

```bash
# 1. 安装二进制文件
sudo cp zhixin-backend /opt/zhixin/backend/
sudo mkdir -p /opt/zhixin/backend/logs

# 2. 创建环境变量文件
cat > /opt/zhixin/backend/.env.cluster << EOF
MONGODB_URI=mongodb://localhost:27017/zhixin
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=zhixin
MYSQL_PASSWORD=your-password
MYSQL_DATABASE=zhixin
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-super-secret-jwt-key
ADMIN_SECRET=admin-secret-key
EOF

# 3. 安装 systemd 服务
sudo cp zhixin-backend@.service /etc/systemd/system/
sudo cp zhixin-nginx.service /etc/systemd/system/

# 4. 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable zhixin-backend@{1,2,3}
sudo systemctl start zhixin-backend@{1,2,3}
sudo systemctl enable zhixin-nginx
sudo systemctl start zhixin-nginx

# 5. 检查服务状态
sudo systemctl status zhixin-backend@{1,2,3}
sudo systemctl status zhixin-nginx
```

### 方式三：多服务器部署

```bash
# 服务器规划
# - 192.168.1.10: Nginx + Backend 1
# - 192.168.1.11: Backend 2 + Backend 3
# - 192.168.1.20: MongoDB + MySQL + Redis

# 修改 nginx-cluster.conf 中的 upstream 地址
upstream backend_api {
    ip_hash;
    server 192.168.1.10:8080;
    server 192.168.1.11:8080;
    server 192.168.1.11:8081;
}

upstream backend_ws {
    ip_hash;
    server 192.168.1.10:8080;
    server 192.168.1.11:8080;
    server 192.168.1.11:8081;
}
```

## ⚙️ Nginx 配置说明

关键配置项：

```nginx
# WebSocket 支持
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";

# ip_hash 会话保持
upstream backend_ws {
    ip_hash;
    server backend_1:8080;
    server backend_2:8080;
    server backend_3:8080;
}

# 健康检查
server backend_1:8080 max_fails=3 fail_timeout=30s;
```

## 📊 监控与运维

### 健康检查

```bash
# 检查所有后端实例
curl http://localhost:8080/health  # 实例1
curl http://localhost:8081/health  # 实例2
curl http://localhost:8082/health  # 实例3

# 通过 Nginx
curl http://localhost/health
```

### 日志查看

```bash
# Docker 日志
docker-compose -f docker-compose-cluster.yml logs -f backend_1

# systemd 日志
journalctl -u zhixin-backend@1 -f
journalctl -u zhixin-nginx -f
```

### 扩展操作

```bash
# 添加新实例
# 1. 修改 nginx-cluster.conf 添加新 upstream
# 2. 重载 Nginx: nginx -s reload
# 3. 启动新实例

# 缩减实例
# 1. 停止实例
# 2. 修改 nginx-cluster.conf 移除 upstream
# 3. 重载 Nginx
```

## 🔒 安全建议

1. **网络隔离**: 使用私有网络，Nginx 只暴露公网 IP
2. **Redis 安全**: 配置密码认证，禁用危险命令
3. **MongoDB 安全**: 启用认证，创建专用应用用户
4. **MySQL 安全**: 遵循最小权限原则
5. **TLS 加密**: 生产环境务必启用 HTTPS

## 📈 性能调优

### 推荐的实例数量

| 用户规模 | 后端实例 | 内存需求 |
|---------|---------|---------|
| < 1,000 | 2 | 2 GB |
| 1,000 - 10,000 | 3-4 | 4-8 GB |
| 10,000 - 100,000 | 6-8 | 16-32 GB |
| > 100,000 | 10+ | 按需扩展 |

### 关键参数调整

```bash
# /etc/security/limits.conf
* soft nofile 65535
* hard nofile 65535

# sysctl.conf
net.core.somaxconn = 65535
net.ipv4.tcp_max_syn_backlog = 65535
```

## 🐛 故障排查

### 消息不同步

1. 检查 Redis 连接: `redis-cli ping`
2. 检查 Pub/Sub: `redis-cli subscribe chat:new_message`
3. 查看后端日志中的集群初始化信息

### WebSocket 断开

1. 检查 Nginx 配置的 WebSocket 头部
2. 确认 `proxy_read_timeout` 足够长
3. 检查客户端重连逻辑

### 负载不均

1. 确认使用 `ip_hash` 而非 `round_robin`
2. 检查是否有 IP 变化（如移动网络）
3. 考虑使用 `sticky` 模块基于 Cookie 会话保持
