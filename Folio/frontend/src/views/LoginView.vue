<template>
  <div class="min-h-screen bg-dark-950 flex items-center justify-center">
    <div class="card max-w-md w-full text-center">
      <div class="mb-8">
        <div
          class="w-16 h-16 bg-gradient-to-br from-accent-red to-accent-blue rounded-2xl flex items-center justify-center mx-auto mb-6">
          <span class="text-3xl"></span>
        </div>
        <h1 class="text-display-2 mb-4 text-gradient">Welcome to Folio</h1>
        <p class="text-body text-dark-300">Your personal reading journal</p>
      </div>

      <div class="space-y-6">
        <a :href="googleAuthUrl" class="btn-primary w-full flex items-center justify-center gap-3 py-4">
          <svg class="w-5 h-5" viewBox="0 0 24 24">
            <path fill="currentColor"
              d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
            <path fill="currentColor"
              d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
            <path fill="currentColor"
              d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
            <path fill="currentColor"
              d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
          </svg>
          Sign in with Google
        </a>

        <div class="relative">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-dark-700" />
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-4 bg-dark-900 text-dark-400">or</span>
          </div>
        </div>

        <button @click="tryAsGuest" :disabled="loading"
          class="btn-secondary w-full flex items-center justify-center gap-3 py-4">
          <span v-if="loading" class="animate-spin rounded-full h-5 w-5 border-2 border-dark-400 border-t-white"></span>
          <span v-else class="text-xl">U</span>
          {{ loading ? 'Creating guest account...' : 'Try as Guest' }}
        </button>

        <div class="mt-8 p-4 bg-dark-800/50 rounded-xl border border-dark-700">
          <h3 class="text-sm font-semibold text-white mb-2">Guest Account Benefits</h3>
          <ul class="text-xs text-dark-300 space-y-1">
            <li>• Save books locally on your device</li>
            <li>• Try all features risk-free</li>
            <li>• Upgrade anytime to sync across devices</li>
          </ul>
        </div>

        <p class="text-caption text-dark-500">
          By signing in, you agree to our Terms of Service and Privacy Policy
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost'
const googleAuthUrl = `${apiUrl}/api/auth/google`

const tryAsGuest = async () => {
  loading.value = true
  try {
    await authStore.createGuestUser()
    router.push('/')
  } catch (error) {
    console.error('Failed to create guest account:', error)
    // Could show a toast notification here
  } finally {
    loading.value = false
  }
}
</script>
