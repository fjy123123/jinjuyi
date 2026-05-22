@echo off
REM ===============================================
REM Chat System Pro - Windows 部署脚本
REM 支持: Windows 10/11, Windows Server
REM ===============================================

title Chat System Pro 部署工具

chcp 65001 >nul
echo.
echo ==============================================
echo   Chat System Pro - Windows 部署工具
echo ==============================================
echo.

REM 检查管理员权限
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [INFO] 请求管理员权限...
    powershell -Command "Start-Process '%~f0' -Verb RunAs"
    exit /b
)

REM 设置颜色
set "ESC=[92m"
set "END=[0m"

:MENU
cls
echo ==============================================
echo   Chat System Pro - Windows 部署工具
echo ==============================================
echo.
echo   1. 一键部署 (首次安装)
echo   2. 启动服务
echo   3. 停止服务
echo   4. 重启服务
echo   5. 查看状态
echo   6. 查看日志
echo   7. 初始化数据库
echo   8. 更新服务
echo   0. 退出
echo.
set /p choice=请选择操作 [0-8]: 

if "%choice%"=="1" goto INSTALL
if "%choice%"=="2" goto START
if "%choice%"=="3" goto STOP
if "%choice%"=="4" goto RESTART
if "%choice%"=="5" goto STATUS
if "%choice%"=="6" goto LOGS
if "%choice%"=="7" goto INITDB
if "%choice%"=="8" goto UPDATE
if "%choice%"=="0" goto EXIT
goto MENU

:INSTALL
echo.
echo [INFO] 开始一键部署...
echo.

REM 检查 Docker Desktop
docker version >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] 未检测到 Docker Desktop!
    echo.
    echo 请先安装 Docker Desktop for Windows
    echo 下载地址: https://www.docker.com/products/docker-desktop/
    echo.
    pause
    goto MENU
)

echo [INFO] Docker 已安装

REM 检查 .env 文件
if not exist ".env" (
    echo [INFO] 复制 .env.example 为 .env
    copy .env.example .env
)

echo.
echo [INFO] 启动服务...
docker-compose up -d

if %errorLevel% equ 0 (
    echo.
    echo ==============================================
    echo   部署成功!
    echo ==============================================
    echo.
    echo   访问地址:
    echo     前端: http://localhost
    echo     API:  http://localhost:8080
    echo.
    echo   查看日志: docker-compose logs -f
    echo.
) else (
    echo.
    echo [ERROR] 部署失败!
    echo.
)

pause
goto MENU

:START
echo.
echo [INFO] 启动服务...
docker-compose up -d
if %errorLevel% equ 0 (
    echo [SUCCESS] 服务已启动!
) else (
    echo [ERROR] 启动失败!
)
pause
goto MENU

:STOP
echo.
echo [INFO] 停止服务...
docker-compose down
echo [SUCCESS] 服务已停止
pause
goto MENU

:RESTART
echo.
echo [INFO] 重启服务...
docker-compose down
timeout /t 3 /nobreak >nul
docker-compose up -d
echo [SUCCESS] 服务已重启
pause
goto MENU

:STATUS
echo.
echo [INFO] 服务状态:
echo.
docker-compose ps
echo.
pause
goto MENU

:LOGS
echo.
echo [INFO] 查看日志 (按 Ctrl+C 退出):
echo.
docker-compose logs -f
goto MENU

:INITDB
echo.
echo [INFO] 初始化数据库...
docker-compose exec backend chat-system-pro init-db
echo [SUCCESS] 数据库初始化完成
pause
goto MENU

:UPDATE
echo.
echo [INFO] 更新服务...

echo [INFO] 拉取最新代码...
git pull

echo [INFO] 拉取最新镜像...
docker-compose pull

echo [INFO] 重启服务...
docker-compose down
docker-compose up -d

echo [SUCCESS] 更新完成!
pause
goto MENU

:EXIT
echo.
echo 再见!
echo.
timeout /t 1 /nobreak >nul
exit /b 0
