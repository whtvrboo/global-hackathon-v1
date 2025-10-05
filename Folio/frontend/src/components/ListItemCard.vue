<template>
    <div class="flex flex-col md:flex-row gap-6 md:gap-8 p-6 bg-dark-900 rounded-xl">
        <div class="relative flex-shrink-0">
            <div class="cursor-pointer" @click="navigateToBook" @mouseenter="showBookPopupImmediate"
                @mouseleave="hideBookPopup">
                <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                    class="w-32 md:w-40 h-auto object-cover rounded-lg shadow-lg mx-auto hover:shadow-xl transition-shadow" />
                <div v-else
                    class="w-32 md:w-40 h-48 md:h-56 bg-dark-800 rounded-lg flex items-center justify-center hover:bg-dark-700 transition-colors">
                    <span class="text-3xl text-dark-400">📚</span>
                </div>
            </div>

            <!-- Book Hover Popup -->
            <BookHoverPopup :show="showBookPopup" :book="item.book" position="bottom-right"
                @close="showBookPopup = false" @logBook="handleLogBook" />
        </div>

        <div class="flex-1">
            <h3 class="text-2xl font-bold text-white mb-1 cursor-pointer hover:text-primary transition-colors"
                @click="navigateToBook">{{ item.book.title }}</h3>
            <p v-if="item.book.authors" class="text-lg text-dark-300 mb-4">by {{ item.book.authors.join(', ') }}</p>

            <blockquote v-if="item.notes" class="relative p-4 bg-dark-800 border-l-4 rounded-r-lg"
                :style="{ borderColor: themeColor }">
                <p class="text-white/90 italic">{{ item.notes }}</p>
            </blockquote>

            <p v-if="item.book.description" class="mt-4 text-sm text-dark-400 line-clamp-3">
                {{ item.book.description }}
            </p>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import BookHoverPopup from './BookHoverPopup.vue'

const router = useRouter()

const props = defineProps({
    item: Object,
    themeColor: {
        type: String,
        default: '#6366f1'
    }
})

const emit = defineEmits(['logBook'])

const showBookPopup = ref(false)
let hoverTimeout = null

const navigateToBook = () => {
    router.push(`/books/${props.item.book.id}`)
}

const hideBookPopup = () => {
    // Add a small delay to prevent flickering when moving between elements
    hoverTimeout = setTimeout(() => {
        showBookPopup.value = false
    }, 100)
}

const showBookPopupImmediate = () => {
    // Clear any pending close timeout when showing popup
    if (hoverTimeout) {
        clearTimeout(hoverTimeout)
        hoverTimeout = null
    }
    showBookPopup.value = true
}

const handleLogBook = (book) => {
    emit('logBook', book)
}
</script>
