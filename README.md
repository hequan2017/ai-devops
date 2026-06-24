<div align="center">
  <h1>AI运维平台（AI-DevOps）</h1>
  <p>服务器 · Docker · Kubernetes · 发版工作流 一体化运维管理平台</p>
</div>

<br/>

## 简介

AI运维平台覆盖**物理服务器（戴尔/华为）、Docker、Kubernetes** 的统一运维管理，提供 **WebSSH 终端、IPMI 电源控制、性能定时采集、代码发版工作流审批**。后端 Go（Gin + GORM），前端 Vue3 + Element Plus，**默认 SQLite 文件数据库，开箱即用**。

> **编译环境**：需 **go 1.24+**（推荐 go 1.26，依赖 casbin v3 / pgx / x/crypto 等使用 go1.21–1.24 标准库）。Docker 镜像 `golang:alpine` 自带 go 1.24。

## 功能模块

### 1. 服务器管理（戴尔 / 华为）
- 资产台账：厂商/型号/SN/机房机柜/配置/负责人
- **WebSSH 终端**：浏览器内通过 SSH 密钥/密码登录（xterm.js + WebSocket）
- **IPMI 电源管理**：状态/开机/软关机/强制关机/硬重启（ipmitool）
- **性能采集**：CPU/内存/磁盘/网速，**每分钟 SSH 主动采集**（cron 定时），历史查询

### 2. Docker 管理
- 引擎接入点（unix socket / 远程 tcp）+ 连接测试
- 容器：列表/启停/重启/删除/**实时日志流**(WebSocket follow)
- 镜像、网络、数据卷

### 3. Kubernetes 管理
- 集群接入点（Token，支持跳过自签 TLS）+ 连接测试
- Namespace / Pod（日志/删除）/ Deployment / Service / Node

### 4. 发版工作流（类 Jenkins 审批）
- 草稿 → 待审批 → 通过/拒绝 → 执行 → 完成/失败
- 多角色审批，执行阶段运行发布脚本并记录结果

### 5. Dashboard
- 仪表盘运维概览，快捷跳转各功能模块

### 平台基础
RBAC + 动态菜单、用户/角色/菜单/API、操作审计、登录日志、字典参数、代码生成器。

## 技术栈

| 端 | 技术 |
| --- | --- |
| 后端 | Go 1.24+ · Gin · GORM · SQLite(`glebarez/sqlite`) · cron(`robfig/cron`) · WebSSH(`gorilla/websocket`+`x/crypto/ssh`) · Docker/K8s 经 REST 直连 |
| 前端 | Vue 3 · Vite · Element Plus · Pinia · xterm.js |

## 目录结构

```
server/
├── model/biz/        # Server / DockerHost / K8sCluster / Release / ServerMetric
├── service/biz/      # 业务 + docker_client/k8s_client/ipmi/terminal/metric
├── api/v1/biz/       # 业务接口
├── router/biz/       # 业务路由
├── task/             # 定时任务（性能采集）
├── source/system/    # 初始化数据（菜单/API/权限）
└── config.yaml       # 配置（db-type 默认 sqlite）
web/src/
├── api/              # server/docker/k8s/release
├── view/             # server/docker/k8s/release
└── components/terminal/  # xTerminal 终端组件
```

## 数据库（SQLite）

`server/config.yaml`：
```yaml
system:
  db-type: sqlite
sqlite:
  path: "db"
  db-name: "aidevops"   # 实际文件：db/aidevops.db
```
首次启动自动建目录与库。切 MySQL/PgSql 仅需改 `db-type`。

## Docker 部署（推荐）

```bash
# 后端（持久化 SQLite 挂载 db 目录；IPMI 需容器内 apk add ipmitool）
docker build -t ai-devops-server ./server
docker run -d --name ai-devops-server -p 8888:8888 \
  -v $(pwd)/data:/go/src/github.com/flipped-aurora/gin-vue-admin/server/db \
  ai-devops-server

# 前端
docker build -t ai-devops-web ./web
docker run -d --name ai-devops-web -p 8080:80 ai-devops-web
```

## 首次初始化

1. 启动后端（8888）与前端（8080）
2. 前端自动跳「初始化」→ 选 **sqlite**（已预填）→ 初始化
3. 登录 **admin / 123456**
4. 在服务器管理录入资产并填写 **SSH 凭证**（用于 WebSSH 与性能采集）
5. 性能采集自动每分钟运行

## 本地开发

```bash
# 后端（需 go 1.24+）
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
