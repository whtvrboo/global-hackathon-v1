<template>
    <transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0"
        enter-to-class="opacity-100" leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100"
        leave-to-class="opacity-0">
        <div v-if="show && list" @click="$emit('close')"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
            <div @click.stop
                class="bg-gradient-to-br from-dark-900 to-dark-950 border border-dark-700 rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">

                <!-- Header -->
                <div class="sticky top-0 glass-strong border-b border-dark-700 p-6 z-10">
                    <div class="flex items-start justify-between">
                        <div class="flex-1">
                            <h2 class="text-heading-1 mb-3">{{ list.name }}</h2>
                            <p v-if="list.description" class="text-body text-dark-300 mb-4">
                                {{ list.description }}
                            </p>
                            <!-- List Creator -->
                            <div class="flex items-center gap-3 mb-4">
                                <img v-if="list.user?.picture" :src="list.user.picture" :alt="list.user.name"
                                    class="w-10 h-10 rounded-full border-2 border-dark-600" />
                                <div v-else
                                    class="w-10 h-10 rounded-full bg-dark-800 border-2 border-dark-600 flex items-center justify-center">
                                    <span class="text-dark-400">U</span>
                                </div>
                                <div>
                                    <div class="font-semibold text-white">{{ list.user?.name || 'Unknown' }}</div>
                                    <div class="text-caption">@{{ list.user?.username || 'unknown' }}</div>
                                </div>
                            </div>
                            <!-- Social Stats -->
                            <div class="flex items-center gap-6">
                                <div class="flex items-center gap-2">
                                    <span class="text-accent-orange"></span>
                                    <span class="text-body"><strong>{{ list.items_count || 0 }}</strong> books</span>
                                </div>
                                <button @click="toggleLike" class="flex items-center gap-2 transition-colors"
                                    :class="list.is_liked ? 'text-accent-red' : 'text-dark-400 hover:text-accent-red'">
                                    <svg class="w-5 h-5" :fill="list.is_liked ? 'currentColor' : 'none'"
                                        stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                                        </path>
                                    </svg>
                                    <span class="text-sm font-medium">{{ list.likes_count || 0 }}</span>
                                </button>
                                <button @click="showComments = !showComments"
                                    class="flex items-center gap-2 text-dark-400 hover:text-white transition-colors">
                                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                                        </path>
                                    </svg>
                                    <span class="text-sm font-medium">{{ list.comments_count || 0 }}</span>
                                </button>
                            </div>
                        </div>
                        <button @click="$emit('close')"
                            class="p-2 text-dark-400 hover:text-white transition-colors rounded-lg hover:bg-dark-800">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M6 18L18 6M6 6l12 12"></path>
                            </svg>
                        </button>
                    </div>
                </div>

                <!-- Books in List -->
                <div class="p-6">
                    <h3 class="text-heading-3 mb-6">Books in this list</h3>

                    <div v-if="loading" class="space-y-4">
                        <div v-for="i in 3" :key="i" class="animate-pulse flex gap-4">
                            <div class="w-20 h-32 bg-dark-800 rounded-xl"></div>
                            <div class="flex-1 space-y-2">
                                <div class="h-4 bg-dark-800 rounded w-3/4"></div>
                                <div class="h-3 bg-dark-800 rounded w-1/2"></div>
                                <div class="h-3 bg-dark-800 rounded w-full"></div>
                            </div>
                        </div>
                    </div>

                    <div v-else-if="list.items && list.items.length > 0" class="space-y-4">
                        <div v-for="(item, index) in list.items" :key="item.id"
                            class="card card-hover flex gap-4 cursor-pointer" @click="openBook(item.book)">
                            <div class="relative">
                                <div class="text-2xl font-bold text-dark-600 absolute -top-2 -left-2 z-10">
                                    {{ index + 1 }}
                                </div>
                                <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                                    class="w-20 h-32 object-cover rounded-xl shadow-lg" />
                                <div v-else class="w-20 h-32 bg-dark-800 rounded-xl flex items-center justify-center">
                                    <span class="text-2xl text-dark-400"></span>
                                </div>
                            </div>
                            <div class="flex-1">
                                <h4 class="text-heading-4 mb-1 line-clamp-2">{{ item.book.title }}</h4>
                                <p v-if="item.book.authors" class="text-body-sm text-dark-300 mb-2">
                                    by {{ item.book.authors.join(', ') }}
                                </p>
                                <p v-if="item.notes" class="text-caption text-dark-400 mb-3 line-clamp-2">
                                    {{ item.notes }}
                                </p>
                                <p v-if="item.book.description" class="text-body-sm text-dark-300 line-clamp-2">
                                    {{ item.book.description }}
                                </p>
                            </div>
                        </div>
                    </div>

                    <div v-else class="text-center py-12 text-dark-400">
                        <p class="text-4xl mb-4"></p>
                        <p>This list is empty</p>
                    </div>

                    <!-- Comments Section -->
                    <div v-if="showComments" class="mt-12 pt-8 border-t border-dark-700">
                        <h3 class="text-heading-3 mb-6">Comments</h3>

                        <!-- Comment Input -->
                        <div class="mb-6">
                            <TextArea v-model="newComment" placeholder="What do you think about this list?" :rows="3" />
                            <div class="flex justify-end mt-2">
                                <button @click="submitComment" :disabled="!newComment.trim()"
                                    class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed">
                                    Post Comment
                                </button>
                            </div>
                        </div>

                        <!-- Comments List -->
                        <div v-if="loadingComments" class="space-y-4">
                            <div v-for="i in 2" :key="i" class="animate-pulse flex gap-3">
                                <div class="w-10 h-10 bg-dark-800 rounded-full"></div>
                                <div class="flex-1 space-y-2">
                                    <div class="h-3 bg-dark-800 rounded w-1/4"></div>
                                    <div class="h-4 bg-dark-800 rounded w-3/4"></div>
                                </div>
                            </div>
                        </div>

                        <div v-else-if="comments.length > 0" class="space-y-4">
                            <div v-for="comment in comments" :key="comment.id" class="flex gap-3">
                                <img v-if="comment.user.picture" :src="comment.user.picture" :alt="comment.user.name"
                                    class="w-10 h-10 rounded-full" />
                                <div v-else class="w-10 h-10 rounded-full bg-dark-800 flex items-center justify-center">
                                    <span class="text-dark-400 text-sm">U</span>
                                </div>
                                <div class="flex-1">
                                    <div class="flex items-center gap-2 mb-1">
                                        <span class="font-semibold text-sm">{{ comment.user.name }}</span>
                                        <span class="text-caption">• {{ timeAgo(comment.created_at) }}</span>
                                    </div>
                                    <p class="text-body-sm">{{ comment.content }}</p>
                                </div>
                            </div>
                        </div>

                        <div v-else class="text-center py-8 text-dark-400">
                            <p>No comments yet. Be the first to share your thoughts!</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </transition>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import TextArea from './ui/TextArea.vue'

const props = defineProps({
    show: Boolean,
    listId: String
})

const emit = defineEmits(['close', 'updated', 'openBook'])

const router = useRouter()
const authStore = useAuthStore()
const toast = useToastStore()

const list = ref(null)
const comments = ref([])
const loading = ref(false)
const loadingComments = ref(false)
const showComments = ref(false)
const newComment = ref('')

watch(() => props.show, async (newVal) => {
    if (newVal && props.listId) {
        showComments.value = false
        await fetchList()
    }
})

async function fetchList() {
    loading.value = true
    try {
        const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/lists/${props.listId}`, {
            headers: {
                'Authorization': `Bearer ${authStore.token}`
            }
        })

        if (res.ok) {
            list.value = await res.json()
        } else {
            throw new Error('Failed to load list')
        }
    } catch (error) {
        console.error('Error fetching list:', error)
        toast.error('Failed to load list')
        emit('close')
    } finally {
        loading.value = false
    }
}

async function fetchComments() {
    loadingComments.value = true
    try {
        const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/lists/${props.listId}/comments`, {
            headers: {
                'Authorization': `Bearer ${authStore.token}`
            }
        })

        if (res.ok) {
            const data = await res.json()
            comments.value = data.comments
        }
    } catch (error) {
        console.error('Error fetching comments:', error)
    } finally {
        loadingComments.value = false
    }
}

async function toggleLike() {
    try {
        const method = list.value.is_liked ? 'DELETE' : 'POST'
        const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/lists/${props.listId}/like`, {
            method,
            headers: {
                'Authorization': `Bearer ${authStore.token}`
            }
        })

        if (res.ok) {
            list.value.is_liked = !list.value.is_liked
            list.value.likes_count = list.value.is_liked ?
                (list.value.likes_count + 1) : (list.value.likes_count - 1)
            emit('updated')
        }
    } catch (error) {
        console.error('Error toggling like:', error)
        toast.error('Failed to update like')
    }
}

async function submitComment() {
    if (!newComment.value.trim()) return

    try {
        const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/lists/${props.listId}/comments`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authStore.token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ content: newComment.value })
        })

        if (res.ok) {
            const comment = await res.json()
            comments.value.push(comment)
            list.value.comments_count = (list.value.comments_count || 0) + 1
            newComment.value = ''
            toast.success('Comment posted!')
            emit('updated')
        } else {
            throw new Error('Failed to post comment')
        }
    } catch (error) {
        console.error('Error posting comment:', error)
        toast.error('Failed to post comment')
    }
}

watch(() => showComments.value, (newVal) => {
    if (newVal && comments.value.length === 0) {
        fetchComments()
    }
})

function openBook(book) {
    router.push(`/books/${book.id}`)
    emit('close') // Close the modal when navigating
}

function timeAgo(dateString) {
    if (!dateString) return ''
    const date = new Date(dateString)
    const now = new Date()
    const seconds = Math.floor((now - date) / 1000)

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
