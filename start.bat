@echo off
setlocal enabledelayedexpansion

:: 启用 ANSI 转义序列支持 (Windows 10+)
for /F "delims=#" %%E in ('"prompt #$E# & for %%E in (1) do rem"') do set "ESC=%%E"

:: 设置控制台颜色
set "GREEN=%ESC%[0;32m"
set "BLUE=%ESC%[0;34m"
set "YELLOW=%ESC%[1;33m"
set "RED=%ESC%[0;31m"
set "NC=%ESC%[0m"

echo  %BLUE%================================================%NC%
echo         Project EduNexus - 一键启动脚本
echo  %BLUE%================================================%NC%

:: 检查依赖环境
echo %YELLOW%检查依赖环境...%NC%

where go >nul 2>nul
if %errorlevel% neq 0 (
    echo %RED%错误: 未找到 'go' 命令，请先安装 Go 语言环境。%NC%
    exit /b 1
)

where npm >nul 2>nul
if %errorlevel% neq 0 (
    echo %RED%错误: 未找到 'npm' 命令，请先安装 Node.js 和 npm。%NC%
    exit /b 1
)

:: 检查项目目录
if not exist "edunexus\" (
    if not exist "wails.json" (
        echo %RED%错误: 请在项目根目录或 edunexus 目录下运行此脚本。%NC%
        exit /b 1
    )
) else (
    cd edunexus
)

echo.
echo %GREEN%[1/2] 正在构建前端 (Building Frontend)...%NC%
cd frontend
call npm install
call npm run build
cd ..

echo.
echo %GREEN%[2/2] 正在启动 EduNexus 客户端...%NC%
echo %YELLOW%正在使用 go run . 启动... (按 Ctrl+C 停止)%NC%

:: 运行 wails 项目
go run .

endlocal
