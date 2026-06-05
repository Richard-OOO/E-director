<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'

import { useI18n } from '@/i18n/useI18n'

const route = useRoute()
const { lang, t, toggleLang } = useI18n()

const currentPath = computed(() => route.path)
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
        <RouterLink v-if="currentPath === '/' || currentPath === '/register'" to="/login" class="btn-action">
          {{ t.nav.login }}
        </RouterLink>
        <RouterLink v-if="currentPath === '/login'" to="/register" class="btn-action">
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
