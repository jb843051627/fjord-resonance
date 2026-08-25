Web 后端项目，负责远海声学浮标水听器校准批次与采样质量管理。
# fjord-resonance

本项目是一个单体 Go 服务，使用磁盘 SQLite 保存浮标、传感器、校准批次、声学样本、质量结果和告警。默认监听 `:8080`，数据库路径由 `FJORD_RESONANCE_DB` 指定，未设置时使用当前目录的 `fjord-resonance.db`。

启动：

```bash
GOTOOLCHAIN=local go run .
```

页面位于 `/`，健康检查位于 `/healthz`。API 以 `/api/v1/` 开头。
