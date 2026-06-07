<script setup lang="ts">
const props = defineProps<{
  progress: number
  statusText: string
  detailText: string
  errorMessage: string
}>()
</script>

<template>
  <section class="stage-view">
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
  </section>
</template>
