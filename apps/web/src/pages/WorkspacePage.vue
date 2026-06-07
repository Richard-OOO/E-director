<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

import { useRoute, useRouter } from 'vue-router'

import ArchiveGrid from '@/components/landing/ArchiveGrid.vue'
import {
  createProject,
  fetchMe,
  getProject,
  listProjects,
  logout as logoutSession,
  projectEventsUrl,
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

const debugLog = (label: string, payload?: unknown) => {
  console.debug(`[e-director:workspace] ${label}`, payload)
}

const summarizeEditorChapters = (chapters: LandingChapter[]) => ({
  chapter_count: chapters.length,
  scene_count: chapters.reduce((sum, chapter) => sum + chapter.scenes.length, 0),
  yaml_scene_count: chapters.reduce((sum, chapter) => sum + chapter.scenes.filter((scene) => scene.yaml?.trim()).length, 0),
  first_chapter: chapters[0]
    ? {
        id: chapters[0].id,
        title: chapters[0].title,
        scene_count: chapters[0].scenes.length,
        first_scene_yaml_length: chapters[0].scenes[0]?.yaml?.length ?? 0,
      }
    : null,
})

const summarizeStreamedChapters = (chapters: StreamedChapterPayload[]) => ({
  chapter_count: chapters.length,
  chapters: chapters.map((chapter) => ({
    chapter_id: chapter.chapter_id,
    chapter_index: chapter.chapter_index,
    scene_count: chapter.scenes?.length ?? 0,
    yaml_lengths: chapter.scenes?.map((scene) => scene.yaml_content?.length ?? 0) ?? [],
  })),
})

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
  debugLog('hasLandingChapterYAML', { hasYAML, chapters: summarizeEditorChapters(chapters) })
  return hasYAML
}

const selectWorkspaceStage = (stage: LandingStage) => {
  debugLog('selectWorkspaceStage requested', { stage, currentStage: currentStage.value, editorChapters: summarizeEditorChapters(editorChapters.value) })
  if (stage === 'editor' && !hasLandingChapterYAML(editorChapters.value)) {
    debugLog('selectWorkspaceStage blocked editor', { reason: 'no editor YAML yet' })
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
  debugLog('syncRouteState', { path: route.path, routeStage: state.stage, nextStage, editorChapters: summarizeEditorChapters(editorChapters.value) })
  view.value = state.view
  currentStage.value = nextStage
}

const resetProcessingState = () => {
  processingProgress.value = 0
  processingStatus.value = 'Waiting for generation to start...'
  processingDetail.value = 'Connecting to the generation event stream.'
  processingError.value = ''
  streamedChapters.value = []
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

const upsertStreamedChapter = (payload: Record<string, unknown>) => {
  const chapterId = payloadString(payload, 'chapter_id')
  debugLog('upsertStreamedChapter', {
    chapterId,
    payloadKeys: Object.keys(payload),
    sceneCount: Array.isArray(payload.scenes) ? payload.scenes.length : 0,
    yamlLengths: Array.isArray(payload.scenes) ? payload.scenes.map((scene) => (typeof scene === 'object' && scene !== null && typeof (scene as Record<string, unknown>).yaml_content === 'string' ? ((scene as Record<string, string>).yaml_content.length) : 0)) : [],
  })
  if (!chapterId) return
  const chapter = payload as StreamedChapterPayload
  const existingIndex = streamedChapters.value.findIndex((item) => item.chapter_id === chapterId)
  if (existingIndex >= 0) {
    streamedChapters.value.splice(existingIndex, 1, chapter)
    debugLog('streamedChapters updated existing', summarizeStreamedChapters(streamedChapters.value))
    return
  }
  streamedChapters.value = [...streamedChapters.value, chapter].sort((a, b) => (a.chapter_index || 0) - (b.chapter_index || 0))
  debugLog('streamedChapters inserted', summarizeStreamedChapters(streamedChapters.value))
}

const hasGeneratedSceneYAML = (chapter: StreamedChapterPayload) => {
  const hasYAML = (chapter.scenes ?? []).some((scene) => typeof scene.yaml_content === 'string' && scene.yaml_content.trim())
  debugLog('hasGeneratedSceneYAML', {
    chapter_id: chapter.chapter_id,
    chapter_index: chapter.chapter_index,
    hasYAML,
    yamlLengths: chapter.scenes?.map((scene) => scene.yaml_content?.length ?? 0) ?? [],
  })
  return hasYAML
}

const isFirstChapterReady = () => {
  const firstChapter = streamedChapters.value.find((chapter) => chapter.chapter_index === 1)
  const ready = !!firstChapter && hasGeneratedSceneYAML(firstChapter)
  debugLog('isFirstChapterReady', { ready, firstChapter, streamedChapters: summarizeStreamedChapters(streamedChapters.value) })
  return ready
}

const syncStreamedChaptersToEditor = () => {
  debugLog('syncStreamedChaptersToEditor start', summarizeStreamedChapters(streamedChapters.value))
  if (!isFirstChapterReady()) {
    debugLog('syncStreamedChaptersToEditor blocked', { reason: 'first chapter has no YAML yet' })
    return
  }
  const chapters = streamedChaptersToLandingChapters(streamedChapters.value)
  if (!chapters.length) {
    debugLog('syncStreamedChaptersToEditor blocked', { reason: 'mapped landing chapters are empty' })
    return
  }
  editorChapters.value = chapters
  debugLog('syncStreamedChaptersToEditor applied', summarizeEditorChapters(editorChapters.value))
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
    processingProgress.value = Math.max(8, progressFromPayload(progressPayload, processingProgress.value))
    processingStatus.value = chapterTitle ? `Generating ${chapterTitle}` : 'Generating chapter YAML'
    processingDetail.value = totalChapters > 0 ? `Chapter ${chapterIndex || completedChapters + 1} of ${totalChapters} is running.` : 'The model is building scene structure and design notes.'
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
  debugLog('subscribeProjectEvents', { projectId, jobId, url })
  projectEventSource = new EventSource(url, { withCredentials: true })

  const handleEvent = async (message: MessageEvent<string>) => {
    try {
      debugLog('sse raw message', { type: message.type, data: message.data })
      const event = JSON.parse(message.data) as GenerationEvent
      debugLog('sse parsed event', {
        event_type: event.event_type,
        sequence: event.sequence,
        payloadKeys: Object.keys(event.payload ?? {}),
        payload: event.payload,
      })
      updateProcessingFromEvent(event)

      if (event.event_type === 'generation_completed') {
        closeProjectEvents()
        await refreshProjects()
        const snapshot = await getProject(projectId)
        debugLog('generation_completed snapshot data', snapshot)
        const chapters = snapshotToLandingChapters(snapshot)
        debugLog('generation_completed mapped chapters', summarizeEditorChapters(chapters))
        if (hasLandingChapterYAML(chapters)) {
          editorChapters.value = chapters
          view.value = 'new'
          selectWorkspaceStage('editor')
        } else if (currentStage.value !== 'editor') {
          debugLog('generation_completed no editable YAML scenes', { snapshot, chapters: summarizeEditorChapters(chapters) })
          processingDetail.value = 'Generation completed, but no editable YAML scenes were returned.'
        }
      }

      if (event.event_type === 'generation_failed') {
        closeProjectEvents()
        await refreshProjects()
      }
    } catch (error) {
      debugLog('sse handle error', error)
      processingError.value = error instanceof Error ? error.message : 'Unable to parse generation event'
    }
  }

  const eventTypes = ['generation_started', 'chapter_started', 'schema_summary', 'chapter_completed', 'generation_completed', 'generation_failed']
  eventTypes.forEach((eventType) => {
    projectEventSource?.addEventListener(eventType, (message) => {
      void handleEvent(message as MessageEvent<string>)
    })
  })

  projectEventSource.onerror = (error) => {
    debugLog('sse error', { error, readyState: projectEventSource?.readyState, progress: processingProgress.value, processingError: processingError.value })
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
  debugLog('loadProjectIntoEditor start', { projectId })
  const snapshot = await getProject(projectId)
  debugLog('loadProjectIntoEditor snapshot', snapshot)
  editorChapters.value = snapshotToLandingChapters(snapshot)
  debugLog('loadProjectIntoEditor mapped chapters', summarizeEditorChapters(editorChapters.value))
  view.value = 'new'
  selectWorkspaceStage('editor')
}

const handleCreateProject = async (payload: CreateProjectPayload) => {
  isCreatingProject.value = true
  createError.value = ''
  resetProcessingState()
  debugLog('handleCreateProject start', { title: payload.title, language: payload.language, source_type: payload.source_type, hasFile: !!payload.file, contentLength: payload.content?.length ?? 0 })
  try {
    const result = await createProject(payload)
    debugLog('createProject result', result)
    selectWorkspaceStage('processing')
    await refreshProjects()
    subscribeProjectEvents(result.project_id, result.job_id)
  } catch (error) {
    debugLog('handleCreateProject error', error)
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

watch(() => route.path, syncRouteState, { immediate: true })

const refreshProjects = async () => {
  const data = await listProjects()
  debugLog('refreshProjects result', data)
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
            @create-project="handleCreateProject"
            @prev-stage="prevStage"
            @next-stage="nextStage"
            @select-stage="selectWorkspaceStage"
            @update-chapters="handleEditorChaptersUpdate"
          />
        </section>

        <ArchiveGrid v-else-if="view === 'archive'" :items="archiveItems" :tag="t.archive.tag" @open-project="openProject" />

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
