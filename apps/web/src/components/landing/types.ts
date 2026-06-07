export type LandingStage = 'import' | 'processing' | 'editor' | 'export'

export type LandingView = 'new' | 'archive' | 'prompts'

export type LandingScene = {
  id: string
  num: string
  title: string
  intent: string
  comment: string
  yaml?: string
}

export type LandingChapter = {
  id: string
  num: string
  title: string
  open: boolean
  scenes: LandingScene[]
}

export type LandingArchiveItem = {
  project_id: string
  title: string
  status: string
  chapter_count: number
  scene_count: number
  created_at: string
}

export type PromptCard = {
  id: number
  num: string
  text: string
  tags: string[]
}
