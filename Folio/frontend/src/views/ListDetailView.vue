<template>
    <div class="min-h-screen bg-dark-950">

        <!-- Loading -->
        <div v-if="loading" class="container-mobile max-w-4xl mx-auto section-padding">
            <div class="text-center py-16">
                <div
                    class="inline-block animate-spin rounded-full h-16 w-16 border-2 border-dark-600 border-t-accent-red mb-4">
                </div>
                <p class="text-body text-dark-300">Loading list...</p>
            </div>
        </div>

        <!-- List Content -->
        <div v-else-if="list" class="min-h-screen">
            <!-- Hero Section -->
            <div class="relative overflow-hidden" :style="{ backgroundColor: list.theme_color || '#6366f1' }">
                <!-- Header Image -->
                <div v-if="list.header_image_url" class="absolute inset-0">
                    <img :src="list.header_image_url" :alt="list.name" class="w-full h-full object-cover opacity-20" />
                </div>

                <!-- Gradient Overlay -->
                <div class="absolute inset-0 bg-gradient-to-b from-black/50 to-transparent"></div>

                <!-- Content -->
                <div class="relative container-mobile max-w-4xl mx-auto section-padding py-16">
                    <div class="text-center text-white">
                        <h1 class="text-4xl md:text-6xl font-bold mb-4">{{ list.name }}</h1>
                        <p v-if="list.description" class="text-xl md:text-2xl text-white/90 mb-8 max-w-3xl mx-auto">
                            {{ list.description }}
                        </p>

                        <!-- Creator Info -->
                        <div class="flex items-center justify-center gap-4 mb-6">
                            <router-link :to="`/profile/${list.creator.username}`"
                                class="flex items-center gap-4 hover:opacity-80 transition-opacity">
                                <img v-if="list.creator.picture" :src="list.creator.picture" :alt="list.creator.name"
                                    class="w-12 h-12 rounded-full border-2 border-white/20" />
                                <div v-else
                                    class="w-12 h-12 rounded-full bg-white/20 border-2 border-white/20 flex items-center justify-center">
                                    <span class="text-white">U</span>
                                </div>
                                <div class="text-left">
                                    <div class="font-semibold">{{ list.creator.name }}</div>
                                    <div class="text-white/70">@{{ list.creator.username }}</div>
                                </div>
                            </router-link>
                        </div>

                        <!-- Liked By Section -->
                        <div v-if="list.liked_by && list.liked_by.length > 0"
                            class="flex items-center justify-center gap-3 mb-8">
                            <div class="flex -space-x-2">
                                <router-link v-for="user in list.liked_by.slice(0, 5)" :key="user.id"
                                    :to="`/profile/${user.username}`"
                                    class="relative hover:z-10 hover:scale-110 transition-transform">
                                    <img v-if="user.picture" :src="user.picture" :alt="user.name"
                                        class="w-8 h-8 rounded-full border-2 border-white/30" />
                                    <div v-else
                                        class="w-8 h-8 rounded-full bg-white/20 border-2 border-white/30 flex items-center justify-center">
                                        <span class="text-xs text-white">U</span>
                                    </div>
                                </router-link>
                            </div>
                            <span class="text-sm text-white/80">
                                Liked by <span class="font-semibold">{{ list.liked_by[0].name }}</span>
                                <span v-if="list.likes_count > 1">
                                    and {{ list.likes_count - 1 }} other{{ list.likes_count > 2 ? 's' : '' }}
                                </span>
                            </span>
                        </div>

                        <!-- Stats -->
                        <div class="flex items-center justify-center gap-8 text-white/80">
                            <span class="flex items-center gap-2">
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253">
                                    </path>
                                </svg>
                                {{ list.items_count }} books
                            </span>
                            <span class="flex items-center gap-2">
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                                    </path>
                                </svg>
                                {{ list.likes_count || 0 }} likes
                            </span>
                            <span class="flex items-center gap-2">
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                                    </path>
                                </svg>
                                {{ list.comments_count || 0 }} comments
                            </span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- List Items -->
            <div class="container-mobile max-w-4xl mx-auto section-padding">
                <div class="space-y-6">
                    <div v-for="(item, index) in list.items" :key="item.id" class="card card-hover group">
                        <div class="flex gap-6">
                            <!-- Number -->
                            <div class="flex-shrink-0 w-12 h-12 rounded-full flex items-center justify-center text-white font-bold text-lg"
                                :style="{ backgroundColor: list.theme_color || '#6366f1' }">
                                {{ index + 1 }}
                            </div>

                            <!-- Book Cover -->
                            <div class="flex-shrink-0">
                                <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                                    class="w-20 h-32 object-cover rounded-xl shadow-lg group-hover:scale-105 transition-transform duration-200" />
                                <div v-else class="w-20 h-32 bg-dark-800 rounded-xl flex items-center justify-center">
                                    <span class="text-2xl text-dark-400"></span>
                                </div>
                            </div>

                            <!-- Book Info -->
                            <div class="flex-1 min-w-0">
                                <h3 class="text-heading-3 mb-2 line-clamp-2">{{ item.book.title }}</h3>
                                <p v-if="item.book.authors?.length" class="text-body text-dark-300 mb-3">
                                    by {{ item.book.authors.join(', ') }}
                                </p>

                                <!-- Notes -->
                                <div v-if="item.notes" class="mb-4 p-4 bg-dark-800/50 rounded-lg border-l-4"
                                    :style="{ borderLeftColor: list.theme_color || '#6366f1' }">
                                    <p class="text-body text-dark-200 italic">"{{ item.notes }}"</p>
                                </div>

                                <!-- Description -->
                                <p v-if="item.book.description" class="text-caption text-dark-400 line-clamp-3">
                                    {{ item.book.description }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Actions -->
                <div class="mt-12 pt-8 border-t border-dark-800">
                    <div class="flex flex-col sm:flex-row gap-4 justify-center">
                        <button v-if="authStore.isAuthenticated" @click="toggleLike"
                            class="btn-primary flex items-center gap-2"
                            :class="{ 'bg-accent-red': isLiked, 'hover:bg-accent-red/80': isLiked }">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                                </path>
                            </svg>
                            {{ isLiked ? 'Liked' : 'Like List' }}
                        </button>

                        <button v-if="authStore.isAuthenticated" @click="showComments = !showComments"
                            class="btn-secondary flex items-center gap-2">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                                </path>
                            </svg>
                            Comment
                        </button>

                        <button @click="shareList" class="btn-outline flex items-center gap-2">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.367 2.684 3 3 0 00-5.367-2.684z">
                                </path>
                            </svg>
                            Share
                        </button>
                    </div>
                </div>

                <!-- Comments Section -->
                <div v-if="showComments" class="mt-8 pt-8 border-t border-dark-800">
                    <h3 class="text-heading-2 mb-6">Comments</h3>

                    <!-- Comment Input -->
                    <div v-if="authStore.isAuthenticated" class="mb-8">
                        <div class="flex gap-3">
                            <div class="flex-shrink-0">
                                <img v-if="authStore.user?.picture" :src="authStore.user.picture"
                                    :alt="authStore.user.name" class="w-10 h-10 rounded-full border border-dark-700" />
                                <div v-else
                                    class="w-10 h-10 rounded-full bg-dark-800 border border-dark-700 flex items-center justify-center">
                                    <span class="text-sm text-dark-400">U</span>
                                </div>
                            </div>
                            <div class="flex-1 flex gap-2">
                                <input v-model="commentInput" type="text"
                                    placeholder="Share your thoughts on this list..."
                                    class="flex-1 px-4 py-3 bg-dark-900 border border-dark-700 rounded-xl text-white placeholder-dark-400 focus:outline-none focus:border-accent-blue transition-colors"
                                    @keyup.enter="submitComment" />
                                <button @click="submitComment" :disabled="!commentInput.trim() || submittingComment"
                                    class="px-6 py-3 bg-accent-blue text-white rounded-xl hover:bg-accent-blue/80 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
                                    Post
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- Comments List -->
                    <div v-if="comments.length > 0" class="space-y-6">
                        <div v-for="comment in comments" :key="comment.id" class="flex gap-4">
                            <div class="flex-shrink-0">
                                <img v-if="comment.user.picture" :src="comment.user.picture" :alt="comment.user.name"
                                    class="w-10 h-10 rounded-full border border-dark-700" />
                                <div v-else
                                    class="w-10 h-10 rounded-full bg-dark-800 border border-dark-700 flex items-center justify-center">
                                    <span class="text-sm text-dark-400">U</span>
                                </div>
                            </div>
                            <div class="flex-1 bg-dark-900 rounded-xl px-4 py-3">
                                <div class="flex items-center justify-between mb-2">
                                    <div class="flex items-center gap-2">
                                        <span class="font-semibold text-white">{{ comment.user.name }}</span>
                                        <span class="text-sm text-dark-400">@{{ comment.user.username }}</span>
                                        <span class="text-sm text-dark-500">•</span>
                                        <span class="text-sm text-dark-400">{{ timeAgo(comment.created_at) }}</span>
                                    </div>
                                </div>
                                <p class="text-dark-200">{{ comment.content }}</p>
                            </div>
                        </div>
                    </div>

                    <!-- No Comments -->
                    <div v-else class="text-center py-8 text-dark-400">
                        <div class="text-4xl mb-4">💬</div>
                        <p>No comments yet. Be the first to share your thoughts!</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Error State -->
        <div v-else class="container-mobile max-w-4xl mx-auto section-padding">
            <div class="card text-center py-16">
                <div class="text-6xl mb-6"></div>
                <h3 class="text-heading-2 mb-4">List not found</h3>
                <p class="text-body text-dark-300 mb-8">This list might be private or doesn't exist.</p>
                <router-link to="/discover" class="btn-primary">Discover Lists</router-link>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import axios from 'axios'

const route = useRoute()
const authStore = useAuthStore()
const toastStore = useToastStore()

const list = ref(null)
const loading = ref(true)
const showComments = ref(false)
const comments = ref([])
const commentInput = ref('')
const submittingComment = ref(false)
const isLiked = ref(false)

const listId = computed(() => route.params.id)

onMounted(async () => {
    await loadList()
    if (authStore.isAuthenticated) {
        await loadComments()
    }
})

const loadList = async () => {
    loading.value = true
    try {
        const response = await axios.get(`/api/lists/${listId.value}`)
        list.value = response.data

        // Check if user has liked this list
        if (authStore.isAuthenticated) {
            // This would need to be implemented in the backend
            // For now, we'll assume false
            isLiked.value = false
        }
    } catch (error) {
        console.error('Error loading list:', error)
        list.value = null
    } finally {
        loading.value = false
    }
}

const loadComments = async () => {
    try {
        const response = await axios.get(`/api/lists/${listId.value}/comments`)
        comments.value = response.data.comments || []
    } catch (error) {
        console.error('Error loading comments:', error)
        comments.value = []
    }
}

const toggleLike = async () => {
    if (!authStore.isAuthenticated) {
        toastStore.error('Please log in to like lists')
        return
    }

    try {
        if (isLiked.value) {
            await axios.delete(`/api/lists/${listId.value}/like`)
            isLiked.value = false
            list.value.likes_count = Math.max((list.value.likes_count || 0) - 1, 0)
            toastStore.info('Unliked')
        } else {
            await axios.post(`/api/lists/${listId.value}/like`)
            isLiked.value = true
            list.value.likes_count = (list.value.likes_count || 0) + 1
            toastStore.success('Liked!')
        }
    } catch (error) {
        console.error('Error toggling like:', error)
        toastStore.error('Failed to update like')
    }
}

const submitComment = async () => {
    if (!authStore.isAuthenticated) {
        toastStore.error('Please log in to comment')
        return
    }

    const content = commentInput.value.trim()
    if (!content) return

    try {
        submittingComment.value = true
        const response = await axios.post(`/api/lists/${listId.value}/comments`, {
            content
        })

        comments.value.push(response.data)
        list.value.comments_count = (list.value.comments_count || 0) + 1
        commentInput.value = ''
        toastStore.success('Comment posted!')
    } catch (error) {
        console.error('Error submitting comment:', error)
        toastStore.error('Failed to post comment')
    } finally {
        submittingComment.value = false
    }
}

const shareList = async () => {
    const url = window.location.href
    try {
        await navigator.clipboard.writeText(url)
        toastStore.success('Link copied to clipboard!')
    } catch (error) {
        // Fallback for older browsers
        const textArea = document.createElement('textarea')
        textArea.value = url
        document.body.appendChild(textArea)
        textArea.select()
        document.execCommand('copy')
        document.body.removeChild(textArea)
        toastStore.success('Link copied to clipboard!')
    }
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
</script>
