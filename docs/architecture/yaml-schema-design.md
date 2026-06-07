# E-director 剧本 YAML Schema 设计说明文档

## 1. 文档目的

本文档定义 E-director 项目中“小说文本转结构化剧本”使用的固定 YAML Schema。该 Schema 用于约束大模型输出格式，使模型在处理不同小说章节时，都按照同一套稳定字段生成场景级剧本 YAML。

本项目的核心流程是：用户上传或输入小说文档，系统按章节切分文本，再调用大模型把章节内容转换为若干个可拍摄、可编辑、可用于 AI 视频生成的场景 YAML。为了避免模型每次自行设计字段导致前后格式不一致，本项目采用固定 YAML Schema 作为输出契约。

## 2. 设计目标

固定 YAML Schema 需要同时满足叙事表达、影视化制作和工程落库三类需求。

1. 保持结构稳定：不同章节、不同场景都使用同一套字段，方便前端编辑、后端存储和后续导出。
2. 支持影视化表达：把小说中的人物、地点、动作、台词、镜头和情绪拆成明确字段。
3. 支持 AI 视频生成：将自然语言叙事转化为可被视频生成模型理解的画面提示、运动提示和风格提示。
4. 支持人工编辑：YAML 字段名称清晰，用户可以在编辑器中直接修改。
5. 支持章节级生成：每个场景都携带章节和场景标识，便于流式生成、局部重试和进度展示。

## 3. 顶层结构

每个场景的 YAML 内容必须使用以下顶层结构：

```yaml
schema_version: "1.0"
chapter:
  chapter_id: "chapter_001"
  chapter_title: "第一章"
  chapter_index: 1
scene:
  scene_id: "chapter_001_scene_001"
  scene_index: 1
  title: "场景标题"
  summary: "场景摘要"
  location:
    name: "地点名称"
    type: "interior"
    time_of_day: "night"
    atmosphere: "紧张、潮湿、压抑"
  characters: []
  action: []
  dialogues: []
  camera: []
  visual_elements: []
  audio:
    music: "低频悬疑氛围"
    sound_effects: []
  emotion_tags: []
  ai_video_prompt:
    positive: "用于 AI 视频生成的正向画面提示"
    negative: "需要避免的画面内容"
  transition:
    type: "cut"
    description: "场景转场说明"
```

## 4. 完整 Schema 定义

### 4.1 `schema_version`

```yaml
schema_version: "1.0"
```

字段说明：当前 YAML Schema 的版本号。

设计原因：版本号用于保证后续功能迭代时仍能兼容旧项目。例如未来增加角色关系、分镜时长或镜头轨迹字段时，可以通过版本号区分解析逻辑。

### 4.2 `chapter`

```yaml
chapter:
  chapter_id: "chapter_001"
  chapter_title: "第一章 雨夜来信"
  chapter_index: 1
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `chapter_id` | string | 是 | 章节唯一标识 |
| `chapter_title` | string | 是 | 章节标题 |
| `chapter_index` | number | 是 | 章节顺序，从 1 开始 |

设计原因：小说通常包含多个章节，系统会按章节调用大模型生成。章节字段可以保证场景在数据库、流式事件和前端编辑器中都能追溯到原始章节，避免多章节生成时出现归属混乱。

### 4.3 `scene`

```yaml
scene:
  scene_id: "chapter_001_scene_001"
  scene_index: 1
  title: "雨夜街口"
  summary: "主角在雨夜街口等待一辆迟迟未到的车。"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `scene_id` | string | 是 | 场景唯一标识 |
| `scene_index` | number | 是 | 场景在当前章节中的顺序 |
| `title` | string | 是 | 场景标题 |
| `summary` | string | 是 | 场景内容摘要 |

设计原因：场景是视频化生产的最小单位。每个场景可以独立编辑、重试生成、展示进度和导出。`scene_id` 与 `scene_index` 同时存在，既支持稳定引用，也支持按顺序渲染。

### 4.4 `location`

```yaml
location:
  name: "老城区公交站"
  type: "exterior"
  time_of_day: "night"
  atmosphere: "雨夜、冷清、压抑"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `name` | string | 是 | 场景发生地点 |
| `type` | string | 是 | `interior` 或 `exterior` |
| `time_of_day` | string | 是 | 时间段，如 `morning`、`day`、`dusk`、`night` |
| `atmosphere` | string | 是 | 环境氛围描述 |

设计原因：地点、内外景和时间是影视制作中的基础信息。它们会直接影响布景、灯光、镜头色调和 AI 视频生成效果。将地点信息结构化，可以避免模型只给出抽象剧情而缺少可视化环境。

### 4.5 `characters`

```yaml
characters:
  - name: "林舟"
    role: "protagonist"
    appearance: "黑色风衣，头发被雨水打湿"
    emotional_state: "疲惫、警惕"
    goal: "等待匿名信中提到的接头人"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `name` | string | 是 | 人物名称 |
| `role` | string | 是 | 人物在场景中的职能，如 `protagonist`、`supporting`、`antagonist` |
| `appearance` | string | 是 | 当前场景中可见的外貌或服装 |
| `emotional_state` | string | 是 | 当前情绪状态 |
| `goal` | string | 是 | 当前场景中的行动目标 |

设计原因：小说人物关系经常跨章节变化，大模型容易在长文本中混淆角色状态。`characters` 字段要求模型明确当前场景中谁出现、处于什么状态、想要完成什么目标，这有助于保持剧情连续性，也便于演员表演、分镜设计和视频角色一致性控制。

### 4.6 `action`

```yaml
action:
  - actor: "林舟"
    description: "他低头看了一眼手机，屏幕上仍然没有新消息。"
    purpose: "表现等待过程中的焦虑"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `actor` | string | 是 | 动作执行者 |
| `description` | string | 是 | 可拍摄的动作描述 |
| `purpose` | string | 是 | 该动作的叙事目的 |

设计原因：小说常用心理描写推动情节，但视频需要通过动作表达人物状态。`action` 字段把抽象心理转化为可拍摄行为，并补充动作背后的叙事目的，帮助导演、剪辑和 AI 生成模型理解为什么要展示这个动作。

### 4.7 `dialogues`

```yaml
dialogues:
  - speaker: "林舟"
    text: "你到底想让我看见什么？"
    emotion: "压低声音、怀疑"
    delivery: "slow"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `speaker` | string | 是 | 说话人物 |
| `text` | string | 是 | 台词文本 |
| `emotion` | string | 是 | 台词情绪 |
| `delivery` | string | 是 | 语速或表达方式，如 `slow`、`fast`、`whisper`、`shout` |

设计原因：台词是剧本表达的核心。将台词拆分为说话人、内容、情绪和表达方式，可以服务配音、字幕、表演指导和 AI 语音生成。即使原文没有直接台词，也可以用空数组表示，保持结构一致。

### 4.8 `camera`

```yaml
camera:
  - shot_type: "wide_shot"
    movement: "slow_push_in"
    subject: "林舟站在空荡公交站下"
    description: "雨幕覆盖街道，镜头从远处缓慢推近人物。"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `shot_type` | string | 是 | 景别，如 `wide_shot`、`medium_shot`、`close_up` |
| `movement` | string | 是 | 镜头运动，如 `static`、`pan`、`tilt`、`slow_push_in`、`tracking` |
| `subject` | string | 是 | 镜头主体 |
| `description` | string | 是 | 镜头画面描述 |

设计原因：小说文本不会天然包含镜头语言，但 AI 视频生成和影视化改编都需要明确的画面调度。`camera` 字段把叙事内容转译为景别、运动、主体和画面描述，使生成结果更接近分镜脚本，而不是普通摘要。

### 4.9 `visual_elements`

```yaml
visual_elements:
  - name: "雨水"
    category: "weather"
    description: "密集雨线在路灯下形成白色光幕"
  - name: "红色信封"
    category: "prop"
    description: "被雨水浸湿的红色信封贴在站牌背面"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `name` | string | 是 | 视觉元素名称 |
| `category` | string | 是 | 类型，如 `prop`、`weather`、`lighting`、`set`、`symbol` |
| `description` | string | 是 | 视觉表现说明 |

设计原因：AI 视频生成对具体画面元素非常敏感。`visual_elements` 用于提取小说中的关键道具、天气、光线、空间和象征物，保证生成结果不会遗漏叙事重点。

### 4.10 `audio`

```yaml
audio:
  music: "低频弦乐，制造悬疑感"
  sound_effects:
    - "持续雨声"
    - "远处车辆驶过的水声"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `music` | string | 是 | 背景音乐或音乐风格 |
| `sound_effects` | string[] | 是 | 环境音和关键音效 |

设计原因：声音是短片氛围的重要组成部分。小说中的环境和情绪可以通过音乐与音效进一步表达。该字段也能服务后续配乐、音效生成和剪辑节奏设计。

### 4.11 `emotion_tags`

```yaml
emotion_tags:
  - "suspense"
  - "loneliness"
  - "fear"
```

字段说明：当前场景的情绪标签数组。

设计原因：小说情绪通常隐藏在环境描写和心理活动中。`emotion_tags` 把隐含情绪显式化，方便前端筛选、AI 配乐、AI 视频风格控制和后期节奏规划。

### 4.12 `ai_video_prompt`

```yaml
ai_video_prompt:
  positive: "cinematic rainy night bus stop, lonely man in black coat, wet street, neon reflection, suspense mood"
  negative: "cartoon style, low quality, distorted face, extra limbs, unreadable text"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `positive` | string | 是 | 面向 AI 视频生成的正向提示词 |
| `negative` | string | 是 | 需要避免的画面问题或风格 |

设计原因：结构化剧本最终要服务 AI 视频生成。`ai_video_prompt` 将前面的地点、人物、动作、镜头、视觉元素和情绪压缩成生成模型更容易理解的提示词。正向提示控制画面目标，负向提示减少低质量、风格偏移或人物错误。

### 4.13 `transition`

```yaml
transition:
  type: "cut"
  description: "切到站牌背面，红色信封被雨水打湿。"
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `type` | string | 是 | 转场方式，如 `cut`、`fade`、`dissolve`、`match_cut`、`time_jump` |
| `description` | string | 是 | 转场说明 |

设计原因：视频叙事不是孤立画面的堆叠，场景之间需要衔接。`transition` 字段帮助系统表达场景之间的节奏和连续性，便于后续剪辑或分镜合成。

## 5. 字段约束规则

1. 每个场景 YAML 必须包含 `schema_version`、`chapter` 和 `scene` 三个顶层字段。
2. `chapter.chapter_id` 必须与后端当前章节 ID 一致。
3. `scene.scene_id` 必须在当前章节内唯一，推荐格式为 `chapter_001_scene_001`。
4. `scene.scene_index` 必须从 1 开始递增。
5. 数组字段即使没有内容，也必须输出空数组，例如 `dialogues: []`。
6. 文本字段不得输出 `null`，没有信息时使用空字符串或简短说明。
7. `location.type` 推荐使用 `interior` 或 `exterior`。
8. `transition.type` 推荐使用 `cut`、`fade`、`dissolve`、`match_cut`、`time_jump`。
9. YAML 内容必须可解析，不允许混入 Markdown 代码块标记。
10. 大模型可以根据原文增减数组项数量，但不能新增或改名 Schema 字段。

## 6. 完整示例

```yaml
schema_version: "1.0"
chapter:
  chapter_id: "chapter_001"
  chapter_title: "第一章 雨夜来信"
  chapter_index: 1
scene:
  scene_id: "chapter_001_scene_001"
  scene_index: 1
  title: "雨夜公交站"
  summary: "林舟在雨夜公交站等待匿名信中提到的接头人，并发现站牌背后的红色信封。"
  location:
    name: "老城区公交站"
    type: "exterior"
    time_of_day: "night"
    atmosphere: "冷清、潮湿、悬疑"
  characters:
    - name: "林舟"
      role: "protagonist"
      appearance: "黑色风衣，头发被雨水打湿"
      emotional_state: "疲惫、警惕"
      goal: "确认匿名信中的线索是否真实"
  action:
    - actor: "林舟"
      description: "他把手机屏幕按灭，抬头看向空荡的街道。"
      purpose: "表现漫长等待带来的焦虑"
    - actor: "林舟"
      description: "他绕到站牌背后，发现一只被雨水浸湿的红色信封。"
      purpose: "推动线索出现"
  dialogues:
    - speaker: "林舟"
      text: "你到底想让我看见什么？"
      emotion: "怀疑、压抑"
      delivery: "whisper"
  camera:
    - shot_type: "wide_shot"
      movement: "static"
      subject: "空荡的公交站和雨中的林舟"
      description: "远景展示雨夜街道的冷清，林舟独自站在站台下。"
    - shot_type: "close_up"
      movement: "slow_push_in"
      subject: "红色信封"
      description: "镜头缓慢推近湿透的信封，突出它与灰暗环境的反差。"
  visual_elements:
    - name: "雨水"
      category: "weather"
      description: "密集雨线在路灯下形成白色光幕"
    - name: "红色信封"
      category: "prop"
      description: "信封边缘被雨水泡皱，但红色仍然醒目"
    - name: "路灯倒影"
      category: "lighting"
      description: "黄色灯光在积水中晃动"
  audio:
    music: "低频弦乐，缓慢增强悬疑感"
    sound_effects:
      - "持续雨声"
      - "远处车辆驶过积水的声音"
      - "纸张被揭下的轻微撕裂声"
  emotion_tags:
    - "suspense"
    - "loneliness"
    - "uncertainty"
  ai_video_prompt:
    positive: "cinematic rainy night bus stop, lonely man in black coat, wet street, red envelope, neon reflection, suspense atmosphere, realistic lighting"
    negative: "cartoon style, low quality, distorted hands, extra limbs, unreadable text, oversaturated colors"
  transition:
    type: "cut"
    description: "切到信封内页，上面只有一行被雨水晕开的地址。"
```

## 7. 后端大模型输出关系

后端调用大模型时，外层仍然要求模型返回 JSON，便于 Go 服务稳定解析、落库和推送事件。JSON 中每个场景的 `yaml_content` 字段必须是符合本文档固定 Schema 的 YAML 字符串。

推荐外层 JSON 结构如下：

```json
{
  "chapter_id": "chapter_001",
  "chapter_title": "第一章 雨夜来信",
  "chapter_summary": "本章摘要",
  "scenes": [
    {
      "scene_id": "chapter_001_scene_001",
      "scene_index": 1,
      "title": "雨夜公交站",
      "summary": "场景摘要",
      "yaml_content": "schema_version: \"1.0\"\nchapter:\n  chapter_id: \"chapter_001\"\n..."
    }
  ],
  "carry_context": {
    "previous_chapter_id": "chapter_001",
    "previous_chapter_title": "第一章 雨夜来信",
    "previous_chapter_summary": "本章摘要",
    "character_timeline": [],
    "time_and_place_notes": [],
    "continuity_notes": []
  }
}
```

设计原因：外层 JSON 解决工程解析问题，内层 YAML 解决剧本编辑与导出问题。这样既能保证后端服务稳定，又能让用户在前端直接编辑最终剧本 YAML。

## 8. 为什么采用固定 Schema

原先让大模型为每个章节动态设计 YAML 字段，会带来三个问题：

1. 字段不稳定：不同章节可能出现不同字段名，前端无法提供一致编辑体验。
2. 验证困难：后端难以判断模型输出是否满足题目要求。
3. 演示不可控：相同输入在不同模型或不同轮次下可能产生结构差异。

固定 Schema 更适合当前题目要求。它把“设计原因”集中写在本文档中，而不是让模型每次输出一套新的字段解释。大模型的职责变为：读取小说内容，并把内容填入固定 Schema。

## 9. 总结

本 YAML Schema 的核心设计原因是：把小说中的叙事信息转化为稳定、可编辑、可拍摄、可用于 AI 视频生成的结构化剧本。固定字段保证工程系统可解析，场景级结构保证前端可展示，镜头、动作、台词、视觉元素和情绪字段保证输出具备影视制作价值。