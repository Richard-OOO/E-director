export type ApiEnvelope<T> = {
  code: number
  msg: string
  data: T
}

export type AuthUser = {
  id: string
  email: string
  display_name: string
}

export type CreateProjectPayload = {
  title: string
  language: string
  source_type: string
  content?: string
  file?: File
}

export type CreateProjectResult = {
  project_id: string
  job_id: string
  status: string
}

export type ProjectListItem = {
  project_id: string
  title: string
  language: string
  status: string
  chapter_count: number
  scene_count: number
  current_job_id: string
  created_at: string
  updated_at: string
}

export type ProjectSnapshotProject = {
  id?: string
  user_id?: string
  title?: string
  language?: string
  source_type?: string
  source_text?: string
  source_text_hash?: string
  chapter_count?: number
  current_job_id?: string
  created_at?: string
  updated_at?: string
}

export type ProjectSnapshotJob = Record<string, unknown>

export type ProjectSnapshotChapter = {
  id?: string
  job_id?: string
  project_id?: string
  chapter_id?: string
  chapter_index?: number
  title?: string
  content?: string
  content_hash?: string
  summary?: string
  status?: string
  scene_count?: number
  completed_scene_count?: number
  failed_scene_count?: number
  design_note_yaml?: string
  error_code?: string
  error_message?: string
  retryable?: boolean
  started_at?: string
  finished_at?: string
  created_at?: string
  updated_at?: string
}

export type ProjectSnapshotScene = {
  id?: string
  job_id?: string
  project_id?: string
  chapter_record_id?: string
  chapter_id?: string
  scene_id?: string
  scene_index?: number
  title?: string
  summary?: string
  status?: string
  generated_yaml?: string
  editable_yaml?: string
  design_reason_yaml?: string
  yaml_hash?: string
  is_edited?: boolean
  edited_by_user_id?: string
  edited_at?: string
  generation_started_at?: string
  generation_finished_at?: string
  error_code?: string
  error_message?: string
  retryable?: boolean
  retry_count?: number
  created_at?: string
  updated_at?: string
}

export type ProjectSnapshot = {
  project?: ProjectSnapshotProject
  job?: ProjectSnapshotJob
  chapters?: ProjectSnapshotChapter[]
  scenes?: ProjectSnapshotScene[]
}

export type PromptItem = {
  id: string
  user_id: string
  key: string
  name: string
  content: string
  default_content: string
  is_default: boolean
  is_editable: boolean
  created_at: string
  updated_at: string
}

type AuthData = {
  user: AuthUser
}

type SendCodeData = {
  cooldown_seconds: number
  expires_in_seconds: number
}

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, '') ?? ''

export function apiUrl(path: string) {
  return `${apiBaseUrl}${path}`
}

async function request<T>(path: string, options: RequestInit = {}) {
  const headers = new Headers(options.headers)

  if (!headers.has('Content-Type') && options.body && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }

  const url = apiUrl(path)
  const response = await fetch(url, { ...options, headers, credentials: 'include' })
  const result = (await response.json()) as ApiEnvelope<T>

  if (!response.ok || result.code !== 200) {
    throw new Error(result.msg || 'Request failed')
  }

  return result.data
}

export const sendRegisterCode = async (payload: { email: string }) => {
  return request<SendCodeData>('/api/v1/auth/send-code', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export const login = async (payload: { email: string; password: string }) => {
  return request<AuthData>('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export const register = async (payload: { name: string; email: string; password: string; verification_code: string }) => {
  return request<AuthData>('/api/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export const logout = () => {
  return request<Record<string, never>>('/api/v1/auth/logout', { method: 'POST' })
}

export const fetchMe = () => request<AuthUser>('/api/v1/auth/me')

export const createProject = (payload: CreateProjectPayload) => {
  if (payload.file) {
    const formData = new FormData()
    formData.set('title', payload.title)
    formData.set('language', payload.language)
    formData.set('source_type', payload.source_type)
    formData.set('file', payload.file)
    return request<CreateProjectResult>('/api/v1/projects', {
      method: 'POST',
      body: formData,
    })
  }

  return request<CreateProjectResult>('/api/v1/projects', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export const listProjects = () => request<{ projects: ProjectListItem[] }>('/api/v1/projects')

export const getProject = (projectId: string) => request<ProjectSnapshot>(`/api/v1/projects/${encodeURIComponent(projectId)}`)

export const updateSceneYAML = (projectId: string, sceneId: string, yaml: string) => {
  return request<{ scene_id: string }>(`/api/v1/projects/${encodeURIComponent(projectId)}/scenes/${encodeURIComponent(sceneId)}`, {
    method: 'PATCH',
    body: JSON.stringify({ yaml }),
  })
}

export const deleteProject = (projectId: string) => {
  return request<{ project_id: string }>(`/api/v1/projects/${encodeURIComponent(projectId)}`, { method: 'DELETE' })
}

export const exportProjectYAML = async (projectId: string) => {
  const url = apiUrl(`/api/v1/projects/${encodeURIComponent(projectId)}/export/yaml`)
  const response = await fetch(url, { credentials: 'include' })
  if (!response.ok) {
    const result = (await response.json().catch(() => null)) as ApiEnvelope<unknown> | null
    throw new Error(result?.msg || 'Export failed')
  }
  const blob = await response.blob()
  const disposition = response.headers.get('Content-Disposition') ?? ''
  const filenameMatch = disposition.match(/filename="?([^";]+)"?/)
  return { blob, filename: filenameMatch?.[1] ?? 'project.yaml.zip' }
}

export const exportCombinedProjectYAML = async (projectId: string) => {
  const url = apiUrl(`/api/v1/projects/${encodeURIComponent(projectId)}/export/yaml/combined`)
  const response = await fetch(url, { credentials: 'include' })
  if (!response.ok) {
    const result = (await response.json().catch(() => null)) as ApiEnvelope<unknown> | null
    throw new Error(result?.msg || 'Export failed')
  }
  const blob = await response.blob()
  const disposition = response.headers.get('Content-Disposition') ?? ''
  const filenameMatch = disposition.match(/filename="?([^";]+)"?/)
  return { blob, filename: filenameMatch?.[1] ?? 'project.yaml' }
}

export const listPrompts = () => request<{ prompts: PromptItem[] }>('/api/v1/prompts')

export const createPrompt = (payload: { name: string; content: string }) => {
  return request<PromptItem>('/api/v1/prompts', { method: 'POST', body: JSON.stringify(payload) })
}

export const updatePrompt = (promptId: string, payload: { name: string; content: string }) => {
  return request<PromptItem>(`/api/v1/prompts/${encodeURIComponent(promptId)}`, { method: 'PATCH', body: JSON.stringify(payload) })
}

export const deletePrompt = (promptId: string) => {
  return request<Record<string, never>>(`/api/v1/prompts/${encodeURIComponent(promptId)}`, { method: 'DELETE' })
}

export const resetPrompt = (promptId: string) => {
  return request<PromptItem>(`/api/v1/prompts/${encodeURIComponent(promptId)}/reset`, { method: 'POST' })
}

export type GenerationEventPayload = Record<string, unknown>

export type GenerationDesignReason = {
  reason_id?: string
  target_path?: string
  field_name?: string
  reason_type?: string
  title?: string
  description?: string
  hover_text?: string
}

export type StreamedChapterScene = {
  scene_id: string
  scene_index: number
  title: string
  summary: string
  yaml_content: string
  design_reasons: GenerationDesignReason[]
}

export type StreamedChapterPayload = GenerationEventPayload & {
  chapter_id: string
  chapter_title: string
  chapter_index: number
  status?: string
  chapter_summary?: string
  chapter_schema_design_note?: {
    summary?: string
    key_reasons?: Array<{ field_name?: string; reason?: string }>
  }
  carry_context_summary?: string
  progress?: {
    completed_chapters?: number
    total_chapters?: number
    completed_scenes?: number
    total_scenes?: number
    overall_progress?: number
  }
  scenes?: StreamedChapterScene[]
}

export type GenerationEvent = {
  event_id: string
  event_type: string
  project_id: string
  job_id: string
  sequence: number
  created_at: string
  payload: GenerationEventPayload
}

export const projectEventsUrl = (projectId: string, jobId: string) => {
  return apiUrl(`/api/v1/projects/${encodeURIComponent(projectId)}/events?job_id=${encodeURIComponent(jobId)}`)
}
