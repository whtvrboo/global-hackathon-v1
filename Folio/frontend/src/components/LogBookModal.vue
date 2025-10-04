<template>
  <transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0"
    enter-to-class="opacity-100" leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100"
    leave-to-class="opacity-0">
    <div v-if="show" @click="$emit('close')"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div @click.stop class="bg-white rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto p-8">
        <h2 class="text-2xl font-bold mb-6">Log "{{ book?.title }}"</h2>

        <!-- Existing Logs Section -->
        <div v-if="existingLogs.length > 0"
          class="mb-6 p-4 bg-gradient-to-r from-blue-50 to-purple-50 border border-blue-200 rounded-lg">
          <div class="flex items-start gap-3">
            <div class="text-blue-500 text-xl">📚</div>
            <div class="flex-1">
              <h3 class="font-semibold text-blue-900 mb-2">You've read this book before!</h3>
              <p class="text-blue-700 text-sm mb-3">
                You have {{ existingLogs.length }} previous log{{ existingLogs.length > 1 ? 's' : '' }} for this book:
              </p>
              <div class="space-y-2 mb-4">
                <div v-for="log in existingLogs" :key="log.id" class="text-sm bg-white p-3 rounded border">
                  <div class="flex items-center justify-between mb-1">
                    <span class="font-medium">{{ statusLabel(log.status) }}</span>
                    <span class="text-gray-500">{{ formatDate(log.created_at) }}</span>
                  </div>
                  <div v-if="log.rating" class="text-yellow-600 mb-1">
                    {{ '★'.repeat(log.rating) }}{{ '☆'.repeat(5 - log.rating) }} ({{ log.rating }}/5)
                  </div>
                  <p v-if="log.review" class="text-gray-600 text-xs line-clamp-2">{{ log.review }}</p>
                </div>
              </div>
              <div class="flex gap-3">
                <button @click="handleLogAnotherReading"
                  class="px-6 py-3 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors flex items-center gap-2">
                  <span>📖</span>
                  <span>Log This New Reading</span>
                </button>
                <button @click="handleViewExistingLogs"
                  class="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-200 transition-colors flex items-center gap-2">
                  <span>👁️</span>
                  <span>View Existing Logs</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Status -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Reading Status *
            </label>
            <div class="grid grid-cols-2 gap-3">
              <button v-for="option in statusOptions" :key="option.value" type="button"
                @click="form.status = option.value" :class="[
                  'p-4 border-2 rounded-lg text-left transition-all',
                  form.status === option.value
                    ? 'border-primary bg-primary/5'
                    : 'border-gray-200 hover:border-gray-300'
                ]">
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
              <button v-for="star in 5" :key="star" type="button" @click="form.rating = star"
                class="text-4xl transition-transform hover:scale-110"
                :class="star <= (form.rating || 0) ? 'text-yellow-500' : 'text-gray-300'">
                ★
              </button>
              <button v-if="form.rating" type="button" @click="form.rating = null"
                class="ml-2 text-sm text-gray-500 hover:text-gray-700">
                Clear
              </button>
            </div>
          </div>

          <!-- Review -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Review
            </label>
            <TextArea v-model="form.review" placeholder="What did you think about this book?" :rows="4" />

            <!-- Spoiler Warning -->
            <div class="mt-3 flex items-center gap-2">
              <input type="checkbox" id="spoiler-flag" v-model="form.spoiler_flag"
                class="w-4 h-4 text-accent-red bg-gray-100 border-gray-300 rounded focus:ring-accent-red focus:ring-2" />
              <label for="spoiler-flag" class="text-sm text-gray-600">
                ⚠️ Contains spoilers
              </label>
            </div>
          </div>

          <!-- Notes (Private) -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Private Notes
            </label>
            <TextArea v-model="form.notes" placeholder="Personal notes, quotes, or thoughts (only visible to you)"
              :rows="3" />
          </div>

          <!-- Dates -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Start Date
              </label>
              <Input v-model="form.start_date" type="date" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Finish Date
              </label>
              <Input v-model="form.finish_date" type="date" />
            </div>
          </div>

          <!-- Visibility -->
          <div class="flex items-center gap-3">
            <input v-model="form.is_public" type="checkbox" id="is_public"
              class="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary" />
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
import { ref, reactive, watch, onMounted } from 'vue'
import axios from 'axios'
import PrimaryButton from './ui/PrimaryButton.vue'
import SecondaryButton from './ui/SecondaryButton.vue'
import Input from './ui/Input.vue'
import TextArea from './ui/TextArea.vue'
import { useToastStore } from '../stores/toast'

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
  is_public: true,
  spoiler_flag: false
})

const loading = ref(false)
const error = ref(null)
const existingLogs = ref([])
const isReread = ref(false)
const toastStore = useToastStore()

const statusOptions = [
  { value: 'want_to_read', label: 'Want to Read', emoji: '📚', description: 'Add to your list' },
  { value: 'reading', label: 'Currently Reading', emoji: '📖', description: 'Reading right now' },
  { value: 'read', label: 'Read', emoji: '✅', description: 'Finished this book' },
  { value: 'dnf', label: 'Did Not Finish', emoji: '🚫', description: 'Stopped reading' }
]

const statusLabel = (status) => {
  const labels = {
    'want_to_read': 'Want to Read',
    'reading': 'Currently Reading',
    'read': 'Read',
    'dnf': 'Did Not Finish'
  }
  return labels[status] || status
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const checkExistingLogs = async () => {
  if (!props.book?.id) return

  try {
    // Get current user's logs for this book
    const response = await axios.get(`/api/users/me/logs`)
    const userLogs = response.data.logs || []

    // Filter logs for this specific book
    existingLogs.value = userLogs.filter(log => log.book_id === props.book.id)
  } catch (error) {
    console.error('Error checking existing logs:', error)
    existingLogs.value = []
  }
}

watch(() => props.show, async (newShow) => {
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
    isReread.value = false

    // Check for existing logs
    await checkExistingLogs()
  }
})

const handleLogAnotherReading = () => {
  isReread.value = true
  // Clear any existing error
  error.value = null
}

const handleViewExistingLogs = () => {
  // Close modal and redirect to profile with book filter
  emit('close')
  // Could emit a custom event to navigate to profile with book filter
  // For now, just close the modal
}

const handleSubmit = async () => {
  if (!form.status) {
    error.value = 'Please select a reading status'
    return
  }

  // If book is already logged and user hasn't chosen to log another reading, show guidance
  if (existingLogs.value.length > 0 && !isReread.value) {
    error.value = 'Please choose an action above to continue.'
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
      is_public: form.is_public,
      is_reread: existingLogs.value.length > 0 && isReread.value
    }

    await axios.post('/api/logs', payload)

    // Show success message
    if (existingLogs.value.length > 0) {
      toastStore.success(`Logged another reading of "${props.book.title}"!`)
    } else {
      toastStore.success(`Successfully logged "${props.book.title}"!`)
    }

    emit('success')
    emit('close')
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to save log'
    toastStore.error('Failed to save log. Please try again.')
    console.error('Error saving log:', err)
  } finally {
    loading.value = false
  }
}
</script>
