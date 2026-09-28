[简体中文](README.md) | [English](README.en.md)

# AI-DevOps

An all-in-one operations platform for servers · Docker · Kubernetes · releases · alerts.

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](server/go.mod)
[![Vue 3](https://img.shields.io/badge/Vue-3-4FC08D?logo=vuedotjs&logoColor=white)](web/package.json)
[![License: BSL 1.1](https://img.shields.io/badge/License-BSL%201.1-blue.svg)](LICENSE)

## Introduction

AI-DevOps is an infrastructure-focused operations platform that unifies the management of **physical servers (Dell/Huawei), Docker and Kubernetes**, providing **WebSSH, IPMI power control, metric collection, alert notifications, batch execution and release approvals**. The backend is Go (Gin + GORM), the frontend is Vue 3 + Element Plus, and it **defaults to SQLite — ready to run out of the box**.

It suits ops teams and individual developers who manage their own physical/virtual machines and container environments and want a lightweight closed loop of "assets + terminals + collection + alerts + releases".

> **Build requirement**: go 1.24+. The `golang:alpine` Docker image ships with go 1.24.

## ✨ Features

### 1. Server Management

- **Asset inventory**: vendor/model/SN/data center rack/configuration/owner plus group tags, with IPMI/BMC address and credentials
- **WebSSH terminal**: SSH key/password login, xterm.js frontend, window resize support (WebSocket JSON protocol)
- **IPMI power control**: status / power on / soft off / hard off / hard reset (via `ipmitool -I lanplus`)
- **Metric collection**: CPU/memory/disk/network throughput, actively collected over SSH every minute by a scheduled task; alert rules are evaluated automatically after each round
- **Batch command execution**: run commands over SSH on multiple selected servers (concurrency limit 10, 30s timeout per command)
- **Port probing**: TCP reachability check with latency (3s timeout)
- **SSH credential management**: centrally managed keys/passwords reused by WebSSH, collection and batch execution

### 2. Docker Management

- Endpoints (unix/tcp) + connection test
- Containers: list/start/stop/restart/delete + live log streaming + **interactive exec terminal** (exec + hijack bridged over WebSocket)
- Images, networks and volumes

### 3. Kubernetes Management

- Cluster endpoints (Token, self-signed TLS skipped) + connection test
- Namespace / Pod (logs/delete) / Deployment / Service / Node management (direct REST calls)

### 4. Release Workflow (Jenkins-style)

- State machine: `draft → pending → approved/rejected → executing → done/failed`; records project/version/branch/target environment and release script, with a full audit trail of applicant, approver, approval/execution timestamps and results

### 5. Alert Management

- Alert rules: metric / operator / threshold / severity
- Automatic evaluation after each collection, **deduplicated per rule while unresolved**, with **webhook notifications** (generic text format for DingTalk / WeCom / Feishu)
- Alert record query + one-click resolution

### 6. Dashboard & Data Hygiene

- Ops overview: counts for servers/Docker/K8s/unresolved alerts/pending approvals with one-click navigation
- Automatic cleanup of collected data (metrics retained 7 days, alert records 90 days, run by a daily task)
- Concurrency and timeout controls against slow commands and connection storms

## 🛠 Tech Stack

| Layer | Technologies |
| --- | --- |
| Backend | Go 1.24+ · Gin · GORM · SQLite · cron · WebSocket (gorilla/websocket) · SSH (x/crypto/ssh) · direct Docker/K8s REST |
| Frontend | Vue 3 · Vite · Element Plus · Pinia · xterm.js (Node.js 20 / pnpm) |

## 🚀 Quick Start

### Local Development

```bash
cd server && go mod tidy && go run .   # requires go 1.24+, listens on 8888
cd web && pnpm install && pnpm dev     # frontend on 8080
```

### Docker Deployment

```bash
docker build -t ai-devops-server ./server
docker run -d --name ai-devops-server -p 8888:8888 \
  -v $(pwd)/data:/go/src/github.com/flipped-aurora/gin-vue-admin/server/db ai-devops-server
docker build -t ai-devops-web ./web
docker run -d --name ai-devops-web -p 8080:80 ai-devops-web
```

> For IPMI support, run `apk add ipmitool` inside the container.

### Database

In `server/config.yaml`: `system.db-type: sqlite` and `sqlite: {path: "db", db-name: "aidevops"}` → database file `db/aidevops.db`, created automatically on first start. Switching to MySQL/PgSQL only requires changing `db-type`.

### First-run Initialization

1. Start the backend (8888) and the frontend (8080)
2. The frontend redirects to the "Initialize" page → pick sqlite → initialize
3. Log in with **admin / 123456**
4. Register assets and SSH credentials under Server Management (used by WebSSH / collection / batch execution)
5. Metric collection runs automatically every minute; configure alert rules + a webhook to close the alerting loop

### Common Commands

```bash
cd server && go test ./... && go vet ./...   # backend checks
cd web && pnpm build                         # frontend build check
```

Deployment assets: `deploy/docker-compose/docker-compose.yaml` (single-host orchestration example), `deploy/kubernetes/` (cluster manifests), `server/Dockerfile` and `web/Dockerfile` (standalone images, both built on official Node.js 20 / Go base images).

## 📁 Directory Structure

```
server/
├── api/v1/biz/       # Server/Docker/K8s/release/alert/terminal APIs
├── model/biz/        # Server/DockerHost/K8sCluster/Release/ServerMetric/AlertRule/AlertRecord
├── service/biz/      # Business logic + docker_client/k8s_client/ipmi/terminal/metric/alert/exec_cmd
├── task/             # Scheduled tasks (metric collection / data cleanup)
├── source/system/    # Seed data (menus/APIs/permissions)
└── config.yaml       # Configuration (db-type defaults to sqlite)
web/src/
├── api/ view/ components/terminal/
deploy/               # Docker, docker-compose and Kubernetes assets
```

## 🤝 Community & Contributing

Issues and pull requests are welcome. See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md), [CONTRIBUTING.md](CONTRIBUTING.md) and [SECURITY.md](SECURITY.md) for the code of conduct, contribution process and security policy.

## 📄 License

This project is a customization of [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) and is licensed under the [Business Source License 1.1](LICENSE) (Licensed Work: gin-vue-admin). Use that complies with the terms — individual, evaluation and development use, as well as non-commercial teaching, research or academic use — is permitted; any use not expressly allowed constitutes Production Use and requires a valid commercial license. Each version converts to the Apache License 2.0 at its Change Date three years after its first public release. See the [LICENSE](LICENSE) file for the exact terms.
