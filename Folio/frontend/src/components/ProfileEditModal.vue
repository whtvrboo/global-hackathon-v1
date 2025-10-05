<template>
    <transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0"
        enter-to-class="opacity-100" leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100"
        leave-to-class="opacity-0">
        <div v-if="show" @click="$emit('close')"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
            <div @click.stop
                class="bg-gradient-to-br from-dark-900 to-dark-950 border border-dark-700 rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto p-8">

                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-heading-1">Edit Profile</h2>
                    <button @click="$emit('close')"
                        class="p-2 text-dark-400 hover:text-white transition-colors rounded-lg hover:bg-dark-800">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12"></path>
                        </svg>
                    </button>
                </div>

                <form @submit.prevent="handleSubmit" class="space-y-6">
                    <!-- Bio -->
                    <div>
                        <label class="block text-sm font-medium text-dark-200 mb-2">
                            Bio
                        </label>
                        <TextArea v-model="form.bio" placeholder="Tell us about yourself and your reading tastes..."
                            :rows="4" />
                    </div>

                    <!-- Banner URL -->
                    <div>
                        <label class="block text-sm font-medium text-dark-200 mb-2">
                            Banner URL
                        </label>
                        <Input v-model="form.banner_url" placeholder="https://example.com/banner.jpg" />
                        <p class="text-xs text-dark-400 mt-1">
                            Optional: Add a custom banner image to your profile
                        </p>
                    </div>

                    <!-- Favorite Books Selection -->
                    <div>
                        <label class="block text-sm font-medium text-dark-200 mb-2">
                            Favorite Books (select up to 4)
                        </label>
                        <p class="text-xs text-dark-400 mb-3">
                            Choose your all-time favorite reads from books you've logged
                        </p>

                        <!-- Search/Filter -->
                        <div class="mb-4">
                            <Input v-model="searchQuery" placeholder="Search your books..." type="text" />
                        </div>

                        <!-- Selected Favorites -->
                        <div v-if="selectedFavorites.length > 0" class="mb-4">
                            <div class="text-sm text-dark-300 mb-2">Selected ({{ selectedFavorites.length }}/4)</div>
                            <div class="flex gap-2 flex-wrap">
                                <div v-for="book in selectedFavorites" :key="book.id" class="relative group">
                                    <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                                        class="w-16 h-24 object-cover rounded-lg border-2 border-accent-red" />
                                    <button @click="removeFavorite(book.id)" type="button"
                                        class="absolute -top-2 -right-2 w-6 h-6 bg-accent-red rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                                        <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor"
                                            viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M6 18L18 6M6 6l12 12"></path>
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        </div>

                        <!-- Available Books -->
                        <div class="max-h-64 overflow-y-auto glass-strong rounded-lg p-4 space-y-2">
                            <div v-if="loading" class="text-center py-4 text-dark-400">
                                Loading your books...
                            </div>
                            <div v-else-if="filteredAvailableBooks.length === 0" class="text-center py-4 text-dark-400">
                                No books found. Log some books to add them as favorites!
                            </div>
                            <button v-else v-for="book in filteredAvailableBooks" :key="book.id" type="button"
                                @click="addFavorite(book)" :disabled="selectedFavorites.length >= 4"
                                class="w-full flex items-center gap-3 p-2 rounded-lg hover:bg-dark-700 transition-colors text-left disabled:opacity-50 disabled:cursor-not-allowed">
                                <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                                    class="w-12 h-16 object-cover rounded" />
                                <div v-else class="w-12 h-16 bg-dark-800 rounded flex items-center justify-center">
                                    <span class="text-lg text-dark-400"></span>
                                </div>
                                <div class="flex-1 min-w-0">
                                    <div class="font-medium text-white text-sm line-clamp-1">{{ book.title }}</div>
                                    <div class="text-xs text-dark-400 line-clamp-1">
                                        {{ book.authors?.join(', ') || 'Unknown Author' }}
                                    </div>
                                </div>
                            </button>
                        </div>
                    </div>

                    <!-- Reading Goal -->
                    <div>
                        <label class="block text-sm font-medium text-dark-200 mb-2">
                            Reading Goal for {{ new Date().getFullYear() }}
                        </label>
                        <Input v-model.number="form.reading_goal" type="number" min="0" placeholder="e.g., 24 books" />
                        <p class="text-xs text-dark-400 mt-1">
                            How many books do you want to read this year?
                        </p>
                    </div>

                    <!-- Action Buttons -->
                    <div class="flex gap-4 pt-6 border-t border-dark-700">
                        <button type="button" @click="$emit('close')" class="btn-secondary flex-1">
                            Cancel
                        </button>
                        <button type="submit" :disabled="submitting" class="btn-primary flex-1">
                            {{ submitting ? 'Saving...' : 'Save Changes' }}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </transition>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import Input from './ui/Input.vue'
import TextArea from './ui/TextArea.vue'

const props = defineProps({
    show: Boolean,
    currentUser: Object,
    userLogs: Array
})

const emit = defineEmits(['close', 'updated'])

const authStore = useAuthStore()
const toast = useToastStore()

const form = ref({
    bio: '',
    banner_url: '',
    reading_goal: 0
})

const selectedFavorites = ref([])
const availableBooks = ref([])
const loading = ref(false)
const submitting = ref(false)
const searchQuery = ref('')

watch(() => props.show, (newVal) => {
    if (newVal) {
        initializeForm()
        loadAvailableBooks()
    }
})

function initializeForm() {
    if (props.currentUser) {
        form.value = {
            bio: props.currentUser.bio || '',
            banner_url: props.currentUser.banner_url || '',
            reading_goal: props.currentUser.reading_goal || 0
        }

        // Load current favorite books
        if (props.currentUser.favorite_book_ids?.length > 0) {
            selectedFavorites.value = props.userLogs
                .filter(log => props.currentUser.favorite_book_ids.includes(log.book_id))
                .map(log => log.book)
                .filter(book => book) // Remove nulls
                .slice(0, 4)
        }
    }
}

function loadAvailableBooks() {
    if (!props.userLogs) return

    // Get all unique books from user's logs (read books with ratings >= 3)
    const bookMap = new Map()
    props.userLogs.forEach(log => {
        if (log.book && log.status === 'read' && (!log.rating || log.rating >= 3)) {
            if (!bookMap.has(log.book.id)) {
                bookMap.set(log.book.id, log.book)
            }
        }
    })

    availableBooks.value = Array.from(bookMap.values())
}

const filteredAvailableBooks = computed(() => {
    const selectedIds = new Set(selectedFavorites.value.map(b => b.id))
    let books = availableBooks.value.filter(book => !selectedIds.has(book.id))

    if (searchQuery.value.trim()) {
        const query = searchQuery.value.toLowerCase()
        books = books.filter(book =>
            book.title?.toLowerCase().includes(query) ||
            book.authors?.some(author => author.toLowerCase().includes(query))
        )
    }

    return books
})

function addFavorite(book) {
    if (selectedFavorites.value.length < 4) {
        selectedFavorites.value.push(book)
    }
}

function removeFavorite(bookId) {
    selectedFavorites.value = selectedFavorites.value.filter(b => b.id !== bookId)
}

async function handleSubmit() {
    submitting.value = true

    try {
        const payload = {
            bio: form.value.bio || null,
            banner_url: form.value.banner_url || null,
            favorite_book_ids: selectedFavorites.value.map(b => b.id),
            reading_goal: form.value.reading_goal || 0
        }

        const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/me/profile`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${authStore.token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        })

        if (res.ok) {
            toast.success('Profile updated successfully!')
            emit('updated')
            emit('close')
        } else {
            throw new Error('Failed to update profile')
        }
    } catch (error) {
        console.error('Error updating profile:', error)
        toast.error('Failed to update profile')
    } finally {
        submitting.value = false
    }
}
</script>
