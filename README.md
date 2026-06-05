# E-director

AI-assisted screenplay creation workspace.

This repository is organized as a polyglot monorepo:

```text
apps/
  gateway/     Go HTTP API Gateway
  ai-worker/   Python AI Worker
  web/         Vue 3 frontend workspace
contracts/
  http/        Browser-facing HTTP API contracts
  schemas/     Request/response JSON schemas
  examples/    HTTP payload examples
docs/
  architecture/  Architecture notes and design docs
infra/
  mysql/       MySQL development infrastructure notes
  redis/       Redis development infrastructure notes
```

## Architecture

- Frontend to Go Gateway: HTTP + JSON.
- Go Gateway to Python AI Worker: HTTP + JSON.
- MySQL stores users, sessions, jobs, and persisted generation records.
- Redis stores transient session state, rate limiting, and job progress caches.
- Python AI Worker processes chapters and scenes in chapter-sized batches to avoid overflowing model context.
- Each generated scene carries both `chapter_id` and `scene_id`.
- Each chapter emits a `schema_design_note`; the final YAML Schema design report is summarized from those chapter notes.

## Local development

Copy the environment example and fill secrets as needed:

```bash
cp .env.example .env
```

Start local services after service Dockerfiles are implemented:

```bash
docker compose up --build
```
