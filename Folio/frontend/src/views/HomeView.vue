<template>
  <div class="min-h-screen bg-dark-950">
    <!-- Header -->
    <header class="glass-strong border-b border-dark-800 sticky top-0 z-50">
      <div class="container-mobile max-w-7xl mx-auto">
        <div class="flex items-center justify-between py-4">
          <!-- Logo -->
          <div class="flex items-center gap-3">
            <div
              class="w-8 h-8 bg-gradient-to-br from-accent-red to-accent-blue rounded-lg flex items-center justify-center">
              <span class="text-white font-bold text-sm">F</span>
            </div>
            <h1 class="text-xl font-bold text-white">Folio</h1>
          </div>

          <!-- Desktop Navigation -->
          <nav class="hidden md:flex items-center gap-6">
            <router-link to="/discover" class="btn-ghost">
              Discover
            </router-link>
            <router-link to="/feed" class="btn-ghost">
              Feed
            </router-link>
            <router-link to="/profile" class="btn-ghost">
              Profile
            </router-link>
            <button @click="authStore.logout()" class="btn-ghost">
              Logout
            </button>
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

      <!-- Content Feed - TikTok-style vertical scroll -->
      <div class="space-y-6 pb-24">
        <!-- Trending Books Section -->
        <section class="container-mobile max-w-7xl mx-auto">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-heading-2">🔥 Trending Now</h2>
            <button class="text-accent-red hover:text-accent-red/80 text-sm font-medium">
              See All
            </button>
          </div>

          <!-- Horizontal Scroll Cards -->
          <div class="flex gap-4 overflow-x-auto pb-4 scrollbar-hide">
            <div v-for="i in 5" :key="i" class="flex-shrink-0 w-48">
              <div class="card card-hover group">
                <div class="aspect-[2/3] bg-dark-800 rounded-xl mb-4 relative overflow-hidden">
                  <div
                    class="w-full h-full bg-gradient-to-br from-accent-red/20 to-accent-blue/20 flex items-center justify-center">
                    <span class="text-4xl">📚</span>
                  </div>
                  <div
                    class="absolute top-3 right-3 w-8 h-8 bg-accent-red rounded-full flex items-center justify-center">
                    <span class="text-white text-xs font-bold">{{ i }}</span>
                  </div>
                </div>
                <h3 class="font-semibold text-white mb-1 line-clamp-2">Sample Book Title {{ i }}</h3>
                <p class="text-caption mb-3">by Author Name</p>
                <div class="flex items-center gap-2">
                  <div class="flex">
                    <span v-for="star in 5" :key="star" class="text-accent-orange text-sm">★</span>
                  </div>
                  <span class="text-caption">4.{{ i }}</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- TikTok-style Discovery Feed -->
        <section class="container-mobile max-w-7xl mx-auto">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-heading-2">✨ Discover Books</h2>
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
            <div v-for="i in 3" :key="i" class="flex justify-center">
              <div class="card card-hover group max-w-sm w-full">
                <div class="aspect-[2/3] bg-dark-800 rounded-xl mb-4 relative overflow-hidden">
                  <div
                    class="w-full h-full bg-gradient-to-br from-accent-red/20 to-accent-blue/20 flex items-center justify-center">
                    <span class="text-4xl">📚</span>
                  </div>
                  <div
                    class="absolute top-3 right-3 w-8 h-8 bg-accent-green rounded-full flex items-center justify-center">
                    <span class="text-white text-xs font-bold">✓</span>
                  </div>
                  <div class="absolute bottom-3 left-3 right-3">
                    <div class="glass-strong rounded-lg p-3">
                      <div class="flex items-center gap-2 mb-1">
                        <div class="flex">
                          <span v-for="star in 5" :key="star" class="text-accent-orange text-sm">★</span>
                        </div>
                        <span class="text-white text-sm font-semibold">4.{{ i }}</span>
                      </div>
                      <p class="text-white text-xs">Based on your preferences</p>
                    </div>
                  </div>
                </div>

                <div class="p-4">
                  <h3 class="font-semibold text-white line-clamp-2 mb-2 text-lg">Sample Book Title {{ i }}</h3>
                  <p class="text-caption mb-3">by Author Name • 2024</p>
                  <p class="text-caption text-dark-300 mb-4 line-clamp-2">
                    A compelling story that will keep you turning pages late into the night...
                  </p>

                  <!-- Engagement Stats -->
                  <div class="flex items-center justify-between mb-4 text-sm">
                    <div class="flex items-center gap-4">
                      <span class="text-dark-400">👥 1.2k readers</span>
                      <span class="text-dark-400">💬 89 reviews</span>
                    </div>
                    <span class="text-accent-green text-xs font-medium">Trending</span>
                  </div>

                  <!-- Action Buttons -->
                  <div class="flex gap-2">
                    <button class="flex-1 btn-primary text-sm py-2">
                      ❤️ Like
                    </button>
                    <button class="flex-1 btn-secondary text-sm py-2">
                      📚 Add to List
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- Personalized Recommendations -->
        <section class="container-mobile max-w-7xl mx-auto">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-heading-2">🎯 For You</h2>
            <button class="text-accent-blue hover:text-accent-blue/80 text-sm font-medium">
              Refresh
            </button>
          </div>

          <!-- Vertical Feed Cards -->
          <div class="space-y-4">
            <div v-for="i in 2" :key="i" class="card card-hover">
              <div class="flex gap-4">
                <div class="relative">
                  <div class="w-20 h-32 bg-dark-800 rounded-xl flex items-center justify-center">
                    <span class="text-2xl">📖</span>
                  </div>
                  <div
                    class="absolute -top-2 -right-2 w-6 h-6 bg-accent-green rounded-full flex items-center justify-center">
                    <span class="text-white text-xs">✓</span>
                  </div>
                </div>
                <div class="flex-1">
                  <div class="flex items-start justify-between mb-2">
                    <h3 class="text-heading-3 line-clamp-2">Recommended Book Title {{ i }}</h3>
                    <button class="text-dark-400 hover:text-white p-1">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
                      </svg>
                    </button>
                  </div>
                  <p class="text-body text-dark-300 mb-3">by Author Name • 2024</p>
                  <p class="text-caption text-dark-400 mb-4 line-clamp-2">
                    Based on your reading history and preferences, we think you'll love this book...
                  </p>
                  <div class="flex items-center gap-3">
                    <button class="btn-primary text-sm px-4 py-2">
                      Add to List
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
          <h2 class="text-heading-2 mb-6">⚡ Quick Actions</h2>
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
                  <span class="text-dark-400">👤</span>
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

        <router-link to="/profile"
          class="flex flex-col items-center gap-1 p-3 text-dark-400 hover:text-white transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
          </svg>
          <span class="text-xs">Profile</span>
        </router-link>
      </div>
    </nav>

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
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import SearchBar from '../components/SearchBar.vue'
import Card from '../components/ui/Card.vue'
import BookDetailModal from '../components/BookDetailModal.vue'
import LogBookModal from '../components/LogBookModal.vue'
import GuestConversionModal from '../components/GuestConversionModal.vue'

const authStore = useAuthStore()

const selectedBookId = ref(null)
const showBookDetail = ref(false)
const showLogModal = ref(false)
const showConversionModal = ref(false)
const bookToLog = ref(null)

const handleBookSelect = (book) => {
  selectedBookId.value = book.id
  showBookDetail.value = true
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
</script>
