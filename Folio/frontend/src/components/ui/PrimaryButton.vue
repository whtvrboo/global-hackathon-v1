<template>
  <button
    :class="[
      'btn-primary inline-flex items-center justify-center gap-2 font-medium transition-all duration-200',
      'focus:outline-none focus:ring-2 focus:ring-primary/50 focus:ring-offset-2',
      'disabled:opacity-50 disabled:cursor-not-allowed',
      sizeClasses,
      loading && 'opacity-75 cursor-wait'
    ]"
    :disabled="disabled || loading"
    v-bind="$attrs"
  >
    <span v-if="loading" class="animate-spin">⏳</span>
    <slot />
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  loading: Boolean,
  disabled: Boolean
})

const sizeClasses = computed(() => {
  const sizes = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2',
    lg: 'px-6 py-3 text-lg'
  }
  return sizes[props.size]
})
</script>

