<template>
  <div class="min-h-screen bg-gradient-to-br from-primary/5 to-secondary/5">
    <!-- Header -->
    <div class="bg-white/80 backdrop-blur border-b sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold">📚 Folio</h1>
          <div class="flex items-center gap-4">
            <router-link to="/discover" class="text-gray-600 hover:text-gray-900 transition-colors">
              Discover
            </router-link>
            <router-link to="/feed" class="text-gray-600 hover:text-gray-900 transition-colors">
              Feed
            </router-link>
            <router-link to="/profile" class="text-gray-600 hover:text-gray-900 transition-colors">
              Profile
            </router-link>
            <button @click="authStore.logout()" class="text-gray-600 hover:text-gray-900 transition-colors">
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 py-12">
      <!-- Hero Section -->
      <div class="text-center mb-12">
        <h2 class="text-5xl font-bold text-gray-900 mb-4">
          Track Your Reading Journey
        </h2>
        <p class="text-xl text-gray-600 mb-8">
          Discover, log, and share the books you love
        </p>

        <!-- Search Bar -->
        <div class="flex justify-center">
          <SearchBar @select="handleBookSelect" />
        </div>
      </div>

      <!-- Welcome Message -->
      <Card v-if="authStore.user" class="max-w-2xl mx-auto text-center">
        <div class="flex items-center justify-center gap-4">
          <img v-if="authStore.user.picture" :src="authStore.user.picture" :alt="authStore.user.name"
            class="w-16 h-16 rounded-full" />
          <div v-else class="w-16 h-16 rounded-full bg-gray-200 flex items-center justify-center text-2xl">
            👤
          </div>
          <div class="text-left">
            <div class="flex items-center gap-2">
              <h3 class="text-xl font-bold">Welcome, {{ authStore.user.name }}!</h3>
              <span v-if="authStore.isGuestUser"
                class="px-2 py-1 text-xs font-medium bg-yellow-100 text-yellow-800 rounded-full">
                Guest
              </span>
            </div>
            <p class="text-gray-600">
              {{ authStore.isGuestUser ?
                'Try out Folio - your data will be saved locally' : 'Start searching for your next great read' }}
            </p>
          </div>
        </div>

        <!-- Guest Conversion Prompt -->
        <div v-if="authStore.isGuestUser" class="mt-4 pt-4 border-t">
          <button @click="showConversionModal = true" class="text-sm text-primary hover:text-primary/80 font-medium">
            Create a full account to save your progress →
          </button>
        </div>
      </Card>
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
