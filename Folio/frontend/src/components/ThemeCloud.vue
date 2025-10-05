<template>
    <div class="bg-dark-900/50 backdrop-blur-sm rounded-lg border border-dark-800 p-6">
        <div class="flex items-center justify-between mb-6">
            <div>
                <h2 class="text-heading-3 text-white">Your Great Themes</h2>
                <p class="mt-1 text-sm text-dark-400">
                    Discover the patterns in your thinking
                </p>
            </div>
            <div v-if="themes.length > 0" class="text-sm text-dark-400">
                {{ themes.length }} theme{{ themes.length !== 1 ? 's' : '' }}
            </div>
        </div>

        <!-- Loading State -->
        <div v-if="isLoading" class="flex items-center justify-center py-12">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-accent-red"></div>
        </div>

        <!-- Empty State -->
        <div v-else-if="themes.length === 0" class="text-center py-12">
            <svg class="mx-auto h-16 w-16 text-dark-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                    d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
            </svg>
            <h3 class="text-lg font-medium text-white mb-2">No themes yet</h3>
            <p class="text-sm text-dark-400 mb-4">
                Start capturing thoughts with tags to see your themes emerge
            </p>
            <button @click="$emit('open-capture')"
                class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-accent-red hover:bg-accent-red/80 rounded-lg transition-colors">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                Capture a thought
            </button>
        </div>

        <!-- Theme Cloud -->
        <div v-else class="space-y-4">
            <!-- Tag Cloud -->
            <div class="flex flex-wrap gap-3">
                <button v-for="theme in themes" :key="theme.tag" @click="selectTheme(theme)"
                    class="group relative inline-flex items-center px-4 py-2 rounded-full border-2 transition-all duration-200 hover:scale-105"
                    :class="[
                        selectedTheme?.tag === theme.tag
                            ? 'border-accent-blue bg-accent-blue/20 text-accent-blue'
                            : 'border-dark-600 bg-dark-800/50 text-white hover:border-accent-blue/50 hover:bg-accent-blue/10',
                        getTagSize(theme.usage_count)
                    ]">
                    <span class="font-medium">#{{ theme.tag }}</span>
                    <span class="ml-2 text-xs opacity-70">{{ theme.usage_count }}</span>

                    <!-- Hover effect -->
                    <div
                        class="absolute inset-0 rounded-full bg-gradient-to-r from-accent-blue/20 to-accent-red/20 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                    </div>
                </button>
            </div>

            <!-- Selected Theme Details -->
            <div v-if="selectedTheme" class="mt-6 p-4 bg-dark-800/50 border border-dark-700 rounded-lg">
                <div class="flex items-center justify-between mb-3">
                    <h3 class="text-lg font-medium text-white">
                        #{{ selectedTheme.tag }}
                    </h3>
                    <button @click="viewThread"
                        class="inline-flex items-center px-3 py-1.5 text-sm font-medium text-accent-blue bg-accent-blue/10 hover:bg-accent-blue/20 rounded-lg transition-colors">
                        View Thread
                        <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                        </svg>
                    </button>
                </div>

                <div class="grid grid-cols-2 gap-4 text-sm">
                    <div>
                        <span class="text-dark-400">Usage:</span>
                        <span class="ml-2 text-white font-medium">{{ selectedTheme.usage_count }} time{{
                            selectedTheme.usage_count !== 1 ? 's' : '' }}</span>
                    </div>
                    <div>
                        <span class="text-dark-400">Last used:</span>
                        <span class="ml-2 text-white">{{ formatDate(selectedTheme.last_used_at) }}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToastStore } from '../stores/toast'

const router = useRouter()
const toast = useToastStore()

// Props
defineProps({
    themes: {
        type: Array,
        default: () => []
    },
    isLoading: {
        type: Boolean,
        default: false
    }
})

// Emits
const emit = defineEmits(['open-capture', 'theme-selected'])

// State
const selectedTheme = ref(null)

// Methods
const selectTheme = (theme) => {
    selectedTheme.value = theme
    emit('theme-selected', theme)
}

const viewThread = () => {
    if (selectedTheme.value) {
        router.push(`/notebook/threads/${encodeURIComponent(selectedTheme.value.tag)}`)
    }
}

const getTagSize = (usageCount) => {
    if (usageCount >= 10) return 'text-lg px-5 py-3'
    if (usageCount >= 5) return 'text-base px-4 py-2'
    if (usageCount >= 3) return 'text-sm px-3 py-1.5'
    return 'text-xs px-2 py-1'
}

const formatDate = (dateString) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now - date
    const diffDays = Math.floor(diffMs / 86400000)

    if (diffDays === 0) return 'Today'
    if (diffDays === 1) return 'Yesterday'
    if (diffDays < 7) return `${diffDays} days ago`
    if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`
    if (diffDays < 365) return `${Math.floor(diffDays / 30)} months ago`

    return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' })
}
</script>

<style scoped>
/* Custom animations for tag cloud */
.group:hover {
    transform: translateY(-1px);
}

/* Smooth transitions for all interactive elements */
button {
    transition: all 0.2s ease;
}

/* Gradient overlay for selected theme */
.group.selected {
    background: linear-gradient(135deg, rgba(59, 130, 246, 0.1), rgba(239, 68, 68, 0.1));
}
</style>
