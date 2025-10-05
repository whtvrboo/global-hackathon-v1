<template>
  <div class="min-h-screen bg-dark-950">
    <!-- Main Content Area -->
    <div class="lg:pl-20 flex flex-col flex-1">
      <!-- Main Content - Unified TikTok/Twitter-like Feed -->
      <main class="flex-1">
        <!-- Hero Search Section -->
        <section class="container-mobile max-w-4xl mx-auto pt-8 pb-6">
          <div class="text-center mb-8">
            <h1 class="text-display-1 mb-4 text-gradient">
              {{ authStore.isAuthenticated ? 'Your Reading Feed' : 'Discover Your Next Read' }}
            </h1>
            <p class="text-body-large text-dark-300 mb-8 max-w-2xl mx-auto">
              {{ authStore.isAuthenticated ?
                'Discover books and see what your friends are reading' :
                'Personalized recommendations just for you' }}
            </p>

            <!-- Guest CTA Section -->
            <div v-if="!authStore.isAuthenticated" class="mb-8">
              <div
                class="bg-gradient-to-r from-accent-red/10 to-accent-blue/10 rounded-2xl p-6 border border-dark-800/50 mb-6">
                <h2 class="text-heading-2 mb-3">Join thousands of readers discovering their next favorite book</h2>
                <p class="text-body text-dark-300 mb-6">
                  Create reading lists, track your progress, and connect with fellow book lovers.
                  Start your reading journey today!
                </p>
                <div class="flex flex-col sm:flex-row gap-3 justify-center">
                  <button @click="$router.push('/login')" class="btn-primary px-8 py-3 text-lg font-semibold">
                    Sign Up Free
                  </button>
                  <button @click="continueAsGuest" class="btn-secondary px-8 py-3 text-lg font-semibold">
                    Continue as Guest
                  </button>
                </div>
                <p class="text-caption text-dark-400 mt-4">
                  No credit card required • Join 10,000+ readers
                </p>
              </div>
            </div>

            <!-- Search Bar -->
            <div class="max-w-2xl mx-auto">
              <SearchBar @select="handleBookSelect" />
            </div>
          </div>
        </section>

        <!-- Unified Feed - TikTok/Twitter Style -->
        <section class="container-mobile max-w-4xl mx-auto pb-24">
          <!-- Feed Header -->
          <div class="flex items-center justify-between mb-6">
            <div>
              <h2 class="text-heading-2">
                {{ authStore.isAuthenticated ? 'For You' : 'Discover Books' }}
              </h2>
              <p v-if="!authStore.isAuthenticated" class="text-sm text-dark-400 mt-1">
                Explore trending books and lists from our community
              </p>
            </div>
            <button @click="refreshFeed" :disabled="feedLoading"
              class="text-accent-blue hover:text-accent-blue/80 text-sm font-medium disabled:opacity-50">
              {{ feedLoading ? 'Loading...' : 'Refresh' }}
            </button>
          </div>

          <!-- Loading State -->
          <div v-if="feedLoading && unifiedFeed.length === 0" class="space-y-6">
            <div v-for="i in 3" :key="i" class="animate-pulse">
              <div class="card">
                <div class="flex gap-4">
                  <div class="w-16 h-24 bg-dark-800 rounded-xl"></div>
                  <div class="flex-1 space-y-3">
                    <div class="h-4 bg-dark-800 rounded w-3/4"></div>
                    <div class="h-3 bg-dark-800 rounded w-1/2"></div>
                    <div class="h-3 bg-dark-800 rounded w-full"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Unified Feed Items -->
          <div v-else class="space-y-1">
            <div v-for="item in unifiedFeed" :key="`${item.type}-${item.id}`"
              class="p-4 hover:bg-dark-800/50 cursor-pointer group border-b border-dark-800/30 transition-colors">

              <!-- Social Activity Item (List Creation/Update) -->
              <div v-if="item.type === 'social'" @click="$router.push(`/lists/${item.id}`)">
                <!-- User Info -->
                <div class="flex items-start gap-3 mb-3">
                  <router-link :to="`/profile/${item.user.username}`"
                    class="relative hover:opacity-80 transition-opacity flex-shrink-0">
                    <img v-if="item.user.picture" :src="item.user.picture" :alt="item.user.name"
                      class="w-10 h-10 rounded-full border border-dark-700" />
                    <div v-else
                      class="w-10 h-10 rounded-full bg-dark-800 border border-dark-700 flex items-center justify-center">
                      <span class="text-sm text-dark-400">👤</span>
                    </div>
                  </router-link>
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-1">
                      <router-link :to="`/profile/${item.user.username}`"
                        class="font-semibold text-white hover:text-accent-blue transition-colors text-sm">
                        {{ item.user.name }}
                      </router-link>
                      <span class="text-dark-500 text-sm">@{{ item.user.username }}</span>
                      <span class="text-dark-500 text-sm">·</span>
                      <span class="text-dark-500 text-sm">{{ timeAgo(item.created_at) }}</span>
                    </div>
                    <p class="text-dark-300 text-sm leading-relaxed">
                      <span class="text-white font-medium">"{{ item.title || item.name }}"</span>
                      <span class="text-dark-400"> - {{ item.description || 'A curated book list' }}</span>
                    </p>
                  </div>
                </div>

                <!-- Enhanced List Preview with Photos -->
                <div class="ml-13 mb-3">
                  <!-- Header Image (if available) -->
                  <div v-if="item.header_image_url" class="mb-3">
                    <img :src="item.header_image_url" :alt="item.title || item.name"
                      class="w-full h-32 object-cover rounded-lg group-hover:opacity-90 transition-opacity" />
                  </div>

                  <!-- Book Covers Grid -->
                  <div v-if="item.preview_books && item.preview_books.length > 0" class="space-y-2">
                    <!-- First row - show up to 4 books -->
                    <div class="flex gap-2">
                      <div v-for="book in item.preview_books.slice(0, 4)" :key="book.id" class="flex-shrink-0">
                        <div class="cursor-pointer" @click.stop="navigateToBook(book.id)">
                          <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                            class="w-12 h-16 object-cover rounded-md shadow-md hover:shadow-lg transition-all duration-200 hover:scale-105" />
                          <div v-else
                            class="w-12 h-16 bg-gradient-to-br from-dark-800 to-dark-700 rounded-md flex items-center justify-center hover:from-dark-700 hover:to-dark-600 transition-all duration-200 hover:scale-105">
                            <span class="text-sm text-dark-400">📖</span>
                          </div>
                        </div>
                      </div>

                      <!-- Show count if more than 4 books -->
                      <div v-if="item.items_count > 4"
                        class="w-12 h-16 bg-gradient-to-br from-accent-blue/20 to-accent-red/20 rounded-md flex items-center justify-center text-white text-xs font-bold border border-dark-600">
                        +{{ item.items_count - 4 }}
                      </div>
                    </div>

                    <!-- Second row - show additional books if more than 4 -->
                    <div v-if="item.preview_books.length > 4" class="flex gap-2">
                      <div v-for="book in item.preview_books.slice(4, 8)" :key="book.id" class="flex-shrink-0">
                        <div class="cursor-pointer" @click.stop="navigateToBook(book.id)">
                          <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                            class="w-10 h-14 object-cover rounded-md shadow-md hover:shadow-lg transition-all duration-200 hover:scale-105" />
                          <div v-else
                            class="w-10 h-14 bg-gradient-to-br from-dark-800 to-dark-700 rounded-md flex items-center justify-center hover:from-dark-700 hover:to-dark-600 transition-all duration-200 hover:scale-105">
                            <span class="text-xs text-dark-400">📖</span>
                          </div>
                        </div>
                      </div>

                      <!-- Show count if more than 8 books -->
                      <div v-if="item.items_count > 8"
                        class="w-10 h-14 bg-gradient-to-br from-accent-green/20 to-accent-blue/20 rounded-md flex items-center justify-center text-white text-xs font-bold border border-dark-600">
                        +{{ item.items_count - 8 }}
                      </div>
                    </div>
                  </div>

                  <!-- Fallback when no preview books but list has items -->
                  <div v-else-if="item.items_count > 0" class="flex items-center gap-2 text-dark-400">
                    <div class="w-8 h-8 bg-dark-800 rounded-full flex items-center justify-center">
                      <span class="text-sm">📚</span>
                    </div>
                    <span class="text-sm">{{ item.items_count }} book{{ item.items_count !== 1 ? 's' : '' }} in this
                      list</span>
                  </div>
                </div>

                <!-- Action Bar -->
                <div class="ml-13 flex items-center gap-6 text-sm text-dark-400">
                  <button @click.stop="toggleLike(item.id)"
                    class="flex items-center gap-1 transition-colors hover:text-accent-red"
                    :class="item.is_liked ? 'text-accent-red' : 'text-dark-400'">
                    <svg class="w-4 h-4" :fill="item.is_liked ? 'currentColor' : 'none'" stroke="currentColor"
                      viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                      </path>
                    </svg>
                    {{ item.likes_count || 0 }}
                  </button>
                  <button class="flex items-center gap-1 transition-colors hover:text-accent-blue">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                      </path>
                    </svg>
                    {{ item.comments_count || 0 }}
                  </button>
                  <span class="text-dark-500">{{ item.items_count }} book{{ item.items_count !== 1 ? 's' : ''
                  }}</span>
                </div>
              </div>

              <!-- Book Recommendation Item -->
              <div v-else-if="item.type === 'book'" @click="handleBookSelect(item.book)">
                <!-- Book Card -->
                <div class="flex gap-4">
                  <div class="relative">
                    <div
                      class="w-20 h-32 bg-dark-800 rounded-xl flex items-center justify-center overflow-hidden cursor-pointer">
                      <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                        class="w-full h-full object-cover" />
                      <div v-else
                        class="w-full h-full bg-gradient-to-br from-accent-red/20 to-accent-blue/20 flex items-center justify-center">
                        <span class="text-2xl">📖</span>
                      </div>
                    </div>
                    <div
                      class="absolute -top-2 -right-2 w-6 h-6 bg-accent-green rounded-full flex items-center justify-center">
                      <span class="text-white text-xs">✓</span>
                    </div>
                  </div>
                  <div class="flex-1">
                    <div class="flex items-start justify-between mb-2">
                      <h3 class="text-heading-3 line-clamp-2">{{ item.book.title }}</h3>
                      <button class="text-dark-400 hover:text-white p-1">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
                        </svg>
                      </button>
                    </div>
                    <p class="text-body text-dark-300 mb-3">by {{ item.book.authors?.[0] || 'Unknown Author' }} • {{
                      item.book.published_date?.split('-')[0] || 'Unknown Year' }}</p>
                    <p class="text-caption text-dark-400 mb-4 line-clamp-2">
                      {{ item.reason?.description || item.book.description || 'No description available.' }}
                    </p>
                    <div class="flex items-center gap-3">
                      <button class="btn-primary text-sm px-4 py-2" @click.stop="handleLogBook(item.book)">
                        Log Book
                      </button>
                      <button class="btn-ghost text-sm px-4 py-2">
                        Preview
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Trending List Item -->
              <div v-else-if="item.type === 'trending-list'" @click="viewList(item)">
                <div class="flex gap-4">
                  <div class="w-16 h-20 bg-dark-800 rounded-xl flex items-center justify-center">
                    <span class="text-2xl">📚</span>
                  </div>
                  <div class="flex-1">
                    <div class="flex items-start justify-between mb-2">
                      <h3 class="text-heading-3 line-clamp-2">{{ item.name }}</h3>
                      <div class="text-right">
                        <div class="text-lg font-bold text-accent-red">{{ item.items_count }}</div>
                        <div class="text-xs text-dark-400">books</div>
                      </div>
                    </div>
                    <p v-if="item.description" class="text-caption text-dark-300 mb-2 line-clamp-2">
                      {{ item.description }}
                    </p>
                    <div class="flex items-center gap-2 mb-3">
                      <router-link :to="`/profile/${item.user.username}`" class="cursor-pointer">
                        <img v-if="item.user.picture" :src="item.user.picture" :alt="item.user.name"
                          class="w-6 h-6 rounded-full" />
                        <div v-else class="w-6 h-6 rounded-full bg-dark-700 flex items-center justify-center">
                          <span class="text-xs text-dark-400">U</span>
                        </div>
                      </router-link>
                      <span class="text-sm text-dark-300">by {{ item.user.name }}</span>
                    </div>
                    <div class="flex gap-2">
                      <button class="flex-1 btn-primary text-sm py-2" @click.stop="viewList(item)">
                        View List
                      </button>
                      <button class="flex-1 btn-secondary text-sm py-2" @click.stop="followList(item)">
                        Follow
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="!feedLoading && unifiedFeed.length === 0" class="card text-center py-16">
            <div class="text-6xl mb-6">📚</div>
            <h3 class="text-heading-2 mb-4">
              {{ authStore.isAuthenticated ? 'No content yet' : 'Ready to discover amazing books?' }}
            </h3>
            <p class="text-body text-dark-300 mb-8">
              {{ authStore.isAuthenticated ?
                'Follow some users or create lists to see activity in your feed!' :
                'Join our community of book lovers and get personalized recommendations!' }}
            </p>

            <!-- Authenticated user actions -->
            <div v-if="authStore.isAuthenticated">
              <button @click="$router.push('/profile')" class="btn-primary">
                Create Your First List
              </button>
            </div>

            <!-- Guest user actions -->
            <div v-else class="space-y-4">
              <div class="flex flex-col sm:flex-row gap-3 justify-center">
                <button @click="$router.push('/login')" class="btn-primary px-8 py-3">
                  Sign Up Free
                </button>
                <button @click="continueAsGuest" class="btn-secondary px-8 py-3">
                  Browse as Guest
                </button>
              </div>
              <div class="text-sm text-dark-400 mt-4">
                <p>✨ Get personalized recommendations</p>
                <p>📚 Create and share reading lists</p>
                <p>👥 Connect with fellow readers</p>
              </div>
            </div>
          </div>
        </section>
      </main>
    </div>

    <!-- Book Detail Modal -->
    <BookDetailModal :show="showBookDetail" :book-id="selectedBookId" @close="showBookDetail = false"
      @log="handleLogBook" />

    <!-- Log Book Modal -->
    <LogBookModal :show="showLogModal" :book="bookToLog" @close="showLogModal = false" @success="handleLogSuccess" />

    <!-- Guest Conversion Modal -->
    <GuestConversionModal :show="showConversionModal" @close="showConversionModal = false" />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import axios from 'axios'
import SearchBar from '../components/SearchBar.vue'
import BookDetailModal from '../components/BookDetailModal.vue'
import LogBookModal from '../components/LogBookModal.vue'
import GuestConversionModal from '../components/GuestConversionModal.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const toastStore = useToastStore()

// Modal states
const selectedBookId = ref(null)
const showBookDetail = ref(false)
const showLogModal = ref(false)
const showConversionModal = ref(false)
const bookToLog = ref(null)

// Data sources
const feed = ref([])
const feedLoading = ref(false)
const trendingBooks = ref([])
const discoveryBooks = ref([])
const personalizedRecommendations = ref([])
const trendingLists = ref([])
const popularUsers = ref([])
const followingUsers = ref({})

// Unified feed computed property that mixes all content types
const unifiedFeed = computed(() => {
  const items = []

  // Add social feed items
  feed.value.forEach(item => {
    items.push({
      ...item,
      type: 'social',
      id: item.id,
      created_at: item.created_at
    })
  })

  // Add book recommendations (mix discovery and personalized)
  const allBooks = [...discoveryBooks.value, ...personalizedRecommendations.value]
  allBooks.forEach(book => {
    items.push({
      ...book,
      type: 'book',
      id: `book-${book.book.id}`,
      created_at: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toISOString() // Random date within last week
    })
  })

  // Add trending lists
  trendingLists.value.forEach(list => {
    items.push({
      ...list,
      type: 'trending-list',
      id: `${list.id}`,
      created_at: new Date(Date.now() - Math.random() * 3 * 24 * 60 * 60 * 1000).toISOString() // Random date within last 3 days
    })
  })

  // Sort by creation date (most recent first)
  return items.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
})

// Event handlers
const handleBookSelect = (book) => {
  router.push(`/books/${book.id}`)
}

const handleLogBook = (book) => {
  bookToLog.value = book
  showBookDetail.value = false
  showLogModal.value = true
}

const handleLogSuccess = () => {
  console.log('Book logged successfully!')
}

const navigateToBook = (bookId) => {
  router.push(`/books/${bookId}`)
}

const timeAgo = (date) => {
  const seconds = Math.floor((new Date() - new Date(date)) / 1000)

  const intervals = {
    year: 31536000,
    month: 2592000,
    week: 604800,
    day: 86400,
    hour: 3600,
    minute: 60
  }

  for (const [unit, secondsInUnit] of Object.entries(intervals)) {
    const interval = Math.floor(seconds / secondsInUnit)
    if (interval >= 1) {
      return `${interval} ${unit}${interval > 1 ? 's' : ''} ago`
    }
  }

  return 'just now'
}

const toggleLike = async (listId) => {
  try {
    const item = feed.value.find(i => i.id === listId)
    if (!item) return

    if (item.is_liked) {
      await axios.delete(`/api/lists/${listId}/like`)
      item.is_liked = false
      item.likes_count = Math.max((item.likes_count || 0) - 1, 0)
      toastStore.info('Unliked')
    } else {
      await axios.post(`/api/lists/${listId}/like`)
      item.is_liked = true
      item.likes_count = (item.likes_count || 0) + 1
      toastStore.success('Liked!')
    }
  } catch (error) {
    console.error('Error toggling like:', error)
    toastStore.error('Failed to update like')
  }
}

const followUser = async (username) => {
  followingUsers.value[username] = true
  try {
    await axios.post(`/api/users/${username}/follow`)

    const user = popularUsers.value.find(u => u.username === username)
    if (user) {
      user.is_following = true
    }

    toastStore.success(`You're now following ${username}!`)

    setTimeout(async () => {
      await loadFeed()
    }, 500)
  } catch (error) {
    console.error('Error following user:', error)
    toastStore.error('Failed to follow user')
  } finally {
    followingUsers.value[username] = false
  }
}

const viewList = (list) => {
  router.push(`/lists/${list.id}`)
}

const followList = (list) => {
  console.log('Following list:', list.name)
  toastStore.info(`Following list: ${list.name}`)
}

const continueAsGuest = () => {
  // Show guest conversion modal or just continue browsing
  showConversionModal.value = true
}

// Data loading functions
const loadFeed = async () => {
  if (!authStore.isAuthenticated) {
    console.log('User not authenticated, skipping feed load')
    return
  }

  try {
    console.log('Loading feed for user:', authStore.user?.username)
    const response = await axios.get('/api/feed')
    console.log('Feed response:', response.data)
    feed.value = response.data.feed || []

    if (feed.value.length === 0) {
      console.log('No feed items found - user may need to follow others or create lists')
    }
  } catch (error) {
    console.error('Error loading feed:', error)
    if (error.response?.status === 401) {
      console.log('Unauthorized - user may need to re-authenticate')
    } else if (error.response?.status === 500) {
      console.log('Server error loading feed')
    }
  }
}

const loadTrendingBooks = async () => {
  try {
    const response = await axios.get('/api/discover', {
      params: { limit: 3 }
    })
    trendingBooks.value = response.data.recommendations || []
  } catch (error) {
    console.error('Error loading trending books:', error)
    trendingBooks.value = []
  }
}

const loadDiscoveryBooks = async () => {
  try {
    const response = await axios.get('/api/discover', {
      params: { limit: 4 }
    })
    discoveryBooks.value = response.data.recommendations || []
  } catch (error) {
    console.error('Error loading discovery books:', error)
    discoveryBooks.value = []
  }
}

const loadPersonalizedRecommendations = async () => {
  try {
    const response = await axios.get('/api/discover', {
      params: { limit: 3 }
    })
    personalizedRecommendations.value = response.data.recommendations || []
  } catch (error) {
    console.error('Error loading personalized recommendations:', error)
    personalizedRecommendations.value = []
  }
}

const loadTrendingLists = async () => {
  try {
    const response = await axios.get('/api/discover/lists', {
      params: { limit: 5 }
    })
    trendingLists.value = response.data.lists || []
  } catch (error) {
    console.error('Error loading trending lists:', error)
    trendingLists.value = []
  }
}

const loadPopularUsers = async () => {
  try {
    const response = await axios.get('/api/users/popular')
    popularUsers.value = response.data.users || []
  } catch (error) {
    console.error('Error loading popular users:', error)
    popularUsers.value = []
  }
}

const refreshFeed = async () => {
  feedLoading.value = true
  try {
    await Promise.all([
      loadFeed(),
      loadTrendingBooks(),
      loadDiscoveryBooks(),
      loadPersonalizedRecommendations(),
      loadTrendingLists(),
      loadPopularUsers()
    ])
  } finally {
    feedLoading.value = false
  }
}

// Initialize data on mount
onMounted(async () => {
  await refreshFeed()
})
</script>