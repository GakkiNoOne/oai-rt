# 多阶段构建 Dockerfile - 完全在 Docker 环境中构建
# 不依赖本地环境，适合 CI/CD 和纯 Docker 环境

# 第一阶段：构建前端
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

# 复制前端依赖文件
COPY frontend/package*.json ./

# 安装前端依赖
RUN npm install

# 复制前端源代码
COPY frontend/ ./

# 构建前端
RUN npm run build

# 第二阶段：构建后端
FROM golang:1.24-alpine AS backend-builder

# 安装构建工具
RUN apk add --no-cache git

WORKDIR /build

# 设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

# 复制 Go 依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制后端源代码
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# 从前端构建阶段复制前端产物
COPY --from=frontend-builder /frontend/dist ./internal/web/dist

# 构建后端（完全静态编译）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server cmd/server/main.go

# 第三阶段：运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=backend-builder /build/server .

# 创建配置目录和数据目录
RUN mkdir -p ./config ./data

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./server"]

