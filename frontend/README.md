# RT Pool 前端管理系统

基于 React + Umi + Ant Design 的 RT管理系统前端项目。

## 功能特性

### 1. RT 管理页面
- ✅ RT 列表展示（分页、排序）
- ✅ 搜索功能（按名称、状态、创建日期）
- ✅ 单条操作
  - 查看详情
  - 编辑 RT
  - 删除 RT
  - 单条刷新（测活）
- ✅ 批量操作
  - 批量选择
  - 批量删除
  - 批量刷新（测活）
  - 批量导入
- ✅ 状态管理（启用/禁用）
- ✅ 代理配置

### 2. 配置管理页面
- ✅ 代理服务器配置（支持多个代理）
- ✅ 自动刷新配置
  - 启用/禁用自动刷新
  - 刷新间隔设置
- ✅ 请求配置
  - 最大并发请求数
  - 请求超时时间
  - 失败重试次数
- ✅ 日志配置
  - 日志保留天数
- ✅ 环境变量展示（只读）

### 3. 用户认证
- ✅ 登录页面
- ✅ 退出登录
- ✅ Token 管理

## 技术栈

- **框架**: React 18 + Umi 4
- **UI 组件库**: Ant Design 5
- **HTTP 客户端**: Axios
- **日期处理**: Day.js
- **语言**: TypeScript

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API 接口定义
│   │   ├── rts.ts       # RT 管理 API
│   │   └── configs.ts   # 配置管理 API
│   ├── layouts/          # 布局组件
│   │   └── index.tsx    # 主布局（侧边栏+顶栏）
│   ├── pages/            # 页面组件
│   │   ├── Login/       # 登录页
│   │   ├── RTs/         # RT 管理页
│   │   └── Configs/     # 配置管理页
│   ├── utils/            # 工具函数
│   │   └── request.ts   # Axios 封装
│   └── global.css        # 全局样式
├── .umirc.ts             # Umi 配置
├── package.json          # 依赖配置
└── tsconfig.json         # TypeScript 配置
```

## 开发指南

### 安装依赖

```bash
cd frontend
pnpm install
```

### 启动开发服务器

```bash
pnpm dev
```

访问 http://localhost:8000

### 构建生产版本

```bash
pnpm build
```

构建产物将生成在 `dist` 目录。

## API 接口

前端通过 `/api` 前缀访问后端接口，开发环境下会代理到 `http://localhost:8080`。

### RT 管理接口
- `GET /api/rts` - 获取 RT 列表
- `POST /api/rts` - 创建 RT
- `PUT /api/rts/:id` - 更新 RT
- `DELETE /api/rts/:id` - 删除 RT
- `POST /api/rts/:id/refresh` - 刷新单个 RT
- `POST /api/rts/batch/import` - 批量导入 RT
- `POST /api/rts/batch/refresh` - 批量刷新 RT
- `POST /api/rts/batch/delete` - 批量删除 RT

### 配置管理接口
- `GET /api/configs/system` - 获取系统配置
- `POST /api/configs/system` - 保存系统配置
- `GET /api/configs/proxy-list` - 获取代理列表

### 认证接口
- `POST /api/auth/login` - 登录
- `POST /api/auth/logout` - 退出登录

## 配置说明

### 代理配置
支持配置多个代理服务器，系统会随机选择使用：
- 支持协议：`http://`、`https://`、`socks5://`
- 格式示例：`http://127.0.0.1:7890`

### 自动刷新
- 可配置是否启用自动刷新
- 刷新间隔：5-1440 分钟
- 自动禁用失效的 RT

### 请求配置
- 最大并发请求数：1-100
- 请求超时时间：5-300 秒
- 失败重试次数：0-10 次

## 注意事项

1. 批量导入时会自动去重，跳过已存在的 RT Token
2. 批量刷新失败的 RT 会被自动禁用
3. 删除操作不可恢复，请谨慎操作
4. 修改配置后会自动更新相关 RT 的配置
