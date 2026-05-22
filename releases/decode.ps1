# PowerShell解码脚本
# 使用方法：
# 1. 从服务器复制 Base64 编码内容
# 2. 粘贴到下面的 $base64 变量中（替换为空字符串）
# 3. 运行此脚本

$base64 = ""  # 把服务器输出的Base64字符串粘贴在这里
$outputFile = "D:\AI\chat-system-pro.zip"

if ([string]::IsNullOrEmpty($base64)) {
    Write-Host "错误：请先把Base64字符串粘贴到脚本的 `$base64` 变量中！" -ForegroundColor Red
    exit 1
}

Write-Host "正在解码Base64并保存到: $outputFile" -ForegroundColor Green

try {
    $bytes = [System.Convert]::FromBase64String($base64)
    [System.IO.File]::WriteAllBytes($outputFile, $bytes)
    Write-Host "✅ 保存成功！" -ForegroundColor Green
    Write-Host "文件位置: $outputFile" -ForegroundColor Cyan
} catch {
    Write-Host "❌ 错误: $_" -ForegroundColor Red
}
