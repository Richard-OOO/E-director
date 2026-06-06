<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'

import { useI18n } from '@/i18n/useI18n'

const { lang, t, toggleLang } = useI18n()

const view = ref<'new' | 'archive' | 'prompts'>('new')
const currentStage = ref<'import' | 'processing' | 'editor' | 'export'>('import')
const sidebarWidth = ref(240)
const activeScene = ref('s1')

const stages = ['import', 'processing', 'editor', 'export'] as const

const currentViewTitle = computed(() => t.value.layout.viewTitles[view.value])

const chapters = ref([
  {
    id: 'c1',
    num: 'CH. 01',
    title: 'The Silent Threshold',
    open: true,
    scenes: [
      { id: 's1', num: '01', title: 'Arrival at the Gate', intent: 'Establish Isolation', comment: 'Use a tight 50mm lens. The AI will prioritize deep shadows to emphasize character solitude.' },
      { id: 's2', num: '02', title: 'The First Encounter', intent: 'Create Tension', comment: 'Slow dolly zoom during the reveal. Maintain a cold color grade.' },
    ],
  },
  {
    id: 'c2',
    num: 'CH. 02',
    title: 'Echoes of the Past',
    open: false,
    scenes: [{ id: 's3', num: '03', title: 'Memory Corridor', intent: 'Surreal Flashback', comment: 'Soft focus edges. The AI suggests a higher frame rate (60fps).' }],
  },
])

const archiveItems = [
  { id: 1, title: 'Obsidian Dreams', date: '2026.05.12', pattern: 'pattern-dots' },
  { id: 2, title: 'Neon Monolith', date: '2026.04.08', pattern: 'pattern-grid' },
  { id: 3, title: 'The Paper Forest', date: '2026.03.22', pattern: 'pattern-lines' },
  { id: 4, title: 'Aperture Sky', date: '2026.02.15', pattern: 'pattern-dots' },
]

const prompts = [
  { id: 1, num: '01', text: 'Act as a professional cinematographer focusing on the Golden Hour aesthetics. Use high dynamic range.', tags: ['Lighting', 'Visuals'] },
  { id: 2, num: '02', text: 'Analyze character motivation through micro-expressions. Direct the AI to capture subtle eye movements.', tags: ['Performance', 'AI'] },
  { id: 3, num: '03', text: 'Implement a non-linear narrative structure. Flashbacks should be triggered by cues.', tags: ['Story', 'Logic'] },
]

const startProcessing = () => {
  currentStage.value = 'processing'
  window.setTimeout(() => {
    currentStage.value = 'editor'
  }, 2000)
}

const scrollToScene = (id: string) => {
  activeScene.value = id
  document.getElementById('scene-' + id)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

let isResizing = false

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
  if (!isResizing) {
    return
  }

  const newWidth = e.clientX
  if (newWidth > 60 && newWidth < 500) {
    sidebarWidth.value = newWidth
  }
}

const prevStage = () => {
  const idx = stages.indexOf(currentStage.value)
  if (idx > 0) {
    currentStage.value = stages[idx - 1]
  }
}

const nextStage = () => {
  const idx = stages.indexOf(currentStage.value)
  if (idx < stages.length - 1) {
    currentStage.value = stages[idx + 1]
  }
}

onMounted(() => {
  window.addEventListener('mousemove', onResize)
  window.addEventListener('mouseup', stopResize)
})

onUnmounted(() => {
  window.removeEventListener('mousemove', onResize)
  window.removeEventListener('mouseup', stopResize)
})
</script>

<template>
  <div class="open-design-shell">
    <aside class="main-sidebar" :style="{ width: sidebarWidth + 'px' }">
      <div class="resizer" @mousedown="startResize"></div>
      <div class="sidebar-logo">
        <span v-if="sidebarWidth > 150">E-DIRECTOR</span>
        <span v-else>ED</span>
      </div>
      <nav class="sidebar-nav">
        <button class="nav-item" :class="{ active: view === 'new' }" @click="view = 'new'" type="button">
          <div class="nav-icon">＋</div>
          <div class="nav-label" v-show="sidebarWidth > 150">{{ t.layout.sidebar.new }}</div>
        </button>
        <button class="nav-item" :class="{ active: view === 'archive' }" @click="view = 'archive'" type="button">
          <div class="nav-icon">▤</div>
          <div class="nav-label" v-show="sidebarWidth > 150">{{ t.layout.sidebar.archive }}</div>
        </button>
        <button class="nav-item" :class="{ active: view === 'prompts' }" @click="view = 'prompts'" type="button">
          <div class="nav-icon">✒</div>
          <div class="nav-label" v-show="sidebarWidth > 150">{{ t.layout.sidebar.prompts }}</div>
        </button>
      </nav>
    </aside>

    <main class="main-content">
      <header class="open-design-header">
        <div class="view-title">{{ currentViewTitle }}</div>
        <div class="nav-meta">
          <button class="lang-switch" type="button" @click="toggleLang">
            <span :class="{ active: lang === 'en' }">EN</span>
            <span class="sep">|</span>
            <span :class="{ active: lang === 'zh' }">中</span>
          </button>
        </div>
      </header>

      <div class="view-scroller">
        <section v-if="view === 'new'" class="workflow-view">
          <div class="view-container">
            <transition name="slide">
              <section v-if="currentStage === 'import'" key="import" class="stage-view">
                <div class="card-import">
                  <h2>{{ t.import.h2 }}</h2>
                  <p>{{ t.import.p }}</p>
                  <div class="upload-zone">
                    <h4>{{ t.import.upload }}</h4>
                    <p class="upload-hint">Word, PDF, TXT</p>
                  </div>
                  <button class="btn-submit" type="button" @click="startProcessing">{{ t.import.cta }}</button>
                </div>
              </section>

              <section v-else-if="currentStage === 'processing'" key="processing" class="stage-view">
                <div class="processing-panel">
                  <h2>{{ t.processing.h2 }}</h2>
                  <div class="loader-bar">
                    <div class="loader-bar-fill"></div>
                  </div>
                  <p>{{ t.processing.p }}</p>
                </div>
              </section>

              <section v-else-if="currentStage === 'editor'" key="editor" class="workspace">
                <aside class="sidebar-editor">
                  <h3 class="panel-title">{{ t.editor.outline }}</h3>
                  <div class="chapter-list">
                    <div v-for="chap in chapters" :key="chap.id">
                      <div class="chapter-header" @click="chap.open = !chap.open">
                        <span>{{ chap.num }} · {{ chap.title }}</span>
                        <span class="chapter-toggle">{{ chap.open ? '−' : '+' }}</span>
                      </div>
                      <div v-show="chap.open" class="scene-list">
                        <div
                          v-for="scene in chap.scenes"
                          :key="scene.id"
                          class="scene-item"
                          :class="{ active: activeScene === scene.id }"
                          @click="scrollToScene(scene.id)"
                        >
                          {{ scene.num }} {{ scene.title }}
                        </div>
                      </div>
                    </div>
                  </div>
                </aside>

                <div class="editor-scroll">
                  <div class="yaml-document">
                    <div v-for="chap in chapters" :key="'doc-' + chap.id">
                      <div v-for="scene in chap.scenes" :key="'sec-' + scene.id" :id="'scene-' + scene.id" class="yaml-section">
                        <div class="yaml-code">
                          <span class="yaml-key">scene:</span> <span class="yaml-val">{{ scene.num }}</span>
                          <br />
                          <span class="yaml-key">title:</span> <span class="yaml-val">{{ scene.title }}</span>
                          <br />
                          <span class="yaml-key">intent:</span> <span class="yaml-val">{{ scene.intent }}</span>
                          <br />
                          <span class="yaml-key">location:</span> <span class="yaml-val">INT. STUDIO - DAY</span>
                        </div>
                        <div class="yaml-commentary">
                          <span class="meta">AI Observation</span>
                          <h6>{{ scene.title }}</h6>
                          <p>{{ scene.comment }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </section>

              <section v-else-if="currentStage === 'export'" key="export" class="stage-view">
                <div class="export-panel">
                  <h2>Ready to Deliver</h2>
                  <div class="export-grid">
                    <div class="export-card">
                      <h3>Director's Cut</h3>
                      <p>Export as a cinematic high-fidelity preview.</p>
                      <button class="btn-submit" type="button">Export PDF</button>
                    </div>
                    <div class="export-card">
                      <h3>YAML Raw Script</h3>
                      <p>Export the structured script for engine import.</p>
                      <button class="btn-submit" type="button">Export YAML</button>
                    </div>
                  </div>
                </div>
              </section>
            </transition>

            <div class="flow-nav-minimal">
              <button class="flow-arrow-min" type="button" @click="prevStage">←</button>
              <div class="flow-line-wrap">
                <button
                  v-for="s in stages"
                  :key="s"
                  class="flow-dot"
                  :class="{ active: currentStage === s }"
                  type="button"
                  @click="currentStage = s"
                ></button>
              </div>
              <button class="flow-arrow-min" type="button" @click="nextStage">→</button>
            </div>
          </div>
        </section>

        <section v-else-if="view === 'archive'" class="archive-view">
          <div class="archive-grid">
            <div v-for="item in archiveItems" :key="item.id" class="archive-card" :class="item.pattern">
              <div class="card-content">
                <span class="card-tag">{{ t.archive.tag }}</span>
                <h3 class="card-title">{{ item.title }}</h3>
                <div class="card-footer">
                  <span>{{ item.date }}</span>
                  <span>{{ t.archive.footer }}</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section v-else-if="view === 'prompts'" class="prompts-view">
          <div class="prompt-grid">
            <div v-for="p in prompts" :key="p.id" class="prompt-card">
              <div class="prompt-num">{{ p.num }}</div>
              <div class="prompt-content">
                <p class="prompt-text">{{ p.text }}</p>
                <div class="prompt-tags">
                  <span v-for="tag in p.tags" :key="tag" class="prompt-tag">{{ tag }}</span>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>
