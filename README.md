<div align="center">

# E-director

### AI 小说转剧本创作工作台

将多章节小说自动改编为结构化、可编辑、可导出的剧本 YAML，帮助作者快速获得可继续打磨的剧本初稿。

<p>
  <img alt="Vue 3" src="https://img.shields.io/badge/Vue-3-42b883?style=for-the-badge&logo=vuedotjs&logoColor=white" />
  <img alt="Go" src="https://img.shields.io/badge/Go-1.22-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="MySQL" src="https://img.shields.io/badge/MySQL-8.4-4479A1?style=for-the-badge&logo=mysql&logoColor=white" />
  <img alt="Redis" src="https://img.shields.io/badge/Redis-7.2-DC382D?style=for-the-badge&logo=redis&logoColor=white" />
  <img alt="Docker" src="https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
</p>

<p>
  <a href="docs/architecture/yaml-schema-design.md"><strong>YAML Schema 文档</strong></a>
  ·
  <a href="docs/architecture/chapter-based-ai-pipeline.md"><strong>AI 生成流水线</strong></a>
  ·
  <a href="docker-compose.yml"><strong>一键启动</strong></a>
  ·
  <a href="https://www.bilibili.com/video/BV1PTEh6yE21/?spm_id_from=333.1387.list.card_archive.click&vd_source=5925ae2f2d04bfcc99578332e52daea4"><strong>演示视频</strong></a>
</p>

<a href="https://www.bilibili.com/video/BV1PTEh6yE21/?spm_id_from=333.1387.list.card_archive.click&vd_source=5925ae2f2d04bfcc99578332e52daea4">
  <img alt="E-director 演示视频" src="https://img.shields.io/badge/点击观看-Bilibili%20演示视频-00AEEF?style=for-the-badge&logo=bilibili&logoColor=white" />
</a>

</div>

<div style="width: 100%; height: 2px; margin: 20px 0; background: linear-gradient(90deg, transparent, #007DFF, transparent);"></div>

## 项目简介

E-director 是一个面向小说作者的 AI 辅助剧本创作工具。作者只需要导入 小说文本，系统就会按章节理解原文内容，并把小说改写为由多个场景组成的结构化剧本。

与普通“生成一段剧本”的工具不同，E-director 更强调可编辑、可追踪和可继续生产。生成结果不是一整段不可控文本，而是稳定的 YAML 结构：每个场景都有明确的章节归属、场景编号、人物状态、动作、台词、镜头、音效、情绪标签和 AI 视频提示词。作者可以在工作台中逐场景查看、修改和导出。

<div style="width: 100%; height: 2px; margin: 20px 0; background: linear-gradient(90deg, transparent, #007DFF, transparent);"></div>

## 体验亮点

<table>
  <tr>
    <td width="50%">
      <h3>从小说到场景剧本</h3>
      <p>系统将长篇小说拆解为章节，再把章节转换为多个可拍摄场景。场景是前端展示、人工编辑、后端保存和导出的最小单位。</p>
    </td>
    <td width="50%">
      <h3>实时生成进度</h3>
      <p>工作台通过 SSE 接收后端生成事件。用户不需要等待黑盒任务结束，可以实时看到章节、场景和状态变化。</p>
    </td>
  </tr>
  <tr>
    <td width="50%">
      <h3>可编辑 YAML</h3>
      <p>每个场景都有独立 YAML 内容，作者可以直接修改台词、动作、镜头、音效、情绪和提示词，再保存为最终剧本初稿。</p>
    </td>
    <td width="50%">
      <h3>固定 Schema</h3>
      <p>剧本输出使用固定 YAML Schema，避免不同章节字段不一致，使编辑、保存、导出和后续工具接入都更加稳定。</p>
    </td>
  </tr>
</table>

## 工作流

```text
导入小说
   ↓
创建项目并切分章节
   ↓
后端按章节调用 AI 生成
   ↓
每章拆分为多个剧本场景
   ↓
前端通过 SSE 实时接收进度
   ↓
作者在工作台编辑 YAML
   ↓
导出结构化剧本初稿
```

这个流程贴近真实改编习惯：先理解章节，再拆场景，再处理人物、动作、台词和镜头。它不会强迫作者一次接受完整结果，而是把剧本拆成可以逐步检查和打磨的小单元。

## YAML 剧本格式

E-director 使用固定 YAML Schema 表达剧本场景。固定结构让输出更稳定，也让前端编辑器、后端存储和导出流程都能围绕同一份数据契约工作。

<p align="center">
  <a href="docs/architecture/yaml-schema-design.md">
    <img alt="Open YAML Schema" src="https://img.shields.io/badge/打开-YAML%20Schema%20设计文档-007DFF?style=for-the-badge&logo=readthedocs&logoColor=white" />
  </a>
</p>

每个场景 YAML 大致包含以下结构：

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

完整字段定义、约束规则、完整示例和设计原因见 [docs/architecture/yaml-schema-design.md](docs/architecture/yaml-schema-design.md)。

<div style="width: 100%; height: 2px; margin: 20px 0; background: linear-gradient(90deg, transparent, #007DFF, transparent);"></div>

## 演示视频

<p align="center">
  <a href="https://www.bilibili.com/video/BV1PTEh6yE21/?spm_id_from=333.1387.list.card_archive.click&vd_source=5925ae2f2d04bfcc99578332e52daea4">
    <img alt="E-director Bilibili 演示视频" src="https://img.shields.io/badge/Bilibili-点击观看完整演示-00AEEF?style=for-the-badge&logo=bilibili&logoColor=white" />
  </a>
</p>

<p align="center">
  <a href="https://www.bilibili.com/video/BV1PTEh6yE21/?spm_id_from=333.1387.list.card_archive.click&vd_source=5925ae2f2d04bfcc99578332e52daea4">
    <strong> https://www.bilibili.com/video/BV1PTEh6yE21/?spm_id_from=333.1387.list.card_archive.click&vd_source=5925ae2f2d04bfcc99578332e52daea4</strong>
  </a>
</p>

> 这是项目演示视频入口，GitHub README 会保留可点击链接，评审点击后可直接跳转播放。

## 前端交互

前端位于 [apps/web](apps/web)，使用 Vue 3 + Vite + TypeScript 实现。页面围绕“作者改编工作台”组织，而不是单纯的表单提交。

- 登录 / 注册：使用邮箱验证码完成账号注册与登录。
- 项目入口：创建小说改编项目，导入小说文本或文档内容。
- 工作台侧栏：展示历史项目、章节和场景，方便在不同内容之间切换。
- 生成进度流：通过 EventSource 连接后端 SSE，持续接收项目生成事件。
- YAML 编辑区：展示当前场景 YAML，支持直接编辑并保存。
- 提示词管理：调整剧本生成提示词和导出风格提示词。
- YAML 导出：将生成并修改后的剧本导出为结构化 YAML。

交互设计上，E-director 尽量减少“等待 AI 完成”的割裂感。生成过程中，用户可以看到章节和场景逐步出现；生成后，用户可以继续在同一个工作台中编辑和导出，不需要切换到其他工具整理结果。

## 后端架构

后端位于 [apps/backend](apps/backend)，使用 Go 1.22 实现。它承担 API、认证、项目管理、文档解析、AI 生成编排、SSE 推送和 YAML 导出等职责。

```text
apps/backend/
  cmd/gateway/          服务入口
  internal/config/      环境变量和配置加载
  internal/server/      Gin 路由注册
  internal/handler/     HTTP API、SSE、项目和提示词接口
  internal/service/     业务逻辑、文档解析、LLM 调用、生成编排
  internal/domain/      领域对象和事件结构
  internal/models/      MySQL / Redis 数据访问
```

后端生成链路大致如下：

```text
Browser Workspace
  ↓ HTTP / JSON / SSE
Nginx
  ├─ /       → Web Frontend
  └─ /api/*  → Go Gateway

Go Gateway
  ├─ MySQL
  │   └─ 用户、项目、章节、场景、生成任务、提示词
  ├─ Redis
  │   └─ 会话、验证码、临时进度缓存
  └─ OpenAI-compatible LLM
      └─ DashScope / OpenAI-compatible API

Go Gateway
  ↓ SSE 推送生成事件
Browser Workspace
```

MySQL 保存用户、会话、项目、章节、场景、生成任务和提示词等持久数据。Redis 用于验证码、会话状态和临时进度缓存。生产部署中，nginx 作为统一入口：`/` 转发到前端，`/api/` 和 `/health` 转发到后端，SSE 路由关闭代理缓冲以保证事件实时到达。

## 系统架构

```text
                 ┌────────────────────────┐
                 │        Browser         │
                 │ Vue 3 + Vite Frontend  │
                 └───────────┬────────────┘
                             │ HTTP / SSE
                             ▼
                 ┌────────────────────────┐
                 │         Nginx          │
                 │ / → web, /api → gateway│
                 └───────────┬────────────┘
                             │ /api
                             ▼
                 ┌────────────────────────┐
                 │       Go Gateway       │
                 │       Gin + GORM       │
                 └───────────┬────────────┘
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
          ▼                  ▼                  ▼
┌────────────────┐  ┌────────────────┐  ┌────────────────────────┐
│     MySQL      │  │     Redis      │  │ OpenAI-compatible LLM  │
│ persistent DB  │  │ session/cache  │  │ DashScope / other API  │
└────────────────┘  └────────────────┘  └────────────────────────┘
```

仓库中保留 [apps/ai-worker](apps/ai-worker) 作为后续独立 AI Worker 的容器占位；当前主要生成链路由 Go 后端直接调用 OpenAI-compatible 接口完成。

<div style="width: 100%; height: 2px; margin: 20px 0; background: linear-gradient(90deg, transparent, #007DFF, transparent);"></div>

## 原创设计

E-director 的核心业务设计和实现集中在以下部分：

- 面向小说改编的“项目 → 章节 → 场景 → YAML”数据模型。
- 场景级 YAML Schema 以及对应的设计说明文档。
- 章节级生成流程，避免长篇小说一次性生成导致上下文过长。
- 同时在当前章节中会加入之前章节的所有重要内容，防止上下文丢失。
- 稳定的 `chapter_id` / `scene_id` 标识，支持实时进度、局部保存和导出。
- 面向作者编辑体验的前端工作台。
- 后端项目、章节、场景、提示词、生成任务和 SSE 事件编排。
- 面向 AI 视频生成的字段扩展：镜头、视觉元素、声音、情绪标签和视频提示词。
- Docker Compose + nginx 的一键部署方案。

第三方依赖主要提供框架、数据库、缓存、构建、文档解析和模型接口能力；小说转剧本的流程设计、数据结构、前后端交互和 YAML Schema 为项目内实现。

## 目录结构

```text
apps/
  backend/      Go HTTP API、生成编排、数据访问、SSE
  web/          Vue 3 前端工作台
  ai-worker/    Python AI Worker 占位服务
contracts/
  http/         HTTP API 契约
  schemas/      请求、响应和事件 JSON Schema
  examples/     API 示例载荷
docs/
  architecture/ 架构说明与 YAML Schema 设计文档
infra/
  nginx/        nginx 反向代理配置
  mysql/        MySQL 说明
  redis/        Redis 说明
```

## 主要依赖

### 后端

| 依赖 | 用途 |
| --- | --- |
| Go 1.22 | 后端运行时和构建环境 |
| Gin | HTTP API 路由和中间件 |
| GORM + MySQL Driver | MySQL ORM 和数据持久化 |
| go-redis | Redis 会话、验证码和缓存访问 |
| x/crypto | 密码哈希和安全相关能力 |
| rsc.io/pdf | PDF 文本提取 |
| OpenAI-compatible / DashScope API | 小说章节到剧本场景的 AI 生成 |

### 前端

| 依赖 | 用途 |
| --- | --- |
| Vue 3 | 前端 UI 框架 |
| Vue Router | 页面路由和工作台导航 |
| Vite | 前端开发服务器和构建工具 |
| TypeScript / vue-tsc | 类型检查 |
| Monaco Editor | YAML 编辑体验 |

### 部署

| 依赖 | 用途 |
| --- | --- |
| Docker / Docker Compose | 一键启动完整服务栈 |
| Nginx | 统一入口、反向代理和 SSE 转发 |
| MySQL 8.4 | 持久化数据库 |
| Redis 7.2 | 会话和进度缓存 |
| Node 22 | 前端构建和静态服务运行时 |

## 一键启动

准备环境变量：

```bash
cp .env.example .env
```

至少需要在 `.env` 中配置模型密钥：

```env
OPENAI_API_KEY=<your-api-key>
# 或者
DASHSCOPE_API_KEY=<your-dashscope-api-key>
```

启动完整服务栈：

```bash
docker compose up --build
```

如果已经构建过镜像，也可以后台启动：

```bash
docker compose up -d
```

默认访问地址：

```text
http://localhost/
```

健康检查：

```bash
curl http://localhost/health
```

常用排查命令：

```bash
docker compose ps
docker compose logs gateway
docker compose logs nginx
docker compose logs web
```

> Windows Docker Desktop 用户如果遇到 `dockerDesktopLinuxEngine` pipe 不存在，通常是 Docker Desktop 尚未启动。启动 Docker Desktop 后重新执行 `docker compose up --build` 即可。

## 环境变量

完整示例见 [.env.example](.env.example)。常用变量如下：

| 变量 | 说明 |
| --- | --- |
| `NGINX_PORT` | nginx 对外端口，默认 `80` |
| `MYSQL_DATABASE` / `MYSQL_USER` / `MYSQL_PASSWORD` | Docker Compose MySQL 配置 |
| `REDIS_ADDR` / `REDIS_DB` | Redis 地址和数据库编号 |
| `TOKEN_SECRET` | 登录会话相关密钥，部署时应替换 |
| `SESSION_COOKIE_SECURE` | HTTPS 部署时建议设为 `true` |
| `CORS_ALLOWED_ORIGINS` | 允许访问后端的前端来源 |
| `SMTP_HOST` / `SMTP_USERNAME` / `SMTP_PASSWORD` | 邮箱验证码发送配置 |
| `OPENAI_API_BASE` | OpenAI-compatible API 地址 |
| `OPENAI_API_KEY` / `DASHSCOPE_API_KEY` | 模型调用密钥 |
| `OPENAI_MODEL` | 生成使用的模型名称 |
| `OPENAI_MAX_TOKENS` | 单次生成最大 token 数 |

请不要把真实 `.env` 密钥提交到仓库。

## 相关文档

- [剧本 YAML Schema 设计文档](docs/architecture/yaml-schema-design.md)
- [章节级 AI 生成流水线说明](docs/architecture/chapter-based-ai-pipeline.md)
- [HTTP 内部契约说明](docs/architecture/http-internal-contract.md)
- [环境变量示例](.env.example)
- [Docker Compose 配置](docker-compose.yml)
- [nginx 反向代理配置](infra/nginx/default.conf)

## 验证方式

后端测试：

```bash
go test -C apps/backend ./...
```

前端构建：

```bash
npm --prefix apps/web run build
```

Docker Compose 配置检查：

```bash
docker compose config
```

完整部署验证：

```bash
docker compose up --build
```

## 后续改进

后续计划引入 RAG 检索增强能力，在并行生成章节时按需召回人物关系、关键事件和前文摘要，从而提升生成速度，并保持剧情上下文连贯。



