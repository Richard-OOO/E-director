<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

import { useI18n } from '@/i18n/useI18n'

const { lang, t, toggleLang } = useI18n()

const intro = computed(() => {
  if (lang.value === 'zh') {
    return {
      hero: {
        titleLead: 'AI 助力',
        titleRest: '导演一切.',
        description: 'E-Director 是一个专业级平台，帮助你通过人工智能编排小说、视觉故事和分镜叙事，把创意转成可执行的导演工作流。',
        cta: '开始导演',
      },
      collage: {
        vision: '视觉',
        motion: '动态',
        result: '成果',
      },
      section: {
        title: '创新视觉叙事的生产方式',
        description: '我们弥合创意文本与技术执行之间的距离，让章节、场景、提示词和 YAML 结构在同一个工作台中流动起来。',
      },
    }
  }

  return {
    hero: {
      titleLead: 'AI Powers',
      titleRest: 'Directing Everything.',
      description: 'E-Director is an editorial-grade platform for orchestrating novels, visual stories, and narrative production through the power of artificial intelligence.',
      cta: 'Get Started',
    },
    collage: {
      vision: 'The Vision',
      motion: 'The Motion',
      result: 'The Result',
    },
    section: {
      title: 'Innovating the Art of Visual Storytelling',
      description: 'We bridge creative vision and technical execution, turning chapters, scenes, prompts, and YAML structure into one directable workspace.',
    },
  }
})
</script>

<template>
  <div class="intro-page">
    <nav class="app-nav container">
      <RouterLink class="logo" to="/">E-Director</RouterLink>
      <div class="nav-right">
        <button class="lang-switch" type="button" @click="toggleLang">
          <span :class="{ active: lang === 'en' }">EN</span>
          <span class="sep">/</span>
          <span :class="{ active: lang === 'zh' }">中</span>
        </button>
        <div class="nav-links">
          <RouterLink class="btn-action" to="/login">{{ t.nav.login }}</RouterLink>
          <RouterLink class="btn-action" to="/register">{{ t.nav.register }}</RouterLink>
        </div>
      </div>
    </nav>

    <main class="container">
      <header class="hero">
        <h1 v-reveal class="reveal">
          {{ intro.hero.titleLead }}<br />
          <em>{{ intro.hero.titleRest }}</em>
        </h1>
        <p v-reveal class="reveal">{{ intro.hero.description }}</p>
        <div v-reveal class="reveal">
          <RouterLink class="btn-submit hero-cta" to="/register">{{ intro.hero.cta }}</RouterLink>
        </div>

        <div class="collage" aria-label="E-Director workflow preview">
          <div v-reveal class="collage-item reveal intro-delay-short">
            <img src="https://images.unsplash.com/photo-1485846234645-a62644f84728?auto=format&fit=crop&w=800&q=80" alt="Film vision planning" />
            <span>{{ intro.collage.vision }}</span>
          </div>
          <div v-reveal class="collage-item reveal intro-delay-medium">
            <img src="https://images.unsplash.com/photo-1492691527719-9d1e07e534b4?auto=format&fit=crop&w=800&q=80" alt="Camera and motion production" />
            <span>{{ intro.collage.motion }}</span>
          </div>
          <div v-reveal class="collage-item reveal intro-delay-long">
            <img src="https://images.unsplash.com/photo-1536440136628-849c177e76a1?auto=format&fit=crop&w=800&q=80" alt="Cinema result screen" />
            <span>{{ intro.collage.result }}</span>
          </div>
        </div>
      </header>

      <section v-reveal class="intro-section reveal">
        <div class="intro-layout">
          <h2>{{ intro.section.title }}</h2>
          <p>{{ intro.section.description }}</p>
        </div>
      </section>
    </main>

    <footer class="footer container">
      <p>{{ t.footer }}</p>
    </footer>
  </div>
</template>
