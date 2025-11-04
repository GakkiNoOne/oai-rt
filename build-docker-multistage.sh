#!/bin/bash

set -e  # 遇到错误立即退出

# 颜色输出
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

IMAGE_NAME="rt-manage:0.0.1"

echo -e "${BLUE}🚀 开始构建 Docker 镜像（多阶段构建）: ${IMAGE_NAME}${NC}"
echo -e "${YELLOW}💡 完全在 Docker 环境中构建，不依赖本地环境${NC}"
echo ""

# 1. 检查 Docker 是否运行
echo -e "${BLUE}📋 检查 Docker...${NC}"
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker 未运行，请先启动 Docker${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker 运行正常${NC}"
echo ""

# 2. 构建 Docker 镜像（多阶段构建）
echo -e "${BLUE}🐳 开始多阶段构建...${NC}"
echo -e "${YELLOW}   阶段 1/3: 构建前端（Node.js）${NC}"
echo -e "${YELLOW}   阶段 2/3: 构建后端（Golang）${NC}"
echo -e "${YELLOW}   阶段 3/3: 创建运行镜像（Alpine）${NC}"
echo ""

# 自动检测架构并选择构建平台
CURRENT_ARCH=$(uname -m)
OS_TYPE=$(uname -s)

echo -e "${BLUE}💻 当前系统: ${OS_TYPE} ${CURRENT_ARCH}${NC}"

# 使用多阶段构建专用的 .dockerignore（保留源代码）
if [ "$CURRENT_ARCH" = "arm64" ] && [ "$OS_TYPE" = "Darwin" ]; then
    # Mac M1/M2，构建 linux/amd64 镜像
    echo -e "${YELLOW}⚠️  检测到 Mac ARM64，将构建 linux/amd64 镜像${NC}"
    docker buildx build --platform linux/amd64 -t ${IMAGE_NAME} --load .
else
    # Linux 或其他系统
    docker build -t ${IMAGE_NAME} .
fi

echo -e "${GREEN}✅ Docker 镜像构建完成${NC}"
echo ""

# 3. 显示镜像信息
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🎉 构建成功！${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}镜像名称:${NC} ${IMAGE_NAME}"
echo -e "${BLUE}镜像大小:${NC}"
docker images ${IMAGE_NAME} --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
echo ""
echo -e "${BLUE}运行命令:${NC}"
echo "  docker run -d -p 8080:8080 -v \$(pwd)/data:/app/data -v \$(pwd)/config:/app/config --name rt-manage ${IMAGE_NAME}"
echo ""
echo -e "${BLUE}查看日志:${NC}"
echo "  docker logs -f rt-manage"
echo ""
echo -e "${YELLOW}💡 提示：${NC}"
echo -e "${YELLOW}   - 此方案完全在 Docker 中构建${NC}"
echo -e "${YELLOW}   - 不依赖本地 Node.js 和 Go 环境${NC}"
echo -e "${YELLOW}   - 适合 CI/CD 和纯 Docker 环境${NC}"
echo ""

