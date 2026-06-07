<script setup lang="ts">
import LandingYamlScene from '@/components/landing/LandingYamlScene.vue'
import type { LandingChapter } from '@/components/landing/types'

const props = defineProps<{
  chapters: LandingChapter[]
}>()
</script>

<template>
  <div class="editor-scroll">
    <div class="yaml-document">
      <template v-for="chapter in props.chapters" :key="chapter.id">
        <LandingYamlScene
          v-for="scene in chapter.scenes"
          :key="scene.id"
          :id="`scene-${scene.id}`"
          :num="scene.num"
          :title="scene.title"
          :intent="scene.intent"
          :comment="scene.comment"
          :yaml="scene.yaml"
        />
      </template>
    </div>
  </div>
</template>
