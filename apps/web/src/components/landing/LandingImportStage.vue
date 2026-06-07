<script setup lang="ts">
import { ref } from 'vue'

import type { CreateProjectPayload } from '@/auth'

const props = defineProps<{
  isSubmitting: boolean
  errorMessage: string
}>()

const emit = defineEmits<{
  (event: 'create-project', payload: CreateProjectPayload): void
}>()

const title = ref('The Silent Threshold')
const language = ref('zh-CN')
const content = ref('第一章：沉默的门槛\n主角抵达一座废弃摄影棚，门外的光线被厚重云层压低。\n第二章：过去的回声\n一条长廊唤醒旧记忆，镜头在现实和回忆之间缓慢切换。')

const submit = () => {
  emit('create-project', {
    title: title.value,
    language: language.value,
    source_type: 'plain_text',
    content: content.value,
  })
}
</script>

<template>
  <section class="stage-view">
    <form class="card-import landing-import-form" @submit.prevent="submit">
      <span class="meta">Source Intake</span>
      <h2>Import Source</h2>
      <p>Drop a manuscript or paste source text to begin generating the project.</p>

      <label class="landing-field">
        <span>Project title</span>
        <input v-model.trim="title" type="text" placeholder="Project title" :disabled="props.isSubmitting" required />
      </label>

      <label class="landing-field">
        <span>Language</span>
        <select v-model="language" :disabled="props.isSubmitting">
          <option value="zh-CN">中文</option>
          <option value="en-US">English</option>
        </select>
      </label>

      <label class="upload-zone landing-source-zone">
        <h4>Paste source text</h4>
        <p class="upload-hint">Plain text manuscript</p>
        <textarea v-model.trim="content" :disabled="props.isSubmitting" rows="8" required></textarea>
      </label>

      <p v-if="props.errorMessage" class="landing-form-error">{{ props.errorMessage }}</p>
      <button class="btn-submit" type="submit" :disabled="props.isSubmitting">
        {{ props.isSubmitting ? 'Creating Project...' : 'Start Analysis' }}
      </button>
    </form>
  </section>
</template>
