<template>
  <button
    :class="[
      'inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-all duration-200',
      'bg-gray-100 hover:bg-gray-200 text-gray-900',
      'focus:outline-none focus:ring-2 focus:ring-gray-300 focus:ring-offset-2',
      'disabled:opacity-50 disabled:cursor-not-allowed',
      sizeClasses
    ]"
    :disabled="disabled"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  size: { type: String, default: 'md' },
  disabled: Boolean
})

const sizeClasses = computed(() => ({
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2',
  lg: 'px-6 py-3 text-lg'
}[props.size]))
</script>

