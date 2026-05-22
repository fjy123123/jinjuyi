# 移动端快速开始

## 🚀 3步打包

### 1️⃣ 下载 HBuilderX

下载地址: https://www.dcloud.io/hbuilderx.html

### 2️⃣ 配置项目

1. 打开 `mobile/` 目录
2. 编辑 `config.js`，设置您的API地址:
```javascript
API_BASE_URL: 'http://your-server-ip:8080'
```

### 3️⃣ 云打包

1. 菜单: **发行 → 原生App-云打包**
2. 选择平台: **Android** 或 **iOS**
3. 填写信息 → 点击打包
4. 等待5-15分钟，下载APK/IPA文件

---

## 📱 APK打包(Android)

### Windows用户

```powershell
cd mobile
build-apk.bat
```

### Mac/Linux用户

```bash
cd mobile
chmod +x build-apk.sh
./build-apk.sh
```

---

## 🍎 iOS打包

### 要求

- Mac电脑
- Xcode 12+
- Apple Developer账号 ($99/年)

### 打包步骤

```bash
cd mobile
chmod +x build-ios.sh
./build-ios.sh
```

---

## ⚙️ 配置说明

### API地址配置

编辑 `config.js`:

```javascript
export default {
  development: {
    API_BASE_URL: 'http://your-server-ip:8080',
    WS_BASE_URL: 'ws://your-server-ip:8080/ws',
    UPLOAD_URL: 'http://your-server-ip:8080/api/v1/upload'
  },
  production: {
    API_BASE_URL: 'https://your-domain.com',
    WS_BASE_URL: 'wss://your-domain.com/ws',
    UPLOAD_URL: 'https://your-domain.com/api/v1/upload'
  }
}
```

### 修改App名称

编辑 `manifest.json`:

```json
{
  "name": "您的App名称"
}
```

---

## 📞 常见问题

### Q: 云打包失败?
A: 检查API地址是否可访问，网络是否正常。

### Q: 如何修改App图标?
A: HBuilderX → 发行 → 原生App-云打包 → 配置图标。

### Q: iOS需要什么证书?
A: Apple Developer账号，创建Development或Distribution证书。

---

## 📚 更多文档

- [完整打包指南](./README.md)
- [UniApp文档](https://uniapp.dcloud.net.cn/)
- [项目主文档](../README.md)
