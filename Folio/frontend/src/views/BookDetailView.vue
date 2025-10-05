<template>
    <div v-if="loading" class="text-center p-12">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
        <p class="mt-4 text-dark-300">Loading book details...</p>
    </div>
    <div v-else-if="error" class="text-center p-12">
        <div class="text-4xl mb-4">😕</div>
        <p class="text-dark-300">{{ error }}</p>
    </div>
    <div v-else-if="bookDetails" class="bg-dark-950 min-h-screen lg:pl-20">
        <!-- Book Hero Section -->
        <div
            class="relative w-full h-[40vh] text-white flex items-center justify-center text-center p-6 bg-gradient-to-br from-accent-blue via-accent-purple to-accent-red">
            <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/50 to-transparent"></div>
            <div class="relative z-10 max-w-6xl mx-auto">
                <div class="flex flex-col lg:flex-row items-center gap-8">
                    <!-- Book Cover -->
                    <div class="flex-shrink-0">
                        <img v-if="bookDetails.cover_url" :src="bookDetails.cover_url" :alt="bookDetails.title"
                            class="w-48 h-72 object-cover rounded-xl shadow-2xl border-2 border-white/20" />
                        <div v-else
                            class="w-48 h-72 bg-white/20 rounded-xl flex items-center justify-center text-6xl border-2 border-white/20 shadow-2xl">
                            📚
                        </div>
                    </div>

                    <!-- Book Info -->
                    <div class="flex-1 text-left lg:text-center">
                        <h1 class="text-3xl md:text-5xl font-bold leading-tight mb-3">
                            {{ bookDetails.title }}
                        </h1>
                        <p v-if="bookDetails.authors?.length" class="text-lg md:text-xl text-white/90 mb-4">
                            by {{ bookDetails.authors.join(', ') }}
                        </p>

                        <!-- Rating -->
                        <div v-if="bookDetails.rating"
                            class="flex items-center justify-center lg:justify-start gap-3 mb-4">
                            <div class="flex">
                                <span v-for="i in 5" :key="i" class="text-2xl"
                                    :class="i <= Math.round(bookDetails.rating) ? 'text-accent-orange' : 'text-white/40'">
                                    ★
                                </span>
                            </div>
                            <span class="text-xl font-bold">{{ bookDetails.rating }}</span>
                            <span v-if="bookDetails.ratings_count" class="text-sm text-white/70">
                                ({{ bookDetails.ratings_count }} ratings)
                            </span>
                        </div>

                        <!-- Categories -->
                        <div v-if="bookDetails.categories?.length"
                            class="flex flex-wrap justify-center lg:justify-start gap-2 mb-6">
                            <span v-for="category in bookDetails.categories" :key="category"
                                class="px-3 py-1 text-xs bg-white/20 text-white rounded-full border border-white/30 backdrop-blur-sm">
                                {{ category }}
                            </span>
                        </div>

                        <!-- Action Buttons -->
                        <div class="flex flex-col sm:flex-row gap-3 justify-center lg:justify-start">
                            <button @click="showLogModal = true"
                                class="px-6 py-3 bg-white text-dark-950 rounded-lg font-semibold hover:bg-white/90 transition-all shadow-lg">
                                📖 Log This Book
                            </button>
                            <button @click="showAddToList = true"
                                class="px-6 py-3 bg-white/20 text-white rounded-lg border border-white/30 hover:bg-white/30 transition-all backdrop-blur-sm font-semibold">
                                📋 Add to List
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Book Content -->
        <div class="container mx-auto max-w-7xl py-8 px-4">
            <!-- Breadcrumb -->
            <nav class="mb-6">
                <ol class="flex items-center space-x-2 text-sm text-dark-400">
                    <li><router-link to="/discover"
                            class="hover:text-accent-red transition-colors">Discover</router-link>
                    </li>
                    <li class="flex items-center">
                        <svg class="w-4 h-4 mx-2" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                                d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
                                clip-rule="evenodd"></path>
                        </svg>
                        <span class="truncate text-white">{{ bookDetails.title }}</span>
                    </li>
                </ol>
            </nav>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <!-- Main Content - Reviews Focus -->
                <div class="lg:col-span-2">
                    <!-- Public Reviews Section - Enhanced and Prominent -->
                    <div class="card mb-6">
                        <div class="flex items-center justify-between mb-6">
                            <h2 class="text-heading-2 flex items-center gap-3">
                                <span class="text-3xl">💬</span>
                                <span>Community Reviews</span>
                                <span class="text-lg text-accent-red font-bold">{{ publicReviews.length }}</span>
                            </h2>
                            <button @click="showLogModal = true"
                                class="px-6 py-3 bg-accent-red text-white rounded-lg hover:bg-accent-red/90 transition-colors font-medium shadow-lg">
                                ✍️ Write Review
                            </button>
                        </div>

                        <!-- Reviews Loading -->
                        <div v-if="loadingReviews" class="text-center py-12">
                            <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-accent-red">
                            </div>
                            <p class="mt-3 text-dark-400">Loading reviews...</p>
                        </div>

                        <!-- Reviews List -->
                        <div v-else-if="publicReviews.length > 0" class="space-y-6">
                            <div v-for="review in publicReviews" :key="review.id"
                                class="bg-dark-800 rounded-xl p-6 border border-dark-700 hover:border-dark-600 transition-all duration-200 hover:shadow-lg">
                                <div class="flex items-start gap-4">
                                    <img v-if="review.user.picture" :src="review.user.picture" :alt="review.user.name"
                                        class="w-16 h-16 rounded-full object-cover border-2 border-dark-600" />
                                    <div v-else
                                        class="w-16 h-16 bg-dark-700 rounded-full flex items-center justify-center text-dark-300 border-2 border-dark-600 font-bold text-xl">
                                        {{ review.user.name.charAt(0).toUpperCase() }}
                                    </div>

                                    <div class="flex-1">
                                        <div class="flex items-center gap-3 mb-3">
                                            <h4 class="font-bold text-xl text-white">{{ review.user.name }}</h4>
                                            <span class="text-sm text-dark-400">@{{ review.user.username }}</span>
                                            <span class="text-sm text-dark-400">•</span>
                                            <span class="text-sm text-dark-400">{{ formatDate(review.created_at)
                                                }}</span>
                                        </div>

                                        <!-- Rating & Status -->
                                        <div class="flex items-center gap-3 mb-4">
                                            <div v-if="review.rating" class="flex">
                                                <span v-for="i in 5" :key="i" class="text-2xl"
                                                    :class="i <= review.rating ? 'text-accent-orange' : 'text-dark-500'">
                                                    ★
                                                </span>
                                            </div>
                                            <span
                                                class="px-3 py-1 rounded-full text-xs font-semibold bg-accent-blue/20 text-accent-blue border border-accent-blue/30">
                                                {{ getStatusText(review.status) }}
                                            </span>
                                        </div>

                                        <!-- Review Content -->
                                        <div v-if="review.review" class="text-dark-100 mb-4">
                                            <p class="whitespace-pre-line text-base leading-relaxed">{{ review.review }}
                                            </p>
                                        </div>

                                        <!-- Notes -->
                                        <div v-if="review.notes"
                                            class="text-sm text-dark-300 bg-dark-700 p-4 rounded-lg border border-dark-600">
                                            <p class="whitespace-pre-line">{{ review.notes }}</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- No Reviews -->
                        <div v-else class="text-center py-20 text-dark-400">
                            <div class="text-8xl mb-6">📝</div>
                            <h3 class="text-2xl font-semibold mb-3 text-white">No reviews yet</h3>
                            <p class="mb-8 text-lg">Be the first to share your thoughts about this book!</p>
                            <button @click="showLogModal = true"
                                class="px-8 py-4 bg-accent-red text-white rounded-lg hover:bg-accent-red/90 transition-colors font-medium text-lg shadow-lg">
                                ✍️ Write the First Review
                            </button>
                        </div>
                    </div>

                    <!-- Lists containing this book - Now more prominent -->
                    <div v-if="listsWithBook.length > 0" class="card">
                        <h2 class="text-heading-2 mb-6 flex items-center gap-3">
                            <span class="text-3xl">📋</span>
                            <span>Featured in Lists</span>
                            <span class="text-lg text-accent-blue font-bold">{{ listsWithBook.length }}</span>
                        </h2>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div v-for="list in listsWithBook" :key="list.id" @click="$router.push(`/lists/${list.id}`)"
                                class="cursor-pointer p-4 bg-dark-800 rounded-xl hover:bg-dark-700 transition-all duration-200 border border-dark-700 hover:border-dark-600 hover:shadow-lg group">
                                <div class="flex items-start gap-3">
                                    <div class="flex-shrink-0">
                                        <div
                                            class="w-12 h-12 bg-accent-blue/20 rounded-lg flex items-center justify-center text-accent-blue text-xl">
                                            📚
                                        </div>
                                    </div>
                                    <div class="flex-1 min-w-0">
                                        <h4
                                            class="font-semibold text-white text-base mb-1 group-hover:text-accent-blue transition-colors">
                                            {{ list.name }}</h4>
                                        <p v-if="list.description" class="text-sm text-dark-400 line-clamp-2 mb-2">{{
                                            list.description }}</p>
                                        <div class="flex items-center gap-2 text-xs text-dark-500">
                                            <span>by {{ list.creator.username }}</span>
                                            <span>•</span>
                                            <span>{{ list.items_count || 0 }} books</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Sidebar -->
                <div class="lg:col-span-1">
                    <!-- Community Stats -->
                    <div class="card mb-6">
                        <h3 class="text-heading-4 mb-4 flex items-center gap-2">
                            <span class="text-xl">📊</span>
                            <span>Community Stats</span>
                        </h3>
                        <div class="space-y-4">
                            <div class="flex justify-between items-center">
                                <span class="text-dark-300 text-sm">Want to Read</span>
                                <span class="font-bold text-accent-blue">{{ communityStats.want_to_read || 0 }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-dark-300 text-sm">Currently Reading</span>
                                <span class="font-bold text-accent-green">{{ communityStats.reading || 0 }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-dark-300 text-sm">Have Read</span>
                                <span class="font-bold text-accent-purple">{{ communityStats.read || 0 }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-dark-300 text-sm">Average Rating</span>
                                <span class="font-bold text-accent-orange">{{ communityStats.avg_rating ?
                                    communityStats.avg_rating.toFixed(1) : 'N/A' }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- Book Details - Compact -->
                    <div class="card mb-6">
                        <h3 class="text-heading-4 mb-4 flex items-center gap-2">
                            <span class="text-xl">ℹ️</span>
                            <span>Book Details</span>
                        </h3>
                        <div class="space-y-3">
                            <div v-if="bookDetails.published_date" class="flex justify-between">
                                <span class="font-medium text-dark-300 text-sm">Published:</span>
                                <span class="text-white text-sm">{{ bookDetails.published_date }}</span>
                            </div>
                            <div v-if="bookDetails.page_count" class="flex justify-between">
                                <span class="font-medium text-dark-300 text-sm">Pages:</span>
                                <span class="text-white text-sm">{{ bookDetails.page_count }}</span>
                            </div>
                            <div v-if="bookDetails.publisher" class="flex justify-between">
                                <span class="font-medium text-dark-300 text-sm">Publisher:</span>
                                <span class="text-white text-sm">{{ bookDetails.publisher }}</span>
                            </div>
                            <div v-if="bookDetails.language" class="flex justify-between">
                                <span class="font-medium text-dark-300 text-sm">Language:</span>
                                <span class="text-white text-sm">{{ bookDetails.language }}</span>
                            </div>
                        </div>

                        <!-- Description -->
                        <div v-if="bookDetails.description" class="mt-4 pt-4 border-t border-dark-700">
                            <h4 class="font-medium text-white text-sm mb-2">Description</h4>
                            <p class="text-dark-200 leading-relaxed whitespace-pre-line text-xs line-clamp-4">
                                {{ bookDetails.description }}
                            </p>
                        </div>
                    </div>

                    <!-- ISBN Info -->
                    <div v-if="bookDetails.isbn_10 || bookDetails.isbn_13" class="card mb-6">
                        <h3 class="text-heading-4 mb-4 flex items-center gap-2">
                            <span class="text-xl">📚</span>
                            <span>ISBN</span>
                        </h3>
                        <div class="space-y-3">
                            <div v-if="bookDetails.isbn_10" class="flex justify-between items-center">
                                <span class="font-medium text-dark-300 text-sm">ISBN-10:</span>
                                <span class="font-mono text-white text-sm">{{ bookDetails.isbn_10 }}</span>
                            </div>
                            <div v-if="bookDetails.isbn_13" class="flex justify-between items-center">
                                <span class="font-medium text-dark-300 text-sm">ISBN-13:</span>
                                <span class="font-mono text-white text-sm">{{ bookDetails.isbn_13 }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- Related Books - Moved to sidebar -->
                    <div v-if="relatedBooks.length > 0" class="card">
                        <h3 class="text-heading-4 mb-4 flex items-center gap-2">
                            <span class="text-xl">🔗</span>
                            <span>Related Books</span>
                        </h3>
                        <div class="space-y-3">
                            <div v-for="book in relatedBooks" :key="book.id" @click="$router.push(`/books/${book.id}`)"
                                class="cursor-pointer group">
                                <div class="flex gap-3">
                                    <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                                        class="w-12 h-16 object-cover rounded shadow-md group-hover:shadow-lg transition-all duration-300" />
                                    <div v-else
                                        class="w-12 h-16 bg-dark-800 rounded flex items-center justify-center text-lg text-dark-400 shadow-md">
                                        📚
                                    </div>
                                    <div class="flex-1 min-w-0">
                                        <h4
                                            class="text-sm font-medium text-white line-clamp-2 group-hover:text-accent-blue transition-colors">
                                            {{ book.title }}
                                        </h4>
                                        <p v-if="book.authors?.length" class="text-xs text-dark-400 line-clamp-1">
                                            {{ book.authors[0] }}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Log Book Modal -->
        <LogBookModal :show="showLogModal" :book="bookDetails" @close="showLogModal = false"
            @success="handleLogSuccess" />

        <!-- Add to List Modal -->
        <div v-if="showAddToList" @click="closeAddToListModal"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
            <div @click.stop class="card rounded-2xl shadow-2xl max-w-md w-full p-6">
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-lg font-semibold">{{ selectedListForNotes ? 'Add Your Note' : 'Add to List' }}</h3>
                    <button @click="closeAddToListModal" class="text-dark-400 hover:text-dark-300 transition-colors">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12">
                            </path>
                        </svg>
                    </button>
                </div>

                <!-- View for adding curator notes -->
                <div v-if="selectedListForNotes">
                    <p class="text-sm text-dark-300 mb-4">
                        You're adding <strong>{{ bookDetails.title }}</strong> to <strong>{{ selectedListForNotes.name
                        }}</strong>.
                    </p>
                    <TextArea v-model="curatorNotes"
                        placeholder="Why is this book on the list? What does it mean to you?" :rows="5" />
                    <div class="flex justify-end gap-3 mt-4">
                        <SecondaryButton @click="selectedListForNotes = null">Back</SecondaryButton>
                        <PrimaryButton @click="addBookToList(selectedListForNotes.id)">Add Book</PrimaryButton>
                    </div>
                </div>

                <!-- View for selecting a list -->
                <div v-else>
                    <div v-if="loadingLists" class="text-center py-4">
                        <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-accent-red"></div>
                        <p class="mt-2 text-sm text-dark-400">Loading lists...</p>
                    </div>

                    <div v-else-if="userLists.length > 0" class="space-y-3 max-h-64 overflow-y-auto">
                        <div v-for="list in userLists" :key="list.id" @click="selectListForNotes(list)"
                            class="p-3 bg-dark-800 border border-dark-700 rounded-lg hover:bg-dark-800 cursor-pointer transition-colors">
                            <div class="flex items-center justify-between">
                                <div>
                                    <h4 class="font-medium">{{ list.name }}</h4>
                                    <p v-if="list.description" class="text-sm text-dark-400 line-clamp-1">{{
                                        list.description }}</p>
                                    <p class="text-xs text-dark-400">{{ list.items_count || 0 }} books</p>
                                </div>
                                <div class="text-dark-400">
                                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M12 4v16m8-8H4"></path>
                                    </svg>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div v-else class="text-center py-4">
                        <p class="text-dark-300 mb-4">You don't have any lists yet.</p>
                        <button @click="createNewList" class="btn-primary text-sm">
                            Create Your First List
                        </button>
                    </div>

                    <div class="mt-4 pt-4 border-t border-dark-700">
                        <button @click="createNewList" class="w-full btn-secondary text-sm">
                            Create New List
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- List Creation Modal -->
    <ListModal :show="showCreateListModal" :book="bookDetails" @close="showCreateListModal = false"
        @success="handleListCreated" />
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import PrimaryButton from '../components/ui/PrimaryButton.vue'
import OutlineButton from '../components/ui/OutlineButton.vue'
import SecondaryButton from '../components/ui/SecondaryButton.vue'
import TextArea from '../components/ui/TextArea.vue'
import LogBookModal from '../components/LogBookModal.vue'
import ListModal from '../components/ListModal.vue'
import { useToastStore } from '../stores/toast'

const route = useRoute()
const router = useRouter()
const toastStore = useToastStore()

const bookDetails = ref(null)
const loading = ref(true)
const error = ref(null)
const showLogModal = ref(false)
const showAddToList = ref(false)
const showCreateListModal = ref(false)
const userLists = ref([])
const loadingLists = ref(false)
const selectedListForNotes = ref(null)
const curatorNotes = ref('')
const publicReviews = ref([])
const loadingReviews = ref(false)
const communityStats = ref({})
const relatedBooks = ref([])
const listsWithBook = ref([])

const bookId = ref(route.params.id)

watch(() => route.params.id, (newId) => {
    if (newId) {
        bookId.value = newId
        fetchBookDetails()
    }
})

onMounted(() => {
    if (bookId.value) {
        fetchBookDetails()
    }
})

const fetchBookDetails = async () => {
    loading.value = true
    error.value = null
    bookDetails.value = null

    try {
        const response = await axios.get(`/api/books/${bookId.value}`)
        bookDetails.value = response.data

        // Fetch additional data in parallel
        await Promise.all([
            fetchPublicReviews(),
            fetchCommunityStats(),
            fetchRelatedBooks(),
            fetchListsWithBook()
        ])
    } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load book details'
        console.error('Error fetching book details:', err)
    } finally {
        loading.value = false
    }
}

const fetchPublicReviews = async () => {
    try {
        loadingReviews.value = true
        const response = await axios.get(`/api/books/${bookId.value}/reviews`)
        publicReviews.value = response.data.reviews || []
    } catch (error) {
        console.error('Error loading reviews:', error)
        publicReviews.value = []
    } finally {
        loadingReviews.value = false
    }
}

const fetchCommunityStats = async () => {
    try {
        const response = await axios.get(`/api/books/${bookId.value}/stats`)
        communityStats.value = response.data
    } catch (error) {
        console.error('Error loading community stats:', error)
        communityStats.value = {}
    }
}

const fetchRelatedBooks = async () => {
    try {
        // Get books by same authors or in same categories
        if (bookDetails.value?.authors?.length > 0) {
            const response = await axios.get(`/api/books/search?q=${encodeURIComponent(bookDetails.value.authors[0])}&limit=4`)
            relatedBooks.value = response.data.results?.filter(book => book.id !== bookId.value) || []
        }
    } catch (error) {
        console.error('Error loading related books:', error)
        relatedBooks.value = []
    }
}

const fetchListsWithBook = async () => {
    try {
        const response = await axios.get(`/api/books/${bookId.value}/lists`)
        listsWithBook.value = response.data.lists || []
    } catch (error) {
        console.error('Error loading lists with book:', error)
        listsWithBook.value = []
    }
}

const loadUserLists = async () => {
    try {
        loadingLists.value = true
        const response = await axios.get('/api/me/lists')
        userLists.value = response.data.lists || []
    } catch (error) {
        console.error('Error loading user lists:', error)
        userLists.value = []
    }
}

const addBookToList = async (listId) => {
    try {
        await axios.post(`/api/lists/${listId}/items`, {
            book_id: bookId.value,
            notes: curatorNotes.value || null
        })
        closeAddToListModal()
        toastStore.success('Book added to list successfully!')
    } catch (error) {
        if (error.response?.status === 409) {
            toastStore.warning('This book is already in the list!')
        } else {
            console.error('Error adding book to list:', error)
            toastStore.error('Failed to add book to list')
        }
    }
}

const createNewList = () => {
    showAddToList.value = false
    showCreateListModal.value = true
}

const handleListCreated = () => {
    showCreateListModal.value = false
    showAddToList.value = true
}

const selectListForNotes = (list) => {
    selectedListForNotes.value = list
}

const closeAddToListModal = () => {
    showAddToList.value = false
    selectedListForNotes.value = null
    curatorNotes.value = ''
}

const handleLogSuccess = () => {
    showLogModal.value = false
    // Refresh reviews and stats
    fetchPublicReviews()
    fetchCommunityStats()
}

const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString()
}

const getStatusText = (status) => {
    const statusMap = {
        'want_to_read': 'Want to Read',
        'reading': 'Currently Reading',
        'read': 'Read',
        'dnf': 'Did Not Finish'
    }
    return statusMap[status] || status
}

// Watch for add to list modal to load user lists
watch(() => showAddToList.value, async (newShow) => {
    if (newShow) {
        await loadUserLists()
    }
})
</script>
