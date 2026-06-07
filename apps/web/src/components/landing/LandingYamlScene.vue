<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  id: string
  num: string
  title: string
  intent: string
  yaml?: string
}>()

const emit = defineEmits<{
  (event: 'update-yaml', value: string): void
}>()

const draftYaml = ref(props.yaml ?? '')

watch(
  () => props.yaml,
  (value) => {
    draftYaml.value = value ?? ''
  },
)

const updateYaml = (event: Event) => {
  const target = event.target as HTMLTextAreaElement
  draftYaml.value = target.value
  emit('update-yaml', target.value)
}
</script>

<template>
  <div class="yaml-section" :id="props.id">
    <textarea
      class="yaml-code yaml-code-editor"
      :value="draftYaml"
      spellcheck="false"
      @input="updateYaml"
    ></textarea>
  </div>
</template>
