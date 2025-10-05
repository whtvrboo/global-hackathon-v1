<template>
    <div
        class="group relative bg-dark-900/50 backdrop-blur-sm rounded-lg border border-dark-700 p-6 hover:border-accent-blue/50 transition-all duration-200">
        <!-- Type Badge -->
        <div class="flex items-start justify-between mb-3">
            <span :class="[
                'inline-flex items-center px-2 py-1 text-xs font-medium rounded-full',
                annotation.type === 'highlight'
                    ? 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30'
                    : 'bg-accent-blue/20 text-accent-blue border border-accent-blue/30'
            ]">
                <svg v-if="annotation.type === 'highlight'" class="w-3 h-3 mr-1" fill="currentColor"
                    viewBox="0 0 20 20">
                    <path
                        d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
                <svg v-else class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
                {{ annotation.type === 'highlight' ? 'Highlight' : 'Note' }}
            </span>

            <!-- Actions Menu -->
            <div class="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                <button @click="editAnnotation"
                    class="p-1.5 text-dark-400 hover:text-accent-blue hover:bg-accent-blue/10 rounded transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                </button>
                <button @click="deleteAnnotation"
                    class="p-1.5 text-dark-400 hover:text-red-400 hover:bg-red-400/10 rounded transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </div>
        </div>

        <!-- Content -->
        <div class="mb-4">
            <p class="text-white leading-relaxed whitespace-pre-wrap">{{ annotation.content }}</p>

            <!-- Context (for highlights) -->
            <div v-if="annotation.context" class="mt-3 p-3 bg-dark-800/50 border-l-4 border-dark-600 rounded-r">
                <p class="text-sm text-dark-300 italic">"{{ annotation.context }}"</p>
            </div>
        </div>

        <!-- Book Context -->
        <div v-if="annotation.book"
            class="flex items-center space-x-3 mb-4 p-3 bg-dark-800/30 rounded-lg border border-dark-700">
            <img v-if="annotation.book.cover_url" :src="annotation.book.cover_url" :alt="annotation.book.title"
                class="w-12 h-16 object-cover rounded shadow-sm" />
            <div class="flex-1 min-w-0">
                <h4 class="font-medium text-white truncate">{{ annotation.book.title }}</h4>
                <p v-if="annotation.book.authors && annotation.book.authors.length > 0"
                    class="text-sm text-dark-400 truncate">
                    {{ annotation.book.authors.join(', ') }}
                </p>
                <div class="flex items-center space-x-2 mt-1">
                    <span v-if="annotation.page_number" class="text-xs text-dark-500 bg-dark-700 px-2 py-0.5 rounded">
                        p. {{ annotation.page_number }}
                    </span>
                    <span class="text-xs text-dark-500">
                        {{ formatDate(annotation.created_at) }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Tags -->
        <div v-if="annotation.tags && annotation.tags.length > 0" class="flex flex-wrap gap-2">
            <span v-for="tag in annotation.tags" :key="tag"
                class="inline-flex items-center px-2 py-1 text-xs text-accent-blue bg-accent-blue/10 border border-accent-blue/20 rounded-full">
                #{{ tag }}
            </span>
        </div>

        <!-- Timeline Connector -->
        <div class="absolute left-6 -bottom-3 w-0.5 h-6 bg-gradient-to-b from-accent-blue/50 to-transparent"></div>
    </div>
</template>

<script setup>
import { useToastStore } from '../stores/toast'

const toast = useToastStore()

// Props
const props = defineProps({
    annotation: {
        type: Object,
        required: true
    }
})

// Emits
const emit = defineEmits(['edit', 'delete'])

// Methods
const editAnnotation = () => {
    emit('edit', props.annotation)
}

const deleteAnnotation = async () => {
    if (!confirm('Are you sure you want to delete this annotation?')) {
        return
    }

    try {
        const token = localStorage.getItem('token')
        const response = await fetch(`/api/annotations/${props.annotation.id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })

        if (response.ok) {
            toast.success('Annotation deleted')
            emit('delete', props.annotation.id)
        } else {
            throw new Error('Failed to delete annotation')
        }
    } catch (error) {
        console.error('Error deleting annotation:', error)
        toast.error('Failed to delete annotation')
    }
}

const formatDate = (dateString) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now - date
    const diffDays = Math.floor(diffMs / 86400000)

    if (diffDays === 0) {
        const diffHours = Math.floor(diffMs / 3600000)
        if (diffHours === 0) {
            const diffMinutes = Math.floor(diffMs / 60000)
            return diffMinutes < 1 ? 'Just now' : `${diffMinutes}m ago`
        }
        return `${diffHours}h ago`
    }
    if (diffDays === 1) return 'Yesterday'
    if (diffDays < 7) return `${diffDays}d ago`
    if (diffDays < 30) return `${Math.floor(diffDays / 7)}w ago`
    if (diffDays < 365) return `${Math.floor(diffDays / 30)}mo ago`

    return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' })
}
</script>

<style scoped>
/* Smooth hover transitions */
.group:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

/* Custom scrollbar for long content */
.whitespace-pre-wrap {
    word-break: break-word;
}

/* Timeline connector animation */
.group:hover .absolute {
    background: linear-gradient(to bottom, rgba(59, 130, 246, 0.8), transparent);
}
</style>
