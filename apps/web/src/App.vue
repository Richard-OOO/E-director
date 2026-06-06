<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'

import { fetchMe } from '@/auth'
import { useI18n } from '@/i18n/useI18n'

const route = useRoute()
const { lang, t, toggleLang } = useI18n()

const currentPath = computed(() => route.path)
const hasSession = ref(false)

watch(
  () => route.fullPath,
  async () => {
    try {
      await fetchMe()
      hasSession.value = true
    } catch {
      hasSession.value = false
    }
  },
  { immediate: true },
)
</script>

<template>
  <nav class="container app-nav">
    <RouterLink to="/" class="logo">E-Director</RouterLink>

    <div class="nav-right">
      <button class="lang-switch" title="Switch Language" type="button" @click="toggleLang">
        <span :class="{ active: lang === 'en' }">EN</span>
        <span class="sep">/</span>
        <span :class="{ active: lang === 'zh' }">中</span>
      </button>

      <div class="nav-links">
        <RouterLink v-if="currentPath !== '/'" to="/">{{ t.nav.home }}</RouterLink>
        <RouterLink v-if="hasSession && currentPath !== '/workspace'" to="/workspace" class="btn-action">
          {{ t.nav.workspace }}
        </RouterLink>
        <RouterLink v-if="!hasSession && (currentPath === '/' || currentPath === '/register')" to="/login" class="btn-action">
          {{ t.nav.login }}
        </RouterLink>
        <RouterLink v-if="!hasSession && currentPath === '/login'" to="/register" class="btn-action">
          {{ t.nav.register }}
        </RouterLink>
      </div>
    </div>
  </nav>

  <RouterView />

  <footer class="footer container">
    <p>{{ t.footer }}</p>
  </footer>
</template>
