<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white border-b sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <router-link to="/" class="text-2xl font-bold">📚 Folio</router-link>
          <div class="flex items-center gap-4">
            <router-link to="/feed" class="text-gray-600 hover:text-gray-900">Feed</router-link>
            <button @click="authStore.logout()" class="text-gray-600 hover:text-gray-900">
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8">
      <!-- Loading -->
      <div v-if="loading" class="text-center py-12">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>

      <!-- Profile Content -->
      <div v-else-if="user" class="space-y-8">
        <!-- User Header -->
        <Card>
          <div class="flex items-start gap-6">
            <img
              v-if="user.picture"
              :src="user.picture"
              :alt="user.name"
              class="w-24 h-24 rounded-full"
            />
            <div class="flex-1">
              <h1 class="text-3xl font-bold mb-1">{{ user.name }}</h1>
              <p class="text-gray-600 mb-4">@{{ user.username }}</p>
              <p v-if="user.bio" class="text-gray-700">{{ user.bio }}</p>
            </div>
          </div>

          <!-- Stats -->
          <div class="grid grid-cols-3 gap-6 mt-6 pt-6 border-t">
            <div class="text-center">
              <div class="text-3xl font-bold text-primary">{{ stats.total }}</div>
              <div class="text-sm text-gray-600">Total Books</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-success">{{ stats.read }}</div>
              <div class="text-sm text-gray-600">Books Read</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-secondary">{{ stats.reading }}</div>
              <div class="text-sm text-gray-600">Currently Reading</div>
            </div>
          </div>
        </Card>

        <!-- Filter Tabs -->
        <div class="flex gap-2 border-b">
          <button
            v-for="tab in tabs"
            :key="tab.value"
            @click="currentTab = tab.value"
            :class="[
              'px-4 py-2 font-medium border-b-2 transition-colors',
              currentTab === tab.value
                ? 'border-primary text-primary'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            ]"
          >
            {{ tab.label }}
          </button>
        </div>

        <!-- Books Grid -->
        <div v-if="filteredLogs.length > 0" class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <BookCard
            v-for="log in filteredLogs"
            :key="log.id"
            :book="{
              ...log.book,
              id: log.book_id,
              status: log.status,
              rating: log.rating
            }"
            @click="selectedBook = log"
          />
        </div>

        <!-- Empty State -->
        <div v-else class="text-center py-12">
          <div class="text-6xl mb-4">📚</div>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">No books yet</h3>
          <p class="text-gray-600 mb-6">Start building your reading journal</p>
          <PrimaryButton @click="$router.push('/')">
            Search for Books
          </PrimaryButton>
        </div>
      </div>
    </div>

    <!-- Book Detail Modal (for viewing logged books) -->
    <transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0"
      leave-active-class="transition ease-in duration-150"
      leave-to-class="opacity-0"
    >
      <div
        v-if="selectedBook"
        @click="selectedBook = null"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
      >
        <Card @click.stop class="max-w-2xl w-full max-h-[90vh] overflow-y-auto p-8">
          <button
            @click="selectedBook = null"
            class="absolute top-4 right-4 p-2 text-gray-400 hover:text-gray-600"
          >
            ✕
          </button>

          <div class="flex gap-6">
            <img
              v-if="selectedBook.book?.cover_url"
              :src="selectedBook.book.cover_url"
              :alt="selectedBook.book.title"
              class="w-32 h-48 object-cover rounded-lg"
            />
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
  reading: logs.value.filter(l => l.status === 'reading').length
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

