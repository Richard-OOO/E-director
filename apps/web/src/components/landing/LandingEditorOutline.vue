<script setup lang="ts">
import type { LandingChapter } from '@/components/landing/types'

const props = defineProps<{
  chapters: LandingChapter[]
  activeScene: string
}>()

const emit = defineEmits<{
  (event: 'select-chapter', chapterId: string): void
  (event: 'select-scene', sceneId: string): void
}>()
</script>

<template>
  <aside class="sidebar-editor">
    <h3 class="panel-title">Outline</h3>
    <div class="outline-status">
      <span class="ai-export-dot"></span>
      <span>AI exporting</span>
    </div>
    <div class="chapter-list">
      <div v-for="chapter in props.chapters" :key="chapter.id" class="outline-chapter-block">
        <button class="chapter-header" type="button" @click="emit('select-chapter', chapter.id)">
          <span>{{ chapter.num }} · {{ chapter.title }}</span>
          <span class="chapter-toggle">↴</span>
        </button>
        <div class="scene-list">
          <button
            v-for="scene in chapter.scenes"
            :key="scene.id"
            class="scene-item"
            :class="{ active: props.activeScene === scene.id }"
            type="button"
            @click="emit('select-scene', scene.id)"
          >
            <span class="scene-item-num">{{ scene.num }}</span>
            <span class="scene-item-title">{{ scene.title }}</span>
          </button>
        </div>
      </div>
    </div>
  </aside>
</template>
