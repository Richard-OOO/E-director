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
} from '@/auth'
import LandingFooterBar from '@/components/landing/LandingFooterBar.vue'
import LandingHeader from '@/components/landing/LandingHeader.vue'
import LandingSidebar from '@/components/landing/LandingSidebar.vue'
import LandingWorkflow from '@/components/landing/LandingWorkflow.vue'
import PromptGrid from '@/components/landing/PromptGrid.vue'
import { snapshotToLandingChapters } from '@/components/landing/snapshot'
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

const selectWorkspaceStage = (stage: LandingStage) => {
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
  view.value = state.view
  currentStage.value = state.stage
}

const resetProcessingState = () => {
  processingProgress.value = 0
  processingStatus.value = 'Waiting for generation to start...'
  processingDetail.value = 'Connecting to the generation event stream.'
  processingError.value = ''
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
    processingProgress.value = progressFromPayload(progressPayload, processingProgress.value)
    processingStatus.value = chapterTitle ? `Completed ${chapterTitle}` : 'Chapter YAML completed'
    processingDetail.value = totalChapters > 0 ? `${completedChapters} of ${totalChapters} chapters are ready.` : 'Generated scenes have been saved and synced.'
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
  projectEventSource = new EventSource(projectEventsUrl(projectId, jobId), { withCredentials: true })

  const handleEvent = async (message: MessageEvent<string>) => {
    try {
      const event = JSON.parse(message.data) as GenerationEvent
      updateProcessingFromEvent(event)

      if (event.event_type === 'generation_completed') {
        closeProjectEvents()
        await refreshProjects()
        await loadProjectIntoEditor(projectId)
      }

      if (event.event_type === 'generation_failed') {
        closeProjectEvents()
        await refreshProjects()
      }
    } catch (error) {
      processingError.value = error instanceof Error ? error.message : 'Unable to parse generation event'
    }
  }

  const eventTypes = ['generation_started', 'chapter_started', 'chapter_completed', 'generation_completed', 'generation_failed']
  eventTypes.forEach((eventType) => {
    projectEventSource?.addEventListener(eventType, (message) => {
      void handleEvent(message as MessageEvent<string>)
    })
  })

  projectEventSource.onerror = () => {
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
            @create-project="handleCreateProject"
            @prev-stage="prevStage"
            @next-stage="nextStage"
            @select-stage="selectWorkspaceStage"
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
