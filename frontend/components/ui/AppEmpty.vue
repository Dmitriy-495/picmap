<template>
  <div class="text-center py-24">
    <div class="text-7xl mb-6">{{ icon }}</div>
    <h2 class="text-3xl font-bold mb-3">{{ title }}</h2>
    <p v-if="description" class="text-gray-400 mb-10 text-lg">{{ description }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
interface Props {
  icon?: string
  title: string
  description?: string
}

withDefaults(defineProps<Props>(), {
  icon: '🗺️',
})
</script>
