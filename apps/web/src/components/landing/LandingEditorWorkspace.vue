<script setup lang="ts">
import { ref, watch } from 'vue'

import LandingEditorDocument from '@/components/landing/LandingEditorDocument.vue'
import LandingEditorOutline from '@/components/landing/LandingEditorOutline.vue'
import { landingDemoChapters } from '@/components/landing/fixtures'
import type { LandingChapter } from '@/components/landing/types'

const props = defineProps<{
  chapters: LandingChapter[]
}>()

const chapters = ref<LandingChapter[]>(props.chapters.length > 0 ? props.chapters : landingDemoChapters)

const activeScene = ref(chapters.value[0]?.scenes[0]?.id ?? '')

watch(
  () => props.chapters,
  (next) => {
    chapters.value = next.length > 0 ? next : landingDemoChapters
    activeScene.value = chapters.value[0]?.scenes[0]?.id ?? ''
  },
)

const selectScene = (sceneId: string) => {
  activeScene.value = sceneId
  document.getElementById(`scene-${sceneId}`)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
}
</script>

<template>
  <section class="workspace">
    <LandingEditorOutline :chapters="chapters" :active-scene="activeScene" @select-scene="selectScene" />
    <LandingEditorDocument :chapters="chapters" />
  </section>
</template>
