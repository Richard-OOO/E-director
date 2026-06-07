<script setup lang="ts">
import type { LandingChapter } from '@/components/landing/types'

const props = defineProps<{
  chapters: LandingChapter[]
  activeScene: string
}>()

const emit = defineEmits<{
  (event: 'select-scene', sceneId: string): void
}>()
</script>

<template>
  <aside class="sidebar-editor">
    <h3 class="panel-title">Outline</h3>
    <div class="chapter-list">
      <div v-for="chapter in props.chapters" :key="chapter.id">
        <div class="chapter-header">
          <span>{{ chapter.num }} · {{ chapter.title }}</span>
          <span class="chapter-toggle">{{ chapter.open ? '−' : '+' }}</span>
        </div>
        <div class="scene-list">
          <button
            v-for="scene in chapter.scenes"
            :key="scene.id"
            class="scene-item"
            :class="{ active: props.activeScene === scene.id }"
            type="button"
            @click="emit('select-scene', scene.id)"
          >
            {{ scene.num }} {{ scene.title }}
          </button>
        </div>
      </div>
    </div>
  </aside>
</template>
