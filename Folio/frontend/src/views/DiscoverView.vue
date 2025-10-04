<template>
    <div class="h-screen flex flex-col bg-gradient-to-br from-primary/5 via-secondary/5 to-primary/5 overflow-hidden">
        <!-- Header -->
        <div class="bg-white/80 backdrop-blur border-b flex-shrink-0">
            <div class="max-w-7xl mx-auto px-4 py-4">
                <div class="flex items-center justify-between">
                    <h1 class="text-2xl font-bold flex items-center gap-2">
                        <span>🎯</span>
                        <span>Discover</span>
                    </h1>
                    <div class="flex items-center gap-4">
                        <router-link to="/" class="text-gray-600 hover:text-gray-900 transition-colors">
                            Search
                        </router-link>
                        <router-link to="/profile" class="text-gray-600 hover:text-gray-900 transition-colors">
                            Profile
                        </router-link>
                    </div>
                </div>
            </div>
        </div>

        <!-- Main Content -->
        <div class="flex-1 flex flex-col items-center justify-center p-4 relative overflow-hidden">
            <!-- Loading Initial Cards -->
            <div v-if="loading && cards.length === 0" class="text-center">
                <div class="inline-block animate-spin rounded-full h-16 w-16 border-b-4 border-primary mb-4"></div>
                <p class="text-lg text-gray-600">Finding amazing books for you...</p>
            </div>

            <!-- Card Stack -->
            <div v-else class="relative w-full max-w-md aspect-[3/4] mx-auto">
                <!-- Background cards (for depth effect) -->
                <div v-for="(card, index) in visibleCards" :key="card.id" :style="{
                    transform: `scale(${1 - index * 0.05}) translateY(${index * 10}px)`,
                    zIndex: visibleCards.length - index,
                    opacity: 1 - index * 0.3
                }" class="absolute inset-0">
                    <SwipeableBookCard v-if="index === 0" :book="card" @swipe-left="handlePass(card)"
                        @swipe-right="handleLike(card)" @info="showBookInfo(card)" />
                    <div v-else class="h-full bg-white rounded-3xl shadow-xl"></div>
                </div>

                <!-- Empty State -->
                <div v-if="cards.length === 0 && !loading"
                    class="absolute inset-0 flex items-center justify-center text-center p-8">
                    <div>
                        <div class="text-6xl mb-4">🎉</div>
                        <h3 class="text-2xl font-bold mb-2">You've seen them all!</h3>
                        <p class="text-gray-600 mb-6">Check back later for more recommendations</p>
                        <PrimaryButton @click="loadMore">
                            Load More
                        </PrimaryButton>
                    </div>
                </div>
            </div>

            <!-- Stats -->
            <div class="mt-8 flex items-center justify-center gap-8 text-center">
                <div>
                    <div class="text-3xl font-bold text-primary">{{ stats.liked }}</div>
                    <div class="text-sm text-gray-600">Liked</div>
                </div>
                <div>
                    <div class="text-3xl font-bold text-gray-600">{{ stats.passed }}</div>
                    <div class="text-sm text-gray-600">Passed</div>
                </div>
                <div>
                    <div class="text-3xl font-bold text-secondary">{{ stats.total }}</div>
                    <div class="text-sm text-gray-600">Seen</div>
                </div>
            </div>

            <!-- Instructions (show on first visit) -->
            <transition enter-active-class="transition ease-out duration-300" enter-from-class="opacity-0 translate-y-4"
                leave-active-class="transition ease-in duration-200" leave-to-class="opacity-0 translate-y-4">
                <div v-if="showInstructions"
                    class="absolute bottom-32 left-0 right-0 mx-4 p-6 bg-white/95 backdrop-blur rounded-2xl shadow-xl">
                    <button @click="showInstructions = false"
                        class="absolute top-2 right-2 p-2 text-gray-400 hover:text-gray-600">
                        ✕
                    </button>
                    <h3 class="font-bold text-lg mb-3">How it works 👋</h3>
                    <ul class="space-y-2 text-sm text-gray-700">
                        <li class="flex items-center gap-2">
                            <span>👈</span>
                            <span>Swipe left to pass</span>
                        </li>
                        <li class="flex items-center gap-2">
                            <span>👉</span>
                            <span>Swipe right to like & save</span>
                        </li>
                        <li class="flex items-center gap-2">
                            <span>ℹ️</span>
                            <span>Tap info for full details</span>
                        </li>
                    </ul>
                </div>
            </transition>
        </div>

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
import axios from 'axios'
import SwipeableBookCard from '../components/SwipeableBookCard.vue'
import BookDetailModal from '../components/BookDetailModal.vue'
import LogBookModal from '../components/LogBookModal.vue'
import PrimaryButton from '../components/ui/PrimaryButton.vue'

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

const handleLike = async (book) => {
    stats.value.liked++
    stats.value.total++

    // Remove card from stack
    cards.value = cards.value.filter(c => c.id !== book.id)

    // Record swipe
    try {
        await axios.post('/api/discover/swipe', {
            book_id: book.id,
            action: 'like'
        })
    } catch (error) {
        // Ignore errors for swipe recording
    }

    // Show log modal
    bookToLog.value = book
    setTimeout(() => {
        showLogModal.value = true
    }, 300)

    // Load more if running low
    if (cards.value.length < 5) {
        loadRecommendations()
    }

    hideInstructionsIfNeeded()
}

const handlePass = async (book) => {
    stats.value.passed++
    stats.value.total++

    // Remove card from stack
    cards.value = cards.value.filter(c => c.id !== book.id)

    // Record swipe
    try {
        await axios.post('/api/discover/swipe', {
            book_id: book.id,
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

const showBookInfo = (book) => {
    selectedBookId.value = book.id
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
