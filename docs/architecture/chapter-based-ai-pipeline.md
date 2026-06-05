# 章节级 AI Pipeline 设计

## 1. 背景

E-director 需要处理三章节以上的长篇小说。如果把整本小说一次性输入大模型，容易撑开上下文窗口，也会导致章节、人物和场景之间互相混淆。因此 AI Worker 必须以章节为主边界处理输入。

## 2. 核心原则

1. 不把整本小说一次性喂给大模型。
2. 每次优先处理一个章节。
3. 章节过长时继续拆分为 chunk。
4. 每个场景必须带 `chapter_id` 和 `scene_id`。
5. 每个章节生成一份 `schema_design_note`。
6. 最终 YAML Schema 设计文档从章节级说明中汇总，而不是重新读取全文。

## 3. Pipeline 流程

```text
Novel Text
  │
  ▼
Chapter Splitter
  │
  ├── chapter_001
  ├── chapter_002
  └── chapter_003
      │
      ▼
Chapter-level LLM Context
      │
      ├── chapter summary
      ├── scene index extraction
      ├── scene YAML generation
      ├── scene design reasons
      └── chapter schema design note
          │
          ▼
Final Schema Report Builder
```

## 4. 章节级处理

每个章节进入模型时，Prompt 应包含：

- 小说标题；
- 当前章节 ID；
- 当前章节标题；
- 当前章节正文；
- 固定输出要求；
- 场景 ID 命名规则；
- YAML 字段说明；
- 设计原因输出要求。

模型输出应包含：

```json
{
  "chapter_id": "chapter_001",
  "chapter_title": "第一章 雨夜",
  "chapter_summary": "本章摘要",
  "scenes": [],
  "chapter_schema_design_note": {
    "summary": "本章 YAML 结构设计原因概述",
    "key_reasons": []
  }
}
```

## 5. 场景级标签

每个场景必须包含：

```yaml
chapter_id: chapter_001
chapter_title: 第一章 雨夜
scene_id: chapter_001_scene_001
scene_index: 1
```

这样即使 Python Worker 并发处理多个章节和多个场景，Go Gateway 与前端也能根据 ID 稳定排序和归组。

## 6. 章节设计小说明

每个章节生成一个 `chapter_schema_design_note`：

```yaml
schema_design_note:
  summary: "本章强调雨夜视觉氛围和人物情绪变化，因此增强 camera、visual_elements 和 emotion_tags 字段。"
  key_reasons:
    - field_name: "camera"
      reason: "用于把环境描写转译为镜头语言，服务后续 AI 视频生成。"
    - field_name: "emotion_tags"
      reason: "用于显式表达角色的心理变化和场景情绪基调。"
```

## 7. 最终文档生成

最终 YAML Schema 设计文档只消费以下数据：

```text
chapter_001.schema_design_note
chapter_002.schema_design_note
chapter_003.schema_design_note
...
```

这样可以：

1. 降低最终总结阶段上下文压力；
2. 保证设计原因来自每个章节的真实内容；
3. 让最终文档可以解释字段为什么存在；
4. 更贴合比赛题目要求。

## 8. 并发策略

建议配置：

- `max_concurrent_chapters`: 控制章节并发数；
- `max_concurrent_scenes`: 控制章节内部场景并发数；
- Redis 保存 job、chapter、scene 状态；
- gRPC stream 持续把章节和场景事件推给 Go Gateway；
- Go Gateway 转换为 SSE 推给 Vue 前端。

## 9. 错误处理

当单个场景失败时，应返回带 `chapter_id` 和 `scene_id` 的错误事件，避免整个任务不可恢复：

```json
{
  "code": "SCENE_GENERATION_FAILED",
  "chapter_id": "chapter_001",
  "scene_id": "chapter_001_scene_003",
  "retryable": true
}
```

## 10. 比赛答辩表述

E-director 的 AI Worker 采用章节级上下文边界，不将整本小说一次性输入模型。每章单独提取场景并生成 YAML，每个场景携带章节标签和场景标签，避免并发生成造成归属混乱。每章还生成一份 YAML Schema 设计小说明，最终文档从这些小说明中归纳汇总，从而解释 Schema 的设计原因。
