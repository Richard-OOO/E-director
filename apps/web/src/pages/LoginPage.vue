<script setup lang="ts">
import { computed, ref } from 'vue'

import { useI18n } from '@/i18n/useI18n'

const { lang, t } = useI18n()

const form = ref({
  email: '',
  password: '',
})

const statusMessage = ref('')
const isSubmitting = ref(false)

const submitLabel = computed(() => (isSubmitting.value ? t.value.login.submiting : t.value.login.submit))

const handleSubmit = async () => {
  isSubmitting.value = true
  statusMessage.value = ''

  try {
    const response = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: form.value.email,
        password: form.value.password,
      }),
    })

    const result = await response.json()

    if (!response.ok || result.code !== 0) {
      throw new Error(result.message || 'Login failed')
    }

    statusMessage.value = t.value.login.success
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : 'Login failed'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="auth-screen">
    <section v-reveal class="auth-card reveal">
      <h2>{{ t.login.h2 }}</h2>
      <p>{{ t.login.p }}</p>

      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="login-email">{{ t.login.email }}</label>
          <input
            id="login-email"
            v-model.trim="form.email"
            type="email"
            :placeholder="lang === 'en' ? 'name@example.com' : 'you@example.com'"
            autocomplete="email"
            required
          />
        </div>
        <div class="form-group">
          <label for="login-password">{{ t.login.password }}</label>
          <input
            id="login-password"
            v-model="form.password"
            type="password"
            placeholder="••••••••"
            autocomplete="current-password"
            required
          />
        </div>
        <button class="btn-submit" type="submit" :disabled="isSubmitting">
          {{ submitLabel }}
        </button>
        <p v-if="statusMessage" class="auth-status">{{ statusMessage }}</p>
        <p class="auth-footer-copy">
          {{ t.login.footer }}
          <RouterLink to="/register">{{ t.login.registerLink }}</RouterLink>
        </p>
      </form>
    </section>
  </main>
</template>
