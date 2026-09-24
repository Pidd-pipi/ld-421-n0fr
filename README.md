# 实验室设备管理系统

面向高校和科研院所实验室的全栈 Web 应用，支持设备资产登记、借用归还、维护保养、使用预约和设备状态看板。

## 一键启动

```bash
cp .env.example .env
docker compose up -d --build
```

启动完成后访问：

- 前端：http://127.0.0.1:18801
- 后端 API：http://127.0.0.1:19301
- 健康检查：http://127.0.0.1:19301/healthz

预置账号（密码见括号）：

- `admin / admin123`（管理员）
- `labmanager / lab123`（实验室管理员）
- `researcher / res123`（研究员）
- `student / stu123`（学生）

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Go 1.22 + Gin + GORM + MySQL 8.0 |
| 认证 | JWT（golang-jwt/jwt/v5） |
| 前端 | React 18 + TypeScript + Vite |
| UI | Ant Design 5 |
| 图表 | ECharts |
| 状态管理 | Zustand |
| 部署 | Docker Compose + Nginx |

## 目录结构

```
.
├── backend/               # Go 后端
│   ├── cmd/server/        # 入口
│   ├── internal/          # config/model/repository/service/handler/router/middleware/dto/...
│   ├── database/          # migrations + seeds
│   ├── api/               # OpenAPI 文档
│   └── Dockerfile
├── frontend/              # React 前端
│   ├── src/api|stores|types|components|hooks|pages|router|utils|constants
│   └── Dockerfile
├── database/init.sql      # MySQL 初始化脚本
└── docker-compose.yml
```

## 核心实体

设备、借用记录、维护记录、预约记录、设备分类，均贯穿前后端。

## 枚举位置

- 后端：`backend/internal/constants/enums.go`
- 前端：`frontend/src/types/enums.ts`

## 权限角色

`Admin` / `LabManager` / `Researcher` / `Student`。Student 不能审批借用/预约。

## 本地开发

后端：

```bash
cd backend
go mod tidy
go run ./cmd/server
```

前端：

```bash
cd frontend
npm install
npm run dev
```

## License

MIT
