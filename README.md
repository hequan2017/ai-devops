<div align="center">
  <h1>AI运维平台（AI-DevOps）</h1>
  <p>服务器 · Docker · Kubernetes · 发版 · 告警 一体化运维平台</p>
</div>

<br/>

## 简介

面向基础设施的运维管理平台，覆盖**物理服务器（戴尔/华为）、Docker、Kubernetes**统一管理，提供 **WebSSH、IPMI 电源、性能采集、告警通知、批量执行、发版审批**。后端 Go（Gin + GORM），前端 Vue3 + Element Plus，**默认 SQLite，开箱即用**。

> **编译**：需 **go 1.24+**（推荐 go 1.26；依赖 casbin v3/pgx/x/crypto 等用 go1.21–1.24 标准库）。Docker 镜像 `golang:alpine` 自带 go 1.24。

## 功能模块

### 1. 服务器管理
- 资产台账：厂商/型号/SN/机房机柜/配置/负责人 + **分组标签**
- **WebSSH 终端**：密钥/密码登录，xterm.js，**窗口自适应 resize**
- **IPMI 电源**：状态/开机/软关机/强制关机/硬重启（ipmitool）
- **性能采集**：CPU/内存/磁盘/网速，**每分钟 SSH 主动采集**（cron），并发上限 10
- **批量命令执行**：多选服务器 SSH 批量跑命令（并发 10 + 单命令 30s 超时）
- **端口探活**：TCP 连通性 + 延迟

### 2. Docker 管理
- 接入点（unix/tcp）+ 连接测试
- 容器：列表/启停/重启/删除 + **实时日志流** + **exec 交互终端**
- 镜像、网络、数据卷

### 3. Kubernetes 管理
- 集群接入点（Token，跳过自签 TLS）+ 连接测试
- Namespace / Pod（日志/删除）/ Deployment / Service / Node

### 4. 发版工作流（类 Jenkins）
- 草稿 → 待审批 → 通过/拒绝 → 执行 → 完成/失败

### 5. 告警管理
- 告警规则：指标/操作符/阈值/级别
- 采集后**自动检测**（去重）+ **webhook 通知**（钉钉/企微/飞书通用）
- 告警记录 + 处理

### 6. Dashboard
- 运维概览：服务器/Docker/K8s/告警未处理/待审批 计数，一键跳转

### 健壮性优化
- 采集数据自动清理（ServerMetric 7 天 / AlertRecord 90 天）
- 并发与超时控制（防慢命令/连接风暴）

## 技术栈

| 端 | 技术 |
| --- | --- |
| 后端 | Go 1.24+ · Gin · GORM · SQLite · cron · WebSocket(`gorilla/websocket`) · SSH(`x/crypto/ssh`) · Docker/K8s REST 直连 |
| 前端 | Vue 3 · Vite · Element Plus · Pinia · xterm.js |

## 目录结构

```
server/
├── model/biz/        # Server/DockerHost/K8sCluster/Release/ServerMetric/AlertRule/AlertRecord
├── service/biz/      # 业务 + docker_client/k8s_client/ipmi/terminal/metric/alert/exec_cmd
├── task/             # 定时任务（采集 / 数据清理）
├── source/system/    # 初始化数据（菜单/API/权限）
└── config.yaml       # 配置（db-type 默认 sqlite）
web/src/
├── api/ view/ components/terminal/
```

## 数据库（SQLite）

`server/config.yaml`：`system.db-type: sqlite`，`sqlite: {path: "db", db-name: "aidevops"}` → `db/aidevops.db`，首次启动自动建库。切 MySQL/PgSql 仅改 `db-type`。

## Docker 部署

```bash
docker build -t ai-devops-server ./server
docker run -d --name ai-devops-server -p 8888:8888 \
  -v $(pwd)/data:/go/src/github.com/flipped-aurora/gin-vue-admin/server/db ai-devops-server
docker build -t ai-devops-web ./web
docker run -d --name ai-devops-web -p 8080:80 ai-devops-web
```
> IPMI 需容器内 `apk add ipmitool`。

## 首次初始化

1. 启动后端（8888）与前端（8080）
2. 前端自动跳「初始化」→ 选 sqlite → 初始化
3. 登录 **admin / 123456**
4. 服务器管理录入资产 + SSH 凭证（用于 WebSSH / 采集 / 批量执行）
5. 性能采集每分钟自动运行；配置告警规则 + webhook 实现告警闭环

## 本地开发

```bash
cd server && go mod tidy && go run .   # 需 go 1.24+
cd web && pnpm install && pnpm dev
```

## 默认账号

admin / 123456

## License

MIT
