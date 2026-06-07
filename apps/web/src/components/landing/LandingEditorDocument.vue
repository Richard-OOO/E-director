<script setup lang="ts">
import { nextTick, watch } from 'vue'

import LandingYamlScene from '@/components/landing/LandingYamlScene.vue'
import type { LandingChapter } from '@/components/landing/types'

const props = defineProps<{
  chapters: LandingChapter[]
  activeScene: string
}>()

const emit = defineEmits<{
  (event: 'select-scene', sceneId: string): void
  (event: 'update-scene-yaml', sceneId: string, yaml: string): void
  (event: 'save-scene-yaml', sceneId: string, yaml: string): void
}>()

const scrollToScene = (sceneId: string) => {
  if (!sceneId) return

  void nextTick(() => {
    const target = document.getElementById(`scene-${sceneId}`)
    if (!target) return

    target.scrollIntoView({
      block: 'center',
      behavior: 'smooth',
    })
  })
}

watch(
  () => props.activeScene,
  (sceneId) => {
    scrollToScene(sceneId)
  },
  { immediate: true },
)
</script>

<template>
  <div class="editor-scroll">
    <div class="yaml-document">
      <section v-for="chapter in props.chapters" :key="chapter.id" :id="`chapter-${chapter.id}`" class="yaml-chapter-section">
        <div class="yaml-chapter-heading">
          <span>{{ chapter.num }}</span>
          <h3>{{ chapter.title }}</h3>
        </div>
        <LandingYamlScene
          v-for="scene in chapter.scenes"
          :key="scene.id"
          :id="`scene-${scene.id}`"
          :num="scene.num"
          :title="scene.title"
          :intent="scene.intent"
          :yaml="scene.yaml"
          :class="{ active: props.activeScene === scene.id }"
          @focusin="emit('select-scene', scene.id)"
          @click="emit('select-scene', scene.id)"
          @update-yaml="emit('update-scene-yaml', scene.id, $event)"
          @save-yaml="emit('save-scene-yaml', scene.id, $event)"
        />
      </section>
    </div>
  </div>
</template>
