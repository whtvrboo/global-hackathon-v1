<template>
  <div class="relative w-full max-w-2xl">
    <div class="relative">
      <input
        v-model="searchQuery"
        @input="handleSearch"
        @focus="showResults = true"
        type="text"
        placeholder="Search for books..."
        class="w-full px-5 py-3 pl-12 pr-12 text-lg border-2 border-gray-200 rounded-xl
               focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary
               transition-all duration-200 shadow-sm"
      />
      <div class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">
        🔍
      </div>
      <button
        v-if="searchQuery"
        @click="clearSearch"
        class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
      >
        ✕
      </button>
    </div>

    <!-- Search Results Dropdown -->
    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="opacity-0 translate-y-1"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-1"
    >
      <div
        v-if="showResults && (loading || results.length > 0 || error)"
        class="absolute z-50 w-full mt-2 bg-white rounded-xl border border-gray-200 shadow-xl max-h-96 overflow-y-auto"
      >
        <!-- Loading -->
        <div v-if="loading" class="p-8 text-center">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
          <p class="mt-2 text-gray-600">Searching...</p>
        </div>

        <!-- Error -->
        <div v-else-if="error" class="p-6 text-center text-red-600">
          <p>{{ error }}</p>
        </div>

        <!-- Results -->
        <div v-else-if="results.length > 0" class="divide-y divide-gray-100">
          <button
            v-for="book in results"
            :key="book.id"
            @click="selectBook(book)"
            class="w-full p-4 flex gap-4 hover:bg-gray-50 transition-colors text-left"
          >
            <img
              v-if="book.cover_url"
              :src="book.cover_url"
              :alt="book.title"
              class="w-16 h-24 object-cover rounded shadow-sm"
            />
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-gray-900 truncate">{{ book.title }}</h3>
              <p v-if="book.authors?.length" class="text-sm text-gray-600 truncate">
                {{ book.authors.join(', ') }}
              </p>
              <p v-if="book.published_date" class="text-xs text-gray-500 mt-1">
                {{ book.published_date }}
              </p>
            </div>
          </button>
        </div>

        <!-- No Results -->
        <div v-else class="p-6 text-center text-gray-500">
          <p>No books found</p>
        </div>
      </div>
    </transition>

    <!-- Backdrop to close results -->
    <div
      v-if="showResults"
      @click="showResults = false"
      class="fixed inset-0 z-40"
    ></div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import axios from 'axios'

const emit = defineEmits(['select'])

const searchQuery = ref('')
const results = ref([])
const loading = ref(false)
const error = ref(null)
const showResults = ref(false)
let searchTimeout = null

const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  
  if (!searchQuery.value.trim()) {
    results.value = []
    return
  }

  loading.value = true
  error.value = null

  searchTimeout = setTimeout(async () => {
    try {
      const response = await axios.get('/api/search', {
        params: { q: searchQuery.value }
      })
      results.value = response.data.results || []
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to search books'
      console.error('Search error:', err)
    } finally {
      loading.value = false
    }
  }, 300) // Debounce 300ms
}

const selectBook = (book) => {
  emit('select', book)
  showResults.value = false
  searchQuery.value = ''
  results.value = []
}

const clearSearch = () => {
  searchQuery.value = ''
  results.value = []
  error.value = null
}
</script>

