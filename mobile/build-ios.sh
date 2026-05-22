#!/bin/bash

# iOS 打包脚本 (仅Mac)
# 使用说明: ./build-ios.sh

set -e

echo "=============================================="
echo "  Chat System Pro iOS 打包"
echo "=============================================="
echo ""

# 检查是否为Mac
if [ "$(uname)" != "Darwin" ]; then
    echo "[错误] iOS 打包需要在 Mac 上运行
    exit 1
fi

# 配置
APP_NAME="即时聊天"
VERSION="1.0.0"
BUNDLE_ID="com.chatsystem.pro"
OUTPUT_DIR="./output/ios"

# 创建输出目录
mkdir -p $OUTPUT_DIR

echo "[1/5] 检查环境..."

# 检查 Xcode
if ! command -v xcodebuild &> /dev/null; then
    echo "[错误] 未找到 Xcode，请从 Mac App Store 安装 Xcode"
    exit 1
fi

echo "[成功] Xcode 检查通过"

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "[错误] 未找到 Node.js，请先安装"
    exit 1
fi

echo "[成功] Node.js 检查通过"

echo ""
echo "[2/5] 配置项目..."

# 检查配置文件
if [ ! -f "manifest.json" ]; then
    echo "[错误] 未找到 manifest.json"
    exit 1
fi

echo "[成功] manifest.json 存在"

echo ""
echo "[3/5] 证书说明..."
echo ""
echo "iOS 打包需要 Apple Developer 账号"
echo ""
echo "请确保您已完成以下步骤:"
echo ""
echo "1. 注册 Apple Developer 账号 (\$99/年)"
echo "   访问: https://developer.apple.com/"
echo ""
echo "2. 创建证书:"
echo "   - iOS Development (用于测试"
echo "   - iOS Distribution (用于发布)"
echo ""
echo "3. 创建 App ID:"
echo "   Bundle ID: $BUNDLE_ID"
echo ""
echo "4. 创建 Provisioning Profile"
echo ""

echo ""
read -p "是否继续？(y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

echo ""
echo "[4/5] HBuilderX 打包..."
echo ""
echo "请按照以下步骤操作:"
echo ""
echo "1. 下载 HBuilderX for Mac"
echo "   下载地址: https://www.dcloud.io/hbuilderx.html"
echo ""
echo "2. 使用 HBuilderX 打开当前目录"
echo "   cd $(pwd)"
echo ""
echo "3. 配置 API 地址"
echo "   编辑 config.js，修改 API_BASE_URL 为您的服务器地址"
echo ""
echo "4. 云打包 iOS"
echo "   菜单: 发行 → 原生App-云打包 → 选择 iOS"
echo ""
echo "5. 填写打包信息:"
echo "   - App名称: $APP_NAME"
echo "   - Bundle ID: $BUNDLE_ID"
echo "   - 版本号: $VERSION"
echo "   - 上传证书和描述文件"
echo ""
echo "6. 等待打包完成，下载 .ipa 文件"
echo ""

echo "[5/5] Xcode 本地打包（可选）..."
echo ""
echo "如果您想本地打包，请确保已安装 Xcode Command Line Tools"
echo "然后运行:"
echo ""
echo "  1. npm install"
echo "  2. npm run build:app-plus"
echo "  3. 使用 Xcode 打开 unpackage/release 目录"
echo "  4. Product → Archive"
echo ""

echo ""
echo "=============================================="
echo "  iOS 打包准备完成！"
echo "=============================================="
echo ""
echo "下一步: 打开 HBuilderX 进行云打包"
echo ""
