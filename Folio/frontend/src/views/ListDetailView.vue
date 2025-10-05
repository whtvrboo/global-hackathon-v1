<template>
    <div v-if="loading" class="text-center p-12">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
        <p class="mt-4 text-dark-300">Loading list...</p>
    </div>
    <div v-else-if="error" class="text-center p-12">
        <div class="text-4xl mb-4">😕</div>
        <p class="text-dark-300">{{ error }}</p>
    </div>
    <div v-else-if="list" class="bg-dark-950 min-h-screen lg:pl-20">
        <ListHero :title="list.name" :description="list.description" :headerImageUrl="list.header_image_url"
            :themeColor="list.theme_color" :creator="list.creator" />

        <div class="container mx-auto max-w-4xl py-12 px-4">
            <div class="space-y-8">
                <ListItemCard v-for="item in list.items" :key="item.id" :item="item" :themeColor="list.theme_color"
                    @logBook="handleLogBook" />
            </div>

            <div class="mt-12 text-center">
                <PrimaryButton @click="shareList" class="inline-flex items-center gap-2">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M8.684 13.342C8.886 12.938 9 12.482 9 12s-.114-.938-.316-1.342m5.632 2.684C14.114 12.938 14 12.482 14 12s.114-.938.316-1.342m-6.264 5.342l5.632-3.332m-5.632-8.684l5.632 3.332" />
                    </svg>
                    Share This List
                </PrimaryButton>
            </div>
        </div>
    </div>

    <!-- Log Book Modal -->
    <LogBookModal :show="showLogModal" :book="selectedBook" @close="showLogModal = false" @success="handleLogSuccess" />
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import ListHero from '../components/ListHero.vue'
import ListItemCard from '../components/ListItemCard.vue'
import PrimaryButton from '../components/ui/PrimaryButton.vue'
import LogBookModal from '../components/LogBookModal.vue'
import { useToastStore } from '../stores/toast'

const route = useRoute()
const router = useRouter()
const list = ref(null)
const loading = ref(true)
const error = ref(null)
const showLogModal = ref(false)
const selectedBook = ref(null)
const listId = computed(() => route.params.id)
const toastStore = useToastStore()


const fetchList = async () => {
    loading.value = true
    error.value = null
    try {
        const response = await axios.get(`/api/lists/${listId.value}`)
        list.value = response.data
    } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load list'
        console.error('Error fetching list:', err)
    } finally {
        loading.value = false
    }
}

const shareList = () => {
    const url = window.location.href
    navigator.clipboard.writeText(url).then(() => {
        toastStore.success('Link copied to clipboard!')
    }, () => {
        toastStore.error('Failed to copy link.')
    })
}

const handleLogBook = (book) => {
    selectedBook.value = book
    showLogModal.value = true
}

const handleLogSuccess = () => {
    showLogModal.value = false
    selectedBook.value = null
}

onMounted(() => {
    fetchList()
})
</script>
