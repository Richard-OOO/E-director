<script setup lang="ts">
import { useI18n } from '@/i18n/useI18n'

const { lang, t } = useI18n()
</script>

<template>
  <main class="auth-screen">
    <section v-reveal class="auth-card reveal">
      <h2>{{ t.login.h2 }}</h2>
      <p>{{ t.login.p }}</p>

      <form @submit.prevent>
        <div class="form-group">
          <label for="login-email">{{ t.login.email }}</label>
          <input
            id="login-email"
            type="email"
            :placeholder="lang === 'en' ? 'name@example.com' : 'you@example.com'"
            autocomplete="email"
            required
          />
        </div>
        <div class="form-group">
          <label for="login-password">{{ t.login.password }}</label>
          <input id="login-password" type="password" placeholder="••••••••" autocomplete="current-password" required />
        </div>
        <button class="btn-submit" type="submit">{{ t.login.submit }}</button>
        <p class="auth-footer-copy">
          {{ t.login.footer }}
          <RouterLink to="/register">{{ t.login.registerLink }}</RouterLink>
        </p>
      </form>
    </section>
  </main>
</template>
