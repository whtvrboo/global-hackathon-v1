<script setup lang="ts">
import type { Comment } from '@/data/dummyData'

interface Props {
    comments: Comment[]
}

defineProps<Props>()

// Format date for display
const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60)

    if (diffInHours < 1) {
        return 'Just now'
    } else if (diffInHours < 24) {
        return `${Math.floor(diffInHours)}h ago`
    } else if (diffInHours < 168) { // 7 days
        return `${Math.floor(diffInHours / 24)}d ago`
    } else {
        return date.toLocaleDateString()
    }
}

// Seek to timestamp (placeholder - would need parent communication)
const seekToTimestamp = (timestamp: number) => {
    console.log('Seek to timestamp:', timestamp)
    // TODO: Implement seeking functionality
    // This would need to communicate with the parent component
    // to seek all waveforms to the specified timestamp
}
</script>

<template>
    <div class="comment-sidebar bg-base-100 border border-base-300 rounded-lg p-4 h-full">
        <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-lg">Comments</h3>
            <span class="badge badge-outline badge-sm">{{ comments.length }}</span>
        </div>

        <div v-if="comments.length === 0" class="text-center py-8">
            <div class="text-base-content/60 mb-2">
                <svg class="w-8 h-8 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                    </path>
                </svg>
            </div>
            <p class="text-sm text-base-content/60">No comments yet</p>
            <p class="text-xs text-base-content/50 mt-1">Click on a waveform to add a comment</p>
        </div>

        <div v-else class="space-y-3 max-h-96 overflow-y-auto">
            <div v-for="comment in comments" :key="comment.id"
                class="comment-item bg-base-200 rounded p-3 hover:bg-base-300 transition-colors">
                <div class="flex justify-between items-start mb-2">
                    <div class="flex items-center gap-2">
                        <div class="w-6 h-6 bg-primary/20 rounded-full flex items-center justify-center">
                            <span class="text-xs font-medium text-primary">{{ comment.author.charAt(0).toUpperCase()
                                }}</span>
                        </div>
                        <span class="text-sm font-medium">{{ comment.author }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                        <span class="badge badge-sm badge-outline">{{ comment.timestamp.toFixed(1) }}s</span>
                        <button class="btn btn-ghost btn-xs" @click="seekToTimestamp(comment.timestamp)">
                            <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                                <path d="M8 5v10l8-5-8-5z" />
                            </svg>
                        </button>
                    </div>
                </div>
                <p class="text-sm mb-2">{{ comment.text }}</p>
                <span class="text-xs text-base-content/50">
                    {{ formatDate(comment.createdAt) }}
                </span>
            </div>
        </div>
    </div>
</template>

<style scoped>
.comment-sidebar {
    min-height: 400px;
}
</style>
