<script setup lang="ts">
import { computed } from 'vue'

import type { StreamedChapterPayload } from '@/auth'

const props = defineProps<{
  progress: number
  statusText: string
  detailText: string
  errorMessage: string
  streamedChapters: StreamedChapterPayload[]
}>()

const chapterStatus = (chapter: StreamedChapterPayload) => {
  if (chapter.status) return chapter.status.replace(/_/g, ' ')
  return (chapter.scenes ?? []).some((scene) => scene.yaml_content?.trim()) ? 'completed' : 'processing'
}

const chapterSummary = (chapter: StreamedChapterPayload) => {
  return chapter.chapter_summary || chapter.chapter_schema_design_note?.summary || 'Scene structure is being prepared. YAML will appear here when the chapter completes.'
}

const yamlScenes = (chapter: StreamedChapterPayload) => {
  return (chapter.scenes ?? []).filter((scene) => scene.yaml_content?.trim())
}

const totalSceneCount = computed(() => props.streamedChapters.reduce((sum, chapter) => sum + (chapter.scenes?.length ?? 0), 0))
</script>

<template>
  <section class="stage-view processing-stage-view">
    <div class="processing-panel">
      <span class="meta">Analysis Queue</span>
      <h2>{{ props.errorMessage ? 'Generation Failed' : 'Processing Manuscript' }}</h2>
      <div class="loader-bar">
        <div class="loader-bar-fill" :style="{ width: `${Math.max(4, props.progress)}%` }"></div>
      </div>
      <p>{{ props.errorMessage || props.statusText }}</p>
      <p class="processing-detail">{{ props.detailText }}</p>
      <div class="processing-steps">
        <span :class="{ active: props.progress >= 5 }">Chapter split</span>
        <span :class="{ active: props.progress >= 35 }">Scene map</span>
        <span :class="{ active: props.progress >= 70 }">YAML draft</span>
      </div>
    </div>

    <div v-if="props.streamedChapters.length" class="processing-stream">
      <div class="processing-stream-heading">
        <span class="meta">Live Chapters</span>
        <strong>{{ props.streamedChapters.length }} chapters · {{ totalSceneCount }} scenes</strong>
      </div>
      <article v-for="chapter in props.streamedChapters" :key="chapter.chapter_id" class="processing-chapter-card">
        <div class="processing-chapter-head">
          <div>
            <span class="meta">Chapter {{ chapter.chapter_index }}</span>
            <h3>{{ chapter.chapter_title }}</h3>
          </div>
          <span class="processing-status-pill">{{ chapterStatus(chapter) }}</span>
        </div>
        <p>{{ chapterSummary(chapter) }}</p>
        <div class="processing-chapter-progress" v-if="chapter.progress?.overall_progress !== undefined">
          <span>Overall {{ Math.round(chapter.progress.overall_progress ?? 0) }}%</span>
          <span>{{ chapter.progress.completed_chapters ?? 0 }}/{{ chapter.progress.total_chapters ?? '?' }} chapters</span>
        </div>
        <div v-if="chapter.scenes?.length" class="processing-scene-list">
          <div v-for="scene in chapter.scenes" :key="scene.scene_id" class="processing-scene-row">
            <div class="processing-scene-title">
              <strong>{{ scene.title || `Scene ${scene.scene_index}` }}</strong>
              <span>SC. {{ String(scene.scene_index).padStart(2, '0') }}</span>
            </div>
            <p>{{ scene.summary || 'Waiting for scene summary.' }}</p>
            <pre v-if="scene.yaml_content?.trim()">{{ scene.yaml_content }}</pre>
          </div>
        </div>
        <div v-else class="processing-scene-empty">Waiting for scene YAML from this chapter...</div>
        <div v-if="yamlScenes(chapter).length" class="processing-yaml-count">{{ yamlScenes(chapter).length }} YAML scene files ready</div>
      </article>
    </div>
  </section>
</template>
