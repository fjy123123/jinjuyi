# Chat System Pro - 快速下载指南

## 🎯 立即获取

### 📦 源码包（85KB）
```
文件：chat-system-pro-v1.0.0-20260522_044311.tar.gz
位置：/workspace/chat-system-pro/releases/
```

**下载命令：**
```bash
scp root@YOUR_SERVER_IP:/workspace/chat-system-pro/releases/chat-system-pro-v1.0.0-20260522_044311.tar.gz ./
```

---

## 🛠️ 必需工具（一键安装）

### Windows 用户
| 工具 | 下载地址 | 大小 |
|------|---------|------|
| Docker Desktop | https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe | ~500MB |
| Git | https://github.com/git-for-windows/git/releases/download/v2.42.0.windows.1/Git-2.42.0-64-bit.exe | ~50MB |
| Go 1.21 | https://go.dev/dl/go1.21.0.windows-amd64.msi | ~100MB |
| Node.js 18 | https://nodejs.org/dist/v18.17.0/node-v18.17.0-x64.msi | ~30MB |
| VS Code | https://code.visualstudio.com/sha/download?build=stable&os=win32-x64 | ~90MB |
| 7-Zip | https://www.7-zip.org/a/7z2301-x64.exe | ~1.5MB |

**快速下载脚本（PowerShell）：**
```powershell
# 创建下载目录
mkdir C:\Tools -Force
cd C:\Tools

# 下载Docker
Invoke-WebRequest -Uri "https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe" -OutFile "DockerDesktopInstaller.exe"

# 下载Git
Invoke-WebRequest -Uri "https://github.com/git-for-windows/git/releases/download/v2.42.0.windows.1/Git-2.42.0-64-bit.exe" -OutFile "Git-2.42.0-64-bit.exe"

# 下载Go
Invoke-WebRequest -Uri "https://go.dev/dl/go1.21.0.windows-amd64.msi" -OutFile "go1.21.0.windows-amd64.msi"

# 下载Node.js
Invoke-WebRequest -Uri "https://nodejs.org/dist/v18.17.0/node-v18.17.0-x64.msi" -OutFile "node-v18.17.0-x64.msi"

# 下载VS Code
Invoke-WebRequest -Uri "https://code.visualstudio.com/sha/download?build=stable&os=win32-x64" -OutFile "VSCodeSetup-stable.exe"

# 下载7-Zip
Invoke-WebRequest -Uri "https://www.7-zip.org/a/7z2301-x64.exe" -OutFile "7z2301-x64.exe"
```

### macOS 用户
| 工具 | 下载地址 | 说明 |
|------|---------|------|
| Docker Desktop | https://desktop.docker.com/mac/main/arm64/Docker.dmg (M芯片) 或 https://desktop.docker.com/mac/main/amd64/Docker.dmg (Intel) | 选择对应版本 |
| Git | 终端运行: `brew install git` | 需要Homebrew |
| Go 1.21 | https://go.dev/dl/go1.21.0.darwin-arm64.tar.gz | M芯片 |
| Node.js 18 | https://nodejs.org/dist/v18.17.0/node-v18.17.0.pkg | 自动识别芯片 |
| VS Code | https://code.visualstudio.com/sha/download?build=stable&os=darwin | 自动识别芯片 |

**快速安装脚本（终端）：**
```bash
# 安装Homebrew（如果没有）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 一键安装所有工具
brew install --cask docker git go node vscode
```

### Linux 用户（Ubuntu/Debian）
| 工具 | 安装命令 |
|------|---------|
| Docker | `curl -fsSL https://get.docker.com | sudo sh` |
| Docker Compose | `sudo apt install docker-compose` |
| Git | `sudo apt install git` |
| Go | `wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz && sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz` |
| Node.js | `curl -fsSL https://deb.nodesource.com/setup_18.x \| sudo -E bash - && sudo apt install nodejs` |
| Nginx | `sudo apt install nginx` |

**一键安装所有工具（Ubuntu/Debian）：**
```bash
sudo apt update && sudo apt upgrade -y
curl -fsSL https://get.docker.com | sudo sh
sudo apt install git nginx curl wget
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install nodejs
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**一键安装所有工具（CentOS/RHEL）：**
```bash
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo systemctl start docker && sudo systemctl enable docker
sudo yum install git nginx curl wget
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install nodejs
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

---

## 📚 文档下载

所有文档都在源码包内，解压后即可查看：

| 文档 | 说明 | 优先级 |
|------|------|--------|
| QUICKSTART.md | 5分钟快速入门 | ⭐⭐⭐ |
| INSTALL.md | 详细安装文档 | ⭐⭐⭐ |
| API.md | 完整API文档 | ⭐⭐⭐ |
| DOWNLOAD.md | 所有工具下载地址 | ⭐⭐ |
| GETTING_STARTED.md | 下载指南（本文档） | ⭐⭐ |
| DEPLOY.md | 部署指南 | ⭐⭐ |
| README.md | 项目说明 | ⭐ |
| FEATURES.md | 功能说明 | ⭐ |
| PACKAGE.md | 打包清单 | ⭐ |
| DELIVERY.md | 交付清单 | ⭐ |

---

## 🚀 快速开始

### 1. 下载源码
```bash
scp root@YOUR_SERVER_IP:/workspace/chat-system-pro/releases/chat-system-pro-v1.0.0-20260522_044311.tar.gz ./
```

### 2. 安装工具
按照上面的表格安装必需工具

### 3. 解压源码
```bash
tar -xzf chat-system-pro-v1.0.0-20260522_044311.tar.gz
cd chat-system-pro-v1.0.0-20260522_044311
```

### 4. 一键部署
```bash
# Windows
deploy.bat

# Linux/macOS
chmod +x deploy.sh
./deploy.sh
```

### 5. 访问系统
- 前端: http://localhost
- API: http://localhost:8080
- 健康检查: http://localhost:8080/health

---

## 🔧 数据库（Docker会自动安装）

如果手动部署需要单独安装数据库：

### MySQL
```bash
# Docker
docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=your_password mysql:8.0

# 或直接安装
# Ubuntu: sudo apt install mysql-server
# CentOS: sudo yum install mysql-community-server
```

### MongoDB
```bash
# Docker
docker run -d --name mongodb -p 27017:27017 mongo:7

# 或直接安装
# Ubuntu: sudo apt install mongodb
# CentOS: sudo yum install mongodb-org
```

### Redis
```bash
# Docker
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 或直接安装
# Ubuntu: sudo apt install redis-server
# CentOS: sudo yum install redis
```

---

## 📞 技术支持

### 文档
- 快速入门: QUICKSTART.md
- 详细安装: INSTALL.md
- API文档: API.md
- 工具下载: DOWNLOAD.md, GETTING_STARTED.md

### 联系方式
- 邮箱: support@example.com
- 问题反馈: https://github.com/yourrepo/chat-system-pro/issues

---

## ✨ 版本信息

- **项目版本**: v1.0.0
- **源码包**: chat-system-pro-v1.0.0-20260522_044311.tar.gz
- **发布日期**: 2024-01-15
- **Go版本**: 1.21+
- **Node.js版本**: 18+
- **Docker版本**: 20.10+

---

**立即开始：下载源码 → 安装工具 → 一键部署！**
