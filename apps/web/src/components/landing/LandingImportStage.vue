<script setup lang="ts">
import { computed, ref } from 'vue'

import type { CreateProjectPayload } from '@/auth'

const props = defineProps<{
  isSubmitting: boolean
  errorMessage: string
}>()

const emit = defineEmits<{
  (event: 'create-project', payload: CreateProjectPayload): void
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const selectedFileName = ref('')
const selectedFileSize = ref(0)
const selectedSourceType = ref<'docx' | 'pdf' | ''>('')
const fileError = ref('')
const isDragging = ref(false)

const title = computed(() => {
  if (!selectedFileName.value) return 'Imported Document'
  return selectedFileName.value.replace(/\.[^.]+$/, '') || 'Imported Document'
})

const sourceLabel = computed(() => {
  if (selectedFileName.value) return selectedFileName.value
  return isDragging.value ? 'Release to import the document' : 'No document selected'
})

const sourceDetail = computed(() => {
  if (selectedSourceType.value === 'pdf') return `PDF document · ${formatFileSize(selectedFileSize.value)}`
  if (selectedSourceType.value === 'docx') return `DOCX document · ${formatFileSize(selectedFileSize.value)}`
  return 'Accepts .docx and .pdf files'
})

const pickerLabel = computed(() => (selectedFileName.value ? 'Replace file' : 'Choose file'))

const openPicker = () => {
  fileInput.value?.click()
}

const normalizeFileName = (name: string) => name.trim()

const formatFileSize = (bytes: number) => {
  if (!bytes) return '0 KB'
  const units = ['B', 'KB', 'MB', 'GB']
  const exponent = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1)
  const value = bytes / 1024 ** exponent
  return `${value >= 10 || exponent === 0 ? value.toFixed(0) : value.toFixed(1)} ${units[exponent]}`
}

const resetFileInput = () => {
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const acceptFile = (file?: File | null) => {
  fileError.value = ''
  if (!file) return

  const name = normalizeFileName(file.name)
  const lowerName = name.toLowerCase()

  if (!lowerName.endsWith('.docx') && !lowerName.endsWith('.pdf')) {
    fileError.value = 'Please choose a DOCX or PDF document.'
    resetFileInput()
    return
  }

  selectedFile.value = file
  selectedFileName.value = name
  selectedFileSize.value = file.size
  selectedSourceType.value = lowerName.endsWith('.pdf') ? 'pdf' : 'docx'
}

const clearFile = () => {
  selectedFile.value = null
  selectedFileName.value = ''
  selectedFileSize.value = 0
  selectedSourceType.value = ''
  fileError.value = ''
  resetFileInput()
}

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement | null
  acceptFile(target?.files?.[0])
}

const onDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragging.value = false
  acceptFile(event.dataTransfer?.files?.[0])
}

const onDragEnter = (event: DragEvent) => {
  event.preventDefault()
  isDragging.value = true
}

const onDragOver = (event: DragEvent) => {
  event.preventDefault()
  isDragging.value = true
}

const onDragLeave = (event: DragEvent) => {
  event.preventDefault()
  if (event.currentTarget === event.target) {
    isDragging.value = false
  }
}

const submit = () => {
  if (!selectedFile.value || !selectedSourceType.value) {
    fileError.value = 'Choose a DOCX or PDF document before starting analysis.'
    return
  }

  emit('create-project', {
    title: title.value,
    language: 'zh-CN',
    source_type: selectedSourceType.value,
    file: selectedFile.value,
  })
}
</script>

<template>
  <section class="stage-view">
    <form class="card-import landing-import-form" @submit.prevent="submit">
      <span class="meta">Source Intake</span>
      <h2>Import Source</h2>
      <p>Drop a DOCX or PDF to begin. The document name becomes the project title.</p>

      <input ref="fileInput" class="sr-only" type="file" accept=".docx,.pdf" @change="onFileChange" />

      <div
        class="upload-zone landing-source-zone"
        :class="{ 'is-dragging': isDragging, 'has-file': selectedFileName }"
        role="button"
        tabindex="0"
        @dragenter="onDragEnter"
        @dragover="onDragOver"
        @dragleave="onDragLeave"
        @drop="onDrop"
        @click="openPicker"
        @keydown.enter.prevent="openPicker"
        @keydown.space.prevent="openPicker"
      >
        <span class="upload-badge">DOCX / PDF</span>
        <h4>{{ selectedFileName ? 'Document ready' : 'Drop document here' }}</h4>
        <p class="upload-hint">{{ sourceDetail }}</p>
        <strong>{{ sourceLabel }}</strong>
        <button class="upload-cta" type="button" @click.stop="openPicker">{{ pickerLabel }}</button>
      </div>

      <div v-if="selectedFileName" class="landing-file-card">
        <span class="meta">Selected file</span>
        <strong>{{ selectedFileName }}</strong>
        <p>{{ title }}</p>
        <button class="landing-file-remove" type="button" @click="clearFile">Remove</button>
      </div>

      <p v-if="fileError || props.errorMessage" class="landing-form-error">{{ fileError || props.errorMessage }}</p>
      <button class="btn-submit" type="submit" :disabled="props.isSubmitting || !selectedFileName">
        {{ props.isSubmitting ? 'Creating Project...' : 'Start Analysis' }}
      </button>
    </form>
  </section>
</template>
