<script setup lang="ts">
import type { StreamedChapterPayload } from '@/auth'

const props = defineProps<{
  progress: number
  statusText: string
  detailText: string
  errorMessage: string
  streamedChapters: StreamedChapterPayload[]
}>()
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
      <article v-for="chapter in props.streamedChapters" :key="chapter.chapter_id" class="processing-chapter-card">
        <span class="meta">Chapter {{ chapter.chapter_index }}</span>
        <h3>{{ chapter.chapter_title }}</h3>
        <p>{{ chapter.chapter_summary || chapter.chapter_schema_design_note?.summary || 'YAML structure is ready.' }}</p>
        <div class="processing-scene-list">
          <div v-for="scene in chapter.scenes || []" :key="scene.scene_id" class="processing-scene-row">
            <strong>{{ scene.title }}</strong>
            <p>{{ scene.summary }}</p>
            <pre>{{ scene.yaml_content }}</pre>
          </div>
        </div>
      </article>
    </div>
  </section>
</template>
