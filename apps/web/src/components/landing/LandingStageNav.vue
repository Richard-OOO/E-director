<script setup lang="ts">
import type { LandingStage } from '@/components/landing/types'

const props = defineProps<{
  currentStage: LandingStage
  stages: readonly LandingStage[]
}>()

const emit = defineEmits<{
  (event: 'prev-stage'): void
  (event: 'next-stage'): void
  (event: 'select-stage', stage: LandingStage): void
}>()
</script>

<template>
  <div class="flow-nav-minimal flow-nav-fixed">
    <button class="flow-arrow-min" type="button" @click="emit('prev-stage')">←</button>
    <div class="flow-line-wrap">
      <button
        v-for="stage in props.stages"
        :key="stage"
        class="flow-dot"
        :class="{ active: props.currentStage === stage }"
        type="button"
        @click="emit('select-stage', stage)"
      ></button>
    </div>
    <button class="flow-arrow-min" type="button" @click="emit('next-stage')">→</button>
  </div>
</template>
