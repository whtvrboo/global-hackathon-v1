<template>
    <div class="min-h-screen bg-dark-950">
        <!-- Header -->
        <div class="bg-dark-900/50 backdrop-blur-sm border-b border-dark-800">
            <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
                <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-4">
                        <button @click="$router.back()"
                            class="p-2 text-dark-400 hover:text-white hover:bg-dark-800 rounded-lg transition-colors">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M15 19l-7-7 7-7" />
                            </svg>
                        </button>
                        <div>
                            <h1 class="text-heading-1 text-white">#{{ tag }}</h1>
                            <p class="mt-1 text-body text-dark-400">
                                Your thoughts on {{ tag }} across {{ annotations.length }} annotation{{
                                annotations.length !== 1 ? 's' : '' }}
                            </p>
                        </div>
                    </div>
                    <div class="flex items-center space-x-3">
                        <button @click="createListFromThread"
                            class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-accent-red hover:bg-accent-red/80 rounded-lg transition-colors">
                            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M12 4v16m8-8H4" />
                            </svg>
                            Create List
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Main Content -->
        <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <!-- Loading State -->
            <div v-if="isLoading" class="flex items-center justify-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-accent-red"></div>
            </div>

            <!-- Empty State -->
            <div v-else-if="annotations.length === 0" class="text-center py-12">
                <svg class="mx-auto h-16 w-16 text-dark-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                        d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                </svg>
                <h3 class="text-lg font-medium text-white mb-2">No annotations found</h3>
                <p class="text-sm text-dark-400">
                    No thoughts tagged with "{{ tag }}" yet
                </p>
            </div>

            <!-- Timeline -->
            <div v-else class="space-y-6">
                <!-- Timeline Header -->
                <div class="text-center py-4">
                    <div class="inline-flex items-center px-4 py-2 bg-dark-800/50 border border-dark-700 rounded-full">
                        <svg class="w-4 h-4 mr-2 text-accent-blue" fill="none" stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span class="text-sm text-white font-medium">Your intellectual journey</span>
                    </div>
                </div>

                <!-- Annotations Timeline -->
                <div class="relative">
                    <!-- Timeline Line -->
                    <div
                        class="absolute left-6 top-0 bottom-0 w-0.5 bg-gradient-to-b from-accent-blue/50 via-accent-blue/30 to-transparent">
                    </div>

                    <!-- Annotations -->
                    <div class="space-y-8">
                        <AnnotationThreadItem v-for="(annotation, index) in annotations" :key="annotation.id"
                            :annotation="annotation" @edit="editAnnotation" @delete="deleteAnnotation" />
                    </div>
                </div>

                <!-- Timeline Footer -->
                <div class="text-center py-8">
                    <div class="inline-flex items-center px-4 py-2 bg-dark-800/30 border border-dark-700 rounded-full">
                        <svg class="w-4 h-4 mr-2 text-dark-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                        <span class="text-sm text-dark-400">End of your {{ tag }} journey</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Create List Modal (simplified for now) -->
        <div v-if="showCreateListModal"
            class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
            <div class="bg-dark-900 border border-dark-700 rounded-lg p-6 w-full max-w-md">
                <h3 class="text-lg font-medium text-white mb-4">Create List from Thread</h3>
                <p class="text-sm text-dark-400 mb-4">
                    Create a curated list from your thoughts on "{{ tag }}"
                </p>
                <div class="flex items-center space-x-3">
                    <button @click="showCreateListModal = false"
                        class="flex-1 px-4 py-2 text-sm font-medium text-dark-400 bg-dark-800 hover:bg-dark-700 rounded-lg transition-colors">
                        Cancel
                    </button>
                    <button @click="confirmCreateList"
                        class="flex-1 px-4 py-2 text-sm font-medium text-white bg-accent-red hover:bg-accent-red/80 rounded-lg transition-colors">
                        Create List
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToastStore } from '../stores/toast'
import AnnotationThreadItem from '../components/AnnotationThreadItem.vue'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()

// State
const annotations = ref([])
const isLoading = ref(false)
const showCreateListModal = ref(false)

// Computed
const tag = computed(() => {
    return decodeURIComponent(route.params.tag || '')
})

// Lifecycle
onMounted(() => {
    if (tag.value) {
        fetchAnnotationThread()
    }
})

// Methods
const fetchAnnotationThread = async () => {
    isLoading.value = true

    try {
        const token = localStorage.getItem('token')
        const response = await fetch(
            `/api/annotations/thread?tag=${encodeURIComponent(tag.value)}`,
            {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            }
        )

        if (response.ok) {
            const data = await response.json()
            annotations.value = data.annotations || []
        } else {
            throw new Error('Failed to fetch annotation thread')
        }
    } catch (error) {
        console.error('Error fetching annotation thread:', error)
        toast.error('Failed to load thread')
    } finally {
        isLoading.value = false
    }
}

const editAnnotation = (annotation) => {
    // TODO: Implement edit functionality
    console.log('Edit annotation:', annotation)
}

const deleteAnnotation = (annotationId) => {
    // Remove from local state
    annotations.value = annotations.value.filter(a => a.id !== annotationId)
}

const createListFromThread = () => {
    showCreateListModal.value = true
}

const confirmCreateList = () => {
    // Get unique books from annotations
    const books = new Map()
    annotations.value.forEach(annotation => {
        if (annotation.book && !books.has(annotation.book.id)) {
            books.set(annotation.book.id, {
                ...annotation.book,
                annotations: annotations.value.filter(a => a.book?.id === annotation.book.id)
            })
        }
    })

    // Navigate to list creation with pre-filled data
    const listData = {
        title: `My Thoughts on #${tag.value}`,
        description: `A curated collection of my thoughts and insights on ${tag.value}`,
        books: Array.from(books.values())
    }

    // Store in session storage for the list creation page
    sessionStorage.setItem('prefilledListData', JSON.stringify(listData))

    showCreateListModal.value = false
    router.push('/lists/create')
}
</script>

<style scoped>
/* Custom timeline styling */
.relative::before {
    content: '';
    position: absolute;
    left: 6px;
    top: 0;
    bottom: 0;
    width: 2px;
    background: linear-gradient(to bottom,
            rgba(59, 130, 246, 0.5) 0%,
            rgba(59, 130, 246, 0.3) 50%,
            transparent 100%);
}

/* Smooth animations */
.space-y-8>* {
    animation: fadeInUp 0.6s ease-out;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(20px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
