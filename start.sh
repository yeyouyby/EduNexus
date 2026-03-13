#!/bin/bash

set -e

# 设置控制台颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}       Project EduNexus - 一键启动脚本          ${NC}"
echo -e "${BLUE}================================================${NC}"

# 检查依赖环境
echo -e "${YELLOW}检查依赖环境...${NC}"
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}错误: 未找到 'go' 命令，请先安装 Go 语言环境。${NC}"
    exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
    echo -e "${RED}错误: 未找到 'npm' 命令，请先安装 Node.js 和 npm。${NC}"
    exit 1
fi

# 检查项目目录
if [ ! -d "edunexus" ]; then
    if [ ! -f "wails.json" ]; then
        echo -e "${RED}错误: 请在项目根目录或 edunexus 目录下运行此脚本。${NC}"
        exit 1
    fi
else
    cd edunexus
fi

echo -e "\n${GREEN}[1/2] 正在构建前端 (Building Frontend)...${NC}"
cd frontend
npm install
npm run build
cd ..

echo -e "\n${GREEN}[2/2] 正在启动 EduNexus 客户端...${NC}"
echo -e "${YELLOW}正在使用 go run . 启动... (按 Ctrl+C 停止)${NC}"

# 运行 wails 项目
# 因为前端已经 build 到 dist，所以可以直接用 go run . 运行基于 wails 的嵌入应用
# 也可以使用 wails dev，但 go run . 不需要全局安装 wails cli
go run .
