# Mobile App 打包指南

## 📋 目录
- [简介](#简介)
- [技术栈](#技术栈)
- [环境准备](#环境准备)
- [APK打包(Android)](#apk打包android)
- [iOS打包](#ios打包)
- [配置说明](#配置说明)

## 📖 简介

本项目使用 UniApp 框架开发，可以同时生成 Android APK 和 iOS App。

## 🛠️ 技术栈

- **框架**: UniApp (Vue 3)
- **构建工具**: HBuilderX / CLI
- **后端 API**: Go + Gin
- **实时通信**: WebSocket

## 📦 环境准备

### 1. 安装 Node.js

```bash
# Windows: 下载安装 https://nodejs.org/
# Linux/Mac:
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install nodejs
```

### 2. 安装 HBuilderX (推荐)

下载地址: https://www.dcloud.io/hbuilderx.html

### 3. 安装 Android SDK (APK打包)

```bash
# 使用 Android Studio 安装
# 或下载 SDK Command Line Tools
```

### 4. 安装 Xcode (iOS打包，仅Mac)

从 Mac App Store 下载安装 Xcode

## 📱 APK打包(Android)

### 方式1: 使用 HBuilderX (推荐)

1. **打开项目**
   ```bash
   # 在 HBuilderX 中打开 mobile/ 目录
   ```

2. **配置 API 地址**
   ```javascript
   // 修改 mobile/config.js
   export const API_BASE_URL = 'http://your-server-ip:8080'
   ```

3. **发行 → 原生App-云打包**
   ```
   菜单: 发行 → 原生App-云打包
   或快捷键: Ctrl + U
   ```

4. **填写打包信息**
   ```
   - App名称: 即时聊天
   - AppID: 使用DCloud分配的AppID
   - 证书: 使用自有证书或DCloud公用证书
   - 版本号: 1.0.0
   - 版本名称: v1.0.0
   ```

5. **点击打包**
   ```
   等待5-15分钟，下载生成的APK文件
   ```

### 方式2: 使用 CLI 本地打包

1. **安装依赖**
   ```bash
   cd mobile
   npm install
   npm install -g @vue/cli
   vue create -p dcloudio/uni-preset-vue my-app
   ```

2. **配置 manifest.json**
   ```json
   {
     "name": "即时聊天",
     "appid": "__UNI__YOURAPPID",
     "versionName": "1.0.0",
     "versionCode": "100"
   }
   ```

3. **编译到Android**
   ```bash
   npm run build:app-plus
   ```

4. **使用 Android Studio 打包**
   ```bash
   # 在 Android Studio 中打开 unpackage/release 目录
   # Build → Generate Signed Bundle / APK
   ```

### 方式3: 使用脚本自动打包

```bash
# 运行打包脚本
cd mobile
chmod +x build-apk.sh
./build-apk.sh
```

## 🍎 iOS打包

### 前置要求

- Mac电脑
- Xcode 12+
- Apple开发者账号 ($99/年)

### 步骤1: 配置证书

1. **登录 Apple Developer**
   ```
   https://developer.apple.com/account/
   ```

2. **创建证书**
   ```
   Certificates → + → iOS Development / Distribution
   ```

3. **创建 App ID**
   ```
   Identifiers → + → App IDs
   Bundle ID: com.yourcompany.chat
   ```

4. **创建 Provisioning Profile**
   ```
   Profiles → + → iOS App Development / App Store
   ```

### 步骤2: HBuilderX打包

1. **配置 iOS 证书**
   ```
   在 HBuilderX 中: 发行 → 原生App-云打包
   选择 iOS 平台
   上传证书和描述文件
   ```

2. **填写信息**
   ```
   Bundle ID: com.yourcompany.chat
   证书密码: 您的证书密码
   版本号: 1.0.0
   ```

3. **点击打包**
   ```
   等待打包完成，下载 .ipa 文件
   ```

### 步骤3: Xcode本地打包

1. **配置项目**
   ```
   cd mobile/unpackage/release
   open project.xcodeproj
   ```

2. **选择证书**
   ```
   Project → Signing & Capabilities
   Team: 选择您的开发者账号
   ```

3. **Archive**
   ```
   Product → Archive
   然后导出 App Store 或 Ad Hoc 版本
   ```

## ⚙️ 配置说明

### API地址配置

```javascript
// mobile/config.js
export default {
  API_BASE_URL: 'http://your-server-ip:8080',
  WS_BASE_URL: 'ws://your-server-ip:8080/ws',
  UPLOAD_URL: 'http://your-server-ip:8080/api/v1/upload'
}
```

### manifest.json 配置

```json
{
  "name": "即时聊天",
  "appid": "__UNI__YOURAPPID",
  "description": "商业级即时聊天应用",
  "versionName": "1.0.0",
  "versionCode": "100",
  "transformPx": false,
  "app-plus": {
    "usingComponents": true,
    "splashscreen": {
      "alwaysShowBeforeRender": true
    },
    "modules": {},
    "distribute": {
      "android": {
        "permissions": [
          "<uses-permission android:name=\"android.permission.CHANGE_NETWORK_STATE\"/>",
          "<uses-permission android:name=\"android.permission.MOUNT_UNMOUNT_FILESYSTEMS\"/>",
          "<uses-permission android:name=\"android.permission.VIBRATE\"/>",
          "<uses-permission android:name=\"android.permission.READ_LOGS\"/>",
          "<uses-permission android:name=\"android.permission.ACCESS_WIFI_STATE\"/>",
          "<uses-permission android:name=\"android.permission.ACCESS_NETWORK_STATE\"/>",
          "<uses-permission android:name=\"android.permission.CAMERA\"/>",
          "<uses-permission android:name=\"android.permission.RECORD_AUDIO\"/>"
        ]
      },
      "ios": {}
    }
  }
}
```

## 🎯 快速开始

### Windows用户

```powershell
# 1. 下载 HBuilderX
# 2. 打开 mobile/ 目录
# 3. 配置 API 地址
# 4. 发行 → 原生App-云打包
```

### Mac用户

```bash
# 1. 下载 HBuilderX
# 2. 打开 mobile/ 目录
# 3. 配置 API 地址
# 4. 发行 → 原生App-云打包
```

## 📞 常见问题

### Q: 云打包失败怎么办?
A: 检查网络连接，确保API地址可访问，或使用本地打包。

### Q: 如何修改App图标?
A: 在 HBuilderX 中: 发行 → 原生App-云打包 → 配置图标和启动图。

### Q: iOS需要什么证书?
A: 需要 Apple Developer 账号，创建 Development 或 Distribution 证书。

## 📚 更多文档

- [UniApp官方文档](https://uniapp.dcloud.net.cn/)
- [HBuilderX使用指南](https://hx.dcloud.net.cn/)
