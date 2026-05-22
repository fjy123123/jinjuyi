# Chat System Pro - 下载指南

## 📦 源码下载

### 源码包位置
```
路径：/workspace/chat-system-pro/releases/
文件：chat-system-pro-v1.0.0-20260522_044311.tar.gz
大小：85KB
```

### 下载方法

#### 方法1: SCP命令（推荐Linux/macOS）
```bash
# 从服务器下载到本地
scp user@your-server-ip:/workspace/chat-system-pro/releases/chat-system-pro-v1.0.0-20260522_044311.tar.gz ./

# 示例
scp root@192.168.1.100:/workspace/chat-system-pro/releases/chat-system-pro-v1.0.0-20260522_044311.tar.gz ./
```

#### 方法2: SFTP工具
使用以下任一工具连接服务器：
- FileZilla (https://filezilla-project.org/)
- WinSCP (https://winscp.net/)
- Cyberduck (https://cyberduck.io/)

连接到服务器后导航到：
```
/workspace/chat-system-pro/releases/
```
下载 `chat-system-pro-v1.0.0-20260522_044311.tar.gz` 文件

#### 方法3: 挂载远程文件系统
使用 sshfs 挂载服务器目录：
```bash
# Linux/macOS
sudo apt install sshfs  # Ubuntu/Debian
# 或
brew install sshfs     # macOS

# 创建挂载点
mkdir -p ~/remote-server

# 挂载
sshfs user@your-server-ip:/ ~/remote-server

# 访问文件
cd ~/remote-server/workspace/chat-system-pro/releases/
```

---

## 🛠️ 必需工具快速下载

### 核心工具（必装）

#### Docker Desktop
- **Windows**: https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe
- **macOS Intel**: https://desktop.docker.com/mac/main/amd64/Docker.dmg
- **macOS M1/M2**: https://desktop.docker.com/mac/main/arm64/Docker.dmg

#### Git
- **Windows**: https://github.com/git-for-windows/git/releases/download/v2.42.0.windows.1/Git-2.42.0-64-bit.exe
- **macOS**: `brew install git`
- **Linux**: `sudo apt install git`

#### Go 1.21
- **Linux x86_64**: https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
- **macOS**: https://go.dev/dl/go1.21.0.darwin-arm64.tar.gz
- **Windows**: https://go.dev/dl/go1.21.0.windows-amd64.msi

#### Node.js 18 LTS
- **Windows**: https://nodejs.org/dist/v18.17.0/node-v18.17.0-x64.msi
- **macOS**: https://nodejs.org/dist/v18.17.0/node-v18.17.0.pkg
- **Linux**: https://nodejs.org/dist/v18.17.0/node-v18.17.0-linux-x64.tar.xz

### 推荐工具

#### VS Code
- **官网**: https://code.visualstudio.com/
- **Windows x64**: https://code.visualstudio.com/sha/download?build=stable&os=win32-x64

#### FinalShell (SSH工具)
- **官网**: https://www.finalshell.com/
- **下载**: https://www.finalshell.com/img/finalshell.zip

#### Postman (API测试)
- **官网**: https://www.postman.com/downloads/
- **Windows**: https://dl.pstmn.io/download/latest/win64

### 数据库工具

#### MySQL Workbench
- **官网**: https://dev.mysql.com/downloads/workbench/
- **下载**: https://dev.mysql.com/get/Downloads/MySQLWorkbench/mysql-workbench-community-8.0.34-windows-x86_64.msi

#### RedisInsight
- **官网**: https://redis.com/redis-enterprise/redis-insight/
- **Windows**: https://downloads.redisinsight.redis.com/latest/RedisInsight-win64.exe

---

## 📋 快速安装清单

### Windows 10/11 用户
```
☐ 1. Docker Desktop: https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe
☐ 2. Git: https://github.com/git-for-windows/git/releases/download/v2.42.0.windows.1/Git-2.42.0-64-bit.exe
☐ 3. Go: https://go.dev/dl/go1.21.0.windows-amd64.msi
☐ 4. Node.js: https://nodejs.org/dist/v18.17.0/node-v18.17.0-x64.msi
☐ 5. VS Code: https://code.visualstudio.com/sha/download?build=stable&os=win32-x64
☐ 6. FinalShell: https://www.finalshell.com/img/finalshell.zip
☐ 7. 7-Zip: https://www.7-zip.org/a/7z2301-x64.exe
```

### macOS Intel/Apple Silicon 用户
```
☐ 1. Docker Desktop: https://desktop.docker.com/mac/main/arm64/Docker.dmg (M系列) 或 https://desktop.docker.com/mac/main/amd64/Docker.dmg (Intel)
☐ 2. Git: 终端运行: brew install git
☐ 3. Go: https://go.dev/dl/go1.21.0.darwin-arm64.tar.gz (M系列) 或 https://go.dev/dl/go1.21.0.darwin-amd64.tar.gz (Intel)
☐ 4. Node.js: https://nodejs.org/dist/v18.17.0/node-v18.17.0.pkg
☐ 5. VS Code: https://code.visualstudio.com/sha/download?build=stable&os=darwin
☐ 6. Homebrew: /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
☐ 7. The Unarchiver: App Store搜索下载
```

### Ubuntu/Debian Linux 用户
```bash
# 一行命令安装所有必需工具
sudo apt update && sudo apt install -y docker.io docker-compose git nginx redis-server mysql-server curl wget
```

### CentOS/RHEL Linux 用户
```bash
# 一行命令安装Docker
sudo yum install -y yum-utils device-mapper-persistent-data lvm2
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo systemctl start docker && sudo systemctl enable docker

# 安装其他工具
sudo yum install git nginx redis mysql-server
```

---

## 🔍 验证安装

安装完成后，运行以下命令验证：

```bash
# 验证Docker
docker --version
docker compose version

# 验证Git
git --version

# 验证Go
go version

# 验证Node.js
node --version
npm --version

# 验证MySQL
mysql --version

# 验证Redis
redis-server --version
```

---

## 📚 完整文档

下载源码后，查看以下文档：

1. **快速开始**: QUICKSTART.md
2. **详细安装**: INSTALL.md
3. **API文档**: API.md
4. **工具地址**: DOWNLOAD.md（本文档）
5. **项目说明**: README.md

---

## 🌐 在线资源

### 官方文档
- Docker Docs: https://docs.docker.com/
- Go Docs: https://go.dev/doc/
- Gin Framework: https://gin-gonic.com/
- React Docs: https://react.dev/

### 学习资源
- Docker入门: https://docker-curriculum.com/
- Go语言教程: https://tour.golang.org/
- React教程: https://reactjs.org/tutorial/tutorial.html

### 社区支持
- Docker Community: https://forums.docker.com/
- Go Forum: https://forum.golang.org/
- Reddit r/docker: https://www.reddit.com/r/docker/

---

## ❓ 常见问题

### Q: 下载速度慢？
使用国内镜像源或选择离您近的下载节点

### Q: 安装失败？
查看详细错误信息，参考INSTALL.md中的故障排除部分

### Q: 版本不兼容？
确保使用推荐版本，详细信息见INSTALL.md

### Q: 需要帮助？
- 查看文档: INSTALL.md, API.md
- 查看FAQ: INSTALL.md 常见问题章节
- 联系支持: support@example.com

---

**祝您使用愉快！**
