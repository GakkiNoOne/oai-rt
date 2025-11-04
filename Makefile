.PHONY: build run front backend clean clean-all rebuild test deps init install-frontend docker docker-simple podman podman-simple up-docker up-podman down help

# 默认目标：完整构建并运行
build: front backend
	@echo "✅ 构建完成！可执行文件：bin/server"
	@echo "🚀 启动服务器..."
	@./bin/server

# 开发运行（不重新构建）
run:
	go run cmd/server/main.go

# 构建前端
front:
	@echo "📦 开始构建前端..."
	@cd frontend && npm run build
	@echo "📋 正在复制前端文件..."
	@mkdir -p internal/web/dist
	@cp -r frontend/dist/* internal/web/dist/
	@echo "✅ 前端构建完成"

# 构建后端
backend:
	@echo "🔨 开始构建后端..."
	@mkdir -p bin
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server cmd/server/main.go
	@echo "✅ 后端构建完成"

# 安装前端依赖
install-frontend:
	@echo "📥 安装前端依赖..."
	@cd frontend && npm install
	@echo "✅ 前端依赖安装完成"

# 清理构建产物（保留数据）
clean:
	@echo "🧹 清理构建产物..."
	@rm -rf bin/
	@rm -rf frontend/dist/
	@rm -rf internal/web/dist/
	@echo "✅ 清理完成"

# 完全清理（包括数据）
clean-all: clean
	@echo "🧹 清理所有数据..."
	@rm -rf data/
	@echo "✅ 完全清理完成"

# 运行测试
test:
	@echo "🧪 运行测试..."
	@go test -v ./...

# 安装/更新 Go 依赖
deps:
	@echo "📥 更新 Go 依赖..."
	@go mod download
	@go mod tidy
	@echo "✅ Go 依赖更新完成"

# 初始化项目目录
init:
	@echo "📁 初始化项目目录..."
	@mkdir -p data config bin
	@echo "✅ 目录初始化完成"

# 快速重新构建（假设前端依赖已安装）
rebuild: clean front backend

# 构建 Docker 镜像（多阶段构建，不依赖本地环境）
docker:
	@./build-docker-multistage.sh

# 构建 Docker 镜像（简化版，本地构建后打包）
docker-simple:
	@./build-docker.sh

# 构建 Podman 镜像（简化版，本地构建后打包）
podman:
	@./build-podman.sh

# 构建 Podman 镜像（简化版，别名）
podman-simple:
	@./build-podman.sh

# 使用 docker-compose 启动
up-docker:
	@echo "🚀 使用 Docker Compose 启动服务..."
	@docker compose up -d

# 使用 podman-compose 启动
up-podman:
	@echo "🚀 使用 Podman Compose 启动服务..."
	@podman-compose up -d

# 停止服务
down:
	@echo "🛑 停止服务..."
	@if command -v docker &> /dev/null; then \
		docker compose down 2>/dev/null || true; \
	fi
	@if command -v podman-compose &> /dev/null; then \
		podman-compose down 2>/dev/null || true; \
	fi

# 帮助信息
help:
	@echo "RT-Manage 项目构建工具"
	@echo ""
	@echo "可用命令："
	@echo "  make                  - 完整构建并运行服务器 [默认]"
	@echo "  make run              - 开发模式运行"
	@echo "  make front            - 构建前端"
	@echo "  make backend          - 构建后端"
	@echo "  make install-frontend - 安装前端依赖"
	@echo "  make clean            - 清理构建产物"
	@echo "  make clean-all        - 完全清理（包括数据）"
	@echo "  make rebuild          - 清理后重新构建"
	@echo "  make docker           - 构建 Docker 镜像（多阶段，纯 Docker 环境）"
	@echo "  make docker-simple    - 构建 Docker 镜像（本地构建，快速）"
	@echo "  make podman           - 构建 Podman 镜像（本地构建，快速）"
	@echo "  make up-docker        - 使用 Docker Compose 启动服务"
	@echo "  make up-podman        - 使用 Podman Compose 启动服务"
	@echo "  make down             - 停止容器服务"
	@echo "  make test             - 运行测试"
	@echo "  make deps             - 更新 Go 依赖"
	@echo "  make init             - 初始化项目目录"
	@echo "  make help             - 显示此帮助信息"

