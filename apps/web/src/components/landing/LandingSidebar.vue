<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  activeView: 'new' | 'archive' | 'prompts'
  width: number
  sidebarLabels: {
    new: string
    archive: string
    prompts: string
  }
}>()

const emit = defineEmits<{
  (event: 'select-view', view: 'new' | 'archive' | 'prompts'): void
  (event: 'resize-start', mouseEvent: MouseEvent): void
}>()

const sidebarStyle = computed(() => ({
  width: `${Math.max(72, props.width)}px`,
}))
</script>

<template>
  <aside class="main-sidebar" :class="{ collapsed: props.width <= 150 }" :style="sidebarStyle">
    <div class="resizer" @mousedown="emit('resize-start', $event)"></div>
    <div class="sidebar-logo">
      <span v-if="props.width > 150">E-DIRECTOR</span>
      <span v-else>ED</span>
    </div>
    <nav class="sidebar-nav">
      <button class="nav-item" :class="{ active: props.activeView === 'new' }" type="button" @click="emit('select-view', 'new')">
        <div class="nav-icon">＋</div>
        <div class="nav-label" v-show="props.width > 150">{{ props.sidebarLabels.new }}</div>
      </button>
      <button class="nav-item" :class="{ active: props.activeView === 'archive' }" type="button" @click="emit('select-view', 'archive')">
        <div class="nav-icon">▤</div>
        <div class="nav-label" v-show="props.width > 150">{{ props.sidebarLabels.archive }}</div>
      </button>
      <button class="nav-item" :class="{ active: props.activeView === 'prompts' }" type="button" @click="emit('select-view', 'prompts')">
        <div class="nav-icon">✒</div>
        <div class="nav-label" v-show="props.width > 150">{{ props.sidebarLabels.prompts }}</div>
      </button>
    </nav>
  </aside>
</template>
