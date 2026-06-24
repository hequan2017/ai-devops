<div align="center">
  <h1>AI运维平台（AI-DevOps）</h1>
  <p>面向物理服务器 / Docker / Kubernetes 的一体化运维管理平台</p>
</div>

<br/>

## 简介

AI运维平台是一套面向基础设施的运维管理系统，覆盖**物理服务器资产（戴尔 / 华为）、Docker 引擎、Kubernetes 集群**的统一管理，并提供**代码发版工作流审批**。后端基于 Go（Gin + GORM），前端基于 Vue3 + Element Plus，**默认使用 SQLite 文件数据库，开箱即用**。

> ⚠️ **编译环境要求**：本项目依赖 casbin v3、pgx、x/crypto、mcp-go 等现代库，其源码使用了 go 1.21–1.24 的标准库（`cmp/maps/slices/iter/crypto/sha3` 等），**必须使用 go 1.24 编译**。go 1.18/1.20 无法编译。Docker 镜像 `golang:alpine` 自带 go 1.24，可直接构建。

## 功能模块

### 1. 服务器管理（戴尔 / 华为）
- 资产台账：厂商、型号、SN、机房机柜、U 位、CPU/内存/磁盘、操作系统、负责人
- 带外管理 IP（iDRAC / iBMC）与业务 IP
- **IPMI 电源管理**：状态查询 / 开机 / 软关机 / 强制关机 / 硬重启（依赖运行环境 `ipmitool`）
- 运行状态可视化：运行中 / 已关机 / 维护中 / 故障

### 2. Docker 管理
- 引擎接入点管理（unix socket / 远程 tcp）+ 连接测试
- 容器：列表 / 启动 / 停止 / 重启 / 删除 / 日志
- 镜像、网络、数据卷浏览

### 3. Kubernetes 管理
- 集群接入点（API Server + Token，支持跳过自签 TLS 校验）+ 连接测试
- Namespace / Pod（日志 / 删除）/ Deployment / Service / Node 浏览

### 4. 发版工作流（类 Jenkins 审批）
- 发版记录：项目 / 版本 / 分支 / 环境 / 发布脚本
- 状态机：草稿 → 待审批 → 通过/拒绝 → 执行 → 完成/失败
- 多角色审批，执行阶段运行发布脚本并记录结果

### 5. 平台基础能力
RBAC 角色权限 + 动态菜单、用户/角色/菜单/API 管理、操作审计、登录日志、字典与参数管理、代码生成器等。

### 🚧 规划中
- **WebSSH 终端**：浏览器内通过密钥/密码登录服务器（xterm.js + WebSocket + SSH），需 go1.24 环境调试，下一阶段交付。

## 技术栈

| 端 | 技术 |
| --- | --- |
| 后端 | Go 1.24 · Gin · GORM · SQLite(`glebarez/sqlite`) · Docker/K8s 经 REST 直连（零重型 SDK） |
| 前端 | Vue 3 · Vite · Element Plus · Pinia |
| 鉴权 | JWT · Casbin |

## 目录结构

```
ai-devops/server/
├── model/biz/        # 业务模型：Server / DockerHost / K8sCluster / Release
├── service/biz/      # 业务逻辑 + docker_client/k8s_client/ipmi 轻量客户端
├── api/v1/biz/       # 业务接口
├── router/biz/       # 业务路由
├── source/system/    # 初始化数据（菜单 / API / 权限，开箱即用）
└── config.yaml       # 配置（db-type 默认 sqlite）
ai-devops/web/src/
├── api/              # server.js / docker.js / k8s.js / release.js
└── view/             # server / docker / k8s / release 页面
```

## 数据库（SQLite）

`server/config.yaml`：
```yaml
system:
  db-type: sqlite
sqlite:
  path: "db"            # 目录
  db-name: "aidevops"   # 实际文件：db/aidevops.db
```
首次启动自动创建目录与库文件。切换 MySQL/PgSql 仅需改 `db-type` 并填配置。

## Docker 部署（推荐）

```bash
# 后端（持久化 SQLite 需挂载 db 目录）
docker build -t ai-devops-server ./server
docker run -d --name ai-devops-server -p 8888:8888 \
  -v $(pwd)/data:/go/src/github.com/flipped-aurora/gin-vue-admin/server/db \
  ai-devops-server

# 前端
docker build -t ai-devops-web ./web
docker run -d --name ai-devops-web -p 8080:80 ai-devops-web
```

> IPMI 电源控制需在后端容器内安装 `ipmitool`（`apk add ipmitool`）。

## 首次初始化

1. 启动后端（8888）与前端（8080）
2. 浏览器访问前端 → 自动跳转「初始化」→ 数据库类型选 **sqlite**（已预填）→ 初始化
3. 登录 **admin / 123456**
4. 左侧菜单：服务器管理 / Docker管理 / K8s管理 / 发版工作流

## 本地开发

```bash
# 后端（需 go 1.24）
cd server && go mod tidy && go run .

# 前端
cd web && pnpm install && pnpm dev
```

## 默认账号

| 用户名 | 密码 |
| --- | --- |
| admin | 123456 |

## License

MIT
