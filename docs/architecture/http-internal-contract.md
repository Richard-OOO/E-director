# Go Gateway 与 Python AI Worker 的 HTTP 内部契约

## 1. 通信边界

E-director 采用两层 HTTP 模型：

```text
Vue Frontend -> Go Gateway: HTTP + JSON
Go Gateway -> Python AI Worker: HTTP + JSON
```

浏览器侧只访问 Go Gateway。Go Gateway 负责统一鉴权、路由、聚合与错误处理，Python AI Worker 仅作为内部 AI 服务存在，不直接暴露给浏览器。

## 2. 契约文件

HTTP 契约建议定义在：

```text
contracts/http/
```

建议同时保留：

- `contracts/http/openapi.yaml`
- `contracts/schemas/*.json`
- `contracts/examples/*.json`

## 3. 请求与响应

所有对外 HTTP API 统一使用 JSON：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

错误时：

```json
{
  "code": 40001,
  "message": "invalid request",
  "data": null
}
```

## 4. 核心稳定性设计

1. Go Gateway 对 AI Worker 调用设置超时和重试边界。
2. Python AI Worker 对每个章节和场景设置处理超时。
3. 所有错误对象必须携带 `chapter_id` 或 `scene_id`，便于定位。
4. Go Gateway 可将 job 状态写入 Redis，支持断点恢复与进度查询。
5. 接口只追加字段，不随意删除已有字段，避免破坏兼容性。

## 5. Why HTTP over gRPC

Browsers do not need gRPC for this project. Keeping the frontend, gateway, and AI worker on HTTP/JSON makes the architecture easier to debug, easier to demonstrate, and faster to ship during the competition.
