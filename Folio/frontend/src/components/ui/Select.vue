<template>
  <select
    :value="modelValue"
    @change="$emit('update:modelValue', $event.target.value)"
    :disabled="disabled"
    class="input w-full px-4 py-2 border border-gray-300 rounded-lg
           focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary
           disabled:bg-gray-50 disabled:cursor-not-allowed
           transition-all duration-200 bg-white"
    v-bind="$attrs"
  >
    <option v-if="placeholder" value="" disabled selected>{{ placeholder }}</option>
    <option v-for="option in options" :key="option" :value="option">
      {{ option }}
    </option>
  </select>
</template>

<script setup>
defineProps({
  modelValue: [String, Number],
  options: { type: Array, required: true },
  placeholder: String,
  disabled: Boolean
})

defineEmits(['update:modelValue'])
</script>

