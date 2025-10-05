<template>
  <div class="relative w-full max-w-2xl">
    <div class="relative">
      <input v-model="searchQuery" @input="handleSearch" @focus="showResults = true" type="text"
        placeholder="Search for books..." class="input input-search w-full px-5 py-4 pl-14 pr-14 text-lg rounded-2xl" />
      <div class="absolute left-5 top-1/2 -translate-y-1/2 text-dark-400">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
        </svg>
      </div>
      <button v-if="searchQuery" @click="clearSearch"
        class="absolute right-5 top-1/2 -translate-y-1/2 text-dark-400 hover:text-white transition-colors">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
        </svg>
      </button>
    </div>

    <!-- Search Results Dropdown -->
    <transition enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 translate-y-2 scale-95" enter-to-class="opacity-100 translate-y-0 scale-100"
      leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100 translate-y-0 scale-100"
      leave-to-class="opacity-0 translate-y-2 scale-95">
      <div v-if="showResults && (loading || results.length > 0 || error)"
        class="absolute z-50 w-full mt-3 glass-strong rounded-2xl border border-dark-700 shadow-2xl shadow-black/20 max-h-96 overflow-y-auto">
        <!-- Loading -->
        <div v-if="loading" class="p-8 text-center">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-2 border-dark-600 border-t-accent-red">
          </div>
          <p class="mt-3 text-dark-300">Searching...</p>
        </div>

        <!-- Error -->
        <div v-else-if="error" class="p-6 text-center text-accent-red">
          <p>{{ error }}</p>
        </div>

        <!-- Results -->
        <div v-else-if="results.length > 0" class="divide-y divide-dark-800">
          <button v-for="book in results" :key="book.id" @click="selectBook(book)"
            class="w-full p-4 flex gap-4 hover:bg-dark-800/50 transition-colors text-left group">
            <div class="relative">
              <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                class="w-16 h-24 object-cover rounded-xl shadow-lg group-hover:scale-105 transition-transform duration-200" />
              <div v-else class="w-16 h-24 bg-dark-800 rounded-xl flex items-center justify-center">
                <span class="text-2xl text-dark-400"></span>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-white truncate mb-1">{{ book.title }}</h3>
              <p v-if="book.authors?.length" class="text-sm text-dark-300 truncate mb-1">
                {{ book.authors.join(', ') }}
              </p>
              <p v-if="book.published_date" class="text-xs text-dark-500">
                {{ book.published_date }}
              </p>
            </div>
          </button>
        </div>

        <!-- No Results -->
        <div v-else class="p-6 text-center text-dark-400">
          <div class="text-4xl mb-2">🔍</div>
          <p>No books found</p>
        </div>
      </div>
    </transition>

    <!-- Backdrop to close results -->
    <div v-if="showResults" @click="showResults = false" class="fixed inset-0 z-40"></div>
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
