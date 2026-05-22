#!/bin/bash

# Chat System Pro - 打包脚本
# 用于生成完整的源代码分发包

set -e

# 定义变量
PROJECT_DIR="/workspace/chat-system-pro"
PROJECT_NAME="chat-system-pro"
VERSION="v1.0.0"
BUILD_DATE=$(date +%Y%m%d_%H%M%S)
PACKAGE_NAME="${PROJECT_NAME}-${VERSION}-${BUILD_DATE}"
OUTPUT_DIR="${PROJECT_DIR}/releases"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

echo "========================================"
echo "  Chat System Pro 打包工具"
echo "========================================"
echo ""

# 显示打包信息
echo "项目名称: $PROJECT_NAME"
echo "版本号: $VERSION"
echo "打包时间: $BUILD_DATE"
echo "输出目录: $OUTPUT_DIR"
echo ""

# 创建临时目录
TEMP_DIR=$(mktemp -d)
echo "创建临时目录: $TEMP_DIR"

# 复制项目文件
echo "复制项目文件..."
cp -r "$PROJECT_DIR" "$TEMP_DIR/$PACKAGE_NAME"

# 进入临时目录
cd "$TEMP_DIR"

# 清理不必要的文件
echo "清理构建产物和临时文件..."
find "$PACKAGE_NAME" -type f -name "*.exe" -delete 2>/dev/null || true
find "$PACKAGE_NAME" -type f -name "*.o" -delete 2>/dev/null || true
find "$PACKAGE_NAME" -type f -name "*.a" -delete 2>/dev/null || true
find "$PACKAGE_NAME" -type d -name "node_modules" -exec rm -rf {} + 2>/dev/null || true
find "$PACKAGE_NAME" -type d -name "dist" -exec rm -rf {} + 2>/dev/null || true
find "$PACKAGE_NAME" -type d -name ".git" -exec rm -rf {} + 2>/dev/null || true
find "$PACKAGE_NAME" -type f -name ".DS_Store" -delete 2>/dev/null || true
find "$PACKAGE_NAME" -type f -name "Thumbs.db" -delete 2>/dev/null || true
rm -rf "$PACKAGE_NAME/backend/chat-server" 2>/dev/null || true

# 创建部署说明文件
cat > "$PACKAGE_NAME/DEPLOY_GUIDE.md" << 'EOF'
# Chat System Pro 快速部署指南

## 快速开始

### 1. 环境要求

- Docker 20.10+
- Docker Compose 2.0+
- 2核CPU, 4GB内存, 50GB磁盘

### 2. 快速部署

```bash
# 解压文件
tar -xzf chat-system-pro-v1.0.0-xxxxxxxx.tar.gz
cd chat-system-pro-v1.0.0-xxxxxxxx

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件

# 启动服务
docker compose up -d

# 查看状态
docker compose ps

# 访问系统
# 前端: http://localhost
# API: http://localhost:8080
```

### 3. 详细文档

- 安装文档: INSTALL.md
- API文档: API.md
- 功能说明: FEATURES.md
- 部署文档: DEPLOY.md

### 4. 常见问题

Q: Docker启动失败？
A: 确保Docker Desktop已启动，或运行 `sudo systemctl start docker`

Q: 数据库连接失败？
A: 等待MySQL初始化完成，约30秒后重试

Q: 如何修改端口？
A: 编辑 docker-compose.yml 中的端口映射

Q: 如何备份数据？
A: 参考 INSTALL.md 中的数据备份章节

### 5. 获取帮助

- 查看详细文档
- 提交Issue: https://github.com/yourrepo/issues
- 邮箱支持: support@example.com

## 默认账号

- 管理员: admin / admin123 (首次登录后请修改)
- 测试用户: test / test123

## 技术支持

商业支持请联系: business@example.com
EOF

# 打包
echo ""
echo "正在打包..."
cd "$TEMP_DIR"

# 创建TAR.GZ包
tar -czf "$OUTPUT_DIR/${PACKAGE_NAME}.tar.gz" "$PACKAGE_NAME"

# 清理临时目录
rm -rf "$TEMP_DIR"

# 显示打包结果
echo ""
echo "========================================"
echo "  打包完成！"
echo "========================================"
echo ""
echo "输出文件:"
ls -lh "$OUTPUT_DIR"
echo ""
echo "使用说明:"
echo "1. 解压文件包: tar -xzf ${PACKAGE_NAME}.tar.gz"
echo "2. 按照 INSTALL.md 进行安装"
echo "3. 或参考 DEPLOY_GUIDE.md 快速开始"
echo ""
echo "打包完成时间: $(date)"
