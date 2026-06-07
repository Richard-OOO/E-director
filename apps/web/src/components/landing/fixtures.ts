import type { LandingChapter, PromptCard } from '@/components/landing/types'

export const landingDemoChapters: LandingChapter[] = [
  {
    id: 'c1',
    num: 'CH. 01',
    title: 'The Silent Threshold',
    open: true,
    scenes: [
      { id: 's1', num: '01', title: 'Arrival at the Gate', intent: 'Establish Isolation' },
      { id: 's2', num: '02', title: 'The First Encounter', intent: 'Create Tension' },
    ],
  },
  {
    id: 'c2',
    num: 'CH. 02',
    title: 'Echoes of the Past',
    open: false,
    scenes: [{ id: 's3', num: '03', title: 'Memory Corridor', intent: 'Surreal Flashback' }],
  },
]

export const landingPromptCards: PromptCard[] = [
  { id: 1, num: '01', text: 'Act as a professional cinematographer focusing on the Golden Hour aesthetics. Use high dynamic range.', tags: ['Lighting', 'Visuals'] },
  { id: 2, num: '02', text: 'Analyze character motivation through micro-expressions. Direct the AI to capture subtle eye movements.', tags: ['Performance', 'AI'] },
  { id: 3, num: '03', text: 'Implement a non-linear narrative structure. Flashbacks should be triggered by cues.', tags: ['Story', 'Logic'] },
  { id: 4, num: '04', text: 'Write concise scene direction with readable beat markers and production-ready structure.', tags: ['Workflow', 'Script'] },
]
