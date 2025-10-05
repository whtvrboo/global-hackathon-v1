<template>
    <div class="min-h-screen bg-dark-950">

        <main class="container-mobile max-w-7xl mx-auto section-padding">
            <div class="text-center mb-8">
                <h1 class="text-heading-1 mb-2 flex items-center justify-center gap-3">
                    <span>{{ authStore.isAuthenticated ? 'Your Reading Feed' : 'Discover Curated Lists' }}</span>
                </h1>
                <p class="text-body text-dark-300 max-w-2xl mx-auto">
                    {{ authStore.isAuthenticated ? 
                        'See when your friends create and update their book lists' : 
                        'Experience the magic of curation. Browse beautiful, hand-crafted book collections from passionate readers.' }}
                </p>
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
                        <router-link :to="`/profile/${item.user.username}`"
                            class="relative hover:opacity-80 transition-opacity">
                            <img v-if="item.user.picture" :src="item.user.picture" :alt="item.user.name"
                                class="w-12 h-12 rounded-full border-2 border-dark-700" />
                            <div v-else
                                class="w-12 h-12 rounded-full bg-dark-800 border-2 border-dark-700 flex items-center justify-center">
                                <span class="text-dark-400">👤</span>
                            </div>
                        </router-link>
                        <div class="flex-1 min-w-0">
                            <p class="text-body text-dark-200 leading-relaxed">
                                <router-link :to="`/profile/${item.user.username}`"
                                    class="font-bold text-white hover:text-accent-blue transition-colors">
                                    {{ item.user.name }}
                                </router-link>
                                <span class="text-dark-300"> just published </span>
                                <span class="font-bold text-white">"{{ item.name }}"</span>
                            </p>
                            <div class="text-xs text-dark-500 mt-1">
                                {{ timeAgo(item.created_at) }}
                            </div>
                        </div>
                    </div>

                    <!-- List Header Image -->
                    <div v-if="item.header_image_url"
                        class="aspect-video bg-gradient-to-br from-accent-blue/20 to-accent-purple/20 rounded-xl mb-4 overflow-hidden">
                        <img :src="item.header_image_url" :alt="item.name"
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                    </div>

                    <!-- List Info -->
                    <div class="mb-4">
                        <h3 class="text-heading-2 mb-2 group-hover:text-accent-blue transition-colors">{{ item.name }}
                        </h3>
                        <p v-if="item.description" class="text-body text-dark-300 line-clamp-2 mb-4">
                            {{ item.description }}
                        </p>
                    </div>

                    <!-- Preview Books -->
                    <div v-if="item.preview_books && item.preview_books.length > 0" class="mb-6">
                        <div class="flex gap-3 overflow-x-auto pb-2">
                            <div v-for="book in item.preview_books" :key="book.id" class="flex-shrink-0 cursor-pointer"
                                @click="navigateToBook(book.id)">
                                <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                                    class="w-20 h-28 object-cover rounded-lg shadow-md hover:shadow-lg transition-shadow" />
                                <div v-else
                                    class="w-20 h-28 bg-dark-800 rounded-lg flex items-center justify-center hover:bg-dark-700 transition-colors">
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

            <!-- Popular Curators Section (for non-authenticated users or when feed is empty) -->
            <div v-if="!authStore.isAuthenticated || feed.length === 0" class="mb-12">
                <h2 class="text-heading-2 mb-6 flex items-center gap-2">
                    <span>Popular Curators</span>
                </h2>
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    <div v-for="user in popularUsers.slice(0, 6)" :key="user.id"
                        class="card card-hover flex items-center gap-4 p-4 cursor-pointer"
                        @click="$router.push(`/profile/${user.username}`)">
                        <img v-if="user.picture" :src="user.picture" :alt="user.name"
                            class="w-16 h-16 rounded-full border-2 border-dark-700" />
                        <div v-else
                            class="w-16 h-16 rounded-full bg-dark-800 border-2 border-dark-700 flex items-center justify-center">
                            <span class="text-2xl text-dark-400">U</span>
                        </div>

                        <div class="flex-1 min-w-0">
                            <div class="font-semibold text-white truncate">{{ user.name }}</div>
                            <div class="text-sm text-dark-400 truncate">@{{ user.username }}</div>
                            <div class="text-xs text-dark-500 mt-1">{{ user.list_count }} lists</div>
                        </div>

                        <button v-if="authStore.isAuthenticated && !user.is_following"
                            @click.stop="followUser(user.username)" :disabled="followingUsers[user.username]"
                            class="btn-primary text-sm px-3 py-2 disabled:opacity-50 whitespace-nowrap">
                            {{ followingUsers[user.username] ? 'Following...' : 'Follow' }}
                        </button>
                        <button v-else-if="authStore.isAuthenticated" @click.stop
                            class="btn-secondary text-sm px-3 py-2 whitespace-nowrap">
                            Following
                        </button>
                    </div>
                </div>
            </div>

            <!-- Popular Lists Section -->
            <div v-if="!authStore.isAuthenticated || feed.length === 0">
                <h2 class="text-heading-2 mb-6 flex items-center gap-2">
                    <span>Popular Lists</span>
                </h2>
            </div>

            <!-- Popular Lists Grid -->
            <div v-if="popularLists.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
                <div v-for="list in popularLists" :key="list.id" @click="$router.push(`/lists/${list.id}`)"
                    class="card card-hover cursor-pointer group">
                    <!-- List Header Image -->
                    <div v-if="list.header_image_url"
                        class="aspect-video bg-gradient-to-br from-accent-blue/20 to-accent-purple/20 rounded-xl mb-4 overflow-hidden">
                        <img :src="list.header_image_url" :alt="list.name"
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                    </div>
                    <div v-else
                        class="aspect-video bg-gradient-to-br from-accent-blue/20 to-accent-purple/20 rounded-xl mb-4 flex items-center justify-center">
                    </div>

                    <!-- List Info -->
                    <div class="p-4">
                        <h3 class="text-heading-3 mb-2 line-clamp-2">{{ list.name }}</h3>
                        <p v-if="list.description" class="text-body text-dark-300 line-clamp-2 mb-4">
                            {{ list.description }}
                        </p>

                        <!-- Creator Info -->
                        <div class="flex items-center gap-3 mb-4">
                            <img v-if="list.creator.picture" :src="list.creator.picture" :alt="list.creator.name"
                                class="w-8 h-8 rounded-full border border-dark-700" />
                            <div v-else
                                class="w-8 h-8 rounded-full bg-dark-800 border border-dark-700 flex items-center justify-center">
                                <span class="text-xs text-dark-400">U</span>
                            </div>
                            <div>
                                <div class="text-sm font-medium text-white">{{ list.creator.name }}</div>
                                <div class="text-xs text-dark-400">@{{ list.creator.username }}</div>
                            </div>
                        </div>

                        <!-- List Stats -->
                        <div class="flex items-center justify-between text-sm text-dark-400">
                            <span>{{ list.items_count }} books</span>
                            <div class="flex items-center gap-4">
                                <span class="flex items-center gap-1">
                                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                                        </path>
                                    </svg>
                                    {{ list.likes_count }}
                                </span>
                                <span class="flex items-center gap-1">
                                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                                        </path>
                                    </svg>
                                    {{ list.comments_count }}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Empty State -->
            <div v-else class="card text-center py-16">
                <div class="text-6xl mb-6"></div>
                <h3 class="text-heading-2 mb-4">No lists yet</h3>
                <p class="text-body text-dark-300 mb-8">
                    {{ authStore.isAuthenticated ?
                        'Be the first to create an amazing book list!' :
                        'Sign up to create and share your own book lists!' }}
                </p>
                <button v-if="authStore.isAuthenticated" @click="$router.push('/profile')" class="btn-primary">
                    Create Your First List
                </button>
                <button v-else @click="$router.push('/login')" class="btn-primary">
                    Sign Up to Create Lists
                </button>
            </div>
        </main>

    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'
import { useToastStore } from '../stores/toast'

const router = useRouter()
const authStore = useAuthStore()

const feed = ref([])
const loading = ref(true)
const popularUsers = ref([])
const popularLists = ref([])
const loadingPopularUsers = ref(false)
const followingUsers = ref({})
const toastStore = useToastStore()

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

const loadPopularLists = async () => {
    try {
        const response = await axios.get('/api/lists/popular', {
            params: { limit: 12 }
        })
        popularLists.value = response.data.lists || []
    } catch (error) {
        console.error('Error loading popular lists:', error)
        popularLists.value = []
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
        // Always load popular users and lists for discover functionality
        await Promise.all([loadPopularUsers(), loadPopularLists()])

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
        }
    } finally {
        loading.value = false
    }
})
</script>
