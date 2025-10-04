<script setup lang="ts">
import { computed } from 'vue'
import type { Version } from '@/composables/useCollaboration'

interface Props {
    versions: Version[]
    currentVersionId: string
}

interface Emits {
    (e: 'switch-version', versionId: string): void
    (e: 'create-new-version'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Computed property for current version
const currentVersion = computed(() => {
    return props.versions.find(v => v.id === props.currentVersionId)
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

const createNewVersion = () => {
    // Add a small delay to show the button was clicked
    emit('create-new-version')
}
</script>

<template>
    <div class="version-selector bg-base-100 border border-base-300 rounded-lg p-4 mb-6">
        <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-3">
                <h3 class="font-semibold text-lg">Version</h3>
                <span class="badge badge-outline badge-sm">{{ versions.length }} version{{ versions.length !== 1 ? 's' :
                    '' }}</span>
            </div>
            <button @click="createNewVersion" class="btn btn-sm btn-primary">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                New Version
            </button>
        </div>

        <div class="flex items-center gap-2">
            <label class="text-sm font-medium">Current:</label>
            <select :value="currentVersionId" @change="switchVersion(($event.target as HTMLSelectElement).value)"
                class="select select-bordered select-sm flex-1">
                <option v-for="version in versions" :key="version.id" :value="version.id">
                    {{ version.name }}
                </option>
            </select>
        </div>

        <!-- Version details -->
        <div v-if="currentVersion" class="mt-3 text-sm text-base-content/70">
            <div class="flex justify-between">
                <span>{{ currentVersion.stems.length }} stem{{ currentVersion.stems.length !== 1 ? 's' : '' }}</span>
                <span>{{ currentVersion.comments.length }} comment{{ currentVersion.comments.length !== 1 ? 's' : ''
                    }}</span>
            </div>
            <div class="text-xs mt-1">
                Created {{ formatDate(currentVersion.createdAt) }}
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Version selector styles */
</style>
