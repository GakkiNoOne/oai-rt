# OAI RT 管理系统
一个用于管理和自动刷新 OpenAI Refresh Token 的系统。

<div align="center">
    <img width="5100" height="2410" alt="Image" src="https://github.com/user-attachments/assets/61e8cada-593d-4466-a678-7c10a5b21939" /></div>
<br/>


## 快速部署

### 1. 准备配置文件

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
  # SQLite 配置（默认）
  type: "sqlite"
  database: "./data/tokens.db"
  
  # MySQL 配置（如需使用请取消注释并注释掉 SQLite 配置）
  # type: "mysql"
  # host: "localhost"
  # port: 3306
  # user: "root"
  # password: "your_mysql_password"
  # database: "rt_manage"
  
  # 通用配置
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

### 2. 使用 Docker Compose 启动

创建 `docker-compose.yml` 文件：

```yaml
services:
  rt-manage:
    image: ghcr.io/gakkinoone/oai-rt:latest
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

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `server.host` | string | `0.0.0.0` | 服务监听地址 |
| `server.port` | int | `8080` | 服务监听端口 |
| `server.mode` | string | `debug` | 运行模式：`debug` / `release` / `test` |

### 数据库配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `database.type` | string | `sqlite` | 数据库类型：`sqlite` 或 `mysql` |
| `database.database` | string | `./data/tokens.db` | SQLite 数据库文件路径（type=sqlite 时使用） |
| `database.host` | string | `localhost` | MySQL 主机地址（type=mysql 时使用） |
| `database.port` | int | `3306` | MySQL 端口（type=mysql 时使用） |
| `database.user` | string | `root` | MySQL 用户名（type=mysql 时使用） |
| `database.password` | string | - | MySQL 密码（type=mysql 时使用） |
| `database.table_prefix` | string | - | 表名前缀（可选） |
| `database.max_idle_conns` | int | `10` | 最大空闲连接数 |
| `database.max_open_conns` | int | `100` | 最大打开连接数 |
| `database.conn_max_lifetime` | int | `3600` | 连接最大生命周期（秒） |

### 认证配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `auth.username` | string | `admin` | 管理后台登录用户名 |
| `auth.password` | string | `admin123` | 管理后台登录密码 |
| `auth.jwt_secret` | string | - | JWT 签名密钥，**生产环境必须修改**, 可以使用https://jwtsecrets.com/去生成 |
| `auth.jwt_expire_hours` | int | `240` | JWT 过期时间（小时） |
| `auth.api_secret` | string | - | 对外 API 密钥，用于 API 接口鉴权 |

## 访问地址

- 管理后台：http://localhost:8080
- 健康检查：http://localhost:8080/health
- API 文档：管理后台登录后可查看

## 对外 API 使用

所有对外 API 接口需要在 HTTP Header 中添加 API Secret：`X-API-Secret`

### 1. 刷新 RT 并获取 AT

支持使用 `biz_id` 或 `email` 查询：

```bash
# 使用 biz_id 查询
curl -X POST http://localhost:8080/public-api/refresh \
  -H "Content-Type: application/json" \
  -H "X-API-Secret: my-api-secret-2025" \
  -d '{
    "biz_id": "user001"
  }'

# 使用 email 查询
curl -X POST http://localhost:8080/public-api/refresh \
  -H "Content-Type: application/json" \
  -H "X-API-Secret: my-api-secret-2025" \
  -d '{
    "email": "user@example.com"
  }'
```

### 2. 获取 AT（不刷新）

支持使用 `biz_id` 或 `email` 查询：

```bash
# 使用 biz_id 查询
curl -X POST http://localhost:8080/public-api/get-at \
  -H "Content-Type: application/json" \
  -H "X-API-Secret: my-api-secret-2025" \
  -d '{
    "biz_id": "user001"
  }'

# 使用 email 查询
curl -X POST http://localhost:8080/public-api/get-at \
  -H "Content-Type: application/json" \
  -H "X-API-Secret: my-api-secret-2025" \
  -d '{
    "email": "user@example.com"
  }'
```

> **注意**：代理、自动刷新等 OpenAI 相关配置可在 Web 管理界面的"配置管理"页面进行设置。

## 数据库表结构

完整的表结构 SQL 文件：[resource/table.sql](resource/table.sql)

## License

MIT
