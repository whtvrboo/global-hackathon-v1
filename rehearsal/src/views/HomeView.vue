<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const email = ref('')
const isLoading = ref(false)

const createTrack = async () => {
    if (!email.value) return

    isLoading.value = true

    try {
        const response = await fetch('/api/create-track', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email: email.value }),
        })

        const data = await response.json()

        if (data.success && data.trackId) {
            // Navigate to the new track page
            router.push(`/track/${data.trackId}`)
        } else {
            console.error('Failed to create track:', data.error)
            alert('Failed to create track. Please try again.')
        }
    } catch (error) {
        console.error('Error creating track:', error)
        alert('Failed to create track. Please try again.')
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

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Email</span>
                    </label>
                    <input v-model="email" type="email" placeholder="Enter your email"
                        class="input input-bordered w-full" />
                </div>

                <div class="card-actions justify-center mt-6">
                    <button @click="createTrack" class="btn btn-primary" :disabled="!email || isLoading">
                        <span v-if="isLoading" class="loading loading-spinner loading-sm"></span>
                        {{ isLoading ? 'Creating...' : 'Create Track' }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Home view styles */
</style>
