<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  id: string
  num: string
  title: string
  intent: string
  comment: string
  yaml?: string
}>()

const isOpen = ref(false)
</script>

<template>
  <div class="yaml-section" :id="props.id">
    <pre v-if="props.yaml" class="yaml-code yaml-code-pre">{{ props.yaml }}</pre>
    <div v-else class="yaml-code">
      <span class="yaml-key">scene:</span> <span class="yaml-val">{{ props.num }}</span>
      <br />
      <span class="yaml-key">title:</span> <span class="yaml-val">{{ props.title }}</span>
      <br />
      <span class="yaml-key">intent:</span> <span class="yaml-val">{{ props.intent }}</span>
      <br />
      <span class="yaml-key">location:</span> <span class="yaml-val">INT. STUDIO - DAY</span>
    </div>

    <div class="yaml-tip-shell">
      <button class="yaml-tip" type="button" :aria-expanded="isOpen" @click="isOpen = !isOpen">Tip</button>
      <div class="yaml-tip-card">
        <span class="meta">Hover Tip</span>
        <strong>{{ props.title }}</strong>
        <p>Click to open the full YAML explanation.</p>
      </div>
      <div v-if="isOpen" class="yaml-commentary-panel is-open">
        <button class="yaml-commentary-close" type="button" aria-label="Close YAML explanation" @click="isOpen = false">×</button>
        <span class="meta">AI Observation</span>
        <h6>{{ props.title }}</h6>
        <p>{{ props.comment }}</p>
      </div>
    </div>
  </div>
</template>
