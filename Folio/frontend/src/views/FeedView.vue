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

        <main class="container-mobile max-w-4xl mx-auto section-padding">
            <div class="mb-8">
                <h1 class="text-heading-1 mb-2">Your Reading Feed</h1>
                <p class="text-body text-dark-300">See what your friends are reading and discover new books</p>
            </div>

            <!-- Loading -->
            <div v-if="loading" class="space-y-6">
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

            <!-- Feed Items -->
            <div v-else-if="feed.length > 0" class="space-y-6">
                <div v-for="item in feed" :key="item.id" class="card card-hover">
                    <!-- User Info -->
                    <div class="flex items-center gap-4 mb-6">
                        <div class="relative">
                            <img v-if="item.user.picture" :src="item.user.picture" :alt="item.user.name"
                                class="w-12 h-12 rounded-full border-2 border-dark-700" />
                            <div v-else
                                class="w-12 h-12 rounded-full bg-dark-800 border-2 border-dark-700 flex items-center justify-center">
                                <span class="text-dark-400">👤</span>
                            </div>
                        </div>
                        <div class="flex-1">
                            <div class="font-semibold text-white">{{ item.user.name }}</div>
                            <div class="text-caption">
                                @{{ item.user.username }} • {{ timeAgo(item.created_at) }}
                            </div>
                        </div>
                        <span class="px-3 py-1 text-xs font-semibold rounded-full glass-strong"
                            :class="statusBadgeClass(item.status)">
                            {{ statusLabel(item.status) }}
                        </span>
                    </div>

                    <!-- Book Info -->
                    <div class="flex gap-6">
                        <div class="relative">
                            <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                                class="w-24 h-36 object-cover rounded-xl shadow-lg" />
                            <div v-else class="w-24 h-36 bg-dark-800 rounded-xl flex items-center justify-center">
                                <span class="text-2xl text-dark-400">📚</span>
                            </div>
                        </div>
                        <div class="flex-1">
                            <h3 class="text-heading-3 mb-2">{{ item.book.title }}</h3>
                            <p v-if="item.book.authors" class="text-body text-dark-300 mb-4">
                                by {{ item.book.authors.join(', ') }}
                            </p>

                            <!-- Rating -->
                            <div v-if="item.rating" class="flex items-center gap-1 mb-4">
                                <div class="flex">
                                    <span v-for="i in 5" :key="i" class="text-lg"
                                        :class="i <= item.rating ? 'text-accent-orange' : 'text-dark-600'">
                                        ★
                                    </span>
                                </div>
                                <span class="text-caption ml-2">{{ item.rating }}/5</span>
                            </div>

                            <!-- Review -->
                            <p v-if="item.review" class="text-body leading-relaxed">
                                {{ item.review }}
                            </p>
                        </div>
                    </div>

                    <!-- Actions -->
                    <div class="flex items-center gap-6 mt-6 pt-6 border-t border-dark-800">
                        <button class="flex items-center gap-2 text-dark-400 hover:text-white transition-colors">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                                </path>
                            </svg>
                            <span class="text-sm">Comment</span>
                        </button>
                        <button class="flex items-center gap-2 text-dark-400 hover:text-accent-red transition-colors">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                                </path>
                            </svg>
                            <span class="text-sm">Like</span>
                        </button>
                    </div>
                </div>
            </div>

            <!-- Empty State -->
            <div v-else class="card text-center py-16">
                <div class="text-6xl mb-6">📖</div>
                <h3 class="text-heading-2 mb-4">Your feed is empty</h3>
                <p class="text-body text-dark-300 mb-8 max-w-md mx-auto">
                    Follow other readers to see their activity here, or start by discovering some great books
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

                <router-link to="/feed" class="flex flex-col items-center gap-1 p-3 text-accent-red transition-colors">
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
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'
import Card from '../components/ui/Card.vue'
import PrimaryButton from '../components/ui/PrimaryButton.vue'

const authStore = useAuthStore()

const feed = ref([])
const loading = ref(true)

const statusLabel = (status) => {
    const labels = {
        'want_to_read': 'Want to Read',
        'reading': 'Reading',
        'read': 'Finished',
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

onMounted(async () => {
    try {
        // Check if user is authenticated
        if (!authStore.isAuthenticated) {
            console.log('User not authenticated, skipping feed load')
            return
        }

        console.log('Loading feed for user:', authStore.user?.username)
        const response = await axios.get('/api/feed')
        console.log('Feed response:', response.data)
        feed.value = response.data.feed || []
    } catch (error) {
        console.error('Error loading feed:', error)
        if (error.response?.status === 401) {
            console.log('Unauthorized - user may need to re-authenticate')
            // Optionally redirect to login or refresh token
        }
    } finally {
        loading.value = false
    }
})
</script>
