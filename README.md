<div align="center">
  <h1>AI运维平台（AI-DevOps）</h1>
  <p>面向戴尔（Dell）/ 华为（Huawei）等物理服务器的运维资产管理系统</p>
</div>

<br/>

## 简介

AI运维平台是一套面向物理服务器资产的运维管理系统，聚焦戴尔、华为等主流品牌服务器的全生命周期管理。后端基于 Go（Gin + GORM），前端基于 Vue3 + Element Plus，**默认使用 SQLite 文件数据库，开箱即用**，并支持 Docker 一键部署。

## 功能特性

- **服务器资产管理**：厂商、型号、SN 序列号、机房、机柜、U 位、CPU/内存/磁盘配置、操作系统、负责人
- **网络信息**：带外管理 IP（iDRAC / iBMC）与业务 IP 统一记录
- **状态可视化**：运行中 / 已关机 / 维护中 / 故障，列表与详情双重展示
- **多维检索**：关键字（名称/SN/IP）+ 厂商 + 状态组合查询
- **权限体系**：JWT + Casbin 的 RBAC 角色权限，动态菜单
- **平台基础能力**：用户/角色/菜单/API 管理、操作审计、登录日志、字典与参数管理、代码生成器等

## 技术栈

| 端 | 技术 |
| --- | --- |
| 后端 | Go 1.24 · Gin · GORM · SQLite(`glebarez/sqlite`，纯 Go 实现，兼容 `CGO_ENABLED=0`) |
| 前端 | Vue 3 · Vite · Element Plus · Pinia |
| 鉴权 | JWT · Casbin |

## 目录结构

```
ai-devops/
├── server/                  # 后端服务
│   ├── model/biz/           # 业务模型：Server 服务器资产
│   ├── service/biz/         # 业务逻辑层
│   ├── api/v1/biz/          # 业务接口层
│   ├── router/biz/          # 业务路由层
│   ├── source/system/       # 初始化数据（菜单 / API / 权限）
│   ├── initialize/          # 数据库与路由初始化
│   └── config.yaml          # 配置文件（db-type 默认 sqlite）
├── web/                     # 前端应用
│   └── src/view/server/     # 服务器管理页面
└── deploy/                  # 部署相关（docker / compose / k8s）
```

## 数据库（SQLite）

默认使用 SQLite 文件数据库，配置位于 `server/config.yaml`：

```yaml
system:
  db-type: sqlite
sqlite:
  path: "db"            # 数据库文件所在目录
  db-name: "aidevops"   # 文件名（不含扩展名）→ 实际文件：db/aidevops.db
```

- 首次启动会**自动创建目录与库文件**，无需手动建库。
- 如需切换为 MySQL / PostgreSQL，仅需修改 `db-type` 并填写对应连接信息。

## Docker 部署（推荐）

### 后端

```bash
docker build -t ai-devops-server ./server

# 持久化 SQLite 文件，需挂载容器内的 db 目录
docker run -d --name ai-devops-server -p 8888:8888 \
  -v $(pwd)/data:/go/src/github.com/flipped-aurora/gin-vue-admin/server/db \
  ai-devops-server
```

> 说明：`server/Dockerfile` 构建时为 `CGO_ENABLED=0`，与纯 Go 实现的 SQLite 驱动完全兼容，**无需安装 CGO/系统 gcc**。

### 前端

```bash
docker build -t ai-devops-web ./web
docker run -d --name ai-devops-web -p 8080:80 ai-devops-web
```

## 首次初始化

1. 分别启动后端（`8888`）与前端（`8080`）。
2. 浏览器访问前端，系统检测到数据库未初始化会自动跳转「初始化」页面。
3. 数据库类型选择 **sqlite**（路径与库名已预填），点击「初始化」。
4. 初始化完成后，使用默认账号登录：**admin / 123456**。
5. 在左侧菜单进入「**服务器管理**」即可开始录入服务器资产。

## 本地开发

```bash
# 后端
cd server
go mod tidy
go run .

# 前端
cd web
pnpm install
pnpm dev
```

## 默认账号

| 用户名 | 密码 |
| --- | --- |
| admin | 123456 |

## License

MIT
