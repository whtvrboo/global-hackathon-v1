<template>
    <Transition name="slide">
        <div v-if="isOpen"
            class="fixed inset-y-0 right-0 w-full sm:w-96 bg-dark-900/95 backdrop-blur-md shadow-2xl border-l border-dark-800 z-40 flex flex-col">
            <!-- Header -->
            <div class="flex items-center justify-between p-4 border-b border-dark-800">
                <h3 class="text-heading-3">Annotations</h3>
                <button @click="close" class="text-dark-400 hover:text-white transition-colors">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            <!-- Tabs -->
            <div class="flex border-b border-dark-800">
                <button @click="activeTab = 'book'" :class="[
                    'flex-1 px-4 py-3 text-sm font-medium transition-colors',
                    activeTab === 'book'
                        ? 'text-accent-blue border-b-2 border-accent-blue'
                        : 'text-dark-400 hover:text-white'
                ]">
                    For This Book
                    <span v-if="bookAnnotations.length > 0" class="ml-1 text-xs">({{ bookAnnotations.length }})</span>
                </button>
                <button @click="activeTab = 'unsorted'" :class="[
                    'flex-1 px-4 py-3 text-sm font-medium transition-colors',
                    activeTab === 'unsorted'
                        ? 'text-accent-blue border-b-2 border-accent-blue'
                        : 'text-dark-400 hover:text-white'
                ]">
                    Unsorted
                    <span v-if="unsortedAnnotations.length > 0" class="ml-1 text-xs">({{ unsortedAnnotations.length
                    }})</span>
                </button>
            </div>

            <!-- Search/Filter -->
            <div class="p-4 border-b border-dark-800">
                <div class="relative">
                    <input v-model="searchQuery" type="text" placeholder="Search annotations..."
                        class="w-full pl-10 pr-4 py-2 text-sm text-white bg-dark-800 border border-dark-600 rounded-lg focus:ring-2 focus:ring-accent-red/50 focus:border-accent-red" />
                    <svg class="absolute left-3 top-2.5 w-5 h-5 text-dark-400" fill="none" stroke="currentColor"
                        viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </div>

                <!-- Type Filter (only for book tab) -->
                <div v-if="activeTab === 'book'" class="flex space-x-2 mt-3">
                    <button @click="typeFilter = 'all'" :class="[
                        'px-3 py-1.5 text-xs font-medium rounded-md transition-colors',
                        typeFilter === 'all'
                            ? 'bg-accent-blue/20 text-accent-blue'
                            : 'bg-dark-800 text-dark-400 hover:bg-dark-700 hover:text-white'
                    ]">
                        All
                    </button>
                    <button @click="typeFilter = 'note'" :class="[
                        'px-3 py-1.5 text-xs font-medium rounded-md transition-colors',
                        typeFilter === 'note'
                            ? 'bg-accent-blue/20 text-accent-blue'
                            : 'bg-dark-800 text-dark-400 hover:bg-dark-700 hover:text-white'
                    ]">
                        Notes
                    </button>
                    <button @click="typeFilter = 'highlight'" :class="[
                        'px-3 py-1.5 text-xs font-medium rounded-md transition-colors',
                        typeFilter === 'highlight'
                            ? 'bg-accent-blue/20 text-accent-blue'
                            : 'bg-dark-800 text-dark-400 hover:bg-dark-700 hover:text-white'
                    ]">
                        Highlights
                    </button>
                </div>
            </div>

            <!-- Annotations List -->
            <div class="flex-1 overflow-y-auto p-4 space-y-3">
                <div v-if="isLoading" class="flex items-center justify-center py-12">
                    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-accent-red"></div>
                </div>

                <div v-else-if="filteredAnnotations.length === 0" class="text-center py-12">
                    <svg class="mx-auto h-12 w-12 text-dark-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <p class="mt-2 text-sm text-dark-400">
                        {{ activeTab === 'book' ? 'No annotations for this book yet' : 'No unsorted notes' }}
                    </p>
                </div>

                <div v-else v-for="annotation in filteredAnnotations" :key="annotation.id"
                    class="group relative p-4 bg-dark-800/50 border border-dark-700 rounded-lg hover:bg-dark-800 hover:border-dark-600 transition-all"
                    :class="{ 'animate-fadeOut': removingIds.includes(annotation.id) }">
                    <!-- Type Badge -->
                    <div class="flex items-start justify-between mb-2">
                        <span :class="[
                            'inline-flex items-center px-2 py-0.5 text-xs font-medium rounded',
                            annotation.type === 'highlight'
                                ? 'bg-yellow-500/20 text-yellow-400'
                                : 'bg-accent-blue/20 text-accent-blue'
                        ]">
                            {{ annotation.type === 'highlight' ? '✨ Highlight' : '📝 Note' }}
                        </span>
                        <span v-if="annotation.page_number" class="text-xs text-dark-400">
                            p. {{ annotation.page_number }}
                        </span>
                    </div>

                    <!-- Content -->
                    <p class="text-sm text-white whitespace-pre-wrap mb-3">
                        {{ annotation.content }}
                    </p>

                    <!-- Tags -->
                    <div v-if="annotation.tags && annotation.tags.length > 0" class="flex flex-wrap gap-1 mb-3">
                        <span v-for="tag in annotation.tags" :key="tag"
                            class="inline-flex items-center px-2 py-0.5 text-xs text-accent-blue bg-accent-blue/10 rounded">
                            #{{ tag }}
                        </span>
                    </div>

                    <!-- Actions -->
                    <div class="flex items-center justify-between">
                        <span class="text-xs text-dark-400">
                            {{ formatDate(annotation.created_at) }}
                        </span>
                        <button @click="injectAnnotation(annotation)"
                            class="flex items-center space-x-1 px-3 py-1.5 text-xs font-medium text-accent-blue bg-accent-blue/10 rounded-md hover:bg-accent-blue/20 transition-colors opacity-0 group-hover:opacity-100">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M12 4v16m8-8H4" />
                            </svg>
                            <span>Inject</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useToastStore } from '../stores/toast'

const props = defineProps({
    isOpen: {
        type: Boolean,
        required: true
    },
    bookId: {
        type: String,
        default: null
    }
})

const emit = defineEmits(['close', 'inject'])

const toast = useToastStore()

// State
const activeTab = ref('book')
const searchQuery = ref('')
const typeFilter = ref('all')
const bookAnnotations = ref([])
const unsortedAnnotations = ref([])
const isLoading = ref(false)
const removingIds = ref([])

// Computed
const filteredAnnotations = computed(() => {
    let annotations = activeTab.value === 'book' ? bookAnnotations.value : unsortedAnnotations.value

    // Apply type filter (only for book tab)
    if (activeTab.value === 'book' && typeFilter.value !== 'all') {
        annotations = annotations.filter(a => a.type === typeFilter.value)
    }

    // Apply search filter
    if (searchQuery.value.trim()) {
        const query = searchQuery.value.toLowerCase()
        annotations = annotations.filter(a =>
            a.content.toLowerCase().includes(query) ||
            (a.tags && a.tags.some(tag => tag.toLowerCase().includes(query)))
        )
    }

    return annotations
})

// Watch for tab changes and book changes
watch(() => props.isOpen, (newVal) => {
    if (newVal) {
        fetchAnnotations()
    }
})

watch(() => props.bookId, () => {
    if (props.isOpen) {
        fetchAnnotations()
    }
})

watch(activeTab, () => {
    searchQuery.value = ''
    typeFilter.value = 'all'
})

// Fetch annotations
const fetchAnnotations = async () => {
    if (!props.bookId) return

    isLoading.value = true

    try {
        const token = localStorage.getItem('token')

        // Fetch book annotations
        const bookResponse = await fetch(
            `${import.meta.env.VITE_API_URL}/api/books/${props.bookId}/annotations`,
            {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            }
        )

        if (bookResponse.ok) {
            const bookData = await bookResponse.json()
            bookAnnotations.value = bookData.annotations || []
        }

        // Fetch unsorted annotations
        const unsortedResponse = await fetch(
            `${import.meta.env.VITE_API_URL}/api/annotations/unassociated`,
            {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            }
        )

        if (unsortedResponse.ok) {
            const unsortedData = await unsortedResponse.json()
            unsortedAnnotations.value = unsortedData.annotations || []
        }
    } catch (error) {
        console.error('Failed to fetch annotations:', error)
        toast.error('Failed to load annotations')
    } finally {
        isLoading.value = false
    }
}

// Inject annotation into parent
const injectAnnotation = async (annotation) => {
    // If this is an unsorted annotation, associate it with the current book
    if (activeTab.value === 'unsorted' && props.bookId) {
        try {
            const token = localStorage.getItem('token')
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/api/annotations/${annotation.id}`,
                {
                    method: 'PATCH',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify({
                        book_id: props.bookId
                    })
                }
            )

            if (response.ok) {
                // Animate removal from unsorted list
                removingIds.value.push(annotation.id)

                setTimeout(() => {
                    // Remove from unsorted and add to book annotations
                    unsortedAnnotations.value = unsortedAnnotations.value.filter(a => a.id !== annotation.id)
                    bookAnnotations.value.push({ ...annotation, is_associated: true })
                    removingIds.value = removingIds.value.filter(id => id !== annotation.id)

                    toast.success('Note associated with book')
                }, 300)
            }
        } catch (error) {
            console.error('Failed to associate annotation:', error)
            toast.error('Failed to associate note')
            return
        }
    }

    // Format content based on type
    let formattedContent = annotation.content
    if (annotation.type === 'highlight') {
        formattedContent = `"${annotation.content}"`
    }

    // Add page reference if available
    if (annotation.page_number) {
        formattedContent += ` (p. ${annotation.page_number})`
    }

    emit('inject', formattedContent)
}

// Format date
const formatDate = (dateString) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now - date
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMs / 3600000)
    const diffDays = Math.floor(diffMs / 86400000)

    if (diffMins < 1) return 'Just now'
    if (diffMins < 60) return `${diffMins}m ago`
    if (diffHours < 24) return `${diffHours}h ago`
    if (diffDays < 7) return `${diffDays}d ago`

    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const close = () => {
    emit('close')
}
</script>

<style scoped>
.slide-enter-active,
.slide-leave-active {
    transition: transform 0.3s ease;
}

.slide-enter-from {
    transform: translateX(100%);
}

.slide-leave-to {
    transform: translateX(100%);
}

@keyframes fadeOut {
    from {
        opacity: 1;
        transform: translateX(0);
    }

    to {
        opacity: 0;
        transform: translateX(-20px);
    }
}

.animate-fadeOut {
    animation: fadeOut 0.3s ease forwards;
}
</style>
