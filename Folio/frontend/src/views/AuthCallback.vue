<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center">
      <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary mb-4"></div>
      <h2 class="text-xl font-semibold text-gray-900">Completing sign in...</h2>
      <p class="text-gray-600 mt-2">Please wait while we log you in.</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

onMounted(async () => {
  const token = route.query.token
  const converted = route.query.converted === 'true'

  if (token) {
    // Store the token
    authStore.setToken(token)
    
    // Fetch user data
    await authStore.fetchUser()
    
    // Show success message if converted from guest
    if (converted) {
      // Could show a toast notification here
      console.log('Successfully converted from guest to full user!')
    }
    
    // Redirect to home
    router.push('/')
  } else {
    // No token, redirect to login
    router.push('/login')
  }
})
</script>

