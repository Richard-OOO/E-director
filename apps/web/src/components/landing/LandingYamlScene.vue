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
  (event: 'save-yaml', value: string): void
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

const saveYaml = () => {
  emit('save-yaml', draftYaml.value)
}
</script>

<template>
  <div class="yaml-section" :id="props.id">
    <div class="yaml-scene-toolbar">
      <span>{{ props.num }}</span>
      <strong>{{ props.title }}</strong>
      <button class="yaml-save-button" type="button" @click.stop="saveYaml">Save</button>
    </div>
    <textarea
      class="yaml-code yaml-code-editor"
      :value="draftYaml"
      spellcheck="false"
      @input="updateYaml"
      @blur="saveYaml"
    ></textarea>
  </div>
</template>
