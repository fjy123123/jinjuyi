# Chat System Pro - 系统自启动配置指南

## 方法一：使用 systemd 服务（推荐）✅

### 1. 创建服务文件

创建 systemd 服务文件：

```bash
sudo nano /etc/systemd/system/chat-system-pro.service
```

复制以下内容（注意修改路径）：

```ini
[Unit]
Description=Chat System Pro - Docker Compose Service
Requires=docker.service
After=docker.service network-online.target
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/fjya/jinjuyi
ExecStart=/usr/bin/docker compose up -d
ExecStop=/usr/bin/docker compose down
ExecReload=/usr/bin/docker compose restart
TimeoutStartSec=0
Restart=on-failure
RestartSec=10s
User=fjya
Group=docker

[Install]
WantedBy=multi-user.target
```

**重要**：请根据实际情况修改：
- `WorkingDirectory`: 你的项目路径
- `User`: 你的用户名
- `ExecStart/ExecStop/ExecReload`: 根据你的 docker-compose 版本使用 `docker compose` 或 `docker-compose`

### 2. 设置权限

```bash
sudo chmod 644 /etc/systemd/system/chat-system-pro.service
sudo chown root:root /etc/systemd/system/chat-system-pro.service
```

### 3. 重新加载 systemd

```bash
sudo systemctl daemon-reload
```

### 4. 启用服务（开机自启动）

```bash
sudo systemctl enable chat-system-pro.service
```

### 5. 立即启动服务

```bash
sudo systemctl start chat-system-pro.service
```

### 6. 检查服务状态

```bash
sudo systemctl status chat-system-pro.service
```

### 7. 常用命令

```bash
# 启动服务
sudo systemctl start chat-system-pro

# 停止服务
sudo systemctl stop chat-system-pro

# 重启服务
sudo systemctl restart chat-system-pro

# 查看状态
sudo systemctl status chat-system-pro

# 查看日志
sudo journalctl -u chat-system-pro -f

# 禁用自启动
sudo systemctl disable chat-system-pro
```

---

## 方法二：使用 rc.local（传统方式）

### 1. 启用 rc-local 服务

```bash
sudo systemctl enable rc-local
```

### 2. 编辑 rc-local 服务文件

```bash
sudo nano /etc/systemd/system/rc-local.service
```

添加以下内容：

```ini
[Unit]
Description=/etc/rc.local Compatibility
After=network.target

[Service]
Type=forking
ExecStart=/etc/rc.local start
TimeoutSuccess=0
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
```

### 3. 创建 rc.local 文件

```bash
sudo nano /etc/rc.local
```

添加以下内容：

```bash
#!/bin/bash
# rc.local

# 等待系统启动完成
sleep 10

# 切换到项目目录
cd /home/fjya/jinjuyi

# 启动 docker compose 服务
/usr/local/bin/docker-compose up -d

exit 0
```

### 4. 设置权限

```bash
sudo chmod +x /etc/rc.local
sudo systemctl start rc-local
sudo systemctl enable rc-local
```

---

## 方法三：使用 crontab @reboot

### 1. 编辑 crontab

```bash
crontab -e
```

### 2. 添加启动命令

在 crontab 文件末尾添加：

```bash
@reboot cd /home/fjya/jinjuyi && /usr/local/bin/docker-compose up -d >> /home/fjya/jinjuyi/boot.log 2>&1
```

### 3. 保存并退出

```
# 如果是 nano 编辑器
Ctrl + O 保存
Ctrl + X 退出
```

---

## 方法四：使用 @reboot 延迟启动

有时候系统启动后 Docker 服务还没就绪，可以延迟启动：

```bash
@reboot sleep 30 && cd /home/fjya/jinjuyi && /usr/bin/docker-compose up -d
```

---

## 验证自启动配置

### 检查服务是否在运行

```bash
# 检查 systemd 服务状态
sudo systemctl is-active chat-system-pro

# 检查 docker 容器状态
docker ps

# 检查端口是否监听
sudo netstat -tlnp | grep -E ':(80|8080)'
```

### 查看启动日志

```bash
# systemd 日志
sudo journalctl -u chat-system-pro -n 50

# rc.local 日志
cat /home/fjya/jinjuyi/boot.log

# Docker 日志
docker compose logs
```

### 测试重启

```bash
# 重启系统
sudo reboot

# 等待几分钟后检查
sudo systemctl status chat-system-pro
docker ps
```

---

## 常见问题

### Q1: 服务启动失败

**检查步骤**：
```bash
# 查看详细日志
sudo journalctl -u chat-system-pro -xe

# 检查 Docker 服务状态
sudo systemctl status docker

# 手动启动测试
cd /home/fjya/jinjuyi
docker compose up -d
```

### Q2: 权限问题

**解决**：
```bash
# 确保用户属于 docker 组
sudo usermod -aG docker $USER

# 重新登录以生效
logout
```

### Q3: 启动顺序问题

**解决**：使用 systemd 的 `After=` 和 `Wants=` 配置依赖关系

### Q4: 网络未就绪

**解决**：使用 `network-online.target` 等待网络就绪，或添加延迟

---

## 推荐配置

**推荐使用方法一（systemd）**，因为：
- ✅ 支持服务依赖管理
- ✅ 支持自动重启失败的服务
- ✅ 支持日志管理
- ✅ 支持服务状态监控
- ✅ 标准的 Linux 服务管理方式

---

## 一键安装脚本

如果你不想手动配置，我提供了一个自动化脚本：

```bash
# 在项目目录执行
chmod +x setup-autostart.sh
sudo ./setup-autostart.sh
```

脚本会自动：
- 检测项目路径
- 创建 systemd 服务文件
- 启用自启动
- 立即启动服务
- 检查服务状态

---

## 安全建议

1. **定期检查日志**：监控服务启动状态
2. **设置日志轮转**：避免日志文件过大
3. **监控磁盘空间**：Docker 日志可能占用大量空间
4. **配置防火墙**：只开放必要的端口
5. **定期更新**：保持系统和 Docker 最新

---

## 参考文档

- [systemd 官方文档](https://www.freedesktop.org/wiki/Software/systemd/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [Ubuntu 服务管理](https://help.ubuntu.com/community/ systemd)

---

祝你部署成功！🎉
