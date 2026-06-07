<script setup lang="ts">
import type { CreateProjectPayload } from '@/auth'
import LandingEditorWorkspace from '@/components/landing/LandingEditorWorkspace.vue'
import LandingExportStage from '@/components/landing/LandingExportStage.vue'
import LandingImportStage from '@/components/landing/LandingImportStage.vue'
import LandingProcessingStage from '@/components/landing/LandingProcessingStage.vue'
import LandingStageNav from '@/components/landing/LandingStageNav.vue'
import type { LandingChapter, LandingStage } from '@/components/landing/types'

const props = defineProps<{
  currentStage: LandingStage
  stages: readonly LandingStage[]
  isCreatingProject: boolean
  createError: string
  chapters: LandingChapter[]
  processingProgress: number
  processingStatus: string
  processingDetail: string
  processingError: string
}>()

const emit = defineEmits<{
  (event: 'create-project', payload: CreateProjectPayload): void
  (event: 'prev-stage'): void
  (event: 'next-stage'): void
  (event: 'select-stage', stage: LandingStage): void
}>()
</script>

<template>
  <div class="view-container landing-workflow-shell">
    <transition name="slide">
      <LandingImportStage
        v-if="props.currentStage === 'import'"
        key="import"
        :is-submitting="props.isCreatingProject"
        :error-message="props.createError"
        @create-project="emit('create-project', $event)"
      />

      <LandingProcessingStage
        v-else-if="props.currentStage === 'processing'"
        key="processing"
        :progress="props.processingProgress"
        :status-text="props.processingStatus"
        :detail-text="props.processingDetail"
        :error-message="props.processingError"
      />

      <LandingEditorWorkspace v-else-if="props.currentStage === 'editor'" key="editor" :chapters="props.chapters" />

      <LandingExportStage v-else-if="props.currentStage === 'export'" key="export" />
    </transition>

    <LandingStageNav
      :current-stage="props.currentStage"
      :stages="props.stages"
      @prev-stage="emit('prev-stage')"
      @next-stage="emit('next-stage')"
      @select-stage="emit('select-stage', $event)"
    />
  </div>
</template>
