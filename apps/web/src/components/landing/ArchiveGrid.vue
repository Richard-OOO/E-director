<script setup lang="ts">
import type { LandingArchiveItem } from '@/components/landing/types'

const props = defineProps<{
  items: LandingArchiveItem[]
  tag: string
}>()

const emit = defineEmits<{
  (event: 'open-project', projectId: string): void
  (event: 'delete-project', projectId: string): void
}>()

const patterns = ['pattern-dots', 'pattern-grid', 'pattern-lines'] as const

const cardPattern = (index: number) => patterns[index % patterns.length]

const deleteProject = (event: MouseEvent, projectId: string) => {
  event.stopPropagation()
  emit('delete-project', projectId)
}
</script>

<template>
  <section class="archive-view">
    <div v-if="props.items.length === 0" class="landing-empty-state">
      <span class="meta">History</span>
      <h3>No projects yet</h3>
      <p>Create a new project first, then generated chapter and scene drafts will appear here.</p>
    </div>

    <div v-else class="archive-grid">
      <article
        v-for="(item, index) in props.items"
        :key="item.project_id"
        class="archive-card"
        :class="cardPattern(index)"
        @click="emit('open-project', item.project_id)"
      >
        <button class="archive-delete" type="button" aria-label="Delete project" @click="deleteProject($event, item.project_id)">×</button>
        <div class="archive-art" aria-hidden="true">
          <span></span>
          <span></span>
          <span></span>
        </div>
        <div class="card-content">
          <span class="card-tag">{{ props.tag }}</span>
          <h3 class="card-title">{{ item.title }}</h3>
          <div class="card-footer">
            <span>{{ item.created_at }}</span>
            <span>{{ item.status }} · {{ item.chapter_count }} ch · {{ item.scene_count }} sc</span>
          </div>
        </div>
      </article>
    </div>
  </section>
</template>
