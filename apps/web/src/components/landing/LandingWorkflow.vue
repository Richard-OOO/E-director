<script setup lang="ts">
import type { CreateProjectPayload, StreamedChapterPayload } from '@/auth'
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
  streamedChapters: StreamedChapterPayload[]
  saveStatus?: string
  activeProjectId: string
  isExporting: boolean
  exportingMode: 'batch' | 'combined' | ''
  exportStatus: string
  canExportYaml: boolean
}>()

const emit = defineEmits<{
  (event: 'create-project', payload: CreateProjectPayload): void
  (event: 'prev-stage'): void
  (event: 'next-stage'): void
  (event: 'select-stage', stage: LandingStage): void
  (event: 'update-chapters', chapters: LandingChapter[]): void
  (event: 'save-scene-yaml', sceneId: string, yaml: string): void
  (event: 'export-yaml', mode: 'batch' | 'combined'): void
}>()

const forwardSaveSceneYAML = (sceneId: string, yaml: string) => {
  emit('save-scene-yaml', sceneId, yaml)
}
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
        :streamed-chapters="props.streamedChapters"
      />

      <div
        v-else-if="props.currentStage === 'editor' && props.chapters.some((chapter) => chapter.scenes.some((scene) => scene.yaml?.trim()))"
        key="editor"
        class="editor-stage-shell"
      >
        <LandingEditorWorkspace
          :chapters="props.chapters"
          @update-chapters="emit('update-chapters', $event)"
          @save-scene-yaml="forwardSaveSceneYAML"
        />
        <div v-if="props.saveStatus" class="editor-save-status">{{ props.saveStatus }}</div>
      </div>

      <LandingProcessingStage
        v-else-if="props.currentStage === 'editor'"
        key="editor-waiting"
        :progress="props.processingProgress"
        :status-text="props.processingStatus"
        :detail-text="props.processingDetail"
        :error-message="props.processingError"
        :streamed-chapters="props.streamedChapters"
      />

      <LandingExportStage
        v-else-if="props.currentStage === 'export'"
        key="export"
        :active-project-id="props.activeProjectId"
        :is-exporting="props.isExporting"
        :exporting-mode="props.exportingMode"
        :export-status="props.exportStatus"
        :can-export-yaml="props.canExportYaml"
        @export-yaml="emit('export-yaml', $event)"
      />
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
