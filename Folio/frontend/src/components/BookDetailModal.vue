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
      <transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div
          v-if="show"
          @click.stop
          class="bg-white rounded-2xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-y-auto"
        >
          <!-- Loading State -->
          <div v-if="loading" class="p-12 text-center">
            <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
            <p class="mt-4 text-gray-600">Loading book details...</p>
          </div>

          <!-- Book Details -->
          <div v-else-if="bookDetails" class="p-8">
            <div class="flex gap-6 mb-6">
              <!-- Book Cover -->
              <div class="flex-shrink-0">
                <img
                  v-if="bookDetails.cover_url"
                  :src="bookDetails.cover_url"
                  :alt="bookDetails.title"
                  class="w-48 h-72 object-cover rounded-lg shadow-lg"
                />
                <div v-else class="w-48 h-72 bg-gray-200 rounded-lg flex items-center justify-center text-6xl">
                  📚
                </div>
              </div>

              <!-- Book Info -->
              <div class="flex-1 min-w-0">
                <h2 class="text-3xl font-bold text-gray-900 mb-2">
                  {{ bookDetails.title }}
                </h2>
                <p v-if="bookDetails.authors?.length" class="text-lg text-gray-600 mb-4">
                  by {{ bookDetails.authors.join(', ') }}
                </p>

                <!-- Metadata -->
                <div class="flex flex-wrap gap-4 text-sm text-gray-600 mb-6">
                  <div v-if="bookDetails.published_date">
                    📅 {{ bookDetails.published_date }}
                  </div>
                  <div v-if="bookDetails.page_count">
                    📖 {{ bookDetails.page_count }} pages
                  </div>
                  <div v-if="bookDetails.publisher">
                    🏢 {{ bookDetails.publisher }}
                  </div>
                </div>

                <!-- Rating -->
                <div v-if="bookDetails.rating" class="flex items-center gap-2 mb-6">
                  <div class="flex">
                    <span
                      v-for="i in 5"
                      :key="i"
                      class="text-2xl"
                      :class="i <= Math.round(bookDetails.rating) ? 'text-yellow-500' : 'text-gray-300'"
                    >
                      ★
                    </span>
                  </div>
                  <span class="text-lg font-semibold">{{ bookDetails.rating }}</span>
                  <span v-if="bookDetails.ratings_count" class="text-sm text-gray-500">
                    ({{ bookDetails.ratings_count }} ratings)
                  </span>
                </div>

                <!-- Categories -->
                <div v-if="bookDetails.categories?.length" class="flex flex-wrap gap-2 mb-6">
                  <span
                    v-for="category in bookDetails.categories"
                    :key="category"
                    class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded-full"
                  >
                    {{ category }}
                  </span>
                </div>

                <!-- Action Buttons -->
                <div class="flex gap-3">
                  <PrimaryButton @click="$emit('log', bookDetails)">
                    📝 Log This Book
                  </PrimaryButton>
                  <OutlineButton>
                    ➕ Add to Watchlist
                  </OutlineButton>
                </div>
              </div>
            </div>

            <!-- Description -->
            <div v-if="bookDetails.description" class="mt-6">
              <h3 class="text-lg font-semibold mb-3">Description</h3>
              <p class="text-gray-700 leading-relaxed whitespace-pre-line">
                {{ bookDetails.description }}
              </p>
            </div>

            <!-- ISBN -->
            <div v-if="bookDetails.isbn_10 || bookDetails.isbn_13" class="mt-6 pt-6 border-t">
              <h3 class="text-sm font-semibold text-gray-500 mb-2">ISBNS</h3>
              <div class="flex gap-4 text-sm text-gray-600">
                <div v-if="bookDetails.isbn_10">ISBN-10: {{ bookDetails.isbn_10 }}</div>
                <div v-if="bookDetails.isbn_13">ISBN-13: {{ bookDetails.isbn_13 }}</div>
              </div>
            </div>
          </div>

          <!-- Error State -->
          <div v-else-if="error" class="p-12 text-center">
            <div class="text-4xl mb-4">😕</div>
            <p class="text-gray-600">{{ error }}</p>
            <SecondaryButton @click="$emit('close')" class="mt-4">
              Close
            </SecondaryButton>
          </div>

          <!-- Close Button -->
          <button
            @click="$emit('close')"
            class="absolute top-4 right-4 p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </transition>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch } from 'vue'
import axios from 'axios'
import PrimaryButton from './ui/PrimaryButton.vue'
import OutlineButton from './ui/OutlineButton.vue'
import SecondaryButton from './ui/SecondaryButton.vue'

const props = defineProps({
  show: Boolean,
  bookId: String
})

const emit = defineEmits(['close', 'log'])

const bookDetails = ref(null)
const loading = ref(false)
const error = ref(null)

watch(() => props.show, async (newShow) => {
  if (newShow && props.bookId) {
    await fetchBookDetails()
  }
})

const fetchBookDetails = async () => {
  loading.value = true
  error.value = null
  bookDetails.value = null

  try {
    const response = await axios.get(`/api/books/${props.bookId}`)
    bookDetails.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to load book details'
    console.error('Error fetching book details:', err)
  } finally {
    loading.value = false
  }
}
</script>

