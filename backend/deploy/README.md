# 部署说明

本项目通过根目录 `docker-compose.yml` 一键部署，服务短名为 `lab-equipment`。

- MySQL：内部服务，命名卷持久化数据
- Backend：Go 1.22 + Gin + GORM，暴露 `19301:8080`
- Frontend：React 18 + Vite，Nginx 暴露 `18801:80`

部署命令：

```bash
cp .env.example .env
docker compose up -d --build
```

健康检查：

- 后端：`http://127.0.0.1:19301/healthz`
- 前端：`http://127.0.0.1:18801`

生产环境请修改 `.env` 中的数据库口令与 `JWT_SECRET`。
