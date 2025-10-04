<script setup lang="ts">
import { computed } from 'vue'
import type { Commit } from '@/composables/useCollaboration'

interface Props {
    versions: Commit[]
    currentVersionId: string
}

interface Emits {
    (e: 'switch-version', versionId: string): void
    (e: 'create-new-version'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Computed property for sorted versions (newest first)
const sortedCommits = computed(() => {
    return [...props.versions].reverse()
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

const switchVersion = (versionId: string) => {
    emit('switch-version', versionId)
}

const createNewCommit = () => {
    emit('create-new-version')
}
</script>

<template>
    <div class="commit-history bg-base-100 border border-base-300 rounded-lg p-4 mb-6">
        <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-lg">Commit History</h3>
            <button @click="createNewCommit" class="btn btn-sm btn-primary">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                New Commit
            </button>
        </div>

        <ul class="space-y-2 max-h-60 overflow-y-auto">
            <li v-for="commit in sortedCommits" :key="commit.id"
                class="commit-item p-3 rounded-md border border-transparent cursor-pointer transition-colors"
                :class="{ 'bg-primary/10 border-primary/50': commit.id === currentVersionId }"
                @click="switchVersion(commit.id)">
                <div class="flex items-center justify-between">
                    <p class="font-mono text-sm font-medium text-primary">{{ commit.id }}</p>
                    <span class="text-xs text-base-content/70">{{ formatDate(commit.createdAt) }}</span>
                </div>
                <p class="text-sm my-1">{{ commit.message }}</p>
                <div class="flex items-center gap-2">
                    <!-- Placeholder for avatar -->
                    <div class="w-4 h-4 bg-base-300 rounded-full"></div>
                    <span class="text-xs font-medium">{{ commit.author.name }}</span>
                </div>
            </li>
        </ul>
    </div>
</template>

<style scoped>
.commit-item:hover {
    background-color: hsl(var(--b2));
}
</style>
