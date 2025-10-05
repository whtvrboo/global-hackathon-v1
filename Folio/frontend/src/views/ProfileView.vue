<template>
  <div class="min-h-screen bg-dark-950">

    <main class="container-mobile max-w-7xl mx-auto section-padding">
      <!-- Loading -->
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-2 border-dark-600 border-t-accent-red">
        </div>
        <p class="mt-4 text-dark-300">Loading your profile...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="!user && !loading" class="text-center py-16">
        <div class="text-6xl mb-6">😞</div>
        <h3 class="text-heading-2 mb-4">Profile not found</h3>
        <p class="text-body text-dark-300 mb-8 max-w-md mx-auto">
          The user profile you're looking for doesn't exist or you don't have permission to view it.
        </p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <button @click="$router.push('/')" class="btn-primary">
            Go Home
          </button>
          <button @click="$router.push('/discover')" class="btn-secondary">
            Discover Books
          </button>
        </div>
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
                U
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
              <p v-if="user.bio" class="text-body text-dark-200 mb-4">{{ user.bio }}</p>

              <!-- Edit Profile Button (Own Profile) -->
              <button v-if="isOwnProfile" @click="showEditProfile = true" class="btn-secondary text-sm">
                ✏️ Edit Profile
              </button>

              <!-- Follow/Unfollow Button (Other's Profile) -->
              <button v-else-if="authStore.isAuthenticated" @click="toggleFollow" :disabled="followLoading"
                class="btn-primary text-sm disabled:opacity-50"
                :class="isFollowing ? 'bg-dark-700 hover:bg-dark-600' : ''">
                {{ followLoading ? 'Loading...' : (isFollowing ? '✓ Following' : '+ Follow') }}
              </button>
            </div>
          </div>

          <!-- Social Stats -->
          <div v-if="user.stats" class="grid grid-cols-2 md:grid-cols-3 gap-6 mt-8 pt-6 border-t border-dark-800">
            <div class="text-center">
              <div class="text-3xl font-bold text-accent-blue">{{ user.stats.followers_count }}</div>
              <div class="text-caption">Followers</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-accent-green">{{ user.stats.following_count }}</div>
              <div class="text-caption">Following</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-accent-purple">{{ user.stats.public_lists }}</div>
              <div class="text-caption">Public Lists</div>
            </div>
          </div>

          <!-- Reading Stats Grid -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-6 mt-6 pt-6 border-t border-dark-800">
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

          <!-- Favorite Books Section -->
          <div v-if="favoriteBooks.length > 0" class="mt-6 pt-6 border-t border-dark-800">
            <h3 class="text-heading-3 mb-4 flex items-center gap-2">
              <span></span>
              <span>Favorite Books</span>
            </h3>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div v-for="book in favoriteBooks" :key="book.id" class="group cursor-pointer"
                @click="showBookDetail(book)">
                <div class="aspect-[2/3] relative overflow-hidden bg-dark-800 rounded-xl mb-2">
                  <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                    class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                  <div v-else class="w-full h-full flex items-center justify-center text-2xl text-dark-400">
                  </div>
                </div>
                <h4 class="text-sm font-medium text-white line-clamp-2 group-hover:text-accent-blue transition-colors">
                  {{ book.title }}
                </h4>
                <p v-if="book.authors?.length" class="text-xs text-dark-400 line-clamp-1">
                  by {{ book.authors.join(', ') }}
                </p>
              </div>
            </div>
          </div>

          <!-- Enhanced Stats with Visualizations -->
          <div v-if="stats.total > 0" class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-6 pt-6 border-t border-dark-800">
            <!-- Average Rating with Star Visualization -->
            <div class="text-center">
              <div class="flex items-center justify-center gap-2 mb-2">
                <div class="text-2xl font-bold text-accent-orange">{{ stats.averageRating.toFixed(1) }}</div>
                <div class="flex">
                  <span v-for="star in 5" :key="star"
                    :class="star <= Math.round(stats.averageRating) ? 'text-accent-orange' : 'text-dark-600'"
                    class="text-lg">★</span>
                </div>
              </div>
              <div class="text-caption">Average Rating</div>
              <div class="w-full bg-dark-700 rounded-full h-2 mt-2">
                <div
                  class="bg-gradient-to-r from-accent-orange to-accent-red h-2 rounded-full transition-all duration-500"
                  :style="{ width: `${(stats.averageRating / 5) * 100}%` }"></div>
              </div>
            </div>

            <!-- Pages Read with Progress Circle -->
            <div class="text-center">
              <div class="relative inline-flex items-center justify-center w-16 h-16 mb-2">
                <svg class="w-16 h-16 transform -rotate-90" viewBox="0 0 36 36">
                  <path d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" fill="none"
                    stroke="#374151" stroke-width="2" />
                  <path d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" fill="none"
                    stroke="#06b6d4" stroke-width="2" stroke-dasharray="100, 100"
                    :stroke-dashoffset="100 - (Math.min(stats.totalPages / 10000, 1) * 100)" />
                </svg>
                <div class="absolute text-sm font-bold text-accent-cyan">{{ Math.floor(stats.totalPages / 1000) }}k
                </div>
              </div>
              <div class="text-2xl font-bold text-accent-cyan">{{ stats.totalPages.toLocaleString() }}</div>
              <div class="text-caption">Pages Read</div>
            </div>

            <!-- This Year Reading with Trend -->
            <div class="text-center">
              <div class="text-2xl font-bold text-accent-pink mb-2">{{ stats.thisYear }}</div>
              <div class="text-caption mb-2">Read This Year</div>
              <div class="flex justify-center items-end gap-1 h-8">
                <div v-for="month in 12" :key="month" class="bg-accent-pink rounded-sm transition-all duration-500"
                  :style="{
                    height: `${Math.random() * 100}%`,
                    width: '4px'
                  }"></div>
              </div>
              <div class="text-xs text-dark-400 mt-1">Monthly Progress</div>
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

        <!-- Lists Tab Content -->
        <div v-if="currentTab === 'lists'">
          <ListManager />
        </div>

        <!-- Books Tab Content -->
        <div v-else>
          <!-- View Toggle -->
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center gap-2">
              <button @click="viewMode = 'grid'" :class="[
                'p-2 rounded-lg transition-colors',
                viewMode === 'grid' ? 'bg-accent-red text-white' : 'bg-dark-800 text-dark-300 hover:bg-dark-700'
              ]">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z">
                  </path>
                </svg>
              </button>
              <button @click="viewMode = 'timeline'" :class="[
                'p-2 rounded-lg transition-colors',
                viewMode === 'timeline' ? 'bg-accent-red text-white' : 'bg-dark-800 text-dark-300 hover:bg-dark-700'
              ]">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </button>
            </div>
            <div class="text-sm text-dark-400">
              {{ filteredLogs.length }} book{{ filteredLogs.length !== 1 ? 's' : '' }}
            </div>
          </div>

          <!-- Grid View -->
          <div v-if="viewMode === 'grid' && filteredLogs.length > 0"
            class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            <BookCard v-for="log in filteredLogs" :key="log.id" :book="{
              ...log.book,
              id: log.book_id,
              status: log.status,
              rating: log.rating
            }" @click="selectedBook = log" />
          </div>

          <!-- Timeline View -->
          <div v-else-if="viewMode === 'timeline' && filteredLogs.length > 0" class="space-y-6">
            <div v-for="group in chronologicalLogs" :key="group.period" class="card">
              <h3 class="text-heading-4 text-white mb-4">{{ group.period }}</h3>
              <div class="space-y-4">
                <div v-for="log in group.logs" :key="log.id" @click="selectedBook = log"
                  class="flex gap-4 p-4 bg-dark-800 rounded-xl hover:bg-dark-700 transition-colors cursor-pointer">
                  <img v-if="log.book.cover_url" :src="log.book.cover_url" :alt="log.book.title"
                    class="w-16 h-24 object-cover rounded-lg" />
                  <div v-else class="w-16 h-24 bg-dark-700 rounded-lg flex items-center justify-center">
                    <span class="text-2xl text-dark-400"></span>
                  </div>
                  <div class="flex-1">
                    <h4 class="text-heading-4 text-white mb-1">{{ log.book.title }}</h4>
                    <p v-if="log.book.authors" class="text-body text-dark-300 mb-2">
                      by {{ log.book.authors.join(', ') }}
                    </p>
                    <div class="flex items-center gap-4 text-sm">
                      <span class="px-2 py-1 rounded-full text-xs font-semibold" :class="statusBadgeClass(log.status)">
                        {{ statusLabel(log.status) }}
                      </span>
                      <span v-if="log.rating" class="text-accent-orange">
                        {{ '★'.repeat(log.rating) }}{{ '☆'.repeat(5 - log.rating) }}
                      </span>
                      <span class="text-dark-400">{{ formatDate(log.created_at) }}</span>
                    </div>
                    <p v-if="log.review" class="text-sm text-dark-300 mt-2 line-clamp-2">
                      {{ log.review }}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="card text-center py-16">
            <div class="text-6xl mb-6"></div>
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


    <!-- Profile Edit Modal -->
    <ProfileEditModal :show="showEditProfile" :current-user="user" :user-logs="logs" @close="showEditProfile = false"
      @updated="refreshProfile" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRoute } from 'vue-router'
import axios from 'axios'
import Card from '../components/ui/Card.vue'
import BookCard from '../components/BookCard.vue'
import PrimaryButton from '../components/ui/PrimaryButton.vue'
import ListManager from '../components/ListManager.vue'
import ProfileEditModal from '../components/ProfileEditModal.vue'
import { useToastStore } from '../stores/toast'

const authStore = useAuthStore()
const route = useRoute()
const toastStore = useToastStore()

const user = ref(null)
const logs = ref([])
const favoriteBooks = ref([])
const loading = ref(true)
const currentTab = ref('lists')
const selectedBook = ref(null)
const viewMode = ref('grid')
const showEditProfile = ref(false)
const isFollowing = ref(false)
const followLoading = ref(false)

const isOwnProfile = computed(() => {
  // Check if the current user (authenticated or guest) is viewing their own profile
  if (!authStore.isAuthenticated || !authStore.user) {
    return false
  }
  return authStore.user.username === route.params.username
})

const tabs = [
  { label: 'Lists', value: 'lists' },
  { label: 'All', value: 'all' },
  { label: 'Read', value: 'read' },
  { label: 'Reading', value: 'reading' },
  { label: 'Want to Read', value: 'want_to_read' }
]

const stats = computed(() => {
  const readLogs = logs.value.filter(l => l.status === 'read')
  const ratings = readLogs.filter(l => l.rating).map(l => l.rating)
  const currentYear = new Date().getFullYear()

  return {
    total: logs.value.length,
    read: readLogs.length,
    reading: logs.value.filter(l => l.status === 'reading').length,
    wantToRead: logs.value.filter(l => l.status === 'want_to_read').length,
    averageRating: ratings.length > 0 ? ratings.reduce((a, b) => a + b, 0) / ratings.length : 0,
    totalPages: readLogs.reduce((total, log) => total + (log.book?.page_count || 0), 0),
    thisYear: readLogs.filter(l => {
      const logYear = new Date(l.created_at).getFullYear()
      return logYear === currentYear
    }).length
  }
})

const filteredLogs = computed(() => {
  if (currentTab.value === 'all') return logs.value
  return logs.value.filter(log => log.status === currentTab.value)
})

const chronologicalLogs = computed(() => {
  const grouped = {}

  filteredLogs.value.forEach(log => {
    const date = new Date(log.created_at)
    const period = getTimePeriod(date)

    if (!grouped[period]) {
      grouped[period] = []
    }
    grouped[period].push(log)
  })

  // Sort periods and logs within each period
  return Object.keys(grouped)
    .sort((a, b) => {
      // Custom sort order for periods
      const order = ['This Week', 'Last Week', 'This Month', 'Last Month', 'This Year', 'Last Year', 'Older']
      return order.indexOf(a) - order.indexOf(b)
    })
    .map(period => ({
      period,
      logs: grouped[period].sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    }))
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

const statusBadgeClass = (status) => {
  const classes = {
    'want_to_read': 'text-accent-blue bg-accent-blue/20',
    'reading': 'text-accent-green bg-accent-green/20',
    'read': 'text-accent-purple bg-accent-purple/20',
    'dnf': 'text-dark-400 bg-dark-600/20'
  }
  return classes[status] || 'text-dark-400 bg-dark-600/20'
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const getTimePeriod = (date) => {
  const now = new Date()
  const diffTime = now - date
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays <= 7) return 'This Week'
  if (diffDays <= 14) return 'Last Week'
  if (diffDays <= 30) return 'This Month'
  if (diffDays <= 60) return 'Last Month'
  if (date.getFullYear() === now.getFullYear()) return 'This Year'
  if (date.getFullYear() === now.getFullYear() - 1) return 'Last Year'
  return 'Older'
}

onMounted(async () => {
  try {
    const username = route.params.username

    // Ensure auth store is initialized for guest users
    if (authStore.isAuthenticated && !authStore.user) {
      await authStore.fetchUser()
    }

    // Fetch user profile
    const profileResponse = await axios.get(`/api/users/${username}`)
    user.value = profileResponse.data
    isFollowing.value = profileResponse.data.is_following

    // Fetch user's logs
    const logsResponse = await axios.get(`/api/users/${username}/logs`)
    logs.value = logsResponse.data.logs || []

    // Load favorite books (for now, get highest rated books)
    if (logs.value.length > 0) {
      const readLogs = logs.value.filter(l => l.status === 'read' && l.rating >= 4)
      favoriteBooks.value = readLogs
        .sort((a, b) => (b.rating || 0) - (a.rating || 0))
        .slice(0, 4)
        .map(log => log.book)
    } else {
      // Initialize empty arrays for new guest users
      favoriteBooks.value = []
    }
  } catch (error) {
    console.error('Error loading profile:', error)
    // Show error state instead of blank screen
    user.value = null
  } finally {
    loading.value = false
  }
})

const toggleFollow = async () => {
  if (followLoading.value) return

  followLoading.value = true
  try {
    if (isFollowing.value) {
      await axios.delete(`/api/users/${user.value.username}/follow`)
      isFollowing.value = false
      toastStore.info('Unfollowed')
    } else {
      await axios.post(`/api/users/${user.value.username}/follow`)
      isFollowing.value = true
      toastStore.success(`You're now following ${user.value.name}!`)
    }
  } catch (error) {
    console.error('Error toggling follow:', error)
    toastStore.error('Failed to update follow status')
  } finally {
    followLoading.value = false
  }
}

const showBookDetail = (book) => {
  selectedBook.value = book
  // showBookDetailModal.value = true // BookDetailModal may not exist yet
}

const refreshProfile = async () => {
  try {
    // Refresh user data from API
    const response = await axios.get('/api/me')
    user.value = response.data
    authStore.setUser(response.data)

    // Reload favorite books based on updated favorite_book_ids
    if (user.value.favorite_book_ids?.length > 0 && logs.value.length > 0) {
      favoriteBooks.value = logs.value
        .filter(log => user.value.favorite_book_ids.includes(log.book_id))
        .map(log => log.book)
        .filter(book => book)
        .slice(0, 4)
    }
  } catch (error) {
    console.error('Error refreshing profile:', error)
  }
}
</script>
