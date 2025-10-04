<template>
  <button
    class="inline-flex items-center justify-center gap-2 px-4 py-2 font-medium rounded-lg border-2 border-gray-300 
           hover:border-primary hover:text-primary transition-all duration-200
           focus:outline-none focus:ring-2 focus:ring-primary/50 focus:ring-offset-2
           disabled:opacity-50 disabled:cursor-not-allowed"
    :disabled="disabled"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>

<script setup>
defineProps({
  disabled: Boolean
})
</script>

