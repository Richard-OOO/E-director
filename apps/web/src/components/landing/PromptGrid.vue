<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { createPrompt, deletePrompt, listPrompts, resetPrompt, updatePrompt, type PromptItem } from '@/auth'

const prompts = ref<PromptItem[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const errorMessage = ref('')
const editorOpen = ref(false)
const editingPrompt = ref<PromptItem | null>(null)
const draftName = ref('')
const draftContent = ref('')
const isContentFocused = ref(false)

const stylePlaceholder = `🔘 剧本节奏：[极快（多动作） / 舒缓（保留环境描写）]
🔘 对话风格：[偏向口语化 / 保留原著文学性]
🔘 额外要求输入框：[“请重点刻画主角的心理挣扎，多加一些微表情的描写”]`

const sortedPrompts = computed(() => [...prompts.value].sort((a, b) => Number(b.is_default) - Number(a.is_default) || a.name.localeCompare(b.name)))

const generationPrompt = computed(() => prompts.value.find((prompt) => prompt.key === 'style_export') ?? null)

const promptBody = (prompt: PromptItem) => prompt.content || prompt.default_content

const isGenerationPrompt = (prompt: PromptItem) => prompt.key === 'style_export'

const contentPlaceholder = computed(() => (isContentFocused.value ? '' : stylePlaceholder))

const refreshPrompts = async () => {
  isLoading.value = true
  errorMessage.value = ''
  try {
    const data = await listPrompts()
    prompts.value = data.prompts
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to load prompts'
  } finally {
    isLoading.value = false
  }
}

const openCreateEditor = () => {
  editingPrompt.value = null
  draftName.value = ''
  draftContent.value = ''
  isContentFocused.value = false
  errorMessage.value = ''
  editorOpen.value = true
}

const openEditEditor = (prompt: PromptItem) => {
  editingPrompt.value = prompt
  draftName.value = prompt.name
  draftContent.value = prompt.content
  isContentFocused.value = false
  errorMessage.value = ''
  editorOpen.value = true
}

const closeEditor = () => {
  editorOpen.value = false
  editingPrompt.value = null
  draftName.value = ''
  draftContent.value = ''
  isContentFocused.value = false
}

const savePrompt = async () => {
  isSaving.value = true
  errorMessage.value = ''
  try {
    const payload = { name: draftName.value, content: draftContent.value }
    if (editingPrompt.value) {
      await updatePrompt(editingPrompt.value.id, payload)
    } else {
      await createPrompt(payload)
    }
    await refreshPrompts()
    closeEditor()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Save failed'
  } finally {
    isSaving.value = false
  }
}

const removePrompt = async (prompt: PromptItem) => {
  if (prompt.is_default) return
  const confirmed = window.confirm(`Delete prompt "${prompt.name}"? This cannot be undone.`)
  if (!confirmed) return
  errorMessage.value = ''
  try {
    await deletePrompt(prompt.id)
    await refreshPrompts()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Delete failed'
  }
}

const restorePrompt = async (prompt: PromptItem) => {
  errorMessage.value = ''
  try {
    await resetPrompt(prompt.id)
    await refreshPrompts()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Reset failed'
  }
}

onMounted(() => {
  void refreshPrompts()
})
</script>

<template>
  <section class="prompts-view">
    <div class="prompt-toolbar">
      <div>
        <span class="meta">Prompt Studio</span>
        <h2>提示词管理</h2>
        <p>生成剧本时固定使用锁定的系统提示词，并套用“导出风格提示词”的当前保存内容。</p>
        <p v-if="generationPrompt" class="prompt-usage">正在用于生成：{{ generationPrompt.name }}</p>
      </div>
      <button class="btn-submit prompt-new-button" type="button" @click="openCreateEditor">New Prompt</button>
    </div>

    <p v-if="errorMessage" class="prompt-error">{{ errorMessage }}</p>
    <p v-if="isLoading" class="prompt-empty">Loading prompts...</p>

    <div v-else class="prompt-grid">
      <article v-for="(prompt, index) in sortedPrompts" :key="prompt.id" class="prompt-card" @click="openEditEditor(prompt)">
        <div class="prompt-num">{{ String(index + 1).padStart(2, '0') }}</div>
        <div class="prompt-content">
          <div class="prompt-card-head">
            <h3>{{ prompt.name }}</h3>
            <div class="prompt-tags">
              <span v-if="isGenerationPrompt(prompt)" class="prompt-tag prompt-tag-active">Used in generation</span>
              <span v-if="prompt.is_default" class="prompt-tag">Default</span>
              <span v-if="prompt.is_editable" class="prompt-tag">Editable</span>
              <span v-else class="prompt-tag">Locked</span>
            </div>
          </div>
          <p class="prompt-text">{{ promptBody(prompt) }}</p>
          <div class="prompt-actions" @click.stop>
            <button class="prompt-action" type="button" :disabled="!prompt.is_editable" @click="openEditEditor(prompt)">Edit</button>
            <button v-if="prompt.is_default" class="prompt-action" type="button" :disabled="!prompt.is_editable" @click="restorePrompt(prompt)">Reset</button>
            <button v-else class="prompt-action danger" type="button" @click="removePrompt(prompt)">Delete</button>
          </div>
        </div>
      </article>
      <p v-if="!sortedPrompts.length" class="prompt-empty">No prompts yet. Create one to start.</p>
    </div>

    <div v-if="editorOpen" class="prompt-modal-backdrop" @click.self="closeEditor">
      <section class="prompt-modal-card">
        <div class="prompt-modal-head">
          <div>
            <span class="meta">{{ editingPrompt ? 'Edit Prompt' : 'New Prompt' }}</span>
            <h3>{{ editingPrompt?.name || 'Create prompt' }}</h3>
          </div>
          <button class="prompt-close" type="button" @click="closeEditor">×</button>
        </div>
        <label class="prompt-field">
          <span>Name</span>
          <input v-model="draftName" type="text" placeholder="提示词名称" :disabled="editingPrompt !== null && !editingPrompt.is_editable" />
        </label>
        <label class="prompt-field">
          <span>Content</span>
          <textarea
            v-model="draftContent"
            :placeholder="contentPlaceholder"
            :disabled="editingPrompt !== null && !editingPrompt.is_editable"
            @focus="isContentFocused = true"
            @blur="isContentFocused = false"
          ></textarea>
        </label>
        <p v-if="editingPrompt && !editingPrompt.is_editable" class="prompt-empty">This default prompt is locked and cannot be edited.</p>
        <div class="prompt-modal-actions">
          <button class="prompt-action" type="button" @click="closeEditor">Cancel</button>
          <button class="btn-submit" type="button" :disabled="isSaving || (editingPrompt !== null && !editingPrompt.is_editable)" @click="savePrompt">
            {{ isSaving ? 'Saving...' : 'Save Prompt' }}
          </button>
        </div>
      </section>
    </div>
  </section>
</template>
