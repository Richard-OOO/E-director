<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

import { useRoute, useRouter } from 'vue-router'

import ArchiveGrid from '@/components/landing/ArchiveGrid.vue'
import {
  createProject,
  deleteProject,
  exportCombinedProjectYAML,
  exportProjectYAML,
  fetchMe,
  getProject,
  listProjects,
  logout as logoutSession,
  projectEventsUrl,
  updateSceneYAML,
  type CreateProjectPayload,
  type GenerationEvent,
  type ProjectListItem,
  type StreamedChapterPayload,
} from '@/auth'
import LandingFooterBar from '@/components/landing/LandingFooterBar.vue'
import LandingHeader from '@/components/landing/LandingHeader.vue'
import LandingSidebar from '@/components/landing/LandingSidebar.vue'
import LandingWorkflow from '@/components/landing/LandingWorkflow.vue'
import PromptGrid from '@/components/landing/PromptGrid.vue'
import { snapshotToLandingChapters, streamedChaptersToLandingChapters } from '@/components/landing/snapshot'
import type { LandingChapter, LandingStage } from '@/components/landing/types'
import { useI18n } from '@/i18n/useI18n'

type LandingView = 'new' | 'archive' | 'prompts'

const route = useRoute()
const router = useRouter()
const { t, toggleLang } = useI18n()

const view = ref<LandingView>('new')
const currentStage = ref<LandingStage>('import')
const sidebarWidth = ref(240)
const isSigningOut = ref(false)
const isCreatingProject = ref(false)
const createError = ref('')
const processingProgress = ref(0)
const processingStatus = ref('Waiting for generation to start...')
const processingDetail = ref('Submit a manuscript to start chapter splitting and YAML generation.')
const processingError = ref('')
const isLoading = ref(true)
const currentUser = ref<{ display_name: string; email: string } | null>(null)
const archiveItems = ref<ProjectListItem[]>([])
const editorChapters = ref<LandingChapter[]>([])
const streamedChapters = ref<StreamedChapterPayload[]>([])
const activeProjectId = ref('')
const saveStatus = ref('')
const isExporting = ref(false)
const exportingMode = ref<'batch' | 'combined' | ''>('')
const exportStatus = ref('')
const isGenerationFinished = ref(false)

const stages = ['import', 'processing', 'editor', 'export'] as const
const routeMap: Record<string, { view: LandingView; stage: LandingStage }> = {
  '/workspace/create': { view: 'new', stage: 'import' },
  '/workspace/progress': { view: 'new', stage: 'processing' },
  '/workspace/editor': { view: 'new', stage: 'editor' },
  '/workspace/export': { view: 'new', stage: 'export' },
  '/workspace/history': { view: 'archive', stage: 'import' },
  '/workspace/prompts': { view: 'prompts', stage: 'import' },
}
const currentViewTitle = computed(() => t.value.layout.viewTitles[view.value])

let isResizing = false
let projectEventSource: EventSource | null = null

const selectWorkspaceView = (nextView: LandingView) => {
  view.value = nextView
  if (nextView === 'archive') {
    void router.push('/workspace/history')
    return
  }
  if (nextView === 'prompts') {
    void router.push('/workspace/prompts')
    return
  }
  void router.push(`/workspace/${stagePath(currentStage.value)}`)
}

const hasLandingChapterYAML = (chapters: LandingChapter[]) => {
  const hasYAML = chapters.some((chapter) => chapter.scenes.some((scene) => typeof scene.yaml === 'string' && scene.yaml.trim()))
  return hasYAML
}

const selectWorkspaceStage = (stage: LandingStage) => {
  if (stage === 'export' && !isGenerationFinished.value) {
    currentStage.value = 'processing'
    view.value = 'new'
    processingDetail.value = 'Please wait until all chapters finish generating before exporting YAML.'
    void router.push('/workspace/progress')
    return
  }
  if (stage === 'editor' && !hasLandingChapterYAML(editorChapters.value)) {
    currentStage.value = 'processing'
    view.value = 'new'
    void router.push('/workspace/progress')
    return
  }
  view.value = 'new'
  currentStage.value = stage
  void router.push(`/workspace/${stagePath(stage)}`)
}

const stagePath = (stage: LandingStage) => {
  if (stage === 'import') return 'create'
  if (stage === 'processing') return 'progress'
  return stage
}

const syncRouteState = () => {
  const state = routeMap[route.path]
  if (!state) return
  const nextStage = state.stage === 'editor' && !hasLandingChapterYAML(editorChapters.value) ? 'processing' : state.stage
  view.value = state.view
  currentStage.value = nextStage
}

const resetProcessingState = () => {
  processingProgress.value = 0
  processingStatus.value = 'Waiting for generation to start...'
  processingDetail.value = 'Connecting to the generation event stream.'
  processingError.value = ''
  streamedChapters.value = []
  activeProjectId.value = ''
  saveStatus.value = ''
  exportStatus.value = ''
  isGenerationFinished.value = false
  editorChapters.value = []
}

const closeProjectEvents = () => {
  projectEventSource?.close()
  projectEventSource = null
}

const payloadNumber = (payload: Record<string, unknown>, key: string, fallback = 0) => {
  const value = payload[key]
  return typeof value === 'number' && Number.isFinite(value) ? value : fallback
}

const payloadString = (payload: Record<string, unknown>, key: string, fallback = '') => {
  const value = payload[key]
  return typeof value === 'string' && value.trim() ? value : fallback
}

const progressFromPayload = (payload: Record<string, unknown>, fallback: number) => {
  const progress = payloadNumber(payload, 'overall_progress', fallback)
  return Math.max(0, Math.min(100, Math.round(progress)))
}

type StreamedScene = NonNullable<StreamedChapterPayload['scenes']>[number]

const statusRank = (status?: string) => {
  const normalized = status?.toLowerCase().replace(/_/g, '-')
  if (!normalized) return 0
  if (normalized.includes('failed') || normalized.includes('error')) return 5
  if (normalized.includes('completed') || normalized.includes('done')) return 4
  if (normalized.includes('processing') || normalized.includes('running') || normalized.includes('started')) return 3
  if (normalized.includes('pending')) return 1
  return 2
}

const mergedStatus = (incoming: string, existing?: string) => {
  if (!existing) return incoming
  if (!incoming) return existing
  return statusRank(incoming) >= statusRank(existing) ? incoming : existing
}

const normalizeScene = (scene: unknown, index: number, existing?: StreamedScene): StreamedScene => {
  const raw = typeof scene === 'object' && scene !== null ? (scene as Record<string, unknown>) : {}
  const designReasons = Array.isArray(raw.design_reasons) ? raw.design_reasons : existing?.design_reasons ?? []
  const sceneIndex = payloadNumber(raw, 'scene_index', existing?.scene_index ?? index + 1)
  const yamlContent = payloadString(raw, 'yaml_content', payloadString(raw, 'editable_yaml', payloadString(raw, 'generated_yaml', payloadString(raw, 'yaml', existing?.yaml_content ?? ''))))
  return {
    scene_id: payloadString(raw, 'scene_id', existing?.scene_id ?? `scene-${sceneIndex}`),
    scene_index: sceneIndex,
    title: payloadString(raw, 'title', payloadString(raw, 'scene_title', existing?.title || `Scene ${sceneIndex}`)),
    summary: payloadString(raw, 'summary', payloadString(raw, 'scene_summary', existing?.summary || payloadString(raw, 'status'))),
    yaml_content: yamlContent,
    design_reasons: designReasons,
  }
}

const mergeScenes = (existingScenes: StreamedScene[] = [], incomingScenes: unknown[] = []) => {
  const merged = [...existingScenes]
  incomingScenes.forEach((scene, index) => {
    const raw = typeof scene === 'object' && scene !== null ? (scene as Record<string, unknown>) : {}
    const sceneId = payloadString(raw, 'scene_id')
    const sceneIndex = payloadNumber(raw, 'scene_index')
    const existingIndex = merged.findIndex((item) => (sceneId && item.scene_id === sceneId) || (sceneIndex > 0 && item.scene_index === sceneIndex))
    const normalized = normalizeScene(scene, index, existingIndex >= 0 ? merged[existingIndex] : undefined)
    if (existingIndex >= 0) {
      merged.splice(existingIndex, 1, normalized)
      return
    }
    merged.push(normalized)
  })
  return merged.sort((a, b) => (a.scene_index || 0) - (b.scene_index || 0))
}

const mergeProgress = (payload: Record<string, unknown>, existing?: StreamedChapterPayload) => {
  const incoming = typeof payload.progress === 'object' && payload.progress !== null ? (payload.progress as StreamedChapterPayload['progress']) : undefined
  if (!incoming) return existing?.progress
  return { ...(existing?.progress ?? {}), ...incoming }
}

const scenePayloadToChapterPatch = (payload: Record<string, unknown>, status: string) => {
  const chapterId = payloadString(payload, 'chapter_id')
  if (!chapterId) return null
  const sceneIndex = payloadNumber(payload, 'scene_index', 1)
  return {
    ...payload,
    status: payloadString(payload, 'chapter_status', 'processing'),
    scenes: [
      {
        ...payload,
        status,
        scene_index: sceneIndex,
        scene_id: payloadString(payload, 'scene_id', `scene-${sceneIndex}`),
      },
    ],
  }
}

const normalizeStreamedChapter = (payload: Record<string, unknown>, existing?: StreamedChapterPayload): StreamedChapterPayload | null => {
  const chapterId = payloadString(payload, 'chapter_id', existing?.chapter_id ?? '')
  if (!chapterId) return null
  const chapterIndex = payloadNumber(payload, 'chapter_index', existing?.chapter_index ?? streamedChapters.value.length + 1)
  const incomingScenes = Array.isArray(payload.scenes) ? payload.scenes : []
  const scenes = incomingScenes.length ? mergeScenes(existing?.scenes ?? [], incomingScenes) : existing?.scenes ?? []
  const rawDesignNote = typeof payload.chapter_schema_design_note === 'object' && payload.chapter_schema_design_note !== null ? (payload.chapter_schema_design_note as StreamedChapterPayload['chapter_schema_design_note']) : existing?.chapter_schema_design_note
  const incomingStatus = payloadString(payload, 'status')

  return {
    ...(existing ?? {}),
    ...payload,
    chapter_id: chapterId,
    chapter_title: payloadString(payload, 'chapter_title', existing?.chapter_title || `Chapter ${chapterIndex}`),
    chapter_index: chapterIndex,
    status: mergedStatus(incomingStatus, existing?.status) || 'processing',
    chapter_summary: payloadString(payload, 'chapter_summary', existing?.chapter_summary ?? ''),
    chapter_schema_design_note: rawDesignNote,
    carry_context_summary: payloadString(payload, 'carry_context_summary', existing?.carry_context_summary ?? ''),
    progress: mergeProgress(payload, existing),
    scenes,
  }
}

const upsertStreamedChapter = (payload: Record<string, unknown>) => {
  const chapterId = payloadString(payload, 'chapter_id')
  if (!chapterId) return
  const existingIndex = streamedChapters.value.findIndex((item) => item.chapter_id === chapterId)
  const chapter = normalizeStreamedChapter(payload, existingIndex >= 0 ? streamedChapters.value[existingIndex] : undefined)
  if (!chapter) return
  if (existingIndex >= 0) {
    streamedChapters.value.splice(existingIndex, 1, chapter)
    return
  }
  streamedChapters.value = [...streamedChapters.value, chapter].sort((a, b) => (a.chapter_index || 0) - (b.chapter_index || 0))
}

const hasGeneratedSceneYAML = (chapter: StreamedChapterPayload) => {
  const hasYAML = (chapter.scenes ?? []).some((scene) => typeof scene.yaml_content === 'string' && scene.yaml_content.trim())
  return hasYAML
}

const isFirstChapterReady = () => {
  const firstChapter = streamedChapters.value.find((chapter) => chapter.chapter_index === 1)
  const ready = !!firstChapter && hasGeneratedSceneYAML(firstChapter)
  return ready
}

const syncStreamedChaptersToEditor = () => {
  if (!streamedChapters.value.some(hasGeneratedSceneYAML)) {
    return
  }
  const chapters = streamedChaptersToLandingChapters(streamedChapters.value)
  if (!chapters.length) {
    return
  }
  editorChapters.value = chapters
  if (currentStage.value === 'processing') {
    selectWorkspaceStage('editor')
  }
}

const updateProcessingFromEvent = (event: GenerationEvent) => {
  const payload = event.payload ?? {}
  const progressPayload = typeof payload.progress === 'object' && payload.progress !== null ? (payload.progress as Record<string, unknown>) : payload
  const completedChapters = payloadNumber(progressPayload, 'completed_chapters')
  const totalChapters = payloadNumber(progressPayload, 'total_chapters')
  const chapterTitle = payloadString(payload, 'chapter_title')
  const chapterIndex = payloadNumber(payload, 'chapter_index')

  if (event.event_type === 'generation_started') {
    const chapterCount = payloadNumber(payload, 'chapter_count')
    processingProgress.value = 5
    processingStatus.value = 'Generation started'
    processingDetail.value = chapterCount > 0 ? `Split into ${chapterCount} chapters. Preparing YAML generation.` : 'Preparing chapter analysis and YAML generation.'
    return
  }

  if (event.event_type === 'chapter_started') {
    upsertStreamedChapter(payload)
    processingProgress.value = Math.max(8, progressFromPayload(progressPayload, processingProgress.value))
    processingStatus.value = chapterTitle ? `Generating ${chapterTitle}` : 'Generating chapter YAML'
    processingDetail.value = totalChapters > 0 ? `Chapter ${chapterIndex || completedChapters + 1} of ${totalChapters} is running.` : 'The model is building scene structure and design notes.'
    return
  }

  if (event.event_type === 'scene_started' || event.event_type === 'scene_delta' || event.event_type === 'scene_completed') {
    const scenePatch = scenePayloadToChapterPatch(payload, event.event_type.replace('scene_', ''))
    if (scenePatch) upsertStreamedChapter(scenePatch)
    if (event.event_type === 'scene_completed') syncStreamedChaptersToEditor()
    processingProgress.value = Math.max(processingProgress.value, progressFromPayload(progressPayload, processingProgress.value))
    const sceneIndex = payloadNumber(payload, 'scene_index')
    const sceneTitle = payloadString(payload, 'scene_title', payloadString(payload, 'title'))
    const sceneLabel = sceneTitle || (sceneIndex ? `Scene ${sceneIndex}` : 'Scene')
    processingStatus.value = chapterTitle ? `${sceneLabel} in ${chapterTitle}` : `${sceneLabel} is updating`
    processingDetail.value = event.event_type === 'scene_completed' ? 'A scene YAML block has been generated and merged into the chapter card.' : 'Scene-level generation updates are streaming into the chapter card.'
    return
  }

  if (event.event_type === 'chapter_completed') {
    upsertStreamedChapter(payload)
    syncStreamedChaptersToEditor()
    processingProgress.value = progressFromPayload(progressPayload, processingProgress.value)
    processingStatus.value = chapterTitle ? `Completed ${chapterTitle}` : 'Chapter YAML completed'
    processingDetail.value = totalChapters > 0 ? `${completedChapters} of ${totalChapters} chapters are ready.` : 'Generated scenes have been saved and synced.'
    return
  }

  if (event.event_type === 'schema_summary') {
    upsertStreamedChapter(payload)
    processingStatus.value = chapterTitle ? `Schema ready for ${chapterTitle}` : 'Chapter schema summary ready'
    processingDetail.value = 'YAML structure explanation is available while the scenes finish syncing.'
    return
  }

  if (event.event_type === 'generation_completed') {
    isGenerationFinished.value = true
    processingProgress.value = 100
    processingStatus.value = 'Generation completed'
    processingDetail.value = 'Loading the generated editor workspace.'
    return
  }

  if (event.event_type === 'generation_failed') {
    const error = typeof payload.error === 'object' && payload.error !== null ? (payload.error as Record<string, unknown>) : payload
    processingProgress.value = progressFromPayload(progressPayload, processingProgress.value)
    processingStatus.value = 'Generation failed'
    processingDetail.value = chapterTitle ? `Stopped while generating ${chapterTitle}.` : 'The generation job stopped before completion.'
    processingError.value = payloadString(error, 'message', 'Generation failed')
  }
}

const subscribeProjectEvents = (projectId: string, jobId: string) => {
  closeProjectEvents()
  const url = projectEventsUrl(projectId, jobId)
  projectEventSource = new EventSource(url, { withCredentials: true })

  const handleEvent = async (message: MessageEvent<string>) => {
    try {
      const event = JSON.parse(message.data) as GenerationEvent
      updateProcessingFromEvent(event)

      if (event.event_type === 'generation_completed') {
        closeProjectEvents()
        await refreshProjects()
        const snapshot = await getProject(projectId)
        const chapters = snapshotToLandingChapters(snapshot)
        if (hasLandingChapterYAML(chapters)) {
          activeProjectId.value = projectId
          isGenerationFinished.value = true
          editorChapters.value = chapters
          view.value = 'new'
          selectWorkspaceStage('editor')
        } else if (currentStage.value !== 'editor') {
          processingDetail.value = 'Generation completed, but no editable YAML scenes were returned.'
        }
      }

      if (event.event_type === 'generation_failed') {
        closeProjectEvents()
        await refreshProjects()
      }
    } catch (error) {
      processingError.value = error instanceof Error ? error.message : 'Unable to parse generation event'
    }
  }

  const eventTypes = ['generation_started', 'chapter_started', 'scene_started', 'scene_delta', 'scene_completed', 'schema_summary', 'chapter_completed', 'generation_completed', 'generation_failed']
  eventTypes.forEach((eventType) => {
    projectEventSource?.addEventListener(eventType, (message) => {
      void handleEvent(message as MessageEvent<string>)
    })
  })

  projectEventSource.onerror = (error) => {
    if (!processingError.value && processingProgress.value < 100) {
      processingDetail.value = 'Waiting for the generation event stream to reconnect...'
    }
  }
}

const startResize = (e: MouseEvent) => {
  isResizing = true
  document.body.style.cursor = 'col-resize'
  e.preventDefault()
}

const stopResize = () => {
  isResizing = false
  document.body.style.cursor = 'default'
}

const onResize = (e: MouseEvent) => {
  if (!isResizing) return
  const newWidth = e.clientX
  if (newWidth > 60 && newWidth < 500) {
    sidebarWidth.value = newWidth
  }
}

const prevStage = () => {
  const idx = stages.indexOf(currentStage.value)
  if (idx > 0) selectWorkspaceStage(stages[idx - 1])
}

const nextStage = () => {
  const idx = stages.indexOf(currentStage.value)
  if (idx < stages.length - 1) selectWorkspaceStage(stages[idx + 1])
}

const loadProjectIntoEditor = async (projectId: string) => {
  const snapshot = await getProject(projectId)
  activeProjectId.value = projectId
  isGenerationFinished.value = true
  editorChapters.value = snapshotToLandingChapters(snapshot)
  view.value = 'new'
  selectWorkspaceStage('editor')
}

const handleCreateProject = async (payload: CreateProjectPayload) => {
  isCreatingProject.value = true
  createError.value = ''
  resetProcessingState()
  try {
    const result = await createProject(payload)
    activeProjectId.value = result.project_id
    selectWorkspaceStage('processing')
    await refreshProjects()
    subscribeProjectEvents(result.project_id, result.job_id)
  } catch (error) {
    createError.value = error instanceof Error ? error.message : 'Create project failed'
    selectWorkspaceStage('import')
  } finally {
    isCreatingProject.value = false
  }
}

const handleSignOut = async () => {
  isSigningOut.value = true
  try {
    await logoutSession()
    await router.push('/login')
  } finally {
    isSigningOut.value = false
  }
}

const openProject = (id: string) => {
  void loadProjectIntoEditor(id)
}

const handleEditorChaptersUpdate = (chapters: LandingChapter[]) => {
  editorChapters.value = chapters
}

const handleSaveSceneYAML = async (sceneId: string, yaml: string) => {
  if (!activeProjectId.value) return
  saveStatus.value = 'Saving...'
  try {
    await updateSceneYAML(activeProjectId.value, sceneId, yaml)
    saveStatus.value = 'Saved'
    await refreshProjects()
  } catch (error) {
    saveStatus.value = error instanceof Error ? error.message : 'Save failed'
  }
}

const downloadBlob = (blob: Blob, filename: string) => {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

const handleExportYAML = async (mode: 'batch' | 'combined' = 'batch') => {
  if (!activeProjectId.value) {
    exportStatus.value = 'Open or generate a project before exporting YAML.'
    return
  }
  if (!isGenerationFinished.value) {
    exportStatus.value = 'Please wait until all chapters finish generating before exporting YAML.'
    selectWorkspaceStage('processing')
    return
  }
  isExporting.value = true
  exportingMode.value = mode
  exportStatus.value = mode === 'combined' ? 'Preparing combined YAML file...' : 'Preparing YAML package...'
  try {
    const { blob, filename } = mode === 'combined' ? await exportCombinedProjectYAML(activeProjectId.value) : await exportProjectYAML(activeProjectId.value)
    downloadBlob(blob, filename)
    exportStatus.value = `Downloaded ${filename}`
  } catch (error) {
    exportStatus.value = error instanceof Error ? error.message : 'Export failed'
  } finally {
    isExporting.value = false
    exportingMode.value = ''
  }
}

const handleDeleteProject = async (projectId: string) => {
  await deleteProject(projectId)
  if (activeProjectId.value === projectId) {
    activeProjectId.value = ''
    editorChapters.value = []
  }
  await refreshProjects()
}

watch(() => route.path, syncRouteState, { immediate: true })

const refreshProjects = async () => {
  const data = await listProjects()
  archiveItems.value = data.projects
}

const loadProjects = async () => {
  try {
    const user = await fetchMe()
    currentUser.value = { display_name: user.display_name, email: user.email }
  } catch {
    await router.push('/login')
    return false
  }

  try {
    await refreshProjects()
  } catch {
    // keep local fallback data when backend is unavailable
  }

  return true
}

onMounted(async () => {
  window.addEventListener('mousemove', onResize)
  window.addEventListener('mouseup', stopResize)

  try {
    const ok = await loadProjects()
    if (!ok) return
  } finally {
    isLoading.value = false
  }
})

onUnmounted(() => {
  closeProjectEvents()
  window.removeEventListener('mousemove', onResize)
  window.removeEventListener('mouseup', stopResize)
})
</script>

<template>
  <div class="open-design-shell workspace-layout-shell">
    <LandingSidebar
      :active-view="view"
      :width="sidebarWidth"
      :sidebar-labels="t.layout.sidebar"
      @select-view="selectWorkspaceView"
      @resize-start="startResize"
    />

    <main class="main-content">
      <LandingHeader :title="currentViewTitle" @toggle-lang="toggleLang" />

      <div class="view-scroller">
        <section v-if="view === 'new'" class="workflow-view">
          <LandingWorkflow
            :current-stage="currentStage"
            :stages="stages"
            :is-creating-project="isCreatingProject"
            :create-error="createError"
            :chapters="editorChapters"
            :processing-progress="processingProgress"
            :processing-status="processingStatus"
            :processing-detail="processingDetail"
            :processing-error="processingError"
            :streamed-chapters="streamedChapters"
            :save-status="saveStatus"
            :active-project-id="activeProjectId"
            :is-exporting="isExporting"
            :exporting-mode="exportingMode"
            :export-status="exportStatus"
            :can-export-yaml="isGenerationFinished && !!activeProjectId && hasLandingChapterYAML(editorChapters)"
            @create-project="handleCreateProject"
            @prev-stage="prevStage"
            @next-stage="nextStage"
            @select-stage="selectWorkspaceStage"
            @update-chapters="handleEditorChaptersUpdate"
            @save-scene-yaml="handleSaveSceneYAML"
            @export-yaml="handleExportYAML"
          />
        </section>

        <ArchiveGrid
          v-else-if="view === 'archive'"
          :items="archiveItems"
          :tag="t.archive.tag"
          @open-project="openProject"
          @delete-project="handleDeleteProject"
        />

        <PromptGrid v-else-if="view === 'prompts'" />
      </div>

      <LandingFooterBar
        :is-loading="isLoading"
        :is-signing-out="isSigningOut"
        :user="currentUser"
        @sign-out="handleSignOut"
      />
    </main>
  </div>
</template>
