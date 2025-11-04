#!/bin/bash

set -e  # 遇到错误立即退出

# 颜色输出
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

IMAGE_NAME="rt-manage:0.0.1"

echo -e "${BLUE}🚀 开始构建 Podman 镜像（优化版）: ${IMAGE_NAME}${NC}"
echo -e "${YELLOW}💡 此方案在本地构建，更快更省资源${NC}"
echo ""

# 1. 检查 Podman 是否安装
echo -e "${BLUE}📋 检查 Podman...${NC}"
if ! command -v podman &> /dev/null; then
    echo -e "${RED}❌ Podman 未安装，请先安装 Podman${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Podman 已安装${NC}"
echo ""

# 2. 构建前端
echo -e "${BLUE}📦 构建前端...${NC}"
cd frontend && pnpm install && pnpm run build
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 前端构建失败${NC}"
    exit 1
fi
cd ..

echo -e "${BLUE}📋 复制前端文件...${NC}"
mkdir -p internal/web/dist
cp -r frontend/dist/* internal/web/dist/
echo -e "${GREEN}✅ 前端构建完成${NC}"
echo ""

# 3. 构建后端
echo -e "${BLUE}🔨 构建后端...${NC}"
mkdir -p bin
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server cmd/server/main.go
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 后端构建失败${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 后端构建完成${NC}"
echo ""

# 4. 构建 Podman 镜像
echo -e "${BLUE}🐳 构建 Podman 镜像...${NC}"

# 自动检测架构
CURRENT_ARCH=$(uname -m)
OS_TYPE=$(uname -s)

echo -e "${BLUE}💻 当前系统: ${OS_TYPE} ${CURRENT_ARCH}${NC}"

# 自动选择目标平台
if [ "$CURRENT_ARCH" = "arm64" ] && [ "$OS_TYPE" = "Darwin" ]; then
    # Mac M1/M2，需要重新构建 linux/amd64 版本
    echo -e "${YELLOW}⚠️  检测到 Mac ARM64，需要为 linux/amd64 重新构建...${NC}"
    echo -e "${BLUE}   重新构建 linux/amd64 版本...${NC}"
    
    # 清理并重新构建 amd64 版本（完全静态链接）
    rm -f bin/server
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server cmd/server/main.go
    
    echo -e "${GREEN}✅ linux/amd64 版本构建完成${NC}"
fi

# 构建 Podman 镜像（使用简化版 Dockerfile）
echo -e "${BLUE}🔨 构建镜像...${NC}"
podman build -f Dockerfile.simple -t ${IMAGE_NAME} .

echo -e "${GREEN}✅ Podman 镜像构建完成${NC}"
echo ""

# 5. 显示镜像信息
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🎉 构建成功！${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}镜像名称:${NC} ${IMAGE_NAME}"
echo -e "${BLUE}镜像大小:${NC}"
podman images ${IMAGE_NAME} --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
echo ""
echo -e "${BLUE}运行命令:${NC}"
echo "  podman run -d -p 8080:8080 -v \$(pwd)/data:/app/data -v \$(pwd)/config:/app/config --name rt-manage ${IMAGE_NAME}"
echo ""
echo -e "${BLUE}查看日志:${NC}"
echo "  podman logs -f rt-manage"
echo ""
echo -e "${YELLOW}💡 提示：${NC}"
echo -e "${YELLOW}   - 此方案镜像更小（~50-80MB vs ~1GB+）${NC}"
echo -e "${YELLOW}   - 构建速度更快（无需在容器中安装依赖）${NC}"
echo -e "${YELLOW}   - 静态文件已嵌入到二进制文件中${NC}"
echo ""

