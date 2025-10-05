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
        <div v-else-if="lists.length > 0" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            <div v-for="list in lists" :key="list.id" class="card card-hover cursor-pointer group overflow-hidden"
                @click="viewList(list)">

                <!-- List Header Image -->
                <div v-if="list.header_image_url"
                    class="aspect-video bg-gradient-to-br from-accent-blue/20 to-accent-purple/20 rounded-xl mb-4 overflow-hidden">
                    <img :src="list.header_image_url" :alt="list.name"
                        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                </div>
                <div v-else class="aspect-video rounded-xl mb-4 flex items-center justify-center"
                    :style="{ backgroundColor: list.theme_color || '#6366f1' }">
                    <span class="text-4xl"></span>
                </div>

                <!-- List Info -->
                <div class="p-4">
                    <div class="flex items-start justify-between mb-3">
                        <h4 class="text-heading-4 text-white line-clamp-2 flex-1">{{ list.name }}</h4>
                        <div class="flex items-center gap-2 ml-2">
                            <span v-if="!list.is_public" class="text-xs text-dark-400 bg-dark-800 px-2 py-1 rounded">
                                Private
                            </span>
                            <button @click.stop="editList(list)"
                                class="text-dark-400 hover:text-white transition-colors">
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

                    <!-- Book Covers Collage -->
                    <div v-if="list.items_count > 0" class="mb-4">
                        <div class="flex -space-x-2">
                            <div v-for="(item, index) in getFirstFourBooks(list)" :key="item.id" class="relative">
                                <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                                    class="w-12 h-16 object-cover rounded-lg border-2 border-dark-800 shadow-lg" />
                                <div v-else
                                    class="w-12 h-16 bg-dark-800 rounded-lg border-2 border-dark-800 flex items-center justify-center">
                                    <span class="text-xs text-dark-400"></span>
                                </div>
                                <div v-if="index === 3 && list.items_count > 4"
                                    class="absolute inset-0 bg-black/50 rounded-lg flex items-center justify-center">
                                    <span class="text-xs text-white font-bold">+{{ list.items_count - 4 }}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="flex items-center justify-between text-sm text-dark-400">
                        <span>{{ list.items_count || 0 }} book{{ (list.items_count || 0) !== 1 ? 's' : '' }}</span>
                        <span>{{ formatDate(list.created_at) }}</span>
                    </div>
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
                    <div class="flex items-center gap-4">
                        <button @click="viewPublicList" class="btn-outline text-sm">
                            View Public Page
                        </button>
                        <button @click="selectedList = null" class="text-dark-400 hover:text-white transition-colors">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M6 18L18 6M6 6l12 12"></path>
                            </svg>
                        </button>
                    </div>
                </div>

                <!-- List Items with Drag and Drop -->
                <div v-if="selectedList.items?.length > 0" class="space-y-4">
                    <draggable v-model="selectedList.items" @end="onDragEnd" :animation="200" ghost-class="ghost"
                        chosen-class="chosen" drag-class="drag">
                        <div v-for="(item, index) in selectedList.items" :key="item.id"
                            class="flex gap-4 p-4 bg-dark-800 rounded-xl cursor-move hover:bg-dark-700 transition-colors">

                            <!-- Drag Handle -->
                            <div
                                class="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-dark-400 hover:text-white transition-colors">
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M4 8h16M4 16h16"></path>
                                </svg>
                            </div>

                            <!-- Number -->
                            <div
                                class="flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-white font-bold text-sm bg-accent-blue">
                                {{ index + 1 }}
                            </div>

                            <!-- Book Cover -->
                            <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                                class="w-16 h-24 object-cover rounded-lg" />
                            <div v-else class="w-16 h-24 bg-dark-700 rounded-lg flex items-center justify-center">
                                <span class="text-2xl text-dark-400"></span>
                            </div>

                            <!-- Book Info -->
                            <div class="flex-1">
                                <h4 class="text-heading-4 text-white mb-1">{{ item.book.title }}</h4>
                                <p v-if="item.book.authors" class="text-body text-dark-300 mb-2">
                                    by {{ item.book.authors.join(', ') }}
                                </p>
                                <div v-if="editingItemId !== item.id">
                                    <p v-if="item.notes" class="text-sm text-dark-400 italic flex items-center gap-2">
                                        <span>"{{ item.notes }}"</span>
                                        <button @click="startEditingNotes(item)"
                                            class="text-dark-500 hover:text-white transition-colors">
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.5L14.732 3.732z">
                                                </path>
                                            </svg>
                                        </button>
                                    </p>
                                    <button v-else @click="startEditingNotes(item)"
                                        class="text-sm text-dark-400 hover:text-white transition-colors flex items-center gap-2">
                                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                                        </svg>
                                        Add Note
                                    </button>
                                </div>
                                <div v-else>
                                    <TextArea v-model="editingNotes" class="w-full" :rows="3" />
                                    <div class="flex gap-2 mt-2">
                                        <PrimaryButton @click="saveNotes(item.id)" size="sm">Save</PrimaryButton>
                                        <SecondaryButton @click="cancelEditingNotes" size="sm">Cancel</SecondaryButton>
                                    </div>
                                </div>
                            </div>

                            <!-- Remove Button -->
                            <button @click="removeFromList(item.id)"
                                class="text-dark-400 hover:text-accent-red transition-colors">
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
                                    </path>
                                </svg>
                            </button>
                        </div>
                    </draggable>
                </div>

                <!-- Empty List -->
                <div v-else class="text-center py-12">
                    <div class="text-6xl mb-4"></div>
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
import { useRouter } from 'vue-router'
import axios from 'axios'
import draggable from 'vuedraggable'
import ListModal from './ListModal.vue'
import { useToastStore } from '../stores/toast'
import PrimaryButton from './ui/PrimaryButton.vue'
import SecondaryButton from './ui/SecondaryButton.vue'
import TextArea from './ui/TextArea.vue'

const router = useRouter()
const lists = ref([])
const loading = ref(true)
const selectedList = ref(null)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingList = ref(null)
const toastStore = useToastStore()
const editingItemId = ref(null)
const editingNotes = ref('')

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
        const response = await axios.get('/api/me/lists')
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

const startEditingNotes = (item) => {
    editingItemId.value = item.id
    editingNotes.value = item.notes || ''
}

const cancelEditingNotes = () => {
    editingItemId.value = null
    editingNotes.value = ''
}

const saveNotes = async (itemId) => {
    try {
        await axios.put(`/api/lists/${selectedList.value.id}/items/${itemId}`, {
            notes: editingNotes.value
        })
        toastStore.success('Note updated successfully')
        // Update the item in the local state
        const item = selectedList.value.items.find(i => i.id === itemId)
        if (item) {
            item.notes = editingNotes.value
        }
        cancelEditingNotes()
    } catch (error) {
        console.error('Error saving notes:', error)
        toastStore.error('Failed to save note')
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

const getFirstFourBooks = (list) => {
    // This would need to be populated from the backend
    // For now, return empty array
    return []
}

const onDragEnd = async (event) => {
    if (event.oldIndex === event.newIndex) return

    try {
        const itemIds = selectedList.value.items.map(item => item.id)
        await axios.put(`/api/lists/${selectedList.value.id}/items/order`, {
            item_ids: itemIds
        })
        toastStore.success('List order updated!')
    } catch (error) {
        console.error('Error updating list order:', error)
        toastStore.error('Failed to update list order')
        // Reload the list to restore original order
        await viewList(selectedList.value)
    }
}

const viewPublicList = () => {
    if (selectedList.value) {
        router.push(`/lists/${selectedList.value.id}`)
        selectedList.value = null
    }
}

onMounted(() => {
    loadLists()
})
</script>

<style scoped>
.ghost {
    opacity: 0.5;
    background: #374151;
}

.chosen {
    background: #4B5563;
}

.drag {
    background: #6B7280;
}
</style>
