<template>
  <div class="p-8 space-y-12 max-w-6xl mx-auto">
    <section>
      <h2 class="text-2xl font-bold mb-6">Folio Design System</h2>
      <p class="text-gray-600">Linear-inspired components for a beautiful reading experience</p>
    </section>

    <!-- Buttons -->
    <section class="space-y-4">
      <h3 class="text-lg font-semibold">Buttons</h3>
      <div class="flex flex-wrap gap-3">
        <PrimaryButton>Primary Action</PrimaryButton>
        <SecondaryButton>Secondary Action</SecondaryButton>
        <OutlineButton>Outline Button</OutlineButton>
        <PrimaryButton size="sm">Small</PrimaryButton>
        <PrimaryButton :loading="true">Loading</PrimaryButton>
      </div>
    </section>

    <!-- Cards -->
    <section class="space-y-4">
      <h3 class="text-lg font-semibold">Cards</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <h4 class="font-semibold mb-2">Simple Card</h4>
          <p class="text-gray-600 text-sm">A clean, minimal card component</p>
        </Card>
        <Card hover>
          <h4 class="font-semibold mb-2">Hoverable Card</h4>
          <p class="text-gray-600 text-sm">With subtle hover effect</p>
        </Card>
      </div>
    </section>

    <!-- Inputs -->
    <section class="space-y-4">
      <h3 class="text-lg font-semibold">Form Inputs</h3>
      <div class="space-y-3 max-w-md">
        <Input placeholder="Enter text..." />
        <TextArea placeholder="Enter description..." :rows="3" />
        <Select :options="['Option 1', 'Option 2', 'Option 3']" />
      </div>
    </section>

    <!-- Typography -->
    <section class="space-y-4">
      <h3 class="text-lg font-semibold">Typography</h3>
      <div class="space-y-2">
        <h1 class="text-heading-1">Heading 1</h1>
        <h2 class="text-heading-2">Heading 2</h2>
        <h3 class="text-heading-3">Heading 3</h3>
        <p class="text-body">Body text with proper line height and spacing for readability.</p>
        <p class="text-caption">Caption text for supplementary information.</p>
      </div>
    </section>
  </div>
</template>

<script setup>
import PrimaryButton from './ui/PrimaryButton.vue'
import SecondaryButton from './ui/SecondaryButton.vue'
import OutlineButton from './ui/OutlineButton.vue'
import Card from './ui/Card.vue'
import Input from './ui/Input.vue'
import TextArea from './ui/TextArea.vue'
import Select from './ui/Select.vue'
</script>

