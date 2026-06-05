<script setup lang="ts">
import { computed, ref } from 'vue'

import { useI18n } from '@/i18n/useI18n'

const { lang, t } = useI18n()

const form = ref({
  name: '',
  email: '',
  password: '',
})

const statusMessage = ref('')
const isSubmitting = ref(false)

const submitLabel = computed(() => (isSubmitting.value ? t.value.register.submiting : t.value.register.submit))

const handleSubmit = async () => {
  isSubmitting.value = true
  statusMessage.value = ''

  try {
    const response = await fetch('/api/v1/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: form.value.name,
        email: form.value.email,
        password: form.value.password,
      }),
    })

    const result = await response.json()

    if (!response.ok || result.code !== 0) {
      throw new Error(result.message || 'Register failed')
    }

    statusMessage.value = t.value.register.success
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : 'Register failed'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="auth-screen">
    <section v-reveal class="auth-card reveal">
      <h2>{{ t.register.h2 }}</h2>
      <p>{{ t.register.p }}</p>

      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="register-name">{{ t.register.name }}</label>
          <input
            id="register-name"
            v-model.trim="form.name"
            type="text"
            :placeholder="lang === 'en' ? 'John Doe' : '张三'"
            autocomplete="name"
            required
          />
        </div>
        <div class="form-group">
          <label for="register-email">{{ t.register.email }}</label>
          <input
            id="register-email"
            v-model.trim="form.email"
            type="email"
            :placeholder="lang === 'en' ? 'name@example.com' : 'you@example.com'"
            autocomplete="email"
            required
          />
        </div>
        <div class="form-group">
          <label for="register-password">{{ t.register.password }}</label>
          <input
            id="register-password"
            v-model="form.password"
            type="password"
            placeholder="••••••••"
            autocomplete="new-password"
            required
          />
        </div>
        <button class="btn-submit" type="submit" :disabled="isSubmitting">
          {{ submitLabel }}
        </button>
        <p v-if="statusMessage" class="auth-status">{{ statusMessage }}</p>
        <p class="auth-footer-copy">
          {{ t.register.footer }}
          <RouterLink to="/login">{{ t.register.loginLink }}</RouterLink>
        </p>
      </form>
    </section>
  </main>
</template>
