<template>
    <div class="min-h-screen bg-dark-950">

        <!-- Main Content -->
        <main class="container-mobile max-w-7xl mx-auto section-padding">
            <div class="text-center mb-8">
                <h1 class="text-heading-1 mb-2 flex items-center justify-center gap-3">
                    <span>Discover Curated Lists</span>
                </h1>
                <p class="text-body text-dark-300 max-w-2xl mx-auto">
                    Experience the magic of curation. Browse beautiful, hand-crafted book collections from passionate
                    readers.
                </p>
            </div>

            <!-- Loading -->
            <div v-if="loading" class="text-center py-16">
                <div
                    class="inline-block animate-spin rounded-full h-16 w-16 border-2 border-dark-600 border-t-accent-red mb-4">
                </div>
                <p class="text-body text-dark-300">Finding amazing lists for you...</p>
            </div>

            <!-- Popular Curators Section -->
            <div v-if="!loading && popularUsers.length > 0" class="mb-12">
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
            <div v-if="!loading">
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
import { useAuthStore } from '../stores/auth'
import axios from 'axios'

const authStore = useAuthStore()

const popularLists = ref([])
const popularUsers = ref([])
const followingUsers = ref({})
const loading = ref(true)

onMounted(async () => {
    await Promise.all([loadPopularLists(), loadPopularUsers()])
})

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

const loadPopularUsers = async () => {
    try {
        const response = await axios.get('/api/users/popular')
        popularUsers.value = response.data.users || []
    } catch (error) {
        console.error('Error loading popular users:', error)
        popularUsers.value = []
    } finally {
        loading.value = false
    }
}

const followUser = async (username) => {
    if (!authStore.isAuthenticated) {
        return
    }

    followingUsers.value[username] = true
    try {
        await axios.post(`/api/users/${username}/follow`)

        // Update the user's following status
        const user = popularUsers.value.find(u => u.username === username)
        if (user) {
            user.is_following = true
        }
    } catch (error) {
        console.error('Error following user:', error)
    } finally {
        followingUsers.value[username] = false
    }
}
</script>
