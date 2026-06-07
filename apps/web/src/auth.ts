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
  content: string
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
  ID?: string
  user_id?: string
  UserID?: string
  title?: string
  Title?: string
  language?: string
  Language?: string
  source_type?: string
  SourceType?: string
  source_text?: string
  SourceText?: string
  chapter_count?: number
  ChapterCount?: number
  current_job_id?: string
  CurrentJobID?: string
  created_at?: string
  CreatedAt?: string
  updated_at?: string
  UpdatedAt?: string
}

export type ProjectSnapshotJob = Record<string, unknown>

export type ProjectSnapshotChapter = {
  id?: string
  ID?: string
  chapter_id?: string
  ChapterID?: string
  chapter_index?: number
  ChapterIndex?: number
  title?: string
  Title?: string
  content?: string
  Content?: string
  summary?: string
  Summary?: string
  status?: string
  Status?: string
  scene_count?: number
  SceneCount?: number
  design_note_y_a_m_l?: string
  DesignNoteYAML?: string
}

export type ProjectSnapshotScene = {
  id?: string
  ID?: string
  chapter_id?: string
  ChapterID?: string
  scene_id?: string
  SceneID?: string
  scene_index?: number
  SceneIndex?: number
  title?: string
  Title?: string
  summary?: string
  Summary?: string
  status?: string
  Status?: string
  generated_y_a_m_l?: string
  GeneratedYAML?: string
  editable_y_a_m_l?: string
  EditableYAML?: string
  design_reason_y_a_m_l?: string
  DesignReasonYAML?: string
}

export type ProjectSnapshot = {
  project?: ProjectSnapshotProject
  Project?: ProjectSnapshotProject
  job?: ProjectSnapshotJob
  Job?: ProjectSnapshotJob
  chapters?: ProjectSnapshotChapter[]
  Chapters?: ProjectSnapshotChapter[]
  scenes?: ProjectSnapshotScene[]
  Scenes?: ProjectSnapshotScene[]
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

  if (!headers.has('Content-Type') && options.body) {
    headers.set('Content-Type', 'application/json')
  }

  const response = await fetch(apiUrl(path), { ...options, headers, credentials: 'include' })
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
  return request<CreateProjectResult>('/api/v1/projects', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export const listProjects = () => request<{ projects: ProjectListItem[] }>('/api/v1/projects')

export const getProject = (projectId: string) => request<ProjectSnapshot>(`/api/v1/projects/${encodeURIComponent(projectId)}`)

export type GenerationEventPayload = Record<string, unknown>

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
