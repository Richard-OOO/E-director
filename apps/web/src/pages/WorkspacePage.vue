<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { fetchMe, logout as logoutSession, type AuthUser } from '@/auth'
import { useI18n } from '@/i18n/useI18n'

const router = useRouter()
const { t } = useI18n()

const user = ref<AuthUser | null>(null)
const statusMessage = ref('')
const isLoading = ref(true)

const logout = async () => {
  try {
    await logoutSession()
  } finally {
    await router.push('/login')
  }
}

onMounted(async () => {
  try {
    user.value = await fetchMe()
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : 'Session expired'
    await router.push('/login')
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <main class="container workspace-screen">
    <section v-reveal class="workspace-card reveal">
      <p class="eyebrow">{{ t.workspace.eyebrow }}</p>
      <h1>{{ t.workspace.h1 }}</h1>
      <p>{{ t.workspace.p }}</p>

      <div v-if="isLoading" class="workspace-panel">
        {{ t.workspace.loading }}
      </div>

      <div v-else-if="user" class="workspace-panel">
        <span>{{ t.workspace.signedInAs }}</span>
        <strong>{{ user.display_name || user.email }}</strong>
        <small>{{ user.email }}</small>
      </div>

      <p v-if="statusMessage" class="auth-status">{{ statusMessage }}</p>

      <button class="btn-submit workspace-action" type="button" @click="logout">
        {{ t.workspace.logout }}
      </button>
    </section>
  </main>
</template>
