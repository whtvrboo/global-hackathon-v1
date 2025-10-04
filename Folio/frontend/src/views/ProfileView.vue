<template>
  <div class="min-h-screen bg-dark-950">
    <!-- Header -->
    <header class="glass-strong border-b border-dark-800 sticky top-0 z-50">
      <div class="container-mobile max-w-7xl mx-auto">
        <div class="flex items-center justify-between py-4">
          <!-- Logo -->
          <router-link to="/" class="flex items-center gap-3">
            <div
              class="w-8 h-8 bg-gradient-to-br from-accent-red to-accent-blue rounded-lg flex items-center justify-center">
              <span class="text-white font-bold text-sm">F</span>
            </div>
            <h1 class="text-xl font-bold text-white">Folio</h1>
          </router-link>

          <!-- Desktop Navigation -->
          <nav class="hidden md:flex items-center gap-6">
            <router-link to="/feed" class="btn-ghost">Feed</router-link>
            <button @click="authStore.logout()" class="btn-ghost">Logout</button>
          </nav>

          <!-- Mobile Menu Button -->
          <button class="md:hidden p-2 text-dark-300 hover:text-white transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
            </svg>
          </button>
        </div>
      </div>
    </header>

    <main class="container-mobile max-w-7xl mx-auto section-padding">
      <!-- Loading -->
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-2 border-dark-600 border-t-accent-red">
        </div>
        <p class="mt-4 text-dark-300">Loading your profile...</p>
      </div>

      <!-- Profile Content -->
      <div v-else-if="user" class="space-y-8">
        <!-- User Header -->
        <div class="card card-hover">
          <div class="flex flex-col md:flex-row items-start gap-6">
            <div class="relative">
              <img v-if="user.picture" :src="user.picture" :alt="user.name"
                class="w-24 h-24 rounded-full border-2 border-dark-700" />
              <div v-else
                class="w-24 h-24 rounded-full bg-dark-800 border-2 border-dark-700 flex items-center justify-center text-2xl">
                👤
              </div>
              <div v-if="authStore.isGuestUser"
                class="absolute -top-1 -right-1 w-6 h-6 bg-accent-orange rounded-full flex items-center justify-center">
                <span class="text-white text-xs font-bold">G</span>
              </div>
            </div>

            <div class="flex-1">
              <div class="flex flex-col md:flex-row md:items-center gap-2 mb-2">
                <h1 class="text-heading-1">{{ user.name }}</h1>
                <span v-if="authStore.isGuestUser"
                  class="inline-flex items-center px-3 py-1 text-xs font-medium bg-accent-orange/20 text-accent-orange rounded-full border border-accent-orange/30">
                  Guest Account
                </span>
              </div>
              <p class="text-body text-dark-300 mb-4">@{{ user.username }}</p>
              <p v-if="user.bio" class="text-body text-dark-200">{{ user.bio }}</p>
            </div>
          </div>

          <!-- Stats Grid -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-6 mt-8 pt-6 border-t border-dark-800">
            <div class="text-center">
              <div class="text-3xl font-bold text-accent-red">{{ stats.total }}</div>
              <div class="text-caption">Total Books</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-accent-green">{{ stats.read }}</div>
              <div class="text-caption">Books Read</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-accent-blue">{{ stats.reading }}</div>
              <div class="text-caption">Currently Reading</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-accent-purple">{{ stats.wantToRead }}</div>
              <div class="text-caption">Want to Read</div>
            </div>
          </div>
        </div>

        <!-- Filter Tabs -->
        <div class="flex gap-2 overflow-x-auto scrollbar-hide">
          <button v-for="tab in tabs" :key="tab.value" @click="currentTab = tab.value" :class="[
            'px-6 py-3 font-medium rounded-xl transition-all duration-200 whitespace-nowrap',
            currentTab === tab.value
              ? 'bg-accent-red text-white shadow-lg'
              : 'bg-dark-800 text-dark-300 hover:bg-dark-700 hover:text-white'
          ]">
            {{ tab.label }}
          </button>
        </div>

        <!-- Books Grid -->
        <div v-if="filteredLogs.length > 0" class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <BookCard v-for="log in filteredLogs" :key="log.id" :book="{
            ...log.book,
            id: log.book_id,
            status: log.status,
            rating: log.rating
          }" @click="selectedBook = log" />
        </div>

        <!-- Empty State -->
        <div v-else class="card text-center py-16">
          <div class="text-6xl mb-6">📚</div>
          <h3 class="text-heading-2 mb-4">No books yet</h3>
          <p class="text-body text-dark-300 mb-8 max-w-md mx-auto">
            Start building your reading journal by discovering and logging your first book
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <button @click="$router.push('/discover')" class="btn-primary">
              Discover Books
            </button>
            <button @click="$router.push('/')" class="btn-secondary">
              Go Home
            </button>
          </div>
        </div>
      </div>
    </main>

    <!-- Book Detail Modal (for viewing logged books) -->
    <transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0"
      leave-active-class="transition ease-in duration-150" leave-to-class="opacity-0">
      <div v-if="selectedBook" @click="selectedBook = null"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <Card @click.stop class="max-w-2xl w-full max-h-[90vh] overflow-y-auto p-8">
          <button @click="selectedBook = null" class="absolute top-4 right-4 p-2 text-gray-400 hover:text-gray-600">
            ✕
          </button>

          <div class="flex gap-6">
            <img v-if="selectedBook.book?.cover_url" :src="selectedBook.book.cover_url" :alt="selectedBook.book.title"
              class="w-32 h-48 object-cover rounded-lg" />
            <div class="flex-1">
              <h2 class="text-2xl font-bold mb-2">{{ selectedBook.book?.title }}</h2>
              <p v-if="selectedBook.book?.authors" class="text-gray-600 mb-4">
                {{ selectedBook.book.authors.join(', ') }}
              </p>

              <div class="space-y-3">
                <div>
                  <span class="text-sm font-medium text-gray-500">Status:</span>
                  <span class="ml-2">{{ statusLabel(selectedBook.status) }}</span>
                </div>

                <div v-if="selectedBook.rating">
                  <span class="text-sm font-medium text-gray-500">Rating:</span>
                  <span class="ml-2 text-yellow-500">
                    {{ '★'.repeat(selectedBook.rating) }}
                  </span>
                </div>

                <div v-if="selectedBook.review" class="pt-4 border-t">
                  <p class="text-sm font-medium text-gray-500 mb-2">Review:</p>
                  <p class="text-gray-700">{{ selectedBook.review }}</p>
                </div>

                <div v-if="selectedBook.start_date || selectedBook.finish_date" class="pt-4 border-t">
                  <div v-if="selectedBook.start_date" class="text-sm text-gray-600">
                    Started: {{ new Date(selectedBook.start_date).toLocaleDateString() }}
                  </div>
                  <div v-if="selectedBook.finish_date" class="text-sm text-gray-600">
                    Finished: {{ new Date(selectedBook.finish_date).toLocaleDateString() }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </Card>
      </div>
    </transition>

    <!-- Bottom Navigation (Mobile) -->
    <nav class="md:hidden fixed bottom-0 left-0 right-0 glass-strong border-t border-dark-800 z-40">
      <div class="flex items-center justify-around py-2">
        <router-link to="/"
          class="flex flex-col items-center gap-1 p-3 text-dark-400 hover:text-white transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6">
            </path>
          </svg>
          <span class="text-xs">Home</span>
        </router-link>

        <router-link to="/discover"
          class="flex flex-col items-center gap-1 p-3 text-dark-400 hover:text-white transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
          <span class="text-xs">Discover</span>
        </router-link>

        <router-link to="/feed"
          class="flex flex-col items-center gap-1 p-3 text-dark-400 hover:text-white transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16">
            </path>
          </svg>
          <span class="text-xs">Feed</span>
        </router-link>

        <router-link to="/profile" class="flex flex-col items-center gap-1 p-3 text-accent-red transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
          </svg>
          <span class="text-xs">Profile</span>
        </router-link>
      </div>
    </nav>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'
import Card from '../components/ui/Card.vue'
import BookCard from '../components/BookCard.vue'
import PrimaryButton from '../components/ui/PrimaryButton.vue'

const authStore = useAuthStore()

const user = ref(null)
const logs = ref([])
const loading = ref(true)
const currentTab = ref('all')
const selectedBook = ref(null)

const tabs = [
  { label: 'All', value: 'all' },
  { label: 'Read', value: 'read' },
  { label: 'Reading', value: 'reading' },
  { label: 'Want to Read', value: 'want_to_read' }
]

const stats = computed(() => ({
  total: logs.value.length,
  read: logs.value.filter(l => l.status === 'read').length,
  reading: logs.value.filter(l => l.status === 'reading').length,
  wantToRead: logs.value.filter(l => l.status === 'want_to_read').length
}))

const filteredLogs = computed(() => {
  if (currentTab.value === 'all') return logs.value
  return logs.value.filter(log => log.status === currentTab.value)
})

const statusLabel = (status) => {
  const labels = {
    'want_to_read': 'Want to Read',
    'reading': 'Reading',
    'read': 'Read',
    'dnf': 'DNF'
  }
  return labels[status] || status
}

onMounted(async () => {
  try {
    // Get current user
    user.value = authStore.user

    // Fetch user's logs
    if (user.value?.username) {
      const response = await axios.get(`/api/users/${user.value.username}/logs`)
      logs.value = response.data.logs || []
    }
  } catch (error) {
    console.error('Error loading profile:', error)
  } finally {
    loading.value = false
  }
})
</script>
