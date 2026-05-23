# Oracle Cloud Free Tier - 知信聊天系统部署指南

## 🎉 Oracle Cloud Free Tier 优势

- ✅ **永久免费ARM实例**（24GB内存，4个CPU核心！）
- ✅ **永久免费数据库**（2个Autonomous Database）
- ✅ **永久免费存储**（10GB对象存储）
- ✅ **永不过期**，不会收费

## 📋 目录

1. [注册Oracle Cloud](#1-注册oracle-cloud)
2. [创建免费ARM实例](#2-创建免费arm实例)
3. [配置实例](#3-配置实例)
4. [安装Docker和Docker Compose](#4-安装docker和docker-compose)
5. [部署知信聊天系统](#5-部署知信聊天系统)
6. [配置域名（可选）](#6-配置域名可选)
7. [故障排除](#7-故障排除)

---

## 1. 注册Oracle Cloud

### 访问注册页面

1. 打开 https://www.oracle.com/cloud/free/
2. 点击 **"Start for free"**

### 完成注册

1. **填写邮箱和密码**
   - 使用真实邮箱（会发送验证邮件）
   - 设置强密码（至少8位，包含大小写和数字）

2. **验证邮箱**
   - 登录邮箱
   - 点击Oracle发送的验证链接

3. **填写个人信息**
   - 国家/地区
   - 姓名
   - 公司名称（可选）
   - 电话号码（需要短信验证）

4. **信用卡验证**
   - 仅用于验证，不会收费
   - 确保账户余额充足（至少$1）
   - 会在几分钟后退回

5. **接受服务条款**

### 等待账户激活

注册完成后，通常需要等待10-30分钟账户激活。

---

## 2. 创建免费ARM实例

### 登录Oracle Cloud控制台

1. 访问 https://cloud.oracle.com
2. 点击 **"Sign In"**
3. 使用你的账户登录

### 创建VM实例

1. **进入计算服务**
   - 在控制台左侧菜单点击 **"Compute"** → **"Instances"**
   - 点击 **"Create Instance"**

2. **配置实例**
   
   **名称**
   ```
   chat-system-pro
   ```
   
   **操作系统选择**
   - 点击 **"Change image"**
   - 选择 **"Ubuntu"** 或 **"Oracle Linux"**
   - 推荐选择 **"Ubuntu 22.04"** 或 **"Ubuntu 24.04"**
   - 点击 **"Select Image"**

3. **选择形状（配置）**
   - 点击 **"Change shape"**
   - 选择 **"Ampere"** 标签（ARM架构）
   - 选择 **"VM.Standard.A1.Flex"**
   - 设置：
     - **OCPU count**: 1（免费额度）
     - **Memory (GB)**: 6（免费额度）
   - 点击 **"Select Shape"**

4. **配置网络**
   - **Virtual cloud network**: 选择 **"Create a new virtual cloud network"**
   - **Subnet**: 选择 **"Create a new subnet"**
   - 确保勾选 **"Assign a public IP address"**（重要！）

5. **添加SSH密钥**
   - 选择 **"Generate a key pair for me"**
   - 点击 **"Save Private Key"** 保存私钥
   - 妥善保管私钥，后续连接服务器需要用到

6. **引导卷配置**
   - 默认50GB，足够使用
   - 不需要额外配置

7. **创建实例**
   - 点击 **"Create"**
   - 等待2-5分钟实例创建完成

---

## 3. 配置实例

### 记录重要信息

创建成功后，记录以下信息：

```
实例信息：
- 公共IP地址：___.___.___.___
- 实例名称：chat-system-pro
- 可用性域：xxx:AD-1

下载的SSH密钥：
- private-key.pem（请妥善保管）
```

### 配置防火墙规则

#### Oracle Cloud安全列表

1. 在实例详情页面，点击 **"Subnet"** 链接
2. 点击 **"Default Security List"**
3. 点击 **"Add Ingress Rules"**

**添加入站规则（允许访问）**

```
# HTTP (80)
Source: 0.0.0.0/0
IP Protocol: TCP
Destination Port Range: 80
Description: HTTP

# HTTPS (443)
Source: 0.0.0.0/0
IP Protocol: TCP
Destination Port Range: 443
Description: HTTPS

# API (8080)
Source: 0.0.0.0/0
IP Protocol: TCP
Destination Port Range: 8080
Description: Backend API

# SSH (22)
Source: 0.0.0.0/0
IP Protocol: TCP
Destination Port Range: 22
Description: SSH

# MySQL (3306)
Source: 10.0.0.0/16
IP Protocol: TCP
Destination Port Range: 3306
Description: MySQL

# Redis (6379)
Source: 10.0.0.0/16
IP Protocol: TCP
Destination Port Range: 6379
Description: Redis
```

### 连接服务器

#### Mac/Linux连接

```bash
# 1. 修改私钥权限
chmod 400 /path/to/private-key.pem

# 2. 连接服务器
ssh -i /path/to/private-key.pem ubuntu@你的公共IP

# 示例：
ssh -i ~/Downloads/ssh-key-2024-01-01.key opc@132.145.67.89
```

#### Windows连接（使用PuTTY）

1. 下载PuTTY: https://www.putty.org/
2. 使用PuTTYgen转换私钥格式
3. 连接：`ubuntu@你的公共IP`

---

## 4. 安装Docker和Docker Compose

### 连接后执行

```bash
# 1. 更新系统
sudo apt update && sudo apt upgrade -y

# 2. 安装依赖
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common

# 3. 添加Docker GPG密钥
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# 4. 添加Docker仓库
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 5. 安装Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 6. 启动Docker
sudo systemctl start docker
sudo systemctl enable docker

# 7. 添加当前用户到docker组（避免每次sudo）
sudo usermod -aG docker $USER

# 8. 验证安装
docker --version
docker compose version
```

### 安装Git（如果没有）

```bash
sudo apt install -y git
```

---

## 5. 部署知信聊天系统

### 克隆项目

```bash
# 1. 创建项目目录
mkdir -p ~/chat-system && cd ~/chat-system

# 2. 克隆仓库
git clone https://github.com/fjy123123/jinjuyi.git .

# 3. 查看文件
ls -la
```

### 创建环境变量文件

```bash
# 1. 复制环境变量模板
cp .env.example .env

# 2. 编辑环境变量
nano .env
```

**修改以下内容**：

```env
# MySQL配置（使用强密码！）
MYSQL_ROOT_PASSWORD=你的强密码_MySQL_Root_2024!
MYSQL_PASSWORD=你的强密码_Chat_User_2024!
MYSQL_USER=chat_user
MYSQL_DATABASE=chat_system_pro

# Redis配置
REDIS_PASSWORD=你的强密码_Redis_2024!

# JWT密钥（使用随机字符串）
JWT_SECRET=你的超长随机密钥_至少32位_random_string_here_change_this

# 运行模式
GIN_MODE=release
```

### 创建必要目录

```bash
mkdir -p backend/uploads logs docker/ssl
```

### 创建Nginx配置文件

```bash
# 如果没有nginx.conf，创建一个
cat > docker/nginx.conf << 'EOF'
server {
    listen 80;
    server_name localhost;
    
    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    
    location /api/ {
        proxy_pass http://backend:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    location /ws {
        proxy_pass http://backend:8080/ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 86400;
    }
}
EOF
```

### 启动服务

```bash
# 1. 构建并启动所有服务（后台运行）
docker compose up -d --build

# 2. 查看服务状态
docker compose ps

# 3. 查看日志
docker compose logs -f
```

**预期输出**：
```
NAME                STATUS          PORTS
chat-mysql-pro      running         0.0.0.0:3306->3306/tcp
chat-mongo-pro      running         0.0.0.0:27017->27017/tcp
chat-redis-pro      running         0.0.0.0:6379->6379/tcp
chat-backend-pro    running         0.0.0.0:8080->8080/tcp
chat-frontend-pro   running         0.0.0.0:80->80/tcp
```

### 等待服务启动

```bash
# 等待30秒让数据库初始化
sleep 30

# 检查后端是否正常运行
curl http://localhost:8080/health
```

**预期响应**：
```json
{"status":"ok"}
```

---

## 6. 配置域名（可选）

### 购买域名

推荐域名注册商：
- **Namecheap**: https://namecheap.com（便宜）
- **Cloudflare Registrar**: https://dash.cloudflare.com（隐私保护好）
- **GoDaddy**: https://godaddy.com（老牌）

### 配置DNS

1. **获取Oracle Cloud实例的公共IP**
   - 在Oracle Cloud控制台查看实例详情
   - 记录公共IP地址

2. **在域名提供商处添加DNS记录**

```
# A记录 - 主域名
Type: A
Name: @
Value: 你的公共IP
TTL: 3600

# A记录 - www子域名
Type: A
Name: www
Value: 你的公共IP
TTL: 3600
```

3. **等待DNS生效**
   - 通常5分钟到24小时
   - 使用 https://dnschecker.org 检查全球DNS传播

### 配置Nginx SSL（使用Let's Encrypt免费证书）

```bash
# 1. 安装Certbot
sudo apt install -y certbot python3-certbot-nginx

# 2. 获取SSL证书
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# 3. 按提示输入邮箱并同意条款

# 4. 自动续期测试
sudo certbot renew --dry-run
```

---

## 7. 故障排除

### 问题1：Docker权限错误

```bash
# 解决方案：重新登录
logout
# 重新连接SSH
ssh -i /path/to/key.pem ubuntu@你的IP
```

### 问题2：端口被占用

```bash
# 检查端口占用
sudo netstat -tlnp | grep :80
sudo netstat -tlnp | grep :8080

# 停止占用进程或修改docker-compose.yml中的端口
```

### 问题3：数据库连接失败

```bash
# 查看MySQL日志
docker compose logs mysql

# 等待MySQL完全启动（可能需要1-2分钟）
docker compose restart mysql
sleep 60
docker compose logs mysql
```

### 问题4：构建失败

```bash
# 清理Docker缓存
docker compose down
docker system prune -af

# 重新构建
docker compose build --no-cache
docker compose up -d
```

### 问题5：内存不足

ARM实例默认6GB内存可能不够：

```bash
# 检查内存使用
free -h

# 增加Swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### 问题6：无法访问服务

#### 检查防火墙

```bash
# 检查UFW状态
sudo ufw status

# 如果UFW启用，放行必要端口
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp
```

#### 检查Oracle Cloud安全列表

1. 进入Oracle Cloud控制台
2. 实例详情 → Subnet → 安全列表
3. 确认入站规则已添加

#### 检查服务状态

```bash
# 查看所有容器
docker ps -a

# 重启所有服务
docker compose restart

# 查看日志
docker compose logs --tail=100
```

---

## 🌐 访问地址

部署成功后：

| 服务 | 地址 |
|------|------|
| **前端** | `http://你的公共IP` |
| **API** | `http://你的公共IP:8080` |
| **API文档** | `http://你的公共IP:8080/api/v1` |
| **健康检查** | `http://你的公共IP:8080/health` |

---

## 📊 资源使用情况

Oracle Cloud Free Tier免费额度：

| 资源 | 免费额度 | 知信使用 |
|------|---------|---------|
| **ARM CPU** | 4 OCPU | 1 OCPU ✅ |
| **内存** | 24GB | 6GB ✅ |
| **存储** | 200GB | ~10GB ✅ |
| **公网带宽** | 不限 | - |

---

## 💰 费用

**完全免费！** Oracle Cloud Free Tier永久免费：
- ✅ 无月费
- ✅ 无年费
- ✅ 不会过期
- ✅ 不会自动收费

---

## 🔒 安全建议

1. **使用强密码** - 环境变量中的密码要足够复杂
2. **定期更新** - `sudo apt update && sudo apt upgrade`
3. **配置防火墙** - 只开放必要端口
4. **启用SSL** - 使用HTTPS
5. **定期备份** - 导出数据库重要数据

---

## 📚 相关文档

- **Oracle Cloud Free Tier**: https://www.oracle.com/cloud/free/
- **Docker官方文档**: https://docs.docker.com/
- **知信项目**: https://github.com/fjy123123/jinjuyi

---

## 🆘 获取帮助

如果在部署过程中遇到问题：

1. 查看本教程的故障排除部分
2. 查看Docker日志：`docker compose logs -f`
3. 检查系统日志：`sudo journalctl -u docker -f`
4. 在GitHub提交Issue

---

**祝你部署成功！🎉**