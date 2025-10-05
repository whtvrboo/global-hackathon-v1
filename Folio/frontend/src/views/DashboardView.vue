<template>
  <div class="min-h-screen bg-dark-950">

    <!-- Main Content - TikTok-inspired content-first layout -->
    <main class="min-h-screen">
      <!-- Hero Search Section -->
      <section class="container-mobile max-w-7xl mx-auto pt-8 pb-6">
        <div class="text-center mb-8">
          <h1 class="text-display-1 mb-4 text-gradient">
            Discover Your Next Read
          </h1>
          <p class="text-body-large text-dark-300 mb-8 max-w-2xl mx-auto">
            Personalized recommendations just for you
          </p>

          <!-- Search Bar -->
          <div class="max-w-2xl mx-auto">
            <SearchBar @select="handleBookSelect" />
          </div>
        </div>
      </section>

      <!-- Public Browse Section for Unauthenticated Users -->
      <section v-if="!authStore.user" class="container-mobile max-w-7xl mx-auto mb-8">
        <div class="text-center mb-8">
          <h2 class="text-heading-2 mb-4">Browse Popular Books</h2>
          <p class="text-body text-dark-300">Discover what others are reading</p>
        </div>

        <!-- Public Book Grid -->
        <div v-if="publicBooks.length > 0" class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <div v-for="book in publicBooks.slice(0, 12)" :key="book.id" class="card card-hover group cursor-pointer"
            @click="handleBookSelect(book)">
            <div class="aspect-[2/3] bg-dark-800 rounded-xl mb-3 relative overflow-hidden">
              <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title" class="w-full h-full object-cover" />
              <div v-else
                class="w-full h-full bg-gradient-to-br from-accent-red/20 to-accent-blue/20 flex items-center justify-center">
              </div>
              <div v-if="book.rating" class="absolute top-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
                {{ book.rating.toFixed(1) }}
              </div>
            </div>
            <h3 class="font-semibold text-white text-sm line-clamp-2 mb-1">{{ book.title }}</h3>
            <p v-if="book.authors" class="text-caption text-dark-300 line-clamp-1">
              by {{ book.authors[0] }}
            </p>
          </div>
        </div>

        <!-- Public Browse CTA -->
        <div class="text-center mt-8">
          <div class="card inline-block">
            <h3 class="text-heading-3 mb-4">Ready to start your reading journey?</h3>
            <p class="text-body text-dark-300 mb-6">Create a free account to log books, track your reading, and discover
              personalized recommendations.</p>
            <div class="flex flex-col sm:flex-row gap-4 justify-center">
              <button @click="$router.push('/login')" class="btn-primary">
                Get Started Free
              </button>
              <button @click="startGuestSession" class="btn-secondary">
                Try as Guest
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- Content Feed - TikTok-style vertical scroll -->
      <div class="space-y-6 pb-24">
        <!-- Trending Books Section -->
        <section class="container-mobile max-w-7xl mx-auto">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-heading-2">Trending Now</h2>
            <button class="text-accent-red hover:text-accent-red/80 text-sm font-medium">
              See All
            </button>
          </div>

          <!-- Horizontal Scroll Cards -->
          <div class="flex gap-4 overflow-x-auto pb-4 scrollbar-hide">
            <div v-for="(book, index) in trendingBooks" :key="book.book.id" class="flex-shrink-0 w-48">
              <div class="card card-hover group cursor-pointer" @click="handleBookSelect(book.book)">
                <div class="aspect-[2/3] bg-dark-800 rounded-xl mb-4 relative overflow-hidden">
                  <img v-if="book.book.cover_url" :src="book.book.cover_url" :alt="book.book.title"
                    class="w-full h-full object-cover" />
                  <div v-else
                    class="w-full h-full bg-gradient-to-br from-accent-red/20 to-accent-blue/20 flex items-center justify-center">
                  </div>
                  <div
                    class="absolute top-3 right-3 w-8 h-8 bg-accent-red rounded-full flex items-center justify-center">
                    <span class="text-white text-xs font-bold">{{ index + 1 }}</span>
                  </div>
                  <div v-if="book.book.rating"
                    class="absolute top-2 left-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
                    {{ book.book.rating.toFixed(1) }}
                  </div>
                </div>
                <h3 class="font-semibold text-white mb-1 line-clamp-2">{{ book.book.title }}</h3>
                <p class="text-caption mb-3">by {{ book.book.authors?.[0] || 'Unknown Author' }}</p>
                <div class="flex items-center gap-2">
                  <div class="flex">
                    <span v-for="star in Math.floor(book.book.rating || 0)" :key="star"
                      class="text-accent-orange text-sm">★</span>
                  </div>
                  <span class="text-caption">{{ book.book.rating?.toFixed(1) || 'N/A' }}</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- TikTok-style Discovery Feed -->
        <section class="container-mobile max-w-7xl mx-auto">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-heading-2">Discover Books</h2>
            <div class="flex gap-2">
              <button class="text-accent-blue hover:text-accent-blue/80 text-sm font-medium">
                Refresh
              </button>
              <span class="text-dark-500">•</span>
              <button class="text-dark-400 hover:text-white text-sm font-medium">
                Filters
              </button>
            </div>
          </div>

          <!-- Swipeable Discovery Cards -->
          <div class="space-y-6">
            <div v-for="book in discoveryBooks" :key="book.book.id" class="flex justify-center">
              <div class="card card-hover group max-w-sm w-full cursor-pointer" @click="handleBookSelect(book.book)">
                <div class="aspect-[2/3] bg-dark-800 rounded-xl mb-4 relative overflow-hidden">
                  <img v-if="book.book.cover_url" :src="book.book.cover_url" :alt="book.book.title"
                    class="w-full h-full object-cover" />
                  <div v-else
                    class="w-full h-full bg-gradient-to-br from-accent-red/20 to-accent-blue/20 flex items-center justify-center">
                  </div>
                  <div
                    class="absolute top-3 right-3 w-8 h-8 bg-accent-green rounded-full flex items-center justify-center">
                    <span class="text-white text-xs font-bold">✓</span>
                  </div>
                  <div class="absolute bottom-3 left-3 right-3">
                    <div class="glass-strong rounded-lg p-3">
                      <div class="flex items-center gap-2 mb-1">
                        <div class="flex">
                          <span v-for="star in Math.floor(book.book.rating || 0)" :key="star"
                            class="text-accent-orange text-sm">★</span>
                        </div>
                        <span class="text-white text-sm font-semibold">{{ book.book.rating?.toFixed(1) || 'N/A'
                        }}</span>
                      </div>
                      <p class="text-white text-xs">{{ book.reason.description }}</p>
                    </div>
                  </div>
                </div>

                <div class="p-4">
                  <h3 class="font-semibold text-white line-clamp-2 mb-2 text-lg">{{ book.book.title }}</h3>
                  <p class="text-caption mb-3">by {{ book.book.authors?.[0] || 'Unknown Author' }} • {{
                    book.book.published_date?.split('-')[0] || 'Unknown Year' }}</p>
                  <p class="text-caption text-dark-300 mb-4 line-clamp-2">
                    {{ book.book.description || 'No description available.' }}
                  </p>

                  <!-- Engagement Stats -->
                  <div class="flex items-center justify-between mb-4 text-sm">
                    <div class="flex items-center gap-4">
                      <span class="text-dark-400">{{ book.book.log_count || 0 }} readers</span>
                      <span class="text-dark-400">{{ book.book.page_count || 'N/A' }} pages</span>
                    </div>
                    <span class="text-accent-green text-xs font-medium">{{ book.reason.type }}</span>
                  </div>

                  <!-- Action Buttons -->
                  <div class="flex gap-2">
                    <button class="flex-1 btn-primary text-sm py-2" @click.stop="handleLogBook(book.book)">
                      Log Book
                    </button>
                    <button class="flex-1 btn-secondary text-sm py-2">
                      Like
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- Trending Lists Section -->
        <section class="container-mobile max-w-7xl mx-auto">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-heading-2">Trending Lists</h2>
            <button @click="loadTrendingLists" class="text-accent-blue hover:text-accent-blue/80 text-sm font-medium">
              Refresh
            </button>
          </div>

          <!-- Horizontal Scroll Lists -->
          <div v-if="trendingLists.length > 0" class="flex gap-4 overflow-x-auto pb-4 scrollbar-hide">
            <div v-for="list in trendingLists" :key="list.id" class="flex-shrink-0 w-64">
              <div class="card card-hover group cursor-pointer" @click="viewList(list)">
                <div class="p-4">
                  <div class="flex items-start justify-between mb-3">
                    <div class="flex-1">
                      <h3 class="font-semibold text-white mb-1 line-clamp-2">{{ list.name }}</h3>
                      <p v-if="list.description" class="text-caption text-dark-300 mb-2 line-clamp-2">
                        {{ list.description }}
                      </p>
                    </div>
                    <div class="text-right">
                      <div class="text-lg font-bold text-accent-red">{{ list.items_count }}</div>
                      <div class="text-xs text-dark-400">books</div>
                    </div>
                  </div>

                  <!-- List Creator -->
                  <div class="flex items-center gap-2 mb-3">
                    <img v-if="list.user.picture" :src="list.user.picture" :alt="list.user.name"
                      class="w-6 h-6 rounded-full" />
                    <div v-else class="w-6 h-6 rounded-full bg-dark-700 flex items-center justify-center">
                      <span class="text-xs text-dark-400">U</span>
                    </div>
                    <span class="text-sm text-dark-300">by {{ list.user.name }}</span>
                  </div>

                  <!-- Action Buttons -->
                  <div class="flex gap-2">
                    <button class="flex-1 btn-primary text-sm py-2" @click.stop="viewList(list)">
                      View List
                    </button>
                    <button class="flex-1 btn-secondary text-sm py-2" @click.stop="followList(list)">
                      Follow
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty State for Lists -->
          <div v-else class="card text-center py-8">
            <div class="text-4xl mb-4"></div>
            <h3 class="text-heading-3 mb-2">No trending lists yet</h3>
            <p class="text-caption text-dark-300">Be the first to create a public list!</p>
          </div>
        </section>

        <!-- Personalized Recommendations -->
        <section class="container-mobile max-w-7xl mx-auto">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-heading-2">For You</h2>
            <button class="text-accent-blue hover:text-accent-blue/80 text-sm font-medium">
              Refresh
            </button>
          </div>

          <!-- Vertical Feed Cards -->
          <div class="space-y-4">
            <div v-for="book in personalizedRecommendations" :key="book.book.id" class="card card-hover cursor-pointer"
              @click="handleBookSelect(book.book)">
              <div class="flex gap-4">
                <div class="relative">
                  <div class="w-20 h-32 bg-dark-800 rounded-xl flex items-center justify-center overflow-hidden">
                    <img v-if="book.book.cover_url" :src="book.book.cover_url" :alt="book.book.title"
                      class="w-full h-full object-cover" />
                  </div>
                  <div
                    class="absolute -top-2 -right-2 w-6 h-6 bg-accent-green rounded-full flex items-center justify-center">
                    <span class="text-white text-xs">✓</span>
                  </div>
                </div>
                <div class="flex-1">
                  <div class="flex items-start justify-between mb-2">
                    <h3 class="text-heading-3 line-clamp-2">{{ book.book.title }}</h3>
                    <button class="text-dark-400 hover:text-white p-1">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
                      </svg>
                    </button>
                  </div>
                  <p class="text-body text-dark-300 mb-3">by {{ book.book.authors?.[0] || 'Unknown Author' }} • {{
                    book.book.published_date?.split('-')[0] || 'Unknown Year' }}</p>
                  <p class="text-caption text-dark-400 mb-4 line-clamp-2">
                    {{ book.reason.description }}
                  </p>
                  <div class="flex items-center gap-3">
                    <button class="btn-primary text-sm px-4 py-2" @click.stop="handleLogBook(book.book)">
                      Log Book
                    </button>
                    <button class="btn-ghost text-sm px-4 py-2">
                      Preview
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- Quick Actions Grid -->
        <section class="container-mobile max-w-7xl mx-auto">
          <h2 class="text-heading-2 mb-6">Quick Actions</h2>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <router-link to="/discover" class="card card-hover group text-center p-6">
              <div
                class="w-12 h-12 bg-accent-blue/20 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:bg-accent-blue/30 transition-colors">
                <svg class="w-6 h-6 text-accent-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                </svg>
              </div>
              <h4 class="font-semibold text-white mb-1">Discover</h4>
              <p class="text-caption">Find books</p>
            </router-link>

            <router-link to="/feed" class="card card-hover group text-center p-6">
              <div
                class="w-12 h-12 bg-accent-purple/20 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:bg-accent-purple/30 transition-colors">
                <svg class="w-6 h-6 text-accent-purple" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 6h16M4 10h16M4 14h16M4 18h16"></path>
                </svg>
              </div>
              <h4 class="font-semibold text-white mb-1">Feed</h4>
              <p class="text-caption">See activity</p>
            </router-link>

            <router-link to="/profile" class="card card-hover group text-center p-6">
              <div
                class="w-12 h-12 bg-accent-red/20 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:bg-accent-red/30 transition-colors">
                <svg class="w-6 h-6 text-accent-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                </svg>
              </div>
              <h4 class="font-semibold text-white mb-1">Profile</h4>
              <p class="text-caption">Your books</p>
            </router-link>

            <button class="card card-hover group text-center p-6">
              <div
                class="w-12 h-12 bg-accent-orange/20 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:bg-accent-orange/30 transition-colors">
                <svg class="w-6 h-6 text-accent-orange" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6">
                  </path>
                </svg>
              </div>
              <h4 class="font-semibold text-white mb-1">Log Book</h4>
              <p class="text-caption">Add reading</p>
            </button>
          </div>
        </section>

        <!-- Welcome Message for Guest Users -->
        <section v-if="authStore.user && authStore.isGuestUser" class="container-mobile max-w-7xl mx-auto">
          <div class="card card-hover">
            <div class="flex items-center gap-4 mb-4">
              <div class="relative">
                <img v-if="authStore.user.picture" :src="authStore.user.picture" :alt="authStore.user.name"
                  class="w-12 h-12 rounded-full border-2 border-dark-700" />
                <div v-else
                  class="w-12 h-12 rounded-full bg-dark-800 border-2 border-dark-700 flex items-center justify-center">
                  <span class="text-dark-400">U</span>
                </div>
                <div
                  class="absolute -top-1 -right-1 w-5 h-5 bg-accent-orange rounded-full flex items-center justify-center">
                  <span class="text-white text-xs font-bold">G</span>
                </div>
              </div>
              <div>
                <h3 class="font-semibold text-white">Welcome, {{ authStore.user.name }}!</h3>
                <p class="text-caption">Guest account • Try Folio risk-free</p>
              </div>
            </div>
            <p class="text-body text-dark-300 mb-4">
              Your reading progress is saved locally. Create a full account to sync across devices and never lose your
              data.
            </p>
            <button @click="showConversionModal = true" class="btn-primary">
              Create Full Account
            </button>
          </div>
        </section>
      </div>
    </main>


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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'
import SearchBar from '../components/SearchBar.vue'
import Card from '../components/ui/Card.vue'
import BookDetailModal from '../components/BookDetailModal.vue'
import LogBookModal from '../components/LogBookModal.vue'
import GuestConversionModal from '../components/GuestConversionModal.vue'

const router = useRouter()
const authStore = useAuthStore()

const selectedBookId = ref(null)
const showBookDetail = ref(false)
const showLogModal = ref(false)
const showConversionModal = ref(false)
const bookToLog = ref(null)
const trendingLists = ref([])
const publicBooks = ref([])
const trendingBooks = ref([])
const discoveryBooks = ref([])
const personalizedRecommendations = ref([])

const handleBookSelect = (book) => {
  router.push(`/books/${book.id}`)
}

const handleLogBook = (book) => {
  bookToLog.value = book
  showBookDetail.value = false
  showLogModal.value = true
}

const handleLogSuccess = () => {
  // Could show a success toast here
  console.log('Book logged successfully!')
}

// Load trending lists on component mount
onMounted(async () => {
  await loadTrendingLists()
  await loadTrendingBooks()
  await loadDiscoveryBooks()
  await loadPersonalizedRecommendations()
  // Load public books for unauthenticated users
  if (!authStore.user) {
    await loadPublicBooks()
  }
})

const loadTrendingLists = async () => {
  try {
    const response = await axios.get('/api/discover/lists', {
      params: { limit: 10 }
    })
    trendingLists.value = response.data.lists || []
  } catch (error) {
    console.error('Error loading trending lists:', error)
    trendingLists.value = []
  }
}

const viewList = (list) => {
  // Navigate to list detail view (could be a modal or new page)
  console.log('Viewing list:', list.name)
  // For now, just show an alert
  alert(`Viewing list: ${list.name}`)
}

const followList = (list) => {
  // Follow list functionality
  console.log('Following list:', list.name)
  // For now, just show an alert
  alert(`Following list: ${list.name}`)
}

const loadPublicBooks = async () => {
  try {
    const response = await axios.get('/api/discover', {
      params: { limit: 12 }
    })
    publicBooks.value = response.data.recommendations || []
  } catch (error) {
    console.error('Error loading public books:', error)
    publicBooks.value = []
  }
}

const loadTrendingBooks = async () => {
  try {
    const response = await axios.get('/api/discover', {
      params: { limit: 5 }
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
      params: { limit: 3 }
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
      params: { limit: 2 }
    })
    personalizedRecommendations.value = response.data.recommendations || []
  } catch (error) {
    console.error('Error loading personalized recommendations:', error)
    personalizedRecommendations.value = []
  }
}

const startGuestSession = () => {
  // Redirect to login with guest option
  // This would typically trigger the guest login flow
  console.log('Starting guest session')
  // For now, redirect to login
  window.location.href = '/login'
}
</script>
