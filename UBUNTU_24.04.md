# Ubuntu 24.04 一键部署指南

## 快速开始

### 方式一：一键命令（推荐）

在您的 Ubuntu 24.04 服务器上直接运行：

```bash
curl -fsSL https://raw.githubusercontent.com/fjy123123/jinjuyi/main/ubuntu-24.04-install.sh | bash
```

或者使用 wget：

```bash
wget -qO- https://raw.githubusercontent.com/fjy123123/jinjuyi/main/ubuntu-24.04-install.sh | bash
```

---

### 方式二：手动下载脚本

```bash
# 下载安装脚本
wget https://raw.githubusercontent.com/fjy123123/jinjuyi/main/ubuntu-24.04-install.sh

# 添加执行权限
chmod +x ubuntu-24.04-install.sh

# 运行脚本
./ubuntu-24.04-install.sh
```

---

## 脚本功能

这个一键部署脚本会自动完成以下操作：

1. ✅ 更新系统
2. ✅ 安装基础依赖
3. ✅ 配置防火墙 (UFW
4. ✅ 安装 Docker (官方最新版)
5. ✅ 安装 Docker Compose
6. ✅ 克隆项目代码
7. ✅ 配置环境变量
8. ✅ 启动所有服务
9. ✅ 显示访问地址

---

## 部署方式选择

脚本启动后，会出现以下选项：

```
请选择部署方式:
  1. 完整部署 (推荐)
  2. 仅启动服务
  0. 退出
```

- **选项 1：完整部署（首次安装）
- **选项 2：仅启动服务（已安装过）

---

## 部署后

部署完成后，会显示：

```
===============================================
  部署完成！
===============================================

  访问地址:
    前端: http://你的服务器IP
    API:  http://你的服务器IP:8080

  常用命令 (在 jinjuyi 目录下执行):
    查看日志: docker compose logs -f
    停止服务: docker compose down
    重启服务: docker compose restart
    更新代码: git pull && docker compose pull && docker compose up -d
```

---

## 手动部署（备选方案

如果自动脚本遇到问题，可以手动部署：

### 1. 安装 Docker

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装依赖
sudo apt install -y curl git ca-certificates gnupg lsb-release

# 安装 Docker
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER

# 重新登录以生效 docker 权限
```

### 2. 配置防火墙

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp
sudo ufw --force enable
```

### 3. 下载项目

```bash
git clone https://github.com/fjy123123/jinjuyi.git
cd jinjuyi
```

### 4. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件配置
nano .env
```

### 5. 启动服务

```bash
docker compose up -d
```

---

## 常用命令

```bash
# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f

# 停止服务
docker compose down

# 重启服务
docker compose restart

# 更新代码
git pull
docker compose pull
docker compose up -d
```

---

## 常见问题

### Q: Docker 命令提示无权限？

A: 重新登录后生效 docker 组权限：
```bash
logout
# 重新 SSH 重新登录
```

或者临时使用 sudo：
```bash
sudo docker compose up -d
```

### Q: 防火墙问题？

A: 手动开放端口：
```bash
sudo ufw allow 80/tcp
sudo ufw allow 8080/tcp
```

### Q: 更新系统版本？

A: 确保使用 Ubuntu 24.04 LTS (Noble Numbat)

---

## 技术支持

- GitHub Issues: https://github.com/fjy123123/jinjuyi/issues

---

祝您部署成功！
