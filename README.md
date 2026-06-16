# echo-framework

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Echo](https://img.shields.io/badge/Echo-v4-7A5AF8?logo=go&logoColor=white)](https://echo.labstack.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个基于 Go + Echo 框架的后端基础项目模板，提供日志、配置加载、数据库连接、JWT 中间件和系统用户模块的基础能力，适合快速搭建业务服务。

## 项目特点

- 基于 Echo v4 的 Web 服务框架
- 支持 YAML 配置加载（优先使用项目根目录的 `config.default.yaml`，否则回退到内置默认配置）
- 内置日志、数据库连接池与监控能力
- 提供系统用户基础 CRUD 接口
- 模块化目录结构，便于扩展业务功能

## 技术栈

- Go 1.25
- Echo v4
- GORM
- MySQL
- Viper
- Zap
- JWT

## 项目结构

```text
.
├── cmd/                  # 启动入口
├── internal/
│   ├── app/              # 应用启动与路由注册
│   ├── config/           # 配置加载与默认配置
│   ├── infrastructure/   # DB、日志、监控基础设施
│   ├── middleware/       # JWT、日志中间件
│   ├── model/            # 数据模型
│   └── module/           # 业务模块
└── pkg/                  # 公共工具包
```

## 快速开始

### 1. 环境准备

- Go 1.25+
- MySQL 数据库
- 可选：创建项目根目录下的 `config.yaml` 覆盖默认配置

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 本地启动

```bash
go run ./cmd
```

默认启动后会监听 `:8080`（如配置文件中已修改则以配置文件为准）。

### 4. Docker 启动

构建镜像：

```bash
docker build -t echo-framework .
```

运行容器：

```bash
docker run --rm -p 8000:8000 echo-framework
```

也可以使用 `docker-compose`：

```bash
docker compose up --build
```

### 5. 验证项目

```bash
go test ./...
```

## 默认配置说明

项目默认配置位于：

- `internal/config/config.default.yaml`

如需自定义环境配置，可在项目根目录创建 `config.yaml`，启动时会优先加载该文件。

## API 示例

当前已提供系统用户模块接口：

- GET    /api/sys/users
- GET    /api/sys/users/:id
- POST   /api/sys/users
- PUT    /api/sys/users/:id
- DELETE /api/sys/users/:id

## 主要模块

- `internal/app`：应用初始化、路由注册与服务启动
- `internal/config`：配置加载逻辑与默认参数
- `internal/infrastructure/db`：数据库连接与连接池设置
- `internal/infrastructure/log`：日志初始化与输出配置
- `internal/module/system/user`：用户管理接口实现

## 许可证

本项目采用 MIT License。
