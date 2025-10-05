<template>
  <div v-if="loading" class="text-center p-12">
    <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
    <p class="mt-4 text-dark-300">Loading profile...</p>
  </div>
  <div v-else-if="error" class="text-center p-12">
    <div class="text-4xl mb-4">😞</div>
    <p class="text-dark-300">{{ error }}</p>
  </div>
  <div v-else-if="user" class="bg-dark-950 min-h-screen">
    <!-- Profile Hero Section - Compact -->
    <div
      class="relative w-full h-[25vh] text-white flex items-center justify-center text-center p-4 bg-gradient-to-br from-accent-red via-accent-purple to-accent-blue">
      <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/50 to-transparent"></div>
      <div class="relative z-10 max-w-4xl mx-auto">
        <div class="flex flex-col md:flex-row items-center gap-4">
          <!-- Profile Picture -->
          <div class="relative flex-shrink-0">
            <img v-if="user.picture" :src="user.picture" :alt="user.name"
              class="w-20 h-20 rounded-full border-3 border-white/50 shadow-2xl" />
            <div v-else
              class="w-20 h-20 rounded-full bg-white/20 border-3 border-white/50 flex items-center justify-center text-2xl font-bold shadow-2xl">
              {{ user.name.charAt(0).toUpperCase() }}
            </div>
            <div v-if="authStore.isGuestUser"
              class="absolute -top-1 -right-1 w-5 h-5 bg-accent-orange rounded-full flex items-center justify-center shadow-lg">
              <span class="text-white text-xs font-bold">G</span>
            </div>
          </div>

          <!-- User Info -->
          <div class="text-center md:text-left flex-1">
            <div class="flex flex-col md:flex-row md:items-center gap-2 mb-2">
              <h1 class="text-2xl md:text-3xl font-bold">{{ user.name }}</h1>
              <span v-if="authStore.isGuestUser"
                class="inline-flex items-center px-2 py-1 text-xs font-medium bg-accent-orange/30 text-white rounded-full border border-white/30 backdrop-blur-sm">
                Guest Account
              </span>
            </div>
            <p class="text-base text-white/90 mb-2">@{{ user.username }}</p>
            <p v-if="user.bio" class="text-sm text-white/80 mb-3 max-w-2xl">{{ user.bio }}</p>

            <!-- Action Buttons -->
            <div class="flex flex-col sm:flex-row gap-2 justify-center md:justify-start">
              <button v-if="isOwnProfile" @click="showEditProfile = true"
                class="px-4 py-2 bg-white/20 text-white rounded-lg border border-white/30 hover:bg-white/30 transition-all backdrop-blur-sm text-sm">
                ✏️ Edit Profile
              </button>
              <button v-else-if="authStore.isAuthenticated" @click="toggleFollow" :disabled="followLoading"
                class="px-4 py-2 rounded-lg border transition-all backdrop-blur-sm disabled:opacity-50 text-sm" :class="isFollowing
                  ? 'bg-white/20 text-white border-white/30 hover:bg-white/30'
                  : 'bg-white text-dark-950 border-white hover:bg-white/90'">
                {{ followLoading ? 'Loading...' : (isFollowing ? '✓ Following' : '+ Follow') }}
              </button>
              <button v-if="isOwnProfile" @click="showInviteModal = true"
                class="px-4 py-2 bg-accent-blue/20 text-white rounded-lg border border-accent-blue/30 hover:bg-accent-blue/30 transition-all backdrop-blur-sm text-sm">
                👥 Invite Friends
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Profile Content - Compact -->
    <div class="container mx-auto max-w-7xl py-4 px-4">
      <!-- Stats Overview - Compact -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
        <!-- Social Stats -->
        <div v-if="user.stats" class="card p-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-accent-blue mb-1">{{ user.stats.followers_count }}</div>
            <div class="text-caption text-xs">Followers</div>
          </div>
        </div>
        <div v-if="user.stats" class="card p-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-accent-green mb-1">{{ user.stats.following_count }}</div>
            <div class="text-caption text-xs">Following</div>
          </div>
        </div>
        <div v-if="user.stats" class="card p-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-accent-purple mb-1">{{ user.stats.public_lists }}</div>
            <div class="text-caption text-xs">Lists</div>
          </div>
        </div>
        <div class="card p-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-accent-red mb-1">{{ stats.total }}</div>
            <div class="text-caption text-xs">Books</div>
          </div>
        </div>
      </div>

      <!-- Quick Analytics - Compact -->
      <div v-if="stats.total > 0" class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
        <div class="card p-4">
          <div class="text-center">
            <div class="text-xl font-bold text-accent-orange mb-1">{{ stats.averageRating.toFixed(1) }}</div>
            <div class="text-caption text-xs">Avg Rating</div>
          </div>
        </div>
        <div class="card p-4">
          <div class="text-center">
            <div class="text-xl font-bold text-accent-cyan mb-1">{{ Math.floor(stats.totalPages / 1000) }}k</div>
            <div class="text-caption text-xs">Pages</div>
          </div>
        </div>
        <div class="card p-4">
          <div class="text-center">
            <div class="text-xl font-bold text-accent-pink mb-1">{{ stats.thisYear }}</div>
            <div class="text-caption text-xs">This Year</div>
          </div>
        </div>
        <div class="card p-4">
          <div class="text-center">
            <div class="text-xl font-bold text-accent-green mb-1">{{ stats.read }}</div>
            <div class="text-caption text-xs">Completed</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Filter Tabs - Compact -->
    <div class="container mx-auto max-w-7xl px-4">
      <div class="flex gap-2 overflow-x-auto scrollbar-hide mb-4">
        <button v-for="tab in tabs" :key="tab.value" @click="currentTab = tab.value" :class="[
          'px-4 py-2 font-medium rounded-lg transition-all duration-200 whitespace-nowrap text-sm',
          currentTab === tab.value
            ? 'bg-accent-red text-white shadow-lg'
            : 'bg-dark-800 text-dark-300 hover:bg-dark-700 hover:text-white'
        ]">
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Lists Tab Content -->
    <div v-if="currentTab === 'lists'">
      <ListManager />
    </div>

    <!-- Books Tab Content -->
    <div v-else class="container mx-auto max-w-7xl px-4">
      <!-- View Toggle - Compact -->
      <div class="card p-4 mb-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <button @click="viewMode = 'grid'" :class="[
              'p-2 rounded-lg transition-all duration-200',
              viewMode === 'grid' ? 'bg-accent-red text-white shadow-lg' : 'bg-dark-800 text-dark-300 hover:bg-dark-700'
            ]">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z">
                </path>
              </svg>
            </button>
            <button @click="viewMode = 'timeline'" :class="[
              'p-2 rounded-lg transition-all duration-200',
              viewMode === 'timeline' ? 'bg-accent-red text-white shadow-lg' : 'bg-dark-800 text-dark-300 hover:bg-dark-700'
            ]">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </button>
          </div>
          <div class="text-sm text-dark-400 font-medium">
            {{ filteredLogs.length }} book{{ filteredLogs.length !== 1 ? 's' : '' }}
          </div>
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

      <!-- Timeline View - Compact -->
      <div v-else-if="viewMode === 'timeline' && filteredLogs.length > 0" class="space-y-4">
        <div v-for="group in chronologicalLogs" :key="group.period" class="card p-4">
          <h3 class="text-heading-5 text-white mb-3 flex items-center gap-2">
            <span class="text-lg">📅</span>
            <span>{{ group.period }}</span>
          </h3>
          <div class="space-y-3">
            <div v-for="log in group.logs" :key="log.id" @click="selectedBook = log"
              class="bg-dark-800 rounded-lg p-3 hover:bg-dark-700 transition-all duration-200 cursor-pointer border border-dark-700 hover:border-dark-600">
              <div class="flex gap-3">
                <img v-if="log.book.cover_url" :src="log.book.cover_url" :alt="log.book.title"
                  class="w-12 h-16 object-cover rounded shadow-md flex-shrink-0" />
                <div v-else
                  class="w-12 h-16 bg-dark-700 rounded flex items-center justify-center shadow-md flex-shrink-0">
                  <span class="text-lg text-dark-400">📚</span>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-start justify-between mb-1">
                    <div>
                      <h4 class="text-sm font-medium text-white mb-1 line-clamp-1">{{ log.book.title }}</h4>
                      <p v-if="log.book.authors" class="text-xs text-dark-300 mb-1 line-clamp-1">
                        by {{ log.book.authors.join(', ') }}
                      </p>
                    </div>
                    <span class="text-xs text-dark-400">{{ formatDate(log.created_at) }}</span>
                  </div>

                  <div class="flex items-center gap-2 mb-2">
                    <span class="px-2 py-1 rounded-full text-xs font-semibold" :class="statusBadgeClass(log.status)">
                      {{ statusLabel(log.status) }}
                    </span>
                    <div v-if="log.rating" class="flex items-center gap-1">
                      <span class="text-accent-orange text-xs">{{ '★'.repeat(log.rating) }}{{ '☆'.repeat(5 - log.rating)
                      }}</span>
                    </div>
                  </div>

                  <!-- Review Content - Compact -->
                  <div v-if="log.review" class="bg-dark-700 rounded p-2 mb-1">
                    <p class="text-xs text-dark-100 leading-relaxed line-clamp-2">{{ log.review }}</p>
                  </div>

                  <!-- Notes -->
                  <div v-if="log.notes" class="text-xs text-dark-300 bg-dark-700/50 rounded p-2">
                    <p class="whitespace-pre-line line-clamp-1">{{ log.notes }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State - Compact -->
      <div v-else class="card text-center py-12">
        <div class="text-6xl mb-6">📚</div>
        <h3 class="text-heading-3 mb-3">No books yet</h3>
        <p class="text-body text-dark-300 mb-6 max-w-md mx-auto">
          Start building your reading journal by discovering and logging your first book
        </p>
        <div class="flex flex-col sm:flex-row gap-3 justify-center">
          <button @click="$router.push('/discover')" class="btn-primary px-6 py-2">
            Discover Books
          </button>
          <button @click="$router.push('/')" class="btn-secondary px-6 py-2">
            Go Home
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Invite Friends Modal -->
  <div v-if="showInviteModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
    <div class="bg-dark-900 rounded-xl p-6 w-full max-w-md">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-heading-4 text-white">Invite Friends</h3>
        <button @click="showInviteModal = false" class="text-dark-400 hover:text-white">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>

      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-white mb-2">Share your profile</label>
          <div class="flex gap-2">
            <input :value="profileUrl" readonly
              class="flex-1 bg-dark-800 border border-dark-700 rounded-lg px-3 py-2 text-sm text-white" />
            <button @click="copyProfileUrl"
              class="px-4 py-2 bg-accent-blue text-white rounded-lg hover:bg-accent-blue/80 transition-colors text-sm">
              Copy
            </button>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-white mb-2">Or invite via email</label>
          <div class="flex gap-2">
            <input v-model="inviteEmail" type="email" placeholder="friend@example.com"
              class="flex-1 bg-dark-800 border border-dark-700 rounded-lg px-3 py-2 text-sm text-white placeholder-dark-400" />
            <button @click="sendInvite" :disabled="!inviteEmail || inviteLoading"
              class="px-4 py-2 bg-accent-green text-white rounded-lg hover:bg-accent-green/80 transition-colors text-sm disabled:opacity-50">
              {{ inviteLoading ? 'Sending...' : 'Send' }}
            </button>
          </div>
        </div>

        <div class="text-xs text-dark-400">
          Share your reading journey with friends and discover new books together!
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import Card from '../components/ui/Card.vue'
import BookCard from '../components/BookCard.vue'
import PrimaryButton from '../components/ui/PrimaryButton.vue'
import ListManager from '../components/ListManager.vue'
import ProfileEditModal from '../components/ProfileEditModal.vue'
import { useToastStore } from '../stores/toast'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const toastStore = useToastStore()

const user = ref(null)
const logs = ref([])
const favoriteBooks = ref([])
const loading = ref(true)
const error = ref(null)
const currentTab = ref('lists')
const selectedBook = ref(null)
const viewMode = ref('grid')
const showEditProfile = ref(false)
const showInviteModal = ref(false)
const isFollowing = ref(false)
const followLoading = ref(false)
const inviteEmail = ref('')
const inviteLoading = ref(false)

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

const fetchProfile = async () => {
  loading.value = true
  error.value = null
  try {
    const username = route.params.username
    const profileResponse = await axios.get(`/api/users/${username}`)
    user.value = profileResponse.data
    isFollowing.value = profileResponse.data.is_following

    const logsResponse = await axios.get(`/api/users/${username}/logs`)
    logs.value = logsResponse.data.logs || []

    if (logs.value.length > 0) {
      const readLogs = logs.value.filter(l => l.status === 'read' && l.rating >= 4)
      favoriteBooks.value = readLogs
        .sort((a, b) => (b.rating || 0) - (a.rating || 0))
        .slice(0, 4)
        .map(log => log.book)
    } else {
      favoriteBooks.value = []
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to load profile'
    console.error('Error loading profile:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProfile()
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
  router.push(`/books/${book.id}`)
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

const profileUrl = computed(() => {
  if (typeof window !== 'undefined') {
    return `${window.location.origin}/profile/${user.value?.username}`
  }
  return ''
})

const copyProfileUrl = async () => {
  try {
    await navigator.clipboard.writeText(profileUrl.value)
    toastStore.success('Profile URL copied to clipboard!')
  } catch (error) {
    console.error('Failed to copy URL:', error)
    toastStore.error('Failed to copy URL')
  }
}

const sendInvite = async () => {
  if (!inviteEmail.value) return

  inviteLoading.value = true
  try {
    await axios.post('/api/invites', {
      email: inviteEmail.value,
      message: `Join me on Folio! Check out my reading profile: ${profileUrl.value}`
    })
    toastStore.success(`Invite sent to ${inviteEmail.value}!`)
    inviteEmail.value = ''
    showInviteModal.value = false
  } catch (error) {
    console.error('Error sending invite:', error)
    toastStore.error('Failed to send invite')
  } finally {
    inviteLoading.value = false
  }
}
</script>
