<template>
    <div class="min-h-screen bg-dark-950">
        <!-- Header -->
        <div class="bg-dark-900/50 backdrop-blur-sm border-b border-dark-800">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
                <div class="flex items-center justify-between">
                    <div>
                        <h1 class="text-heading-1">My Notebook</h1>
                        <p class="mt-1 text-body text-dark-400">
                            Organize your unassociated notes and thoughts
                        </p>
                    </div>
                    <div class="flex items-center space-x-4">
                        <span class="text-body text-dark-400">
                            {{ unassociatedNotes.length }} unsorted note{{ unassociatedNotes.length !== 1 ? 's' : '' }}
                        </span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Main Content -->
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <!-- Left Column: Unassociated Notes -->
                <div class="space-y-4">
                    <div class="flex items-center justify-between">
                        <h2 class="text-heading-3">Unsorted Notes</h2>
                        <button v-if="unassociatedNotes.length > 0" @click="selectAll"
                            class="text-sm text-accent-blue hover:text-accent-blue/80">
                            {{ selectedNote ? 'Deselect' : 'Select first' }}
                        </button>
                    </div>

                    <!-- Search -->
                    <div class="relative">
                        <input v-model="searchQuery" type="text" placeholder="Search notes..."
                            class="w-full pl-10 pr-4 py-2 text-sm text-white bg-dark-800 border border-dark-600 rounded-lg focus:ring-2 focus:ring-accent-red/50 focus:border-accent-red" />
                        <svg class="absolute left-3 top-2.5 w-5 h-5 text-dark-400" fill="none" stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                    </div>

                    <!-- Loading State -->
                    <div v-if="isLoading" class="flex items-center justify-center py-12">
                        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-accent-red"></div>
                    </div>

                    <!-- Empty State -->
                    <div v-else-if="filteredNotes.length === 0"
                        class="text-center py-12 bg-dark-900/50 backdrop-blur-sm rounded-lg border-2 border-dashed border-dark-700">
                        <svg class="mx-auto h-12 w-12 text-dark-400" fill="none" stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                        </svg>
                        <h3 class="mt-2 text-sm font-medium text-white">No unsorted notes</h3>
                        <p class="mt-1 text-sm text-dark-400">
                            {{ searchQuery ? 'Try a different search term' : 'All your notes are organized!' }}
                        </p>
                    </div>

                    <!-- Notes List -->
                    <div v-else class="space-y-3">
                        <div v-for="note in filteredNotes" :key="note.id" @click="selectNote(note)"
                            class="p-4 bg-dark-900/50 backdrop-blur-sm rounded-lg border-2 cursor-pointer transition-all"
                            :class="[
                                selectedNote?.id === note.id
                                    ? 'border-accent-blue ring-2 ring-accent-blue ring-opacity-50'
                                    : 'border-dark-700 hover:border-accent-blue/50',
                                removingIds.includes(note.id) && 'animate-fadeOut'
                            ]">
                            <!-- Type Badge -->
                            <div class="flex items-start justify-between mb-2">
                                <span :class="[
                                    'inline-flex items-center px-2 py-0.5 text-xs font-medium rounded',
                                    note.type === 'highlight'
                                        ? 'bg-yellow-500/20 text-yellow-400'
                                        : 'bg-accent-blue/20 text-accent-blue'
                                ]">
                                    {{ note.type === 'highlight' ? '✨ Highlight' : '📝 Note' }}
                                </span>
                                <span v-if="note.page_number" class="text-xs text-dark-400">
                                    p. {{ note.page_number }}
                                </span>
                            </div>

                            <!-- Content -->
                            <p class="text-sm text-white line-clamp-3 whitespace-pre-wrap mb-2">
                                {{ note.content }}
                            </p>

                            <!-- Tags -->
                            <div v-if="note.tags && note.tags.length > 0" class="flex flex-wrap gap-1 mb-2">
                                <span v-for="tag in note.tags" :key="tag"
                                    class="inline-flex items-center px-2 py-0.5 text-xs text-accent-blue bg-accent-blue/10 rounded">
                                    #{{ tag }}
                                </span>
                            </div>

                            <!-- Date -->
                            <p class="text-xs text-dark-400">
                                {{ formatDate(note.created_at) }}
                            </p>
                        </div>
                    </div>
                </div>

                <!-- Right Column: Book Picker -->
                <div class="lg:sticky lg:top-8 lg:self-start">
                    <div class="bg-dark-900/50 backdrop-blur-sm rounded-lg border border-dark-800 p-6">
                        <h2 class="text-heading-3 mb-4">
                            {{ selectedNote ? 'Assign to Book' : 'Select a note to organize' }}
                        </h2>

                        <!-- Selected Note Preview -->
                        <div v-if="selectedNote" class="mb-6 p-4 bg-dark-800/50 border border-dark-700 rounded-lg">
                            <p class="text-sm text-white font-medium mb-2">Selected Note:</p>
                            <p class="text-sm text-dark-300 line-clamp-3">
                                {{ selectedNote.content }}
                            </p>
                        </div>

                        <!-- Book Search -->
                        <div v-if="selectedNote" class="mb-4">
                            <div class="relative">
                                <input v-model="bookSearchQuery" @input="searchBooks" type="text"
                                    placeholder="Search your books..."
                                    class="w-full pl-10 pr-4 py-2 text-sm text-white bg-dark-800 border border-dark-600 rounded-lg focus:ring-2 focus:ring-accent-red/50 focus:border-accent-red" />
                                <svg class="absolute left-3 top-2.5 w-5 h-5 text-dark-400" fill="none"
                                    stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                                </svg>
                            </div>
                        </div>

                        <!-- Books List -->
                        <div v-if="selectedNote" class="space-y-2 max-h-96 overflow-y-auto">
                            <div v-if="isLoadingBooks" class="flex items-center justify-center py-8">
                                <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-accent-red"></div>
                            </div>

                            <div v-else-if="filteredBooks.length === 0" class="text-center py-8 text-dark-400">
                                <p class="text-sm">No books found</p>
                            </div>

                            <button v-else v-for="book in filteredBooks" :key="book.id" @click="assignToBook(book)"
                                class="w-full flex items-center p-3 text-left bg-dark-800 hover:bg-accent-blue/10 border border-dark-700 rounded-lg transition-all hover:shadow-md group">
                                <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                                    class="w-12 h-16 object-cover rounded mr-3" />
                                <div class="flex-1 min-w-0">
                                    <p class="font-medium text-white truncate group-hover:text-accent-blue">
                                        {{ book.title }}
                                    </p>
                                    <p class="text-xs text-dark-400 capitalize">
                                        {{ book.status?.replace('_', ' ') || 'Unknown' }}
                                    </p>
                                </div>
                                <svg class="w-5 h-5 text-dark-400 group-hover:text-accent-blue" fill="none"
                                    stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M9 5l7 7-7 7" />
                                </svg>
                            </button>
                        </div>

                        <!-- No Selection State -->
                        <div v-else class="text-center py-12 text-dark-400">
                            <svg class="mx-auto h-12 w-12 text-dark-400 mb-3" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M7 11.5V14m0-2.5v-6a1.5 1.5 0 113 0m-3 6a1.5 1.5 0 00-3 0v2a7.5 7.5 0 0015 0v-5a1.5 1.5 0 00-3 0m-6-3V11m0-5.5v-1a1.5 1.5 0 013 0v1m0 0V11m0-5.5a1.5 1.5 0 013 0v3m0 0V11" />
                            </svg>
                            <p class="text-sm">Select a note from the left to assign it to a book</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()

// State
const unassociatedNotes = ref([])
const selectedNote = ref(null)
const searchQuery = ref('')
const bookSearchQuery = ref('')
const userBooks = ref([])
const isLoading = ref(false)
const isLoadingBooks = ref(false)
const removingIds = ref([])

// Computed
const filteredNotes = computed(() => {
    if (!searchQuery.value.trim()) return unassociatedNotes.value

    const query = searchQuery.value.toLowerCase()
    return unassociatedNotes.value.filter(note =>
        note.content.toLowerCase().includes(query) ||
        (note.tags && note.tags.some(tag => tag.toLowerCase().includes(query)))
    )
})

const filteredBooks = computed(() => {
    if (!bookSearchQuery.value.trim()) return userBooks.value

    const query = bookSearchQuery.value.toLowerCase()
    return userBooks.value.filter(book =>
        book.title.toLowerCase().includes(query) ||
        (book.authors && book.authors.some(author => author.toLowerCase().includes(query)))
    )
})

// Lifecycle
onMounted(() => {
    fetchUnassociatedNotes()
    fetchUserBooks()
})

// Fetch unassociated notes
const fetchUnassociatedNotes = async () => {
    isLoading.value = true

    try {
        const token = localStorage.getItem('token')
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/api/annotations/unassociated`,
            {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            }
        )

        if (response.ok) {
            const data = await response.json()
            unassociatedNotes.value = data.annotations || []
        } else {
            throw new Error('Failed to fetch notes')
        }
    } catch (error) {
        console.error('Error fetching unassociated notes:', error)
        toast.error('Failed to load notes')
    } finally {
        isLoading.value = false
    }
}

// Fetch user's books
const fetchUserBooks = async () => {
    isLoadingBooks.value = true

    try {
        const token = localStorage.getItem('token')
        const username = JSON.parse(localStorage.getItem('user'))?.username

        if (!username) return

        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/api/users/${username}/logs`,
            {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            }
        )

        if (response.ok) {
            const data = await response.json()
            // Extract unique books from logs
            const booksMap = new Map()
            data.logs?.forEach(log => {
                if (log.book && !booksMap.has(log.book.id)) {
                    booksMap.set(log.book.id, {
                        ...log.book,
                        status: log.status
                    })
                }
            })
            userBooks.value = Array.from(booksMap.values())
        }
    } catch (error) {
        console.error('Error fetching user books:', error)
    } finally {
        isLoadingBooks.value = false
    }
}

// Search books (could be enhanced with API search)
const searchBooks = () => {
    // Currently using client-side filtering
    // Could be enhanced with server-side search if needed
}

// Select note
const selectNote = (note) => {
    selectedNote.value = note
    bookSearchQuery.value = ''
}

// Select all (first note)
const selectAll = () => {
    if (selectedNote.value) {
        selectedNote.value = null
    } else if (filteredNotes.value.length > 0) {
        selectedNote.value = filteredNotes.value[0]
    }
}

// Assign note to book
const assignToBook = async (book) => {
    if (!selectedNote.value) return

    try {
        const token = localStorage.getItem('token')
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/api/annotations/${selectedNote.value.id}`,
            {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    book_id: book.id
                })
            }
        )

        if (response.ok) {
            // Animate removal
            removingIds.value.push(selectedNote.value.id)

            setTimeout(() => {
                // Remove from list
                unassociatedNotes.value = unassociatedNotes.value.filter(
                    n => n.id !== selectedNote.value.id
                )
                removingIds.value = removingIds.value.filter(id => id !== selectedNote.value.id)

                toast.success(`Note assigned to "${book.title}"`)

                // Select next note if available
                if (filteredNotes.value.length > 0) {
                    selectedNote.value = filteredNotes.value[0]
                } else {
                    selectedNote.value = null
                }
            }, 300)
        } else {
            throw new Error('Failed to assign note')
        }
    } catch (error) {
        console.error('Error assigning note:', error)
        toast.error('Failed to assign note to book')
    }
}

// Format date
const formatDate = (dateString) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now - date
    const diffDays = Math.floor(diffMs / 86400000)

    if (diffDays === 0) return 'Today'
    if (diffDays === 1) return 'Yesterday'
    if (diffDays < 7) return `${diffDays} days ago`

    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<style scoped>
@keyframes fadeOut {
    from {
        opacity: 1;
        transform: scale(1);
    }

    to {
        opacity: 0;
        transform: scale(0.95);
    }
}

.animate-fadeOut {
    animation: fadeOut 0.3s ease forwards;
}

.line-clamp-3 {
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
</style>
