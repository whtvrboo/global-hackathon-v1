<template>
  <div
    :class="[
      'bg-white rounded-xl border border-gray-200 p-6 transition-all duration-200',
      hover && 'hover:shadow-md hover:border-gray-300 cursor-pointer',
      shadow && 'shadow-sm'
    ]"
    v-bind="$attrs"
  >
    <slot />
  </div>
</template>

<script setup>
defineProps({
  hover: Boolean,
  shadow: { type: Boolean, default: true }
})
</script>

