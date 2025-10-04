<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const isLoading = ref(false)
const auth = useAuthStore()
const isAuthed = computed(() => auth.isAuthenticated)

const signInWithGithub = () => {
    window.location.href = '/api/auth/github'
}

const createTrack = async () => {
    if (!auth.token) return signInWithGithub()
    isLoading.value = true
    try {
        const response = await fetch('/api/create-track', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...auth.getAuthHeader(),
            },
        })
        const data = await response.json()
        if (data.success && data.trackId) {
            router.push(`/track/${data.trackId}`)
        } else {
            alert('Failed to create track.')
        }
    } catch (error) {
        console.error('Error creating track:', error)
        alert('Failed to create track.')
    } finally {
        isLoading.value = false
    }
}
</script>

<template>
    <div class="min-h-screen bg-base-200 flex items-center justify-center">
        <div class="card w-96 bg-base-100 shadow-xl">
            <div class="card-body">
                <h1 class="card-title text-3xl font-bold text-center mb-6">Rehearsal</h1>
                <p class="text-center text-base-content/70 mb-6">
                    Start collaborating on your music
                </p>

                <div class="card-actions justify-center mt-2">
                    <button @click="signInWithGithub" class="btn btn-outline w-full">
                        Sign in with GitHub
                    </button>
                </div>

                <div class="divider">or</div>

                <div class="card-actions justify-center mt-2">
                    <button @click="createTrack" class="btn btn-primary w-full" :disabled="isLoading">
                        <span v-if="isLoading" class="loading loading-spinner loading-sm"></span>
                        {{ isLoading ? 'Creating...' :
                            (isAuthed ? 'Create New Track'
                                : 'Create Track (requires sign-in)') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Home view styles */
</style>
