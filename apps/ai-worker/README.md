# Python AI Worker

Python AI Worker for E-director.

Responsibilities:

- Split long novels by chapter.
- Feed chapters, not the whole novel, into the LLM.
- Extract scenes per chapter.
- Generate scene YAML with `chapter_id` and `scene_id`.
- Generate per-chapter `schema_design_note`.
- Expose HTTP endpoints for the Go Gateway.
