<template>
    <div class="space-y-6">
        <!-- Header -->
        <div class="flex items-center justify-between">
            <h3 class="text-heading-3">My Lists</h3>
            <button @click="showCreateModal = true" class="btn-primary flex items-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                Create List
            </button>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="space-y-4">
            <div v-for="i in 3" :key="i" class="animate-pulse">
                <div class="card">
                    <div class="h-4 bg-dark-800 rounded w-1/3 mb-2"></div>
                    <div class="h-3 bg-dark-800 rounded w-2/3 mb-4"></div>
                    <div class="h-3 bg-dark-800 rounded w-1/4"></div>
                </div>
            </div>
        </div>

        <!-- Lists Grid -->
        <div v-else-if="lists.length > 0" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            <div v-for="list in lists" :key="list.id" class="card card-hover cursor-pointer" @click="viewList(list)">
                <div class="flex items-start justify-between mb-3">
                    <h4 class="text-heading-4 text-white">{{ list.name }}</h4>
                    <div class="flex items-center gap-2">
                        <span v-if="!list.is_public" class="text-xs text-dark-400 bg-dark-800 px-2 py-1 rounded">
                            Private
                        </span>
                        <button @click.stop="editList(list)" class="text-dark-400 hover:text-white transition-colors">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z">
                                </path>
                            </svg>
                        </button>
                        <button @click.stop="deleteList(list)"
                            class="text-dark-400 hover:text-accent-red transition-colors">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
                                </path>
                            </svg>
                        </button>
                    </div>
                </div>

                <p v-if="list.description" class="text-body text-dark-300 mb-4 line-clamp-2">
                    {{ list.description }}
                </p>

                <div class="flex items-center justify-between text-sm text-dark-400">
                    <span>{{ list.items_count || 0 }} book{{ (list.items_count || 0) !== 1 ? 's' : '' }}</span>
                    <span>{{ formatDate(list.created_at) }}</span>
                </div>
            </div>
        </div>

        <!-- Empty State -->
        <div v-else class="card text-center py-12">
            <div class="text-6xl mb-4">📋</div>
            <h4 class="text-heading-4 mb-2">Start curating your reading journey</h4>
            <p class="text-body text-dark-300 mb-6">
                Create custom lists to organize your books by genre, mood, or any theme you love.
                Share your favorites with the community!
            </p>
            <button @click="showCreateModal = true" class="btn-primary">
                <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                Create Your First List
            </button>
        </div>

        <!-- List Detail Modal -->
        <div v-if="selectedList" @click="selectedList = null"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
            <div @click.stop
                class="bg-dark-900 rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-y-auto p-8">
                <div class="flex items-center justify-between mb-6">
                    <div>
                        <h2 class="text-heading-2 text-white">{{ selectedList.name }}</h2>
                        <p v-if="selectedList.description" class="text-body text-dark-300 mt-2">
                            {{ selectedList.description }}
                        </p>
                    </div>
                    <button @click="selectedList = null" class="text-dark-400 hover:text-white transition-colors">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12"></path>
                        </svg>
                    </button>
                </div>

                <!-- List Items -->
                <div v-if="selectedList.items?.length > 0" class="space-y-4">
                    <div v-for="item in selectedList.items" :key="item.id"
                        class="flex gap-4 p-4 bg-dark-800 rounded-xl">
                        <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                            class="w-16 h-24 object-cover rounded-lg" />
                        <div v-else class="w-16 h-24 bg-dark-700 rounded-lg flex items-center justify-center">
                            <span class="text-2xl text-dark-400">📚</span>
                        </div>
                        <div class="flex-1">
                            <h4 class="text-heading-4 text-white mb-1">{{ item.book.title }}</h4>
                            <p v-if="item.book.authors" class="text-body text-dark-300 mb-2">
                                by {{ item.book.authors.join(', ') }}
                            </p>
                            <p v-if="item.notes" class="text-sm text-dark-400 italic">
                                "{{ item.notes }}"
                            </p>
                        </div>
                        <button @click="removeFromList(item.id)"
                            class="text-dark-400 hover:text-accent-red transition-colors">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
                                </path>
                            </svg>
                        </button>
                    </div>
                </div>

                <!-- Empty List -->
                <div v-else class="text-center py-12">
                    <div class="text-6xl mb-4">📖</div>
                    <h4 class="text-heading-4 mb-2">No books in this list</h4>
                    <p class="text-body text-dark-300">
                        Add books to this list from the book details page
                    </p>
                </div>
            </div>
        </div>

        <!-- Create/Edit List Modal -->
        <ListModal :show="showCreateModal || showEditModal" :list="editingList" @close="closeModals"
            @success="handleListSuccess" />
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import ListModal from './ListModal.vue'
import { useToastStore } from '../stores/toast'

const lists = ref([])
const loading = ref(true)
const selectedList = ref(null)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingList = ref(null)
const toastStore = useToastStore()

const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    })
}

const loadLists = async () => {
    try {
        loading.value = true
        const response = await axios.get('/api/users/me/lists')
        lists.value = response.data.lists || []
    } catch (error) {
        console.error('Error loading lists:', error)
    } finally {
        loading.value = false
    }
}

const viewList = async (list) => {
    try {
        const response = await axios.get(`/api/lists/${list.id}`)
        selectedList.value = response.data
    } catch (error) {
        console.error('Error loading list details:', error)
    }
}

const editList = (list) => {
    editingList.value = list
    showEditModal.value = true
}

const deleteList = async (list) => {
    if (!confirm(`Are you sure you want to delete "${list.name}"? This action cannot be undone.`)) {
        return
    }

    try {
        await axios.delete(`/api/lists/${list.id}`)
        toastStore.success(`List "${list.name}" deleted successfully`)
        await loadLists()
    } catch (error) {
        console.error('Error deleting list:', error)
        toastStore.error('Failed to delete list')
    }
}

const removeFromList = async (itemId) => {
    if (!confirm('Remove this book from the list?')) {
        return
    }

    try {
        await axios.delete(`/api/lists/${selectedList.value.id}/items/${itemId}`)
        toastStore.success('Book removed from list')
        // Reload the list
        await viewList(selectedList.value)
        // Update the main lists
        await loadLists()
    } catch (error) {
        console.error('Error removing book from list:', error)
        toastStore.error('Failed to remove book from list')
    }
}

const closeModals = () => {
    showCreateModal.value = false
    showEditModal.value = false
    editingList.value = null
}

const handleListSuccess = () => {
    loadLists()
    closeModals()
}

onMounted(() => {
    loadLists()
})
</script>
