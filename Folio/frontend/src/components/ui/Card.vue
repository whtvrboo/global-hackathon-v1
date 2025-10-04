<template>
  <div :class="[
    'card',
    hover && 'card-hover',
    glass && 'card-glass',
    shadow && 'shadow-lg shadow-black/10'
  ]" v-bind="$attrs">
    <slot />
  </div>
</template>

<script setup>
defineProps({
  hover: Boolean,
  shadow: { type: Boolean, default: true },
  glass: Boolean
})
</script>
