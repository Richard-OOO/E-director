# E-director 剧本 YAML Schema 设计说明文档

## 1. 设计目标

E-director 的目标是把三章节以上的长篇小说转换为可用于影视化、短剧化或 AI 视频生成的结构化剧本。YAML Schema 不只是为了保存文本结果，而是为了让小说中的叙事信息、画面信息、台词信息和情绪信息都能够被后续系统稳定理解。

因此，本 Schema 的设计目标包括：

1. 保持章节与场景的归属关系，避免并发生成时混淆内容来源。
2. 把小说文字转换为可拍摄、可渲染、可剪辑的结构化信息。
3. 支持前端实时展示每个场景的 YAML 内容。
4. 支持 Hover 气泡展示字段设计原因。
5. 支持后续从章节级设计说明汇总出全局 YAML Schema 设计文档。

## 2. Schema 总览

推荐的剧本 YAML 顶层结构如下：

```yaml
script:
  schema_version: "1.0"
  title: "小说标题"
  language: "zh-CN"
  chapters:
    - chapter_id: "chapter_001"
      chapter_title: "第一章"
      chapter_summary: "本章摘要"
      schema_design_note:
        summary: "本章 YAML 结构设计原因概述"
        key_reasons: []
      scenes: []
```

## 3. 章节结构

每个章节必须带有明确的 `chapter_id`，章节下的每个场景也必须重复携带 `chapter_id`。

```yaml
chapter_id: "chapter_001"
chapter_title: "第一章 雨夜"
chapter_summary: "本章讲述主角在雨夜中的关键相遇。"
```

### 设计原因

长篇小说至少包含三个章节，不同章节可能发生在不同时间、地点和人物关系阶段。AI Worker 会按照章节分别输入大模型，以避免整本小说一次性进入上下文导致内容溢出。章节 ID 可以保证并发处理后，每个场景仍能追溯到原始章节。

## 4. 场景结构

场景是视频化生产的最小单元。

```yaml
scene_id: "chapter_001_scene_001"
chapter_id: "chapter_001"
chapter_title: "第一章 雨夜"
scene_index: 1
title: "雨夜街口"
summary: "主角在雨夜街口等待车辆。"
```

### 设计原因

`scene_id` 用于支持场景级生成、重试、编辑和前端局部刷新。`chapter_id` 与 `scene_id` 同时存在，可以避免多章节并发生成时出现归属混乱。

## 5. 核心字段设计原因

### 5.1 `characters`

```yaml
characters:
  - name: "林舟"
    role: "protagonist"
    emotional_state: "lonely"
```

小说中人物关系往往跨章节发展。`characters` 字段可以明确当前场景中出现的人物、人物职能和情绪状态，避免大模型在长文本生成中混淆角色。

### 5.2 `camera`

```yaml
camera:
  - shot_type: "wide_shot"
    movement: "slow_push_in"
    description: "雨幕笼罩街道，镜头缓慢推向主角。"
```

小说是文字叙事，而视频生成需要画面指令。`camera` 字段把环境描写、空间关系和人物动作转译为镜头语言，例如远景、近景、特写、推镜、摇镜和跟拍。这有利于后续 AI 视频生成模型理解画面构图与镜头调度。

### 5.3 `dialogues`

```yaml
dialogues:
  - character: "林舟"
    text: "今晚的雨，好像不会停了。"
    emotion: "tired"
```

台词是剧本的核心表达方式。`dialogues` 字段把小说中的对话、独白和叙述性语言拆成可表演、可配音、可字幕化的结构。

### 5.4 `emotion_tags`

```yaml
emotion_tags:
  - loneliness
  - tension
```

小说的情绪经常隐藏在环境描写和人物动作中。`emotion_tags` 把隐含情绪显式化，帮助演员表演、AI 配乐、AI 视频风格选择和后期剪辑节奏控制。

### 5.5 `visual_elements`

```yaml
visual_elements:
  - rain
  - street_light
  - empty_road
```

AI 视频生成通常需要明确的视觉要素。`visual_elements` 把小说中的关键物品、天气、光线、空间和符号元素提取出来，帮助生成模型保留画面重点。

### 5.6 `scene_transition`

```yaml
scene_transition:
  type: "cut"
  description: "切到远处驶来的车灯。"
```

视频不是孤立片段，而是连续叙事。`scene_transition` 描述场景之间的衔接方式，例如 cut、fade、dissolve、match_cut 或 time_jump，有助于最终形成完整短片或分镜脚本。

## 6. 动态 Schema 设计策略

E-director 不把所有章节都强行套入完全相同的字段组合，而是采用动态 Schema 思路：

1. 每个章节单独进入大模型上下文。
2. 每个章节先提取场景，再生成场景 YAML。
3. 每个章节生成一份 `schema_design_note`，说明本章为什么需要某些字段。
4. 最终文档只汇总各章节的 `schema_design_note`，而不是再次读取整本小说。

这样既能避免上下文过长，也能让设计说明来自章节内容本身。

## 7. 面向 AI 视频生成的优势

该 Schema 对 AI 视频生成有以下好处：

1. `camera` 提供镜头调度信息。
2. `visual_elements` 提供关键画面元素。
3. `emotion_tags` 提供情绪和风格控制信号。
4. `dialogues` 支持配音、字幕和表演。
5. `scene_transition` 支持连续叙事剪辑。
6. `chapter_id` 与 `scene_id` 支持稳定追踪和局部重试。

## 8. 总结

本 YAML Schema 的核心设计原因是：把小说从自然语言叙事转换成可被工程系统、前端工作台和 AI 视频生成系统共同理解的结构化剧本。它既保留了章节和场景的叙事结构，也显式表达了镜头、台词、情绪、视觉元素和转场信息。
