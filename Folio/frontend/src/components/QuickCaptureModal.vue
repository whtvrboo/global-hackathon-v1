<template>
    <Transition name="modal">
        <div v-if="isOpen"
            class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
            @click.self="close">
            <div class="w-full max-w-lg bg-dark-900/95 backdrop-blur-md rounded-t-3xl sm:rounded-2xl shadow-2xl border border-dark-800 transform transition-all duration-300 ease-out"
                :class="{ 'translate-y-0': isOpen, 'translate-y-full': !isOpen }">
                <!-- Header -->
                <div class="flex items-center justify-between p-4 border-b border-dark-800">
                    <h3 class="text-heading-3">Quick Capture</h3>
                    <button @click="close" class="text-dark-400 hover:text-white transition-colors">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <!-- Content -->
                <div v-if="!showBookPicker" class="p-6 space-y-4">
                    <!-- Book Selection Display -->
                    <div v-if="selectedBook"
                        class="flex items-center justify-between p-3 bg-accent-blue/10 border border-accent-blue/20 rounded-lg">
                        <div class="flex items-center space-x-3">
                            <svg class="w-5 h-5 text-accent-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                            </svg>
                            <div>
                                <p class="text-sm font-medium text-white">Note for:</p>
                                <p class="text-sm text-accent-blue font-semibold">{{
                                    selectedBook.title }}</p>
                            </div>
                        </div>
                        <button @click="showBookPicker = true"
                            class="text-sm text-accent-blue hover:text-accent-blue/80 font-medium transition-colors">
                            Change
                        </button>
                    </div>

                    <div v-else
                        class="flex items-center justify-between p-3 bg-dark-800/50 border border-dark-700 rounded-lg">
                        <div class="flex items-center space-x-3">
                            <svg class="w-5 h-5 text-dark-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                            <p class="text-sm text-dark-300">General Note</p>
                        </div>
                        <button @click="showBookPicker = true"
                            class="text-sm text-accent-blue hover:text-accent-blue/80 font-medium transition-colors">
                            Select Book
                        </button>
                    </div>

                    <!-- Note Type Toggle -->
                    <div class="flex space-x-2 p-1 bg-dark-800 rounded-lg">
                        <button @click="noteType = 'note'" :class="[
                            'flex-1 px-4 py-2 text-sm font-medium rounded-md transition-all',
                            noteType === 'note'
                                ? 'bg-dark-700 text-white shadow-sm'
                                : 'text-dark-400 hover:text-white'
                        ]">
                            Note
                        </button>
                        <button @click="noteType = 'highlight'" :class="[
                            'flex-1 px-4 py-2 text-sm font-medium rounded-md transition-all',
                            noteType === 'highlight'
                                ? 'bg-dark-700 text-white shadow-sm'
                                : 'text-dark-400 hover:text-white'
                        ]">
                            Highlight
                        </button>
                    </div>

                    <!-- Content Textarea -->
                    <div>
                        <textarea ref="contentInput" v-model="content"
                            :placeholder="noteType === 'note' ? 'What\'s on your mind?' : 'Paste or type the highlighted text...'"
                            class="w-full h-32 px-4 py-3 text-white bg-dark-800 border border-dark-600 rounded-lg focus:ring-2 focus:ring-accent-red/50 focus:border-accent-red resize-none"
                            @input="detectPageNumber"></textarea>
                        <p class="mt-1 text-xs text-dark-400">
                            Tip: Type "p.42" or "pg 42" to auto-detect page numbers
                        </p>
                    </div>

                    <!-- Page Number (auto-detected or manual) -->
                    <div v-if="pageNumber !== null"
                        class="flex items-center space-x-2 p-3 bg-accent-green/10 border border-accent-green/20 rounded-lg">
                        <svg class="w-5 h-5 text-accent-green" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span class="text-sm text-accent-green">Page {{ pageNumber }} detected</span>
                        <button @click="pageNumber = null"
                            class="ml-auto text-xs text-accent-green hover:text-accent-green/80">
                            Clear
                        </button>
                    </div>

                    <!-- Tags Input -->
                    <div>
                        <input v-model="tagInput" @keydown.enter.prevent="addTag" @input="highlightTags"
                            placeholder="Add tags (press Enter or use #tag)"
                            class="w-full px-4 py-2 text-sm text-white bg-dark-800 border border-dark-600 rounded-lg focus:ring-2 focus:ring-accent-red/50 focus:border-accent-red" />
                        <div v-if="tags.length > 0" class="flex flex-wrap gap-2 mt-2">
                            <span v-for="(tag, index) in tags" :key="index"
                                class="inline-flex items-center px-3 py-1 text-xs font-medium text-accent-blue bg-accent-blue/10 rounded-full">
                                {{ tag }}
                                <button @click="removeTag(index)"
                                    class="ml-1.5 text-accent-blue hover:text-accent-blue/80">
                                    <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd"
                                            d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                                            clip-rule="evenodd" />
                                    </svg>
                                </button>
                            </span>
                        </div>
                    </div>

                    <!-- Offline Indicator -->
                    <div v-if="!isOnline"
                        class="flex items-center space-x-2 p-3 bg-yellow-500/10 border border-yellow-500/20 rounded-lg">
                        <svg class="w-5 h-5 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                        <span class="text-sm text-yellow-500">Offline - will sync when
                            connected</span>
                    </div>

                    <!-- Action Buttons -->
                    <div class="flex space-x-3 pt-2">
                        <button @click="close"
                            class="flex-1 px-4 py-2.5 text-sm font-medium text-dark-300 bg-dark-800 border border-dark-600 rounded-lg hover:bg-dark-700 transition-colors">
                            Cancel
                        </button>
                        <button @click="saveAnnotation" :disabled="!content.trim() || isSaving"
                            class="flex-1 px-4 py-2.5 text-sm font-medium text-white bg-accent-red rounded-lg hover:bg-accent-red/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
                            {{ isSaving ? 'Saving...' : 'Save' }}
                        </button>
                    </div>
                </div>

                <!-- Book Picker View -->
                <div v-else class="p-6 space-y-4">
                    <button @click="showBookPicker = false"
                        class="flex items-center text-sm text-dark-400 hover:text-white mb-4">
                        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                        </svg>
                        Back
                    </button>

                    <!-- Special Options -->
                    <button @click="selectBook(null)"
                        class="w-full flex items-center p-4 text-left bg-dark-800/50 hover:bg-dark-800 border border-dark-700 rounded-lg transition-colors">
                        <svg class="w-6 h-6 text-dark-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                        </svg>
                        <div>
                            <p class="font-medium text-white">Save as general note</p>
                            <p class="text-xs text-dark-400">Not associated with any book</p>
                        </div>
                    </button>

                    <!-- Recent Books -->
                    <div v-if="recentBooks.length > 0" class="space-y-2">
                        <h4 class="text-xs font-semibold text-dark-400 uppercase tracking-wider">
                            Recent Books</h4>
                        <button v-for="book in recentBooks" :key="book.id" @click="selectBook(book)"
                            class="w-full flex items-center p-3 text-left bg-dark-800 hover:bg-accent-blue/10 border border-dark-700 rounded-lg transition-colors"
                            :class="{ 'ring-2 ring-accent-blue': selectedBook?.id === book.id }">
                            <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                                class="w-10 h-14 object-cover rounded mr-3" />
                            <div class="flex-1 min-w-0">
                                <p class="font-medium text-white truncate">{{ book.title }}</p>
                                <p class="text-xs text-dark-400 capitalize">{{
                                    book.status.replace('_', ' ') }}</p>
                            </div>
                            <span v-if="book.status === 'reading'"
                                class="ml-2 px-2 py-1 text-xs font-medium text-accent-green bg-accent-green/10 rounded">
                                Reading
                            </span>
                        </button>
                    </div>

                    <div v-else class="text-center py-8 text-dark-400">
                        <p>No recent books found</p>
                        <p class="text-sm mt-1">Start reading a book to see it here</p>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useToastStore } from '../stores/toast'

const props = defineProps({
    isOpen: {
        type: Boolean,
        required: true
    }
})

const emit = defineEmits(['close', 'saved'])

const toast = useToastStore()

// State
const content = ref('')
const noteType = ref('note')
const selectedBook = ref(null)
const pageNumber = ref(null)
const tags = ref([])
const tagInput = ref('')
const recentBooks = ref([])
const showBookPicker = ref(false)
const isSaving = ref(false)
const isOnline = ref(navigator.onLine)
const contentInput = ref(null)

// Auto-focus textarea when modal opens
watch(() => props.isOpen, async (newVal) => {
    if (newVal) {
        await nextTick()
        contentInput.value?.focus()
        fetchRecentBooks()
    } else {
        // Reset form
        content.value = ''
        noteType.value = 'note'
        pageNumber.value = null
        tags.value = []
        tagInput.value = ''
        showBookPicker.value = false
    }
})

// Online/Offline detection
const updateOnlineStatus = () => {
    isOnline.value = navigator.onLine
}

onMounted(() => {
    window.addEventListener('online', updateOnlineStatus)
    window.addEventListener('offline', updateOnlineStatus)
})

onUnmounted(() => {
    window.removeEventListener('online', updateOnlineStatus)
    window.removeEventListener('offline', updateOnlineStatus)
})

// Fetch recent books
const fetchRecentBooks = async () => {
    try {
        const token = localStorage.getItem('token')
        const response = await fetch('/api/users/me/recents', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })

        if (response.ok) {
            const data = await response.json()
            recentBooks.value = data.books || []

            // Auto-select the first "reading" book
            const readingBook = recentBooks.value.find(b => b.status === 'reading')
            if (readingBook) {
                selectedBook.value = readingBook
            }
        }
    } catch (error) {
        console.error('Failed to fetch recent books:', error)
    }
}

// Page number detection
const detectPageNumber = () => {
    const pageRegex = /(p|pg)\.?\s*(\d+)/i
    const match = content.value.match(pageRegex)
    if (match) {
        pageNumber.value = parseInt(match[2])
    }
}

// Tag handling
const highlightTags = () => {
    const hashtagRegex = /#(\w+)/g
    const matches = tagInput.value.match(hashtagRegex)
    if (matches) {
        matches.forEach(match => {
            const tag = match.substring(1) // Remove #
            if (!tags.value.includes(tag)) {
                tags.value.push(tag)
            }
        })
        tagInput.value = tagInput.value.replace(hashtagRegex, '').trim()
    }
}

const addTag = () => {
    const tag = tagInput.value.trim().replace(/^#/, '')
    if (tag && !tags.value.includes(tag)) {
        tags.value.push(tag)
        tagInput.value = ''
    }
}

const removeTag = (index) => {
    tags.value.splice(index, 1)
}

// Book selection
const selectBook = (book) => {
    selectedBook.value = book
    showBookPicker.value = false
}

// Save annotation
const saveAnnotation = async () => {
    if (!content.value.trim()) return

    isSaving.value = true

    try {
        const token = localStorage.getItem('token')
        const payload = {
            content: content.value.trim(),
            type: noteType.value,
            book_id: selectedBook.value?.id || null,
            page_number: pageNumber.value,
            tags: tags.value
        }

        const response = await fetch('/api/annotations/capture', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        })

        if (response.ok) {
            const data = await response.json()

            if (!isOnline.value) {
                toast.info('Saved offline - will sync when connected')
            } else {
                toast.success('Annotation saved successfully!')
            }

            emit('saved', data)
            close()
        } else {
            throw new Error('Failed to save annotation')
        }
    } catch (error) {
        console.error('Error saving annotation:', error)
        toast.error('Failed to save annotation')
    } finally {
        isSaving.value = false
    }
}

const close = () => {
    emit('close')
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}

.modal-enter-active .modal-content,
.modal-leave-active .modal-content {
    transition: transform 0.3s ease;
}

.modal-enter-from .modal-content {
    transform: translateY(100%);
}

.modal-leave-to .modal-content {
    transform: translateY(100%);
}

@media (min-width: 640px) {

    .modal-enter-from .modal-content,
    .modal-leave-to .modal-content {
        transform: scale(0.95);
    }
}
</style>
