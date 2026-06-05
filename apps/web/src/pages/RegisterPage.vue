<script setup lang="ts">
import { useI18n } from '@/i18n/useI18n'

const { lang, t } = useI18n()
</script>

<template>
  <main class="auth-screen">
    <section v-reveal class="auth-card reveal">
      <h2>{{ t.register.h2 }}</h2>
      <p>{{ t.register.p }}</p>

      <form @submit.prevent>
        <div class="form-group">
          <label for="register-name">{{ t.register.name }}</label>
          <input
            id="register-name"
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
            type="email"
            :placeholder="lang === 'en' ? 'name@example.com' : 'you@example.com'"
            autocomplete="email"
            required
          />
        </div>
        <div class="form-group">
          <label for="register-password">{{ t.register.password }}</label>
          <input id="register-password" type="password" placeholder="••••••••" autocomplete="new-password" required />
        </div>
        <button class="btn-submit" type="submit">{{ t.register.submit }}</button>
        <p class="auth-footer-copy">
          {{ t.register.footer }}
          <RouterLink to="/login">{{ t.register.loginLink }}</RouterLink>
        </p>
      </form>
    </section>
  </main>
</template>
