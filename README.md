# OAI RT 管理系统

一个用于管理和自动刷新 OpenAI Refresh Token 的系统。

## 快速部署

### 1. 构建镜像

```bash
docker build -t rt-manage:0.0.1 .
```

### 2. 准备配置文件

在项目目录创建 `config` 目录和 `config.yml` 配置文件：

```bash
mkdir -p config data
```

`config/config.yml` 配置示例：

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"  # debug, release, test

database:
  type: "sqlite"  # sqlite 或 mysql
  database: "./data/tokens.db"  # SQLite 数据库文件路径
  table_prefix: ""
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

auth:
  username: "admin"  # 管理后台登录用户名
  password: "admin123"  # 管理后台登录密码
  jwt_secret: "your-secret-key-change-this-in-production"  # JWT 签名密钥（生产环境必须修改）
  jwt_expire_hours: 240  # JWT 过期时间（小时）
  api_secret: "my-api-secret-2025"  # 对外 API 密钥
```

### 3. 使用 Docker Compose 启动

创建 `docker-compose.yml` 文件：

```yaml
services:
  rt-manage:
    image: rt-manage:0.0.1
    container_name: rt-manage
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
      - ./config:/app/config  # 配置文件目录，需将 config.yml 放在 ./config/ 目录下
    environment:
      - TZ=Asia/Shanghai
```

启动服务：

```bash
docker-compose up -d
```

查看日志：

```bash
docker-compose logs -f
```

## 配置说明

### 服务器配置

- `server.host`: 服务监听地址，默认 `0.0.0.0`
- `server.port`: 服务监听端口，默认 `8080`
- `server.mode`: 运行模式，可选 `debug`/`release`/`test`

### 数据库配置

**SQLite 模式**（推荐用于单机部署）：
- `database.type`: 设置为 `sqlite`
- `database.database`: 数据库文件路径，如 `./data/tokens.db`

**MySQL 模式**（推荐用于生产环境）：
- `database.type`: 设置为 `mysql`
- `database.host`: MySQL 主机地址
- `database.port`: MySQL 端口，默认 `3306`
- `database.user`: MySQL 用户名
- `database.password`: MySQL 密码
- `database.database`: 数据库名称

**通用配置**：
- `database.table_prefix`: 表名前缀（可选）
- `database.max_idle_conns`: 最大空闲连接数，默认 `10`
- `database.max_open_conns`: 最大打开连接数，默认 `100`
- `database.conn_max_lifetime`: 连接最大生命周期（秒），默认 `3600`

### 认证配置

- `auth.username`: 管理后台登录用户名
- `auth.password`: 管理后台登录密码
- `auth.jwt_secret`: JWT 签名密钥，**生产环境必须修改**
- `auth.jwt_expire_hours`: JWT 过期时间（小时）
- `auth.api_secret`: 对外 API 密钥，用于 API 接口鉴权

## 访问地址

- 管理后台：http://localhost:8080
- 健康检查：http://localhost:8080/health
- API 文档：管理后台登录后可查看

## 对外 API 使用

所有对外 API 接口需要在 HTTP Header 中添加 API Secret：

```bash
# 刷新 RT 并获取 AT
curl -X POST http://localhost:8080/public-api/refresh \
  -H "Content-Type: application/json" \
  -H "X-API-Secret: my-api-secret-2025" \
  -d '{"biz_id": "user001"}'

# 获取 AT（不刷新）
curl -X POST http://localhost:8080/public-api/get-at \
  -H "Content-Type: application/json" \
  -H "X-API-Secret: my-api-secret-2025" \
  -d '{"biz_id": "user001"}'
```

> **注意**：代理、自动刷新等 OpenAI 相关配置可在 Web 管理界面的"配置管理"页面进行设置。

## License

MIT
