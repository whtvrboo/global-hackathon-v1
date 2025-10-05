<template>
    <div class="min-h-screen bg-dark-950">
        <!-- Loading State -->
        <div v-if="loading" class="flex items-center justify-center min-h-screen">
            <div class="text-center">
                <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-accent-red"></div>
                <p class="mt-4 text-dark-300">Loading book details...</p>
            </div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="flex items-center justify-center min-h-screen">
            <div class="text-center">
                <div class="text-4xl mb-4">😕</div>
                <p class="text-dark-300 mb-4">{{ error }}</p>
                <PrimaryButton @click="$router.push('/discover')">
                    Back to Discover
                </PrimaryButton>
            </div>
        </div>

        <!-- Book Details -->
        <div v-else-if="bookDetails" class="max-w-6xl mx-auto px-4 py-8">
            <!-- Breadcrumb -->
            <nav class="mb-6">
                <ol class="flex items-center space-x-2 text-sm text-dark-400">
                    <li><router-link to="/discover"
                            class="hover:text-accent-red transition-colors">Discover</router-link></li>
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

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <!-- Main Content -->
                <div class="lg:col-span-2">
                    <!-- Book Header -->
                    <div class="card p-8 mb-6">
                        <div class="flex gap-6">
                            <!-- Book Cover -->
                            <div class="flex-shrink-0">
                                <img v-if="bookDetails.cover_url" :src="bookDetails.cover_url" :alt="bookDetails.title"
                                    class="w-48 h-72 object-cover rounded-lg shadow-lg" />
                                <div v-else
                                    class="w-48 h-72 bg-dark-800 rounded-lg flex items-center justify-center text-6xl text-dark-400">
                                    📚
                                </div>
                            </div>

                            <!-- Book Info -->
                            <div class="flex-1 min-w-0">
                                <h1 class="text-3xl font-bold text-white mb-2">
                                    {{ bookDetails.title }}
                                </h1>
                                <p v-if="bookDetails.authors?.length" class="text-lg text-dark-300 mb-4">
                                    by {{ bookDetails.authors.join(', ') }}
                                </p>

                                <!-- Metadata -->
                                <div class="flex flex-wrap gap-4 text-sm text-dark-300 mb-6">
                                    <div v-if="bookDetails.published_date">
                                        <span class="font-medium">Published:</span> {{ bookDetails.published_date }}
                                    </div>
                                    <div v-if="bookDetails.page_count">
                                        <span class="font-medium">Pages:</span> {{ bookDetails.page_count }}
                                    </div>
                                    <div v-if="bookDetails.publisher">
                                        <span class="font-medium">Publisher:</span> {{ bookDetails.publisher }}
                                    </div>
                                    <div v-if="bookDetails.language">
                                        <span class="font-medium">Language:</span> {{ bookDetails.language }}
                                    </div>
                                </div>

                                <!-- Rating -->
                                <div v-if="bookDetails.rating" class="flex items-center gap-2 mb-6">
                                    <div class="flex">
                                        <span v-for="i in 5" :key="i" class="text-2xl"
                                            :class="i <= Math.round(bookDetails.rating) ? 'text-accent-orange' : 'text-dark-400'">
                                            ★
                                        </span>
                                    </div>
                                    <span class="text-lg font-semibold">{{ bookDetails.rating }}</span>
                                    <span v-if="bookDetails.ratings_count" class="text-sm text-dark-400">
                                        ({{ bookDetails.ratings_count }} ratings)
                                    </span>
                                </div>

                                <!-- Categories -->
                                <div v-if="bookDetails.categories?.length" class="flex flex-wrap gap-2 mb-6">
                                    <span v-for="category in bookDetails.categories" :key="category"
                                        class="px-3 py-1 text-sm bg-dark-800 text-dark-200 rounded-full">
                                        {{ category }}
                                    </span>
                                </div>

                                <!-- Action Buttons -->
                                <div class="flex flex-wrap gap-3">
                                    <PrimaryButton @click="showLogModal = true">
                                        Log This Book
                                    </PrimaryButton>
                                    <OutlineButton @click="showAddToList = true">
                                        Add to List
                                    </OutlineButton>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Description -->
                    <div v-if="bookDetails.description" class="card rounded-2xl shadow-sm p-8 mb-6">
                        <h2 class="text-xl font-semibold mb-4">Description</h2>
                        <p class="text-dark-200 leading-relaxed whitespace-pre-line">
                            {{ bookDetails.description }}
                        </p>
                    </div>

                    <!-- Public Reviews Section -->
                    <div class="card rounded-2xl shadow-sm p-8 mb-6">
                        <div class="flex items-center justify-between mb-6">
                            <h2 class="text-xl font-semibold">Community Reviews</h2>
                            <div class="text-sm text-dark-400">
                                {{ publicReviews.length }} reviews
                            </div>
                        </div>

                        <!-- Reviews Loading -->
                        <div v-if="loadingReviews" class="text-center py-8">
                            <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-accent-red">
                            </div>
                            <p class="mt-2 text-sm text-dark-400">Loading reviews...</p>
                        </div>

                        <!-- Reviews List -->
                        <div v-else-if="publicReviews.length > 0" class="space-y-6">
                            <div v-for="review in publicReviews" :key="review.id"
                                class="border-b border-dark-800 pb-6 last:border-b-0">
                                <div class="flex items-start gap-4">
                                    <img v-if="review.user.picture" :src="review.user.picture" :alt="review.user.name"
                                        class="w-10 h-10 rounded-full object-cover" />
                                    <div v-else
                                        class="w-10 h-10 bg-dark-800 rounded-full flex items-center justify-center text-dark-400">
                                        {{ review.user.name.charAt(0).toUpperCase() }}
                                    </div>

                                    <div class="flex-1">
                                        <div class="flex items-center gap-2 mb-2">
                                            <h4 class="font-medium">{{ review.user.name }}</h4>
                                            <span class="text-sm text-dark-400">@{{ review.user.username }}</span>
                                            <span class="text-sm text-dark-400">•</span>
                                            <span class="text-sm text-dark-400">{{ formatDate(review.created_at)
                                            }}</span>
                                        </div>

                                        <!-- Rating -->
                                        <div v-if="review.rating" class="flex items-center gap-1 mb-3">
                                            <div class="flex">
                                                <span v-for="i in 5" :key="i" class="text-sm"
                                                    :class="i <= review.rating ? 'text-accent-orange' : 'text-dark-400'">
                                                    ★
                                                </span>
                                            </div>
                                            <span class="text-sm text-dark-300">{{ getStatusText(review.status)
                                            }}</span>
                                        </div>

                                        <!-- Review Content -->
                                        <div v-if="review.review" class="text-dark-200 mb-2">
                                            <p class="whitespace-pre-line">{{ review.review }}</p>
                                        </div>

                                        <!-- Notes -->
                                        <div v-if="review.notes" class="text-sm text-dark-300">
                                            <p class="whitespace-pre-line">{{ review.notes }}</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- No Reviews -->
                        <div v-else class="text-center py-8 text-dark-400">
                            <div class="text-4xl mb-2">📝</div>
                            <p>No public reviews yet</p>
                            <p class="text-sm">Be the first to share your thoughts!</p>
                        </div>
                    </div>

                    <!-- Related Books -->
                    <div v-if="relatedBooks.length > 0" class="card rounded-2xl shadow-sm p-8">
                        <h2 class="text-xl font-semibold mb-6">Related Books</h2>
                        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                            <div v-for="book in relatedBooks" :key="book.id" @click="$router.push(`/books/${book.id}`)"
                                class="cursor-pointer group">
                                <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                                    class="w-full h-48 object-cover rounded-lg shadow-sm group-hover:shadow-md transition-shadow" />
                                <div v-else
                                    class="w-full h-48 bg-dark-800 rounded-lg flex items-center justify-center text-2xl text-dark-400">
                                    📚
                                </div>
                                <h3 class="mt-2 text-sm font-medium text-white line-clamp-2">{{ book.title }}</h3>
                                <p v-if="book.authors?.length" class="text-xs text-dark-400 line-clamp-1">
                                    {{ book.authors[0] }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Sidebar -->
                <div class="lg:col-span-1">
                    <!-- ISBN Info -->
                    <div v-if="bookDetails.isbn_10 || bookDetails.isbn_13" class="card rounded-2xl shadow-sm p-6 mb-6">
                        <h3 class="text-lg font-semibold mb-4">Book Details</h3>
                        <div class="space-y-2 text-sm">
                            <div v-if="bookDetails.isbn_10">
                                <span class="font-medium text-dark-300">ISBN-10:</span>
                                <span class="ml-2 font-mono">{{ bookDetails.isbn_10 }}</span>
                            </div>
                            <div v-if="bookDetails.isbn_13">
                                <span class="font-medium text-dark-300">ISBN-13:</span>
                                <span class="ml-2 font-mono">{{ bookDetails.isbn_13 }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- Community Stats -->
                    <div class="card rounded-2xl shadow-sm p-6 mb-6">
                        <h3 class="text-lg font-semibold mb-4">Community Stats</h3>
                        <div class="space-y-4">
                            <div class="flex justify-between items-center">
                                <span class="text-sm text-dark-300">Want to Read</span>
                                <span class="font-semibold">{{ communityStats.want_to_read || 0 }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-sm text-dark-300">Currently Reading</span>
                                <span class="font-semibold">{{ communityStats.reading || 0 }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-sm text-dark-300">Have Read</span>
                                <span class="font-semibold">{{ communityStats.read || 0 }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-sm text-dark-300">Average Rating</span>
                                <span class="font-semibold">{{ communityStats.avg_rating ?
                                    communityStats.avg_rating.toFixed(1) : 'N/A' }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- Lists containing this book -->
                    <div v-if="listsWithBook.length > 0" class="card rounded-2xl shadow-sm p-6">
                        <h3 class="text-lg font-semibold mb-4">Featured in Lists</h3>
                        <div class="space-y-3">
                            <div v-for="list in listsWithBook" :key="list.id" @click="$router.push(`/lists/${list.id}`)"
                                class="cursor-pointer p-3 bg-dark-800 rounded-lg hover:bg-dark-800 transition-colors">
                                <h4 class="font-medium text-sm">{{ list.name }}</h4>
                                <p v-if="list.description" class="text-xs text-dark-400 line-clamp-1">{{
                                    list.description }}</p>
                                <p class="text-xs text-dark-400 mt-1">by {{ list.creator.username }}</p>
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
                                d="M6 18L18 6M6 6l12 12"></path>
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

        <!-- List Creation Modal -->
        <ListModal :show="showCreateListModal" :book="bookDetails" @close="showCreateListModal = false"
            @success="handleListCreated" />
    </div>
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
const loading = ref(false)
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
