<template>
  <transition
    enter-active-class="transition ease-out duration-200"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition ease-in duration-150"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="show"
      @click="$emit('close')"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
      <div
        @click.stop
        class="bg-white rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto p-8"
      >
        <h2 class="text-2xl font-bold mb-6">Log "{{ book?.title }}"</h2>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Status -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Reading Status *
            </label>
            <div class="grid grid-cols-2 gap-3">
              <button
                v-for="option in statusOptions"
                :key="option.value"
                type="button"
                @click="form.status = option.value"
                :class="[
                  'p-4 border-2 rounded-lg text-left transition-all',
                  form.status === option.value
                    ? 'border-primary bg-primary/5'
                    : 'border-gray-200 hover:border-gray-300'
                ]"
              >
                <div class="text-2xl mb-1">{{ option.emoji }}</div>
                <div class="font-medium">{{ option.label }}</div>
                <div class="text-xs text-gray-500">{{ option.description }}</div>
              </button>
            </div>
          </div>

          <!-- Rating -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Your Rating
            </label>
            <div class="flex gap-2">
              <button
                v-for="star in 5"
                :key="star"
                type="button"
                @click="form.rating = star"
                class="text-4xl transition-transform hover:scale-110"
                :class="star <= (form.rating || 0) ? 'text-yellow-500' : 'text-gray-300'"
              >
                ★
              </button>
              <button
                v-if="form.rating"
                type="button"
                @click="form.rating = null"
                class="ml-2 text-sm text-gray-500 hover:text-gray-700"
              >
                Clear
              </button>
            </div>
          </div>

          <!-- Review -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Review
            </label>
            <TextArea
              v-model="form.review"
              placeholder="What did you think about this book?"
              :rows="4"
            />
          </div>

          <!-- Notes (Private) -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Private Notes
            </label>
            <TextArea
              v-model="form.notes"
              placeholder="Personal notes, quotes, or thoughts (only visible to you)"
              :rows="3"
            />
          </div>

          <!-- Dates -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Start Date
              </label>
              <Input
                v-model="form.start_date"
                type="date"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Finish Date
              </label>
              <Input
                v-model="form.finish_date"
                type="date"
              />
            </div>
          </div>

          <!-- Visibility -->
          <div class="flex items-center gap-3">
            <input
              v-model="form.is_public"
              type="checkbox"
              id="is_public"
              class="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary"
            />
            <label for="is_public" class="text-sm text-gray-700">
              Make this log public (visible to followers)
            </label>
          </div>

          <!-- Error -->
          <div v-if="error" class="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
            {{ error }}
          </div>

          <!-- Actions -->
          <div class="flex gap-3 justify-end pt-4 border-t">
            <SecondaryButton @click="$emit('close')" type="button">
              Cancel
            </SecondaryButton>
            <PrimaryButton type="submit" :loading="loading">
              Save Log
            </PrimaryButton>
          </div>
        </form>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import axios from 'axios'
import PrimaryButton from './ui/PrimaryButton.vue'
import SecondaryButton from './ui/SecondaryButton.vue'
import Input from './ui/Input.vue'
import TextArea from './ui/TextArea.vue'

const props = defineProps({
  show: Boolean,
  book: Object
})

const emit = defineEmits(['close', 'success'])

const form = reactive({
  status: 'want_to_read',
  rating: null,
  review: '',
  notes: '',
  start_date: '',
  finish_date: '',
  is_public: true
})

const loading = ref(false)
const error = ref(null)

const statusOptions = [
  { value: 'want_to_read', label: 'Want to Read', emoji: '📚', description: 'Add to your list' },
  { value: 'reading', label: 'Currently Reading', emoji: '📖', description: 'Reading right now' },
  { value: 'read', label: 'Read', emoji: '✅', description: 'Finished this book' },
  { value: 'dnf', label: 'Did Not Finish', emoji: '🚫', description: 'Stopped reading' }
]

watch(() => props.show, (newShow) => {
  if (newShow) {
    // Reset form
    form.status = 'want_to_read'
    form.rating = null
    form.review = ''
    form.notes = ''
    form.start_date = ''
    form.finish_date = ''
    form.is_public = true
    error.value = null
  }
})

const handleSubmit = async () => {
  if (!form.status) {
    error.value = 'Please select a reading status'
    return
  }

  loading.value = true
  error.value = null

  try {
    const payload = {
      book_id: props.book.id,
      status: form.status,
      rating: form.rating || undefined,
      review: form.review || undefined,
      notes: form.notes || undefined,
      start_date: form.start_date || undefined,
      finish_date: form.finish_date || undefined,
      is_public: form.is_public
    }

    await axios.post('/api/logs', payload)
    emit('success')
    emit('close')
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to save log'
    console.error('Error saving log:', err)
  } finally {
    loading.value = false
  }
}
</script>

