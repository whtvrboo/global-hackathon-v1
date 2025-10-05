<template>
    <div class="min-h-screen bg-dark-950">
        <!-- Header -->
        <div class="bg-dark-900/50 backdrop-blur-sm border-b border-dark-800">
            <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
                <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-4">
                        <button @click="$router.back()"
                            class="p-2 text-dark-400 hover:text-white hover:bg-dark-800 rounded-lg transition-colors">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M15 19l-7-7 7-7" />
                            </svg>
                        </button>
                        <div>
                            <h1 class="text-heading-1 text-white">Create List</h1>
                            <p class="mt-1 text-body text-dark-400">
                                {{ prefilledData ? 'From your thread insights' : 'Build your curated collection' }}
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Main Content -->
        <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <ListModal :show="true" :prefilled-data="prefilledData" @close="handleClose" @success="handleSuccess" />
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ListModal from '../components/ListModal.vue'

const router = useRouter()
const prefilledData = ref(null)

onMounted(() => {
    // Check for pre-filled data from session storage
    const storedData = sessionStorage.getItem('prefilledListData')
    if (storedData) {
        try {
            prefilledData.value = JSON.parse(storedData)
        } catch (error) {
            console.error('Error parsing prefilled data:', error)
        }
    }
})

const handleClose = () => {
    // Clear session storage
    sessionStorage.removeItem('prefilledListData')
    router.back()
}

const handleSuccess = () => {
    // Clear session storage
    sessionStorage.removeItem('prefilledListData')
    // Navigate to the newly created list or dashboard
    router.push('/')
}
</script>
