#!/bin/bash

# APK 打包脚本
# 使用说明: ./build-apk.sh

set -e

echo "=============================================="
echo "  Chat System Pro APK 打包"
echo "=============================================="
echo ""

# 配置
APP_NAME="即时聊天"
VERSION="1.0.0"
BUILD_NUMBER=$(date +%Y%m%d%H%M)
OUTPUT_DIR="./output/apk"

# 创建输出目录
mkdir -p $OUTPUT_DIR

echo "[1/4] 检查环境..."

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "错误: 未找到 Node.js，请先安装"
    exit 1
fi

echo "✓ Node.js 检查通过"

# 检查 HBuilderX CLI
echo "提示: 本脚本需要配合 HBuilderX 使用"
echo ""

echo "[2/4] 配置项目..."

# 检查配置文件
if [ ! -f "manifest.json" ]; then
    echo "错误: 未找到 manifest.json"
    exit 1
fi

echo "✓ manifest.json 存在"

# 检查 config.js
if [ ! -f "config.js" ]; then
    echo "警告: 未找到 config.js，使用默认配置"
else
    echo "✓ config.js 存在"
fi

echo ""
echo "[3/4] 构建说明..."
echo ""
echo "请按照以下步骤操作："
echo ""
echo "1. 下载并安装 HBuilderX"
echo "   下载地址: https://www.dcloud.io/hbuilderx.html"
echo ""
echo "2. 使用 HBuilderX 打开当前目录"
echo "   cd $(pwd)"
echo ""
echo "3. 配置 API 地址"
echo "   编辑 config.js，修改 API_BASE_URL 为您的服务器地址"
echo ""
echo "4. 云打包 APK"
echo "   菜单: 发行 → 原生App-云打包 → 选择 Android"
echo ""
echo "5. 填写打包信息"
echo "   - App名称: $APP_NAME"
echo "   - 版本号: $VERSION"
echo "   - 版本名称: v$VERSION"
echo ""
echo "6. 等待打包完成，下载 APK"
echo ""

echo "[4/4] 本地打包（可选）..."
echo ""
echo "如果您想本地打包，请确保已安装 Android SDK"
echo "然后运行:"
echo ""
echo "  1. npm install"
echo "  2. npm run build:app-plus"
echo "  3. 使用 Android Studio 打开 unpackage/release 目录"
echo "  4. Build → Generate Signed Bundle / APK"
echo ""

echo ""
echo "=============================================="
echo "  APK 打包准备完成！"
echo "=============================================="
echo ""
echo "下一步: 打开 HBuilderX 进行云打包"
echo ""
