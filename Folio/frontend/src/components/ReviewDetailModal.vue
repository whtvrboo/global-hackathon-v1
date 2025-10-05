<template>
    <transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0"
        enter-to-class="opacity-100" leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100"
        leave-to-class="opacity-0">
        <div v-if="show && log" @click="$emit('close')"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
            <div @click.stop
                class="bg-gradient-to-br from-dark-900 to-dark-950 border border-dark-700 rounded-2xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">

                <!-- Header with Book Info -->
                <div class="sticky top-0 glass-strong border-b border-dark-700 p-6 z-10">
                    <div class="flex items-start gap-6">
                        <img v-if="log.book.cover_url" :src="log.book.cover_url" :alt="log.book.title"
                            class="w-24 h-36 object-cover rounded-xl shadow-xl" />
                        <div v-else class="w-24 h-36 bg-dark-800 rounded-xl flex items-center justify-center">
                            <span class="text-4xl text-dark-400"></span>
                        </div>
                        <div class="flex-1">
                            <h2 class="text-heading-2 mb-2">{{ log.book.title }}</h2>
                            <p v-if="log.book.authors" class="text-body text-dark-300 mb-3">
                                by {{ log.book.authors.join(', ') }}
                            </p>
                            <div class="flex items-center gap-4">
                                <span class="px-3 py-1 text-xs font-semibold rounded-full glass-strong"
                                    :class="statusBadgeClass(log.status)">
                                    {{ statusLabel(log.status) }}
                                </span>
                                <div v-if="log.rating" class="flex items-center gap-1">
                                    <span v-for="i in 5" :key="i" class="text-lg"
                                        :class="i <= log.rating ? 'text-accent-orange' : 'text-dark-600'">
                                        ★
                                    </span>
                                    <span class="text-caption ml-1">{{ log.rating }}/5</span>
                                </div>
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

                <!-- Review Content -->
                <div class="p-8">
                    <!-- User Info -->
                    <div class="flex items-center gap-4 mb-8">
                        <img v-if="log.user.picture" :src="log.user.picture" :alt="log.user.name"
                            class="w-16 h-16 rounded-full border-2 border-dark-600" />
                        <div v-else
                            class="w-16 h-16 rounded-full bg-dark-800 border-2 border-dark-600 flex items-center justify-center">
                            <span class="text-dark-400 text-2xl">U</span>
                        </div>
                        <div class="flex-1">
                            <div class="font-semibold text-heading-3">{{ log.user.name }}</div>
                            <div class="text-caption">
                                @{{ log.user.username }} • {{ formatDate(log.created_at) }}
                            </div>
                        </div>
                    </div>

                    <!-- Spoiler Warning -->
                    <div v-if="log.spoiler_flag && !showSpoiler"
                        class="mb-6 p-6 bg-accent-red/10 border-2 border-accent-red/20 rounded-xl text-center">
                        <div class="text-4xl mb-3">!</div>
                        <h3 class="text-heading-3 text-accent-red mb-2">This review contains spoilers</h3>
                        <p class="text-body text-dark-300 mb-4">
                            Are you sure you want to read it?
                        </p>
                        <button @click="showSpoiler = true"
                            class="btn-secondary bg-accent-red/20 hover:bg-accent-red/30 text-accent-red border-accent-red/30">
                            Show Spoilers
                        </button>
                    </div>

                    <!-- Review Text -->
                    <div v-if="log.review && (!log.spoiler_flag || showSpoiler)"
                        class="prose prose-invert prose-lg max-w-none mb-8">
                        <div class="text-body-lg leading-relaxed whitespace-pre-wrap">
                            {{ log.review }}
                        </div>
                    </div>

                    <!-- Notes (if any) -->
                    <div v-if="log.notes && (!log.spoiler_flag || showSpoiler)"
                        class="mb-8 p-4 bg-dark-800/50 rounded-lg">
                        <h4 class="text-sm font-semibold text-dark-300 mb-2">📝 Personal Notes</h4>
                        <p class="text-body-sm text-dark-300">{{ log.notes }}</p>
                    </div>

                    <!-- Reading Dates -->
                    <div v-if="log.start_date || log.finish_date" class="mb-8 flex gap-4">
                        <div v-if="log.start_date" class="p-4 bg-dark-800/30 rounded-lg flex-1">
                            <div class="text-xs text-dark-400 mb-1">Started</div>
                            <div class="text-body font-medium">{{ formatDate(log.start_date) }}</div>
                        </div>
                        <div v-if="log.finish_date" class="p-4 bg-dark-800/30 rounded-lg flex-1">
                            <div class="text-xs text-dark-400 mb-1">Finished</div>
                            <div class="text-body font-medium">{{ formatDate(log.finish_date) }}</div>
                        </div>
                    </div>

                    <!-- Social Actions -->
                    <div class="flex items-center gap-6 pt-6 border-t border-dark-700">
                        <button @click="toggleLike" class="flex items-center gap-2 transition-colors"
                            :class="log.is_liked ? 'text-accent-red' : 'text-dark-400 hover:text-accent-red'">
                            <svg class="w-6 h-6" :fill="log.is_liked ? 'currentColor' : 'none'" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                                </path>
                            </svg>
                            <span class="text-sm font-medium">{{ log.likes_count || 0 }}</span>
                        </button>

                        <button @click="showComments = !showComments"
                            class="flex items-center gap-2 text-dark-400 hover:text-white transition-colors">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                                </path>
                            </svg>
                            <span class="text-sm font-medium">{{ log.comments_count || 0 }}</span>
                        </button>
                    </div>

                    <!-- Comments Section -->
                    <div v-if="showComments" class="mt-8 pt-8 border-t border-dark-700">
                        <h3 class="text-heading-3 mb-6">Comments</h3>

                        <!-- Comment Input -->
                        <div class="mb-6">
                            <TextArea v-model="newComment" placeholder="Share your thoughts..." :rows="3" />
                            <div class="flex justify-end mt-2">
                                <button @click="submitComment" :disabled="!newComment.trim()"
                                    class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed">
                                    Post Comment
                                </button>
                            </div>
                        </div>

                        <!-- Comments List -->
                        <div v-if="loading" class="space-y-4">
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
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import TextArea from './ui/TextArea.vue'

const props = defineProps({
    show: Boolean,
    logId: String
})

const emit = defineEmits(['close', 'updated'])

const authStore = useAuthStore()
const toast = useToastStore()

const log = ref(null)
const comments = ref([])
const loading = ref(false)
const showSpoiler = ref(false)
const showComments = ref(false)
const newComment = ref('')

watch(() => props.show, async (newVal) => {
    if (newVal && props.logId) {
        showSpoiler.value = false
        showComments.value = false
        await fetchLog()
    }
})

async function fetchLog() {
    loading.value = true
    try {
        const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/logs/${props.logId}`, {
            headers: {
                'Authorization': `Bearer ${authStore.token}`
            }
        })

        if (res.ok) {
            log.value = await res.json()
        } else {
            throw new Error('Failed to load review')
        }
    } catch (error) {
        console.error('Error fetching log:', error)
        toast.error('Failed to load review')
        emit('close')
    } finally {
        loading.value = false
    }
}

async function fetchComments() {
    try {
        const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/logs/${props.logId}/comments`, {
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
    }
}

async function toggleLike() {
    try {
        const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/logs/${props.logId}/like`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authStore.token}`
            }
        })

        if (res.ok) {
            const data = await res.json()
            log.value.is_liked = data.liked
            log.value.likes_count = data.liked ? (log.value.likes_count + 1) : (log.value.likes_count - 1)
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
        const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/logs/${props.logId}/comments`, {
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
            log.value.comments_count = (log.value.comments_count || 0) + 1
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

function statusLabel(status) {
    const labels = {
        want_to_read: 'Want to Read',
        reading: 'Reading',
        read: 'Read',
        dnf: 'DNF'
    }
    return labels[status] || status
}

function statusBadgeClass(status) {
    const classes = {
        want_to_read: 'text-accent-blue border-accent-blue/30',
        reading: 'text-accent-orange border-accent-orange/30',
        read: 'text-accent-green border-accent-green/30',
        dnf: 'text-dark-400 border-dark-600'
    }
    return classes[status] || ''
}

function formatDate(dateString) {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
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
