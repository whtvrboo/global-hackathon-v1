<template>
    <div class="min-h-screen bg-gray-50">
        <!-- Header -->
        <div class="bg-white border-b sticky top-0 z-10">
            <div class="max-w-5xl mx-auto px-4 py-4">
                <div class="flex items-center justify-between">
                    <router-link to="/" class="text-2xl font-bold">📚 Folio</router-link>
                    <div class="flex items-center gap-4">
                        <router-link to="/profile" class="text-gray-600 hover:text-gray-900">Profile</router-link>
                        <button @click="authStore.logout()" class="text-gray-600 hover:text-gray-900">
                            Logout
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div class="max-w-3xl mx-auto px-4 py-8">
            <h1 class="text-3xl font-bold mb-8">Your Reading Feed</h1>

            <!-- Loading -->
            <div v-if="loading" class="space-y-4">
                <div v-for="i in 3" :key="i" class="animate-pulse">
                    <Card>
                        <div class="flex gap-4">
                            <div class="w-16 h-24 bg-gray-200 rounded"></div>
                            <div class="flex-1 space-y-3">
                                <div class="h-4 bg-gray-200 rounded w-3/4"></div>
                                <div class="h-3 bg-gray-200 rounded w-1/2"></div>
                                <div class="h-3 bg-gray-200 rounded w-full"></div>
                            </div>
                        </div>
                    </Card>
                </div>
            </div>

            <!-- Feed Items -->
            <div v-else-if="feed.length > 0" class="space-y-6">
                <Card v-for="item in feed" :key="item.id" hover class="transition-all duration-200">
                    <!-- User Info -->
                    <div class="flex items-center gap-3 mb-4">
                        <img v-if="item.user.picture" :src="item.user.picture" :alt="item.user.name"
                            class="w-10 h-10 rounded-full" />
                        <div class="flex-1">
                            <div class="font-semibold">{{ item.user.name }}</div>
                            <div class="text-sm text-gray-500">
                                @{{ item.user.username }} • {{ timeAgo(item.created_at) }}
                            </div>
                        </div>
                        <span class="px-3 py-1 text-xs font-medium rounded-full" :class="statusBadgeClass(item.status)">
                            {{ statusLabel(item.status) }}
                        </span>
                    </div>

                    <!-- Book Info -->
                    <div class="flex gap-4">
                        <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                            class="w-20 h-30 object-cover rounded-lg shadow-sm" />
                        <div class="flex-1">
                            <h3 class="font-bold text-lg mb-1">{{ item.book.title }}</h3>
                            <p v-if="item.book.authors" class="text-gray-600 text-sm mb-3">
                                by {{ item.book.authors.join(', ') }}
                            </p>

                            <!-- Rating -->
                            <div v-if="item.rating" class="flex items-center gap-1 mb-3">
                                <span v-for="i in 5" :key="i" class="text-xl"
                                    :class="i <= item.rating ? 'text-yellow-500' : 'text-gray-300'">
                                    ★
                                </span>
                            </div>

                            <!-- Review -->
                            <p v-if="item.review" class="text-gray-700 leading-relaxed">
                                {{ item.review }}
                            </p>
                        </div>
                    </div>

                    <!-- Actions -->
                    <div class="flex items-center gap-4 mt-4 pt-4 border-t text-sm">
                        <button class="text-gray-500 hover:text-primary transition-colors">
                            💬 Comment
                        </button>
                        <button class="text-gray-500 hover:text-red-500 transition-colors">
                            ❤️ Like
                        </button>
                    </div>
                </Card>
            </div>

            <!-- Empty State -->
            <Card v-else class="text-center py-12">
                <div class="text-6xl mb-4">📖</div>
                <h3 class="text-xl font-semibold text-gray-900 mb-2">Your feed is empty</h3>
                <p class="text-gray-600 mb-6">Follow other readers to see their activity here</p>
                <PrimaryButton @click="$router.push('/')">
                    Explore Books
                </PrimaryButton>
            </Card>
        </div>
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
        'want_to_read': 'bg-blue-100 text-blue-700',
        'reading': 'bg-green-100 text-green-700',
        'read': 'bg-purple-100 text-purple-700',
        'dnf': 'bg-gray-100 text-gray-700'
    }
    return classes[status] || 'bg-gray-100 text-gray-700'
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
