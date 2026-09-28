[简体中文](README.md) | [English](README.en.md)

# AI运维平台（AI-DevOps）

服务器 · Docker · Kubernetes · 发版 · 告警 一体化运维平台。

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](server/go.mod)
[![Vue 3](https://img.shields.io/badge/Vue-3-4FC08D?logo=vuedotjs&logoColor=white)](web/package.json)
[![License: BSL 1.1](https://img.shields.io/badge/License-BSL%201.1-blue.svg)](LICENSE)

## 简介

AI-DevOps 是面向基础设施的运维管理平台，覆盖**物理服务器（戴尔/华为）、Docker、Kubernetes** 的统一管理，提供 **WebSSH、IPMI 电源、性能采集、告警通知、批量执行、发版审批**能力。后端 Go（Gin + GORM），前端 Vue 3 + Element Plus，**默认 SQLite，开箱即用**。

适合管理自有物理机/虚拟机与容器环境、需要一个轻量「资产 + 终端 + 采集 + 告警 + 发版」闭环的运维团队和个人开发者。

> **编译要求**：go 1.24+。Docker 镜像 `golang:alpine` 自带 go 1.24。

## ✨ 功能特性

### 1. 服务器管理

- **资产台账**：厂商/型号/SN/机房机柜/配置/负责人 + 分组标签，含 IPMI/BMC 地址与账号
- **WebSSH 终端**：SSH 密钥/密码登录，xterm.js 前端，支持窗口自适应 resize（WebSocket JSON 协议）
- **IPMI 电源管理**：状态查询/开机/软关机/强制关机/硬重启（基于 `ipmitool -I lanplus`）
- **性能采集**：CPU/内存/磁盘/网速，定时任务**每分钟 SSH 主动采集**，采集后自动检查告警规则
- **批量命令执行**：多选服务器 SSH 批量跑命令（并发上限 10，单命令 30s 超时）
- **端口探活**：TCP 连通性检测 + 延迟（3s 超时）
- **SSH 凭据管理**：集中管理登录密钥/密码，供 WebSSH、采集、批量执行复用

### 2. Docker 管理

- 接入点（unix/tcp）+ 连接测试
- 容器：列表/启停/重启/删除 + 实时日志流 + **exec 交互终端**（exec + hijack 双向桥接 WebSocket）
- 镜像、网络、数据卷管理

### 3. Kubernetes 管理

- 集群接入点（Token，跳过自签 TLS）+ 连接测试
- Namespace / Pod（日志/删除）/ Deployment / Service / Node 资源管理（REST 直连）

### 4. 发版工作流（类 Jenkins）

- 状态机：`draft → pending → approved/rejected → executing → done/failed`；记录项目/版本/分支/发布环境与发布脚本，申请人、审批人、审批与执行时间、执行结果全程留痕

### 5. 告警管理

- 告警规则：指标/操作符/阈值/级别
- 采集后自动检测，**同规则未处理告警去重**，支持 **webhook 通知**（钉钉/企业微信/飞书通用 text 格式）
- 告警记录查询 + 一键标记处理

### 6. Dashboard 与数据治理

- 运维概览：服务器/Docker/K8s/未处理告警/待审批计数，一键跳转
- 采集数据自动清理（性能指标保留 7 天，告警记录保留 90 天，每日定时执行）
- 并发与超时控制，防慢命令/连接风暴

## 🛠 技术栈

| 端 | 技术 |
| --- | --- |
| 后端 | Go 1.24+ · Gin · GORM · SQLite · cron · WebSocket（gorilla/websocket）· SSH（x/crypto/ssh）· Docker/K8s REST 直连 |
| 前端 | Vue 3 · Vite · Element Plus · Pinia · xterm.js（Node.js 20 / pnpm） |

## 🚀 快速开始

### 本地开发

```bash
cd server && go mod tidy && go run .   # 需 go 1.24+，默认监听 8888
cd web && pnpm install && pnpm dev     # 前端 8080
```

### Docker 部署

```bash
docker build -t ai-devops-server ./server
docker run -d --name ai-devops-server -p 8888:8888 \
  -v $(pwd)/data:/go/src/github.com/flipped-aurora/gin-vue-admin/server/db ai-devops-server
docker build -t ai-devops-web ./web
docker run -d --name ai-devops-web -p 8080:80 ai-devops-web
```

> IPMI 功能需在容器内 `apk add ipmitool`。

### 数据库

`server/config.yaml` 中 `system.db-type: sqlite`，`sqlite: {path: "db", db-name: "aidevops"}` → 数据文件 `db/aidevops.db`，首次启动自动建库。切换 MySQL/PgSql 仅需修改 `db-type`。

### 首次初始化

1. 启动后端（8888）与前端（8080）
2. 前端自动跳「初始化」页 → 选 sqlite → 初始化
3. 登录 **admin / 123456**
4. 服务器管理录入资产 + SSH 凭据（供 WebSSH / 采集 / 批量执行使用）
5. 性能采集每分钟自动运行；配置告警规则 + webhook 即可形成告警闭环

### 常用命令

```bash
cd server && go test ./... && go vet ./...   # 后端验证
cd web && pnpm build                         # 前端构建验证
```

部署资产：`deploy/docker-compose/docker-compose.yaml`（单机编排样例）、`deploy/kubernetes/`（集群清单）、`server/Dockerfile` 与 `web/Dockerfile`（独立镜像，均基于 Node.js 20 / Go 官方基础镜像构建）。

## 📁 目录结构

```
server/
├── api/v1/biz/       # 服务器/Docker/K8s/发版/告警/终端 API
├── model/biz/        # Server/DockerHost/K8sCluster/Release/ServerMetric/AlertRule/AlertRecord
├── service/biz/      # 业务实现 + docker_client/k8s_client/ipmi/terminal/metric/alert/exec_cmd
├── task/             # 定时任务（指标采集 / 数据清理）
├── source/system/    # 初始化数据（菜单/API/权限）
└── config.yaml       # 配置（db-type 默认 sqlite）
web/src/
├── api/ view/ components/terminal/
deploy/               # Docker、docker-compose、Kubernetes 部署资产
```

## 🤝 社区与贡献

欢迎通过 Issue 反馈问题、提交 PR 参与贡献。行为准则、贡献流程与安全策略见 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)、[CONTRIBUTING.md](CONTRIBUTING.md)、[SECURITY.md](SECURITY.md)。

## 📄 License

本项目基于 [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) 定制开发，采用 [Business Source License 1.1](LICENSE) 授权（Licensed Work: gin-vue-admin）。许可证允许满足条款的个人使用、评估与开发使用，以及非商业教学、研究或学术用途；未被明确允许的使用属于 Production Use，需要取得有效商业许可证。每个版本在首次公开发布三年后到达 Change Date，转为 Apache License 2.0，具体以 [LICENSE](LICENSE) 原文为准。
