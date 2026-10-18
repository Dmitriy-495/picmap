<template>
  <div :class="['inline-flex items-center px-4 py-1.5 rounded-full text-xs font-semibold border backdrop-blur-md', colorClass]">
    <slot />
  </div>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'will' | 'now' | 'past' | 'default'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
})

const colorClass = computed(() => {
  switch (props.variant) {
    case 'will':
      return 'bg-blue-500/20 text-blue-400 border-blue-500/30'
    case 'now':
      return 'bg-green-500/20 text-green-400 border-green-500/30'
    case 'past':
      return 'bg-gray-500/20 text-gray-400 border-gray-500/30'
    default:
      return 'bg-gray-500/20 text-gray-400 border-gray-500/30'
  }
})
</script>
