<script setup lang="ts">
import { computed } from 'vue'
import type { Take } from '@/composables/useCollaboration'

interface Props {
    takes: Take[]
    currentTakeId: string
}

interface Emits {
    (e: 'switch-take', takeId: string): void
    (e: 'create-new-take'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Computed property for sorted takes (newest first)
const sortedTakes = computed(() => {
    return [...props.takes].reverse()
})

// Format date for display
const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60)

    if (diffInHours < 1) {
        return 'Just now'
    } else if (diffInHours < 24) {
        return `${Math.floor(diffInHours)}h ago`
    } else if (diffInHours < 168) { // 7 days
        return `${Math.floor(diffInHours / 24)}d ago`
    } else {
        return date.toLocaleDateString()
    }
}

const switchTake = (takeId: string) => {
    emit('switch-take', takeId)
}

const createNewTake = () => {
    emit('create-new-take')
}
</script>

<template>
    <div class="takes-explorer bg-base-100 border border-base-300 rounded-lg p-4 mb-6">
        <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-lg">Song & Takes</h3>
            <button @click="createNewTake" class="btn btn-sm btn-primary">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                Save as New Take
            </button>
        </div>

        <ul class="space-y-2 max-h-60 overflow-y-auto">
            <li v-for="take in sortedTakes" :key="take.id"
                class="take-item p-3 rounded-md border border-transparent cursor-pointer transition-all duration-200 hover:shadow-md"
                :class="{ 'bg-primary/10 border-primary/50 shadow-sm': take.id === currentTakeId }"
                @click="switchTake(take.id)">
                <div class="flex items-center justify-between mb-2">
                    <div class="flex items-center gap-2">
                        <div class="w-2 h-2 rounded-full"
                            :class="take.id === currentTakeId ? 'bg-primary' : 'bg-base-300'"></div>
                        <p class="font-medium text-sm">{{ take.name }}</p>
                        <span v-if="take.id === currentTakeId" class="badge badge-xs badge-primary">Current</span>
                    </div>
                    <span class="text-xs text-base-content/70">{{ formatDate(take.createdAt) }}</span>
                </div>
                <p class="text-sm my-1 text-base-content/80 mb-2">{{ take.description }}</p>
                <div class="flex items-center justify-between">
                    <div class="flex items-center gap-2">
                        <div
                            class="w-5 h-5 bg-gradient-to-br from-primary/20 to-primary/40 rounded-full flex items-center justify-center">
                            <span class="text-xs font-medium text-primary">{{ take.author.name.charAt(0).toUpperCase()
                                }}</span>
                        </div>
                        <span class="text-xs font-medium">{{ take.author.name }}</span>
                    </div>
                    <div class="flex items-center gap-2 text-xs text-base-content/60">
                        <span class="badge badge-xs badge-outline">{{ take.stems.length }} stem{{ take.stems.length !==
                            1 ? 's' : '' }}</span>
                        <span class="badge badge-xs badge-outline">{{ take.comments.length }} comment{{
                            take.comments.length !== 1 ? 's' : '' }}</span>
                    </div>
                </div>
            </li>
        </ul>
    </div>
</template>

<style scoped>
.take-item:hover {
    background-color: hsl(var(--b2));
}
</style>
