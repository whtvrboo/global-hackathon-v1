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
                        <router-link to="/profile" class="btn-ghost">Profile</router-link>
                        <button @click="authStore.logout()" class="btn-ghost">Logout</button>
                    </nav>

                    <!-- Mobile Menu Button -->
                    <button class="md:hidden p-2 text-dark-300 hover:text-white transition-colors">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M4 6h16M4 12h16M4 18h16"></path>
                        </svg>
                    </button>
                </div>
            </div>
        </header>

        <!-- Main Content -->
        <main class="container-mobile max-w-7xl mx-auto section-padding">
            <div class="text-center mb-8">
                <h1 class="text-heading-1 mb-2 flex items-center justify-center gap-3">
                    <span>🎯</span>
                    <span>Discover Books</span>
                </h1>
                <p class="text-body text-dark-300">Swipe through personalized recommendations</p>
            </div>

            <!-- Loading Initial Cards -->
            <div v-if="loading && cards.length === 0" class="text-center py-16">
                <div
                    class="inline-block animate-spin rounded-full h-16 w-16 border-2 border-dark-600 border-t-accent-red mb-4">
                </div>
                <p class="text-body text-dark-300">Finding amazing books for you...</p>
            </div>

            <!-- Card Stack -->
            <div v-else class="relative w-full max-w-sm mx-auto aspect-[3/4] mb-8">
                <!-- Background cards (for depth effect) -->
                <div v-for="(card, index) in visibleCards" :key="card.book.id" :style="{
                    transform: `scale(${1 - index * 0.05}) translateY(${index * 10}px)`,
                    zIndex: visibleCards.length - index,
                    opacity: 1 - index * 0.3
                }" class="absolute inset-0">
                    <SwipeableBookCard v-if="index === 0" :book="card.book" :reason="card.reason"
                        @like="handleLike(card)" @pass="handlePass(card)" />
                    <div v-else class="h-full bg-dark-800 rounded-3xl shadow-xl"></div>
                </div>

                <!-- Empty State -->
                <div v-if="cards.length === 0 && !loading"
                    class="absolute inset-0 flex items-center justify-center text-center p-8">
                    <div class="card">
                        <div class="text-6xl mb-6">🎉</div>
                        <h3 class="text-heading-2 mb-4">You've seen them all!</h3>
                        <p class="text-body text-dark-300 mb-8">Check back later for more recommendations</p>
                        <button @click="loadMore" class="btn-primary">
                            Load More
                        </button>
                    </div>
                </div>
            </div>

            <!-- Stats -->
            <div class="grid grid-cols-3 gap-6 max-w-md mx-auto mb-8">
                <div class="card text-center">
                    <div class="text-3xl font-bold text-accent-green mb-2">{{ stats.liked }}</div>
                    <div class="text-caption">Liked</div>
                </div>
                <div class="card text-center">
                    <div class="text-3xl font-bold text-accent-red mb-2">{{ stats.passed }}</div>
                    <div class="text-caption">Passed</div>
                </div>
                <div class="card text-center">
                    <div class="text-3xl font-bold text-accent-blue mb-2">{{ stats.total }}</div>
                    <div class="text-caption">Seen</div>
                </div>
            </div>

            <!-- Instructions (show on first visit) -->
            <transition enter-active-class="transition ease-out duration-300" enter-from-class="opacity-0 translate-y-4"
                leave-active-class="transition ease-in duration-200" leave-to-class="opacity-0 translate-y-4">
                <div v-if="showInstructions" class="fixed bottom-24 left-4 right-4 z-50">
                    <div class="card glass-strong">
                        <button @click="showInstructions = false"
                            class="absolute top-3 right-3 p-2 text-dark-400 hover:text-white">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M6 18L18 6M6 6l12 12"></path>
                            </svg>
                        </button>
                        <h3 class="font-bold text-lg mb-4 text-white">How it works 👋</h3>
                        <ul class="space-y-3 text-sm text-dark-200">
                            <li class="flex items-center gap-3">
                                <span class="text-2xl">👈</span>
                                <span>Swipe left to pass</span>
                            </li>
                            <li class="flex items-center gap-3">
                                <span class="text-2xl">👉</span>
                                <span>Swipe right to like & save</span>
                            </li>
                            <li class="flex items-center gap-3">
                                <span class="text-2xl">ℹ️</span>
                                <span>Tap info for full details</span>
                            </li>
                        </ul>
                    </div>
                </div>
            </transition>
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
                    class="flex flex-col items-center gap-1 p-3 text-accent-blue transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                    </svg>
                    <span class="text-xs">Discover</span>
                </router-link>

                <router-link to="/feed"
                    class="flex flex-col items-center gap-1 p-3 text-dark-400 hover:text-white transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M4 6h16M4 10h16M4 14h16M4 18h16"></path>
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
        <BookDetailModal :show="showDetailModal" :book-id="selectedBookId" @close="showDetailModal = false"
            @log="handleLogBook" />

        <!-- Log Book Modal -->
        <LogBookModal :show="showLogModal" :book="bookToLog" @close="showLogModal = false"
            @success="handleLogSuccess" />
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'
import SwipeableBookCard from '../components/SwipeableBookCard.vue'
import BookDetailModal from '../components/BookDetailModal.vue'
import LogBookModal from '../components/LogBookModal.vue'
import PrimaryButton from '../components/ui/PrimaryButton.vue'

const authStore = useAuthStore()

const cards = ref([])
const loading = ref(true)
const showInstructions = ref(true)
const stats = ref({
    liked: 0,
    passed: 0,
    total: 0
})

const showDetailModal = ref(false)
const showLogModal = ref(false)
const selectedBookId = ref(null)
const bookToLog = ref(null)

const visibleCards = computed(() => {
    return cards.value.slice(0, 3) // Show 3 cards for depth effect
})

onMounted(async () => {
    await loadRecommendations()

    // Check if user has seen instructions before
    const hasSeenInstructions = localStorage.getItem('folio_seen_discover_instructions')
    if (hasSeenInstructions) {
        showInstructions.value = false
    }
})

const loadRecommendations = async () => {
    loading.value = true

    try {
        const response = await axios.get('/api/discover', {
            params: { limit: 20 }
        })

        const newCards = response.data.recommendations || []
        cards.value.push(...newCards)
    } catch (error) {
        console.error('Error loading recommendations:', error)
    } finally {
        loading.value = false
    }
}

const handleLike = async (recommendation) => {
    stats.value.liked++
    stats.value.total++

    // Remove card from stack
    cards.value = cards.value.filter(c => c.book.id !== recommendation.book.id)

    // Record swipe
    try {
        await axios.post('/api/discover/swipe', {
            book_id: recommendation.book.id,
            action: 'like'
        })
    } catch (error) {
        // Ignore errors for swipe recording
    }

    // Show log modal
    bookToLog.value = recommendation.book
    setTimeout(() => {
        showLogModal.value = true
    }, 300)

    // Load more if running low
    if (cards.value.length < 5) {
        loadRecommendations()
    }

    hideInstructionsIfNeeded()
}

const handlePass = async (recommendation) => {
    stats.value.passed++
    stats.value.total++

    // Remove card from stack
    cards.value = cards.value.filter(c => c.book.id !== recommendation.book.id)

    // Record swipe
    try {
        await axios.post('/api/discover/swipe', {
            book_id: recommendation.book.id,
            action: 'pass'
        })
    } catch (error) {
        // Ignore errors
    }

    // Load more if running low
    if (cards.value.length < 5) {
        loadRecommendations()
    }

    hideInstructionsIfNeeded()
}

const showBookInfo = (recommendation) => {
    selectedBookId.value = recommendation.book.id
    showDetailModal.value = true
}

const handleLogBook = (book) => {
    bookToLog.value = book
    showDetailModal.value = false
    showLogModal.value = true
}

const handleLogSuccess = () => {
    console.log('Book logged successfully!')
}

const loadMore = () => {
    loadRecommendations()
}

const hideInstructionsIfNeeded = () => {
    if (showInstructions.value) {
        setTimeout(() => {
            showInstructions.value = false
            localStorage.setItem('folio_seen_discover_instructions', 'true')
        }, 2000)
    }
}
</script>
