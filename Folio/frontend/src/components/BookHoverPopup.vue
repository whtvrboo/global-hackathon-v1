<template>
    <div v-if="show"
        class="absolute z-50 bg-dark-900 border border-dark-700 rounded-xl shadow-2xl p-4 w-80 transform transition-all duration-200"
        :class="positionClass" @mouseenter="$emit('mouseenter')" @mouseleave="$emit('mouseleave')">
        <!-- Book Info -->
        <div class="flex gap-3 mb-4">
            <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                class="w-16 h-20 object-cover rounded shadow-md flex-shrink-0" />
            <div v-else
                class="w-16 h-20 bg-dark-800 rounded flex items-center justify-center text-lg text-dark-400 shadow-md flex-shrink-0">
                📚
            </div>
            <div class="flex-1 min-w-0">
                <h4 class="font-bold text-white text-sm line-clamp-2 mb-1">{{ book.title }}</h4>
                <p v-if="book.authors?.length" class="text-xs text-dark-300 line-clamp-1 mb-2">
                    by {{ book.authors.join(', ') }}
                </p>
                <div v-if="book.rating" class="flex items-center gap-1">
                    <span class="text-accent-orange text-xs">{{ '★'.repeat(Math.round(book.rating)) }}{{ '☆'.repeat(5 -
                        Math.round(book.rating)) }}</span>
                    <span class="text-xs text-dark-400">{{ book.rating }}</span>
                </div>
            </div>
        </div>


        <!-- Categories -->
        <div v-if="book.categories?.length" class="mb-4">
            <div class="flex flex-wrap gap-1">
                <span v-for="category in book.categories.slice(0, 3)" :key="category"
                    class="px-2 py-1 text-xs bg-dark-800 text-dark-300 rounded-full border border-dark-600">
                    {{ category }}
                </span>
            </div>
        </div>

        <!-- Book Stats -->
        <div class="grid grid-cols-2 gap-3 mb-4">
            <div class="text-center">
                <div class="text-sm font-bold text-accent-blue">{{ book.page_count || 'N/A' }}</div>
                <div class="text-xs text-dark-400">Pages</div>
            </div>
            <div class="text-center">
                <div class="text-sm font-bold text-accent-green">{{ book.published_date || 'N/A' }}</div>
                <div class="text-xs text-dark-400">Published</div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-2">
            <button @click="logBook"
                class="flex-1 px-3 py-2 text-xs font-medium bg-accent-red text-white rounded-lg hover:bg-accent-red/90 transition-all duration-200">
                📖 Log Book
            </button>
            <button @click="viewBook"
                class="flex-1 px-3 py-2 text-xs font-medium bg-dark-800 text-white rounded-lg border border-dark-600 hover:bg-dark-700 transition-all duration-200">
                View Details
            </button>
        </div>

        <!-- Quick Add to List -->
        <div v-if="userLists.length > 0" class="mt-3 pt-3 border-t border-dark-700">
            <div class="text-xs text-dark-400 mb-2">Quick add to list:</div>
            <div class="flex flex-wrap gap-1">
                <button v-for="list in userLists.slice(0, 3)" :key="list.id" @click="addToList(list.id)"
                    class="px-2 py-1 text-xs bg-dark-800 text-dark-300 rounded border border-dark-600 hover:bg-dark-700 hover:text-white transition-all duration-200">
                    {{ list.name }}
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'
import { useToastStore } from '../stores/toast'

const props = defineProps({
    show: Boolean,
    book: Object,
    position: {
        type: String,
        default: 'bottom-right' // bottom-right, bottom-left, top-right, top-left
    }
})

const emit = defineEmits(['close', 'logBook', 'mouseenter', 'mouseleave'])

const router = useRouter()
const authStore = useAuthStore()
const toastStore = useToastStore()

const userLists = ref([])

const positionClass = computed(() => {
    const classes = {
        'bottom-right': 'top-full left-1/2 transform -translate-x-1/2 mt-2',
        'bottom-left': 'top-full right-1/2 transform translate-x-1/2 mt-2',
        'top-right': 'bottom-full left-1/2 transform -translate-x-1/2 mb-2',
        'top-left': 'bottom-full right-1/2 transform translate-x-1/2 mb-2'
    }
    return classes[props.position] || classes['bottom-right']
})


const logBook = () => {
    emit('logBook', props.book)
    emit('close')
}

const viewBook = () => {
    router.push(`/books/${props.book.id}`)
    emit('close')
}

const addToList = async (listId) => {
    try {
        await axios.post(`/api/lists/${listId}/items`, {
            book_id: props.book.id
        })
        toastStore.success('Book added to list!')
    } catch (error) {
        if (error.response?.status === 409) {
            toastStore.warning('This book is already in the list!')
        } else {
            console.error('Error adding book to list:', error)
            toastStore.error('Failed to add book to list')
        }
    }
}

const loadUserLists = async () => {
    if (!authStore.isAuthenticated) return

    try {
        const response = await axios.get('/api/me/lists')
        userLists.value = response.data.lists || []
    } catch (error) {
        console.error('Error loading user lists:', error)
        userLists.value = []
    }
}

onMounted(() => {
    loadUserLists()
})
</script>
