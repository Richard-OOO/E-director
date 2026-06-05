# Go Gateway 与 Python AI Worker 的 gRPC 内部契约

## 1. 通信边界

E-director 采用两层通信模型：

```text
Vue Frontend <-> Go Gateway: HTTP + SSE
Go Gateway <-> Python AI Worker: gRPC Server Streaming
```

浏览器侧保留 HTTP/SSE，便于实时展示和调试；后端内部使用 gRPC，便于强类型契约、跨语言代码生成和流式传输。

## 2. 契约文件

内部 gRPC 契约定义在：

```text
contracts/proto/edirector/ai/v1/generation.proto
```

核心服务：

```proto
service AiGenerationService {
  rpc GenerateScript(GenerateScriptRequest) returns (stream GenerationStreamEvent);
}
```

## 3. 流式事件

AI Worker 返回的流事件包含：

- `generation_started`
- `chapter_started`
- `scene_started`
- `scene_delta`
- `scene_completed`
- `chapter_completed`
- `schema_summary`
- `generation_completed`
- `generation_failed`

Go Gateway 接收这些事件后转换为 SSE 事件推给前端。

## 4. 稳定性设计

为了保证内部通信稳定，后续实现时应遵守：

1. Go Gateway 调用 gRPC 时设置 deadline。
2. Python AI Worker 对每个章节和场景设置超时。
3. 所有错误事件必须携带 `chapter_id` 或 `scene_id`，便于定位。
4. Go Gateway 将每个事件写入 Redis，用于断线恢复。
5. gRPC stream 中断时，Gateway 返回可重连的 SSE 错误事件。
6. Proto 字段只追加不重排，避免破坏兼容性。

## 5. 为什么不用前端 gRPC

浏览器原生不直接支持普通 gRPC。如果前端也使用 gRPC，需要引入 gRPC-Web 和 Envoy 等额外组件。比赛项目更需要清晰可演示，因此前端保持 HTTP/SSE，内部服务使用 gRPC 是更稳妥的组合。
