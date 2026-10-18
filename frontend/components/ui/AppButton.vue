<template>
  <component
    :is="to ? 'NuxtLink' : 'button'"
    :to="to"
    :type="!to ? type : undefined"
    :disabled="!to ? disabled : undefined"
    :class="[
      'inline-flex items-center justify-center gap-2 font-semibold rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed',
      sizeClass,
      variantClass,
    ]"
  >
    <slot />
  </component>
</template>

<script setup lang="ts">
interface Props {
  to?: string
  type?: 'button' | 'submit'
  disabled?: boolean
  size?: 'sm' | 'md' | 'lg'
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
}

const props = withDefaults(defineProps<Props>(), {
  type: 'button',
  disabled: false,
  size: 'md',
  variant: 'primary',
})

const sizeClass = computed(() => {
  switch (props.size) {
    case 'sm': return 'px-4 py-2 text-sm'
    case 'lg': return 'px-8 py-4 text-base'
    default: return 'px-6 py-3 text-sm'
  }
})

const variantClass = computed(() => {
  switch (props.variant) {
    case 'primary': return 'bg-[#00dc82] text-[#020420] hover:bg-[#00c070]'
    case 'secondary': return 'bg-white/10 text-white hover:bg-white/20'
    case 'ghost': return 'bg-transparent text-gray-400 hover:text-white hover:bg-white/5'
    case 'danger': return 'bg-red-500/20 text-red-400 hover:bg-red-500/30'
    default: return 'bg-[#00dc82] text-[#020420] hover:bg-[#00c070]'
  }
})
</script>
