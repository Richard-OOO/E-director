<script setup lang="ts">
import { useI18n } from '@/i18n/useI18n'

const { t } = useI18n()
</script>

<template>
  <main class="container">
    <header class="hero">
      <h1 v-reveal class="reveal" v-html="t.hero.h1" />
      <p v-reveal class="reveal">{{ t.hero.p }}</p>
      <div v-reveal class="reveal">
        <RouterLink to="/register" class="btn-submit hero-cta">{{ t.hero.cta }}</RouterLink>
      </div>

      <div class="collage">
        <div v-reveal class="collage-item reveal delay-1">
          <img
            src="https://images.unsplash.com/photo-1485846234645-a62644f84728?auto=format&fit=crop&w=800&q=80"
            alt="Film production vision"
          />
          <span>{{ t.collage.vision }}</span>
        </div>
        <div v-reveal class="collage-item reveal delay-2">
          <img
            src="https://images.unsplash.com/photo-1492691527719-9d1e07e534b4?auto=format&fit=crop&w=800&q=80"
            alt="Camera motion"
          />
          <span>{{ t.collage.motion }}</span>
        </div>
        <div v-reveal class="collage-item reveal delay-3">
          <img
            src="https://images.unsplash.com/photo-1536440136628-849c177e76a1?auto=format&fit=crop&w=800&q=80"
            alt="Cinema result"
          />
          <span>{{ t.collage.result }}</span>
        </div>
      </div>
    </header>

    <section v-reveal class="intro-section reveal">
      <div class="intro-layout">
        <h2>{{ t.section.h2 }}</h2>
        <p>{{ t.section.p }}</p>
      </div>
    </section>
  </main>
</template>
