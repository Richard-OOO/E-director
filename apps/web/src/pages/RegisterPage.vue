<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'

import { useRouter } from 'vue-router'

import { register, sendRegisterCode } from '@/auth'
import { useI18n } from '@/i18n/useI18n'

const router = useRouter()
const { lang, t } = useI18n()

const form = ref({
  name: '',
  email: '',
  password: '',
  verificationCode: '',
})

const statusMessage = ref('')
const isSubmitting = ref(false)
const isSendingCode = ref(false)
const cooldownSeconds = ref(0)
let cooldownTimer: number | undefined

const submitLabel = computed(() => (isSubmitting.value ? t.value.register.submiting : t.value.register.submit))
const codeButtonLabel = computed(() => {
  if (isSendingCode.value) return t.value.register.sendingCode
  if (cooldownSeconds.value > 0) return `${t.value.register.resendIn} ${cooldownSeconds.value}s`
  return t.value.register.sendCode
})
const isCodeButtonDisabled = computed(() => isSendingCode.value || cooldownSeconds.value > 0 || !form.value.email)

const startCooldown = (seconds: number) => {
  cooldownSeconds.value = seconds
  if (cooldownTimer) {
    window.clearInterval(cooldownTimer)
  }
  cooldownTimer = window.setInterval(() => {
    cooldownSeconds.value = Math.max(cooldownSeconds.value - 1, 0)
    if (cooldownSeconds.value === 0 && cooldownTimer) {
      window.clearInterval(cooldownTimer)
      cooldownTimer = undefined
    }
  }, 1000)
}

const handleSendCode = async () => {
  isSendingCode.value = true
  statusMessage.value = ''

  try {
    const result = await sendRegisterCode({ email: form.value.email })
    statusMessage.value = t.value.register.codeSent
    startCooldown(result.cooldown_seconds)
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : 'Failed to send verification code'
  } finally {
    isSendingCode.value = false
  }
}

const handleSubmit = async () => {
  isSubmitting.value = true
  statusMessage.value = ''

  try {
    await register({
      name: form.value.name,
      email: form.value.email,
      password: form.value.password,
      verification_code: form.value.verificationCode,
    })

    statusMessage.value = t.value.register.success
    await router.push('/workspace')
  } catch (error) {
    statusMessage.value = error instanceof Error ? error.message : 'Register failed'
  } finally {
    isSubmitting.value = false
  }
}

onBeforeUnmount(() => {
  if (cooldownTimer) {
    window.clearInterval(cooldownTimer)
  }
})
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
        <div class="form-group">
          <label for="register-code">{{ t.register.verificationCode }}</label>
          <div class="code-input-row">
            <input
              id="register-code"
              v-model.trim="form.verificationCode"
              type="text"
              inputmode="numeric"
              maxlength="6"
              placeholder="123456"
              autocomplete="one-time-code"
              required
            />
            <button class="btn-secondary" type="button" :disabled="isCodeButtonDisabled" @click="handleSendCode">
              {{ codeButtonLabel }}
            </button>
          </div>
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
