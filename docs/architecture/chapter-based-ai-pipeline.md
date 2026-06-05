# HTTP-based chapter pipeline

## 1. Background

E-director needs to process long novels with three or more chapters. If the whole novel is passed into the model at once, the context window can overflow and chapter, character, and scene boundaries become blurred. The AI Worker therefore processes content in chapter-sized batches.

## 2. Core principles

1. Do not feed the whole novel to the model at once.
2. Prefer one chapter per request.
3. Split a chapter into chunks when it is too long.
4. Every scene must carry `chapter_id` and `scene_id`.
5. Every chapter generates a `schema_design_note`.
6. The final YAML Schema design document is summarized from chapter-level notes, not from rereading the entire novel.

## 3. Pipeline flow

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

## 4. Chapter-level processing

Each chapter request should include:

- novel title
- chapter ID
- chapter title
- chapter content
- fixed output requirements
- scene ID naming rules
- YAML field explanations
- design-reason output requirements

The response should include:

```json
{
  "chapter_id": "chapter_001",
  "chapter_title": "First Chapter",
  "chapter_summary": "Chapter summary",
  "scenes": [],
  "chapter_schema_design_note": {
    "summary": "Why this chapter needs this YAML shape",
    "key_reasons": []
  }
}
```

## 5. Scene labels

Every scene must include:

```yaml
chapter_id: chapter_001
chapter_title: First Chapter
scene_id: chapter_001_scene_001
scene_index: 1
```

This keeps scenes stable even if the worker processes chapters concurrently.

## 6. Chapter design notes

Every chapter should emit a `chapter_schema_design_note`:

```yaml
schema_design_note:
  summary: "This chapter emphasizes rainy-night visuals and emotional change, so camera, visual_elements, and emotion_tags are important."
  key_reasons:
    - field_name: "camera"
      reason: "Converts environmental description into shot language for downstream generation."
    - field_name: "emotion_tags"
      reason: "Makes emotional progression explicit for rendering and editing."
```

## 7. Final document generation

The final YAML Schema design document should consume only chapter-level notes:

```text
chapter_001.schema_design_note
chapter_002.schema_design_note
chapter_003.schema_design_note
...
```

## 8. Concurrency strategy

Recommended controls:

- `max_concurrent_chapters`: chapter concurrency
- `max_concurrent_scenes`: scene concurrency inside a chapter
- Redis stores job, chapter, and scene state
- Go Gateway exposes progress and job status through HTTP endpoints

## 9. Error handling

When a single scene fails, return an error object containing `chapter_id` and `scene_id` so the whole job remains recoverable:

```json
{
  "code": "SCENE_GENERATION_FAILED",
  "chapter_id": "chapter_001",
  "scene_id": "chapter_001_scene_003",
  "retryable": true
}
```

## 10. Competition-facing description

E-director uses chapter-sized context boundaries rather than sending the entire novel to the model at once. Each chapter is processed independently, scenes are tagged with chapter and scene IDs, and the final schema design document is summarized from chapter-level design notes so the system can explain why each field exists.
