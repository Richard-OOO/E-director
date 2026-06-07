<script setup lang="ts">
const props = defineProps<{
  title: string
}>()

const emit = defineEmits<{
  (event: 'toggle-lang'): void
}>()
</script>

<template>
  <header class="open-design-header">
    <div class="view-title">{{ props.title }}</div>
    <div class="nav-meta">
      <button class="lang-switch" type="button" @click="emit('toggle-lang')">
        <span>EN</span>
        <span class="sep">|</span>
        <span>中</span>
      </button>
    </div>
  </header>
</template>
