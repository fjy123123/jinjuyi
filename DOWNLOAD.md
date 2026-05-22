# Chat System Pro - 源码下载与工具地址

## 📦 源码下载

### 项目源码包
```
文件：chat-system-pro-v1.0.0-20260522_044311.tar.gz
大小：85KB
路径：/workspace/chat-system-pro/releases/
```

**下载方式：**
1. 直接从服务器下载：`scp user@server:/workspace/chat-system-pro/releases/chat-system-pro-v1.0.0-20260522_044311.tar.gz .`
2. 使用SFTP工具下载
3. 挂载服务器文件系统下载

---

## 🛠️ 必备工具下载地址

### 1. Docker（容器化平台）

#### Docker Desktop（Windows/macOS）
- **官网**: https://www.docker.com/products/docker-desktop/
- **Windows版**: https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe
- **macOS Intel**: https://desktop.docker.com/mac/main/amd64/Docker.dmg
- **macOS Apple Silicon**: https://desktop.docker.com/mac/main/arm64/Docker.dmg
- **版本要求**: Docker Desktop 4.0+ (Docker 20.10+)

#### Docker Engine（Linux服务器）
- **Ubuntu/Debian**:
  ```bash
  curl -fsSL https://get.docker.com | sudo sh
  ```
- **CentOS/RHEL**:
  ```bash
  sudo yum install -y yum-utils
  sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
  sudo yum install docker-ce docker-ce-cli containerd.io docker-compose-plugin
  ```

#### Docker Compose（独立版）
- **下载地址**: https://github.com/docker/compose/releases
- **直接下载（v2.20.0）**:
  - Linux x86_64: https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-linux-x86_64
  - Linux ARM64: https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-linux-aarch64
  - macOS: https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-darwin-x86_64
  - Windows: https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-windows-x86_64.exe

---

### 2. Go语言（后端编译）

- **官网**: https://go.dev/dl/
- **直接下载**:
  - **Linux x86_64**: https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
  - **Linux ARM64**: https://go.dev/dl/go1.21.0.linux-arm64.tar.gz
  - **macOS x86_64**: https://go.dev/dl/go1.21.0.darwin-amd64.tar.gz
  - **macOS ARM64**: https://go.dev/dl/go1.21.0.darwin-arm64.tar.gz
  - **Windows**: https://go.dev/dl/go1.21.0.windows-amd64.msi
- **版本要求**: Go 1.21+

---

### 3. Node.js（前端开发）

- **官网**: https://nodejs.org/
- **LTS版本直接下载**:
  - **Windows x64**: https://nodejs.org/dist/v18.17.0/node-v18.17.0-x64.msi
  - **macOS x64**: https://nodejs.org/dist/v18.17.0/node-v18.17.0.pkg
  - **macOS ARM64**: https://nodejs.org/dist/v18.17.0/node-v18.17.0.pkg
  - **Linux x64**: https://nodejs.org/dist/v18.17.0/node-v18.17.0-linux-x64.tar.xz
- **版本要求**: Node.js 18+

---

### 4. 数据库软件

#### MySQL 8.0
- **官网**: https://dev.mysql.com/downloads/mysql/
- **直接下载**:
  - **Windows**: https://dev.mysql.com/get/Downloads/MySQLInstaller/mysql-installer-community-8.0.34.0.msi
  - **Linux (APT)**: `sudo apt install mysql-server`
  - **Linux (YUM)**: `sudo yum install mysql-community-server`
- **Docker镜像**: `docker pull mysql:8.0`

#### MongoDB 7.0
- **官网**: https://www.mongodb.com/try/download/community
- **直接下载**:
  - **Windows**: https://fastdl.mongodb.org/windows/mongodb-windows-x86_64-7.0.1-signed.msi
  - **Linux x86_64**: https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu2204-7.0.1.tgz
  - **Linux ARM64**: https://fastdl.mongodb.org/linux/mongodb-linux-aarch64-ubuntu2204-7.0.1.tgz
  - **macOS**: https://fastdl.mongodb.org/osx/mongodb-macos-arm64-7.0.1.tgz
- **Docker镜像**: `docker pull mongo:7`

#### Redis 7.0
- **官网**: https://redis.io/download/
- **直接下载**:
  - **Windows**: https://github.com/microsoftarchive/redis/releases (Redis 3.0 for Windows)
  - **Linux**: `sudo apt install redis-server` 或 `sudo yum install redis`
  - **源码**: https://github.com/redis/redis/archive/7.0.13.tar.gz
- **Docker镜像**: `docker pull redis:7-alpine`

---

### 5. Nginx（Web服务器）

- **官网**: https://nginx.org/en/download.html
- **直接下载**:
  - **Windows**: https://nginx.org/en/download.html (Stable version)
  - **Linux**: 通常通过包管理器安装
    - Ubuntu/Debian: `sudo apt install nginx`
    - CentOS/RHEL: `sudo yum install nginx`
- **Docker镜像**: `docker pull nginx:alpine`

---

### 6. Git（版本控制）

- **官网**: https://git-scm.com/downloads
- **直接下载**:
  - **Windows**: https://github.com/git-for-windows/git/releases/download/v2.42.0.windows.1/Git-2.42.0-64-bit.exe
  - **macOS**: `brew install git` 或 https://git-scm.com/download/mac
  - **Linux**: `sudo apt install git` 或 `sudo yum install git`

---

### 7. 代码编辑器/IDE

#### VS Code
- **官网**: https://code.visualstudio.com/
- **直接下载**:
  - **Windows**: https://code.visualstudio.com/sha/download?build=stable&os=win32-x64
  - **macOS**: https://code.visualstudio.com/sha/download?build=stable&os=darwin
  - **Linux**: https://code.visualstudio.com/sha/download?build=stable&os=linux-x64

#### GoLand（JetBrains）
- **官网**: https://www.jetbrains.com/go/download/
- **学生免费**: https://www.jetbrains.com/community/education/#students

#### WebStorm（JetBrains）
- **官网**: https://www.jetbrains.com/webstorm/download/
- **学生免费**: https://www.jetbrains.com/community/education/#students

---

### 8. 数据库管理工具

#### MySQL Workbench
- **官网**: https://dev.mysql.com/downloads/workbench/
- **直接下载**: https://dev.mysql.com/get/Downloads/MySQLWorkbench/mysql-workbench-community-8.0.34-windows-x86_64.msi

#### Navicat（推荐）
- **官网**: https://www.navicat.com/en/products
- **免费试用**: https://www.navicat.com/en/download/navicat-premium

#### Studio 3T（MongoDB）
- **官网**: https://studio3t.com/download/
- **免费版**: Studio 3T Free

#### RedisInsight
- **官网**: https://redis.com/redis-enterprise/redis-insight/
- **直接下载**: https://downloads.redisinsight.redis.com/latest/RedisInsight-win64.exe

---

### 9. API测试工具

#### Postman
- **官网**: https://www.postman.com/downloads/
- **直接下载**:
  - **Windows**: https://dl.pstmn.io/download/latest/win64
  - **macOS**: https://dl.pstmn.io/download/latest/osx
  - **Linux**: https://dl.pstmn.io/download/latest/linux64

#### Insomnia
- **官网**: https://insomnia.rest/download
- **直接下载**: https://updates.insomnia.rest/downloads/insomnia-v2023.5.8.zip

---

### 10. SSH客户端

#### FinalShell（推荐，支持Windows/macOS/Linux）
- **官网**: https://www.finalshell.com/
- **下载地址**: https://www.finalshell.com/img/finalshell.pkg

#### PuTTY（Windows）
- **官网**: https://www.chiark.greenend.org.uk/~sgtatham/putty/latest.html
- **直接下载**: https://the.earth.li/~sgtatham/putty/latest/w64/putty-64bit-0.79-installer.msi

#### Termius（多平台）
- **官网**: https://termius.com/
- **下载**: https://termius.com/windows

#### MobaXterm（Windows增强终端）
- **官网**: https://mobaxterm.mobatek.net/download.html
- **免费版**: https://mobaxterm.mobatek.net/MobaXterm_Portable_v23.1.zip

---

### 11. SFTP/FTP工具

#### FileZilla
- **官网**: https://filezilla-project.org/download.php
- **直接下载**: https://filezilla-project.org/download.php?platform=win64

#### WinSCP（Windows）
- **官网**: https://winscp.net/eng/download.php
- **直接下载**: https://winscp.net/download/WinSCP-6.1.1-Setup.exe

#### Cyberduck（macOS/Windows）
- **官网**: https://cyberduck.io/download/
- **直接下载**: https://update.cyberduck.io/Cyberduck-8.5.5.39283.pkg

---

### 12. 文本编辑器

#### Notepad++（Windows）
- **官网**: https://notepad-plus-plus.org/downloads/
- **直接下载**: https://github.com/notepad-plus-plus/notepad-plus-plus/releases/download/v8.5.7/npp.8.5.7.Installer.x64.exe

#### Sublime Text（跨平台）
- **官网**: https://www.sublimetext.com/download
- **直接下载**: https://download.sublimetext.com/Sublime%20Text%20Build%203212%20x64 Setup.exe

#### Nano/Vim（Linux终端）
- 通常预装，无需额外下载

---

### 13. 虚拟化/容器工具

#### VirtualBox
- **官网**: https://www.virtualbox.org/wiki/Downloads
- **直接下载**: https://download.virtualbox.org/virtualbox/7.0.10/VirtualBox-7.0.10-158379-Win.exe

#### VMware Workstation Pro（商业软件）
- **官网**: https://www.vmware.com/products/workstation-pro.html

---

### 14. 监控工具

#### Prometheus
- **官网**: https://prometheus.io/download/
- **直接下载**: https://github.com/prometheus/prometheus/releases/download/v2.47.0/prometheus-2.47.0.linux-amd64.tar.gz

#### Grafana
- **官网**: https://grafana.com/grafana/download
- **直接下载**: https://dl.grafana.com/oss/release/grafana-10.1.0.windows-amd64.zip

#### Portainer（Docker管理界面）
- **官网**: https://www.portainer.io/
- **Docker运行**: `docker run -d -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock portainer/portainer`

---

### 15. SSL证书工具

#### Let's Encrypt (Certbot)
- **官网**: https://certbot.eff.org/
- **安装命令**:
  - Ubuntu/Debian: `sudo apt install certbot python3-certbot-nginx`
  - CentOS: `sudo yum install certbot python3-certbot-nginx`

#### OpenSSL（通常预装）
- **官网**: https://www.openssl.org/source/
- **Windows**: https://slproweb.com/products/Win32OpenSSL.html

---

### 16. 压缩/解压工具

#### 7-Zip（Windows）
- **官网**: https://www.7-zip.org/download.html
- **直接下载**: https://www.7-zip.org/a/7z2301-x64.exe

#### The Unarchiver（macOS）
- **官网**: https://theunarchiver.com/
- **App Store**: 在App Store搜索"The Unarchiver"

#### PeaZip（Windows/macOS/Linux）
- **官网**: https://peazip.org/
- **直接下载**: https://github.com/peazip/PeaZip/releases/download/9.6.0/peazip-9.6.0.WINDOWS.x86_64.zip

---

### 17. 系统信息工具

#### CPU-Z（Windows）
- **官网**: https://www.cpuid.com/softwares/cpu-z.html
- **直接下载**: https://download.cpuid.com/cpu-z/cpu-z_1.99.pdf

#### htop（Linux/macOS终端）
- **官网**: https://htop.dev/
- **安装**: `sudo apt install htop` 或 `brew install htop`

#### Neofetch（Linux终端）
- **官网**: https://github.com/dylanaraps/neofetch
- **安装**: `sudo apt install neofetch` 或 `brew install neofetch`

---

## 🔧 快速安装命令

### Ubuntu/Debian
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装Docker
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 安装Git
sudo apt install git

# 安装Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install nodejs

# 安装Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# 安装Nginx
sudo apt install nginx

# 安装Redis
sudo apt install redis-server

# 安装MySQL
sudo apt install mysql-server
```

### CentOS/RHEL
```bash
# 安装Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 启用Docker
sudo systemctl start docker
sudo systemctl enable docker

# 安装其他工具
sudo yum install git nginx redis mysql-server
```

---

## 📥 推荐下载清单

### Windows用户
1. Docker Desktop: https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe
2. Git: https://github.com/git-for-windows/git/releases/download/v2.42.0.windows.1/Git-2.42.0-64-bit.exe
3. VS Code: https://code.visualstudio.com/sha/download?build=stable&os=win32-x64
4. FinalShell: https://www.finalshell.com/img/finalshell.zip
5. 7-Zip: https://www.7-zip.org/a/7z2301-x64.exe
6. Postman: https://dl.pstmn.io/download/latest/win64

### macOS用户
1. Docker Desktop: https://desktop.docker.com/mac/main/arm64/Docker.dmg (Apple Silicon) 或 https://desktop.docker.com/mac/main/amd64/Docker.dmg (Intel)
2. Git: `brew install git`
3. VS Code: https://code.visualstudio.com/sha/download?build=stable&os=darwin
4. FinalShell: https://www.finalshell.com/img/finalshell.pkg
5. The Unarchiver: 在App Store搜索下载
6. Homebrew: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`

### Linux用户
```bash
# 一键安装所有工具
sudo apt update
sudo apt install -y docker.io docker-compose git nginx redis-server mysql-server postman
```

---

## 🌐 镜像加速（可选）

如果网络访问GitHub、Docker Hub等国外网站较慢，可以使用国内镜像：

### Docker镜像加速
```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ]
}
```

### Go模块加速
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### npm镜像
```bash
npm config set registry https://registry.npmmirror.com
```

### GitHub加速
- 使用FastGit: https://fastgit.org/
- 或Gitee镜像

---

## 📞 技术支持

如遇下载问题：
- **官方文档**: 查看各工具官网
- **社区支持**: 加入相关技术社区
- **邮箱支持**: support@example.com

---

**版本**: v1.0.0  
**更新日期**: 2024-01-15
