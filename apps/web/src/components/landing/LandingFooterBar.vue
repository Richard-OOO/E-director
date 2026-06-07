<script setup lang="ts">
const props = defineProps<{
  isLoading: boolean
  isSigningOut: boolean
  user: { display_name: string; email: string } | null
}>()

const emit = defineEmits<{
  (event: 'sign-out'): void
}>()
</script>

<template>
  <div class="landing-footer-bar">
    <div v-if="props.isLoading">Loading session...</div>
    <div v-else-if="props.user">{{ props.user.display_name }} · {{ props.user.email }}</div>
    <div v-else>Session unavailable</div>
    <button class="btn-secondary" type="button" :disabled="props.isSigningOut" @click="emit('sign-out')">
      {{ props.isSigningOut ? 'Signing out...' : 'Sign Out' }}
    </button>
  </div>
</template>
