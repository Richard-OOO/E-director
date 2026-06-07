<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import LandingEditorDocument from '@/components/landing/LandingEditorDocument.vue'
import LandingEditorOutline from '@/components/landing/LandingEditorOutline.vue'
import type { LandingChapter } from '@/components/landing/types'

const props = defineProps<{
  chapters: LandingChapter[]
}>()

const emit = defineEmits<{
  (event: 'update-chapters', chapters: LandingChapter[]): void
}>()

const chapters = ref<LandingChapter[]>(props.chapters)
const activeScene = ref(chapters.value[0]?.scenes[0]?.id ?? '')

const sceneIds = computed(() => chapters.value.flatMap((chapter) => chapter.scenes.map((scene) => scene.id)))

const selectScene = (sceneId: string) => {
  if (!sceneIds.value.includes(sceneId) || activeScene.value === sceneId) return
  activeScene.value = sceneId
}

const updateSceneYaml = (sceneId: string, yaml: string) => {
  if (!sceneIds.value.includes(sceneId)) return

  chapters.value = chapters.value.map((chapter) => ({
    ...chapter,
    scenes: chapter.scenes.map((scene) => (scene.id === sceneId ? { ...scene, yaml } : scene)),
  }))
  emit('update-chapters', chapters.value)
}

watch(
  () => props.chapters,
  (next: LandingChapter[]) => {
    chapters.value = next

    const nextSceneIds = next.flatMap((chapter) => chapter.scenes.map((scene) => scene.id))
    if (!nextSceneIds.includes(activeScene.value)) {
      activeScene.value = nextSceneIds[0] ?? ''
    }
  },
)
</script>

<template>
  <section class="workspace">
    <LandingEditorOutline :chapters="chapters" :active-scene="activeScene" @select-scene="selectScene" />
    <LandingEditorDocument
      :chapters="chapters"
      :active-scene="activeScene"
      @select-scene="selectScene"
      @update-scene-yaml="updateSceneYaml"
    />
  </section>
</template>
