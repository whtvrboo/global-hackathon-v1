<template>
  <input
    :value="modelValue"
    @input="$emit('update:modelValue', $event.target.value)"
    :type="type"
    :placeholder="placeholder"
    :disabled="disabled"
    class="input w-full px-4 py-2 border border-gray-300 rounded-lg
           focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary
           disabled:bg-gray-50 disabled:cursor-not-allowed
           transition-all duration-200"
    v-bind="$attrs"
  />
</template>

<script setup>
defineProps({
  modelValue: [String, Number],
  type: { type: String, default: 'text' },
  placeholder: String,
  disabled: Boolean
})

defineEmits(['update:modelValue'])
</script>

