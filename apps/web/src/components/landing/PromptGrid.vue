<script setup lang="ts">
import { landingPromptCards } from '@/components/landing/fixtures'
</script>

<template>
  <section class="prompts-view">
    <div class="prompt-grid">
      <article v-for="prompt in landingPromptCards" :key="prompt.id" class="prompt-card">
        <div class="prompt-num">{{ prompt.num }}</div>
        <div class="prompt-content">
          <p class="prompt-text">{{ prompt.text }}</p>
          <div class="prompt-tags">
            <span v-for="tag in prompt.tags" :key="tag" class="prompt-tag">{{ tag }}</span>
          </div>
        </div>
      </article>
    </div>
  </section>
</template>
