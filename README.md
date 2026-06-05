# E-director

AI-assisted screenplay creation workspace.

This repository is organized as a polyglot monorepo:

```text
apps/
  gateway/    Go + Gin API Gateway
  ai-worker/  Python AI Worker with gRPC streaming
  web/        Vue 3 frontend workspace
contracts/
  proto/      Backend internal gRPC contracts
  schemas/    Browser-facing HTTP/SSE JSON Schemas
  examples/   Contract examples
docs/
  architecture/ Architecture notes and YAML Schema design docs
infra/
  redis/      Redis development infrastructure notes
```

## Architecture

- Frontend to Go Gateway: HTTP + SSE.
- Go Gateway to Python AI Worker: gRPC server streaming.
- Redis stores sessions, job state, stream replay data, and generated artifacts.
- Python AI Worker feeds the LLM by chapter to avoid overflowing model context.
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
