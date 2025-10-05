<template>
    <div class="min-h-screen bg-dark-950">

        <main class="container-mobile max-w-4xl mx-auto section-padding">
            <div class="mb-8">
                <h1 class="text-heading-1 mb-2">Your Reading Feed</h1>
                <p class="text-body text-dark-300">See when your friends create and update their book lists</p>
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
                <div v-for="item in feed" :key="item.id" @click="$router.push(`/lists/${item.id}`)"
                    class="card card-hover cursor-pointer group">
                    <!-- User Info -->
                    <div class="flex items-center gap-4 mb-6">
                        <div class="relative">
                            <img v-if="item.user.picture" :src="item.user.picture" :alt="item.user.name"
                                class="w-12 h-12 rounded-full border-2 border-dark-700" />
                            <div v-else
                                class="w-12 h-12 rounded-full bg-dark-800 border-2 border-dark-700 flex items-center justify-center">
                                <span class="text-dark-400">U</span>
                            </div>
                        </div>
                        <div class="flex-1">
                            <div class="font-semibold text-white">{{ item.user.name }}</div>
                            <div class="text-caption">
                                @{{ item.user.username }} • {{ timeAgo(item.created_at) }}
                            </div>
                        </div>
                        <span
                            class="px-3 py-1 text-xs font-semibold rounded-full glass-strong text-accent-purple bg-accent-purple/20">
                            Created a List
                        </span>
                    </div>

                    <!-- List Header Image -->
                    <div v-if="item.header_image_url"
                        class="aspect-video bg-gradient-to-br from-accent-blue/20 to-accent-purple/20 rounded-xl mb-4 overflow-hidden">
                        <img :src="item.header_image_url" :alt="item.title"
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                    </div>

                    <!-- List Info -->
                    <div class="mb-4">
                        <h3 class="text-heading-2 mb-2 group-hover:text-accent-blue transition-colors">{{ item.title }}
                        </h3>
                        <p v-if="item.description" class="text-body text-dark-300 line-clamp-2 mb-4">
                            {{ item.description }}
                        </p>
                    </div>

                    <!-- Preview Books -->
                    <div v-if="item.preview_books && item.preview_books.length > 0" class="mb-6">
                        <div class="flex gap-3 overflow-x-auto pb-2">
                            <div v-for="book in item.preview_books" :key="book.id" class="flex-shrink-0">
                                <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                                    class="w-20 h-28 object-cover rounded-lg shadow-md" />
                                <div v-else class="w-20 h-28 bg-dark-800 rounded-lg flex items-center justify-center">
                                    <span class="text-xl text-dark-400"></span>
                                </div>
                            </div>
                            <div v-if="item.items_count > 3"
                                class="w-20 h-28 bg-dark-800 rounded-lg flex items-center justify-center text-dark-300 text-sm font-medium">
                                +{{ item.items_count - 3 }}
                            </div>
                        </div>
                    </div>

                    <!-- List Stats -->
                    <div class="flex items-center justify-between text-sm text-dark-400 pt-4 border-t border-dark-800">
                        <span class="font-medium">{{ item.items_count }} book{{ item.items_count !== 1 ? 's' : ''
                            }}</span>
                        <div class="flex items-center gap-4">
                            <button @click.stop="toggleLike(item.id)" class="flex items-center gap-1 transition-colors"
                                :class="item.is_liked ? 'text-accent-red' : 'text-dark-400 hover:text-accent-red'">
                                <svg class="w-4 h-4" :fill="item.is_liked ? 'currentColor' : 'none'"
                                    stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                                    </path>
                                </svg>
                                {{ item.likes_count }}
                            </button>
                            <span class="flex items-center gap-1">
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                                    </path>
                                </svg>
                                {{ item.comments_count }}
                            </span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Empty State: Who to Follow -->
            <div v-else class="space-y-8">
                <div class="card text-center py-12">
                    <div class="text-6xl mb-6">👋</div>
                    <h3 class="text-heading-2 mb-4">Welcome to Your Feed</h3>
                    <p class="text-body text-dark-300 mb-8 max-w-md mx-auto">
                        Your feed will show list activity from people you follow. Start by following some of the best
                        curators on Folio!
                    </p>
                </div>

                <!-- Who to Follow Module -->
                <div class="card">
                    <div class="flex items-center justify-between mb-6">
                        <h3 class="text-heading-2">Top Curators to Follow</h3>
                    </div>

                    <!-- Loading Popular Users -->
                    <div v-if="loadingPopularUsers" class="space-y-4">
                        <div v-for="i in 3" :key="i"
                            class="animate-pulse flex items-center gap-4 p-4 bg-dark-800 rounded-xl">
                            <div class="w-12 h-12 bg-dark-700 rounded-full"></div>
                            <div class="flex-1 space-y-2">
                                <div class="h-4 bg-dark-700 rounded w-1/3"></div>
                                <div class="h-3 bg-dark-700 rounded w-1/2"></div>
                            </div>
                        </div>
                    </div>

                    <!-- Popular Users List -->
                    <div v-else-if="popularUsers.length > 0" class="space-y-4">
                        <div v-for="user in popularUsers" :key="user.id"
                            class="flex items-center gap-4 p-4 bg-dark-800 rounded-xl hover:bg-dark-700 transition-colors">
                            <img v-if="user.picture" :src="user.picture" :alt="user.name"
                                class="w-12 h-12 rounded-full border-2 border-dark-700" />
                            <div v-else
                                class="w-12 h-12 rounded-full bg-dark-700 border-2 border-dark-600 flex items-center justify-center">
                                <span class="text-dark-400">U</span>
                            </div>

                            <div class="flex-1 min-w-0">
                                <div class="font-semibold text-white">{{ user.name }}</div>
                                <div class="text-sm text-dark-400">@{{ user.username }}</div>
                                <div class="text-xs text-dark-500 mt-1">{{ user.list_count }} public lists</div>
                            </div>

                            <button v-if="!user.is_following" @click="followUser(user.username)"
                                :disabled="followingUsers[user.username]"
                                class="btn-primary text-sm px-4 py-2 disabled:opacity-50">
                                {{ followingUsers[user.username] ? 'Following...' : 'Follow' }}
                            </button>
                            <button v-else class="btn-secondary text-sm px-4 py-2">
                                Following
                            </button>
                        </div>
                    </div>

                    <!-- Empty State -->
                    <div v-else class="text-center py-8 text-dark-400">
                        <p>No curators found yet. Check back soon!</p>
                    </div>
                </div>

                <!-- CTA -->
                <div class="text-center">
                    <button @click="$router.push('/discover')" class="btn-secondary">
                        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                        Explore Lists
                    </button>
                </div>
            </div>
        </main>

    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'
import { useToastStore } from '../stores/toast'

const authStore = useAuthStore()

const feed = ref([])
const loading = ref(true)
const popularUsers = ref([])
const loadingPopularUsers = ref(false)
const followingUsers = ref({})
const toastStore = useToastStore()

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
            // Unlike the list
            await axios.delete(`/api/lists/${listId}/like`)
            item.is_liked = false
            item.likes_count = Math.max((item.likes_count || 0) - 1, 0)
            toastStore.info('Unliked')
        } else {
            // Like the list
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

const loadPopularUsers = async () => {
    loadingPopularUsers.value = true
    try {
        const response = await axios.get('/api/users/popular')
        popularUsers.value = response.data.users || []
    } catch (error) {
        console.error('Error loading popular users:', error)
        popularUsers.value = []
    } finally {
        loadingPopularUsers.value = false
    }
}

const followUser = async (username) => {
    followingUsers.value[username] = true
    try {
        await axios.post(`/api/users/${username}/follow`)

        // Update the user's following status
        const user = popularUsers.value.find(u => u.username === username)
        if (user) {
            user.is_following = true
        }

        toastStore.success(`You're now following ${username}!`)

        // Reload the feed after a short delay
        setTimeout(async () => {
            const response = await axios.get('/api/feed')
            feed.value = response.data.feed || []
        }, 500)
    } catch (error) {
        console.error('Error following user:', error)
        toastStore.error('Failed to follow user')
    } finally {
        followingUsers.value[username] = false
    }
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

        // If feed is empty, load popular users
        if (feed.value.length === 0) {
            await loadPopularUsers()
        }
    } catch (error) {
        console.error('Error loading feed:', error)
        if (error.response?.status === 401) {
            console.log('Unauthorized - user may need to re-authenticate')
        }
        // Load popular users even on error
        await loadPopularUsers()
    } finally {
        loading.value = false
    }
})
</script>
