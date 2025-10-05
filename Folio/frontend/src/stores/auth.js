import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import router from '../router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || null)
  const user = ref(null)
  const loading = ref(false)
  const isGuest = ref(false)

  const isAuthenticated = computed(() => !!token.value)
  const isGuestUser = computed(() => isGuest.value)

  // Set up axios interceptor
  axios.interceptors.request.use((config) => {
    if (token.value) {
      config.headers.Authorization = `Bearer ${token.value}`
      console.log('Adding auth header to request:', config.url)
    } else {
      console.log('No token available for request:', config.url)
    }
    return config
  })

  function setToken(newToken) {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
      // Parse token to set initial guest status
      try {
        const tokenParts = newToken.split('.')
        if (tokenParts.length === 3) {
          const payload = JSON.parse(atob(tokenParts[1]))
          const isGuestFromToken = payload.is_guest === true
          console.log('Setting token - JWT payload:', payload)
          console.log('Setting token - Is guest:', isGuestFromToken)
          isGuest.value = isGuestFromToken
        }
      } catch (error) {
        console.error('Failed to parse token for guest status:', error)
        isGuest.value = false
      }
    } else {
      localStorage.removeItem('token')
      isGuest.value = false
    }
  }

  function setGuestStatus(guestStatus) {
    isGuest.value = guestStatus
  }

  async function fetchUser() {
    if (!token.value) return

    loading.value = true
    try {
      // Parse JWT token to check if user is guest
      const tokenParts = token.value.split('.')
      if (tokenParts.length === 3) {
        try {
          const payload = JSON.parse(atob(tokenParts[1]))
          const isGuestFromToken = payload.is_guest === true
          console.log('JWT payload:', payload)
          console.log('Is guest from token:', isGuestFromToken)
          
          if (isGuestFromToken) {
            // User is a guest, use guest endpoint
            console.log('Using guest endpoint')
            const guestResponse = await axios.get('/api/guest/me')
            user.value = guestResponse.data
            isGuest.value = true
          } else {
            // User is authenticated, use regular endpoint
            console.log('Using regular user endpoint')
            const response = await axios.get('/api/me')
            user.value = response.data
            isGuest.value = false
          }
        } catch (parseError) {
          console.error('Failed to parse JWT token:', parseError)
          // Fallback to regular endpoint
          const response = await axios.get('/api/me')
          user.value = response.data
          isGuest.value = false
        }
      } else {
        // Invalid token format, try regular endpoint
        const response = await axios.get('/api/me')
        user.value = response.data
        isGuest.value = false
      }
    } catch (error) {
      console.error('Failed to fetch user:', error)
      // If regular user fails, try guest endpoint as fallback
      try {
        const guestResponse = await axios.get('/api/guest/me')
        user.value = guestResponse.data
        isGuest.value = true
      } catch (guestError) {
        console.error('Failed to fetch guest user:', guestError)
        // If both fail, clear token
        if (guestError.response?.status === 401) {
          logout()
        }
      }
    } finally {
      loading.value = false
    }
  }

  function logout() {
    setToken(null)
    user.value = null
    isGuest.value = false
    // Redirect to home page after logout
    router.push('/')
  }

  async function createGuestUser() {
    loading.value = true
    try {
      const response = await axios.post('/api/auth/guest')
      setToken(response.data.token)
      user.value = response.data.user
      isGuest.value = true
      return response.data.guest_session_id
    } catch (error) {
      console.error('Failed to create guest user:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // Initialize user on store creation
  if (token.value) {
    fetchUser()
  }

  return {
    token,
    user,
    loading,
    isAuthenticated,
    isGuestUser,
    setToken,
    setGuestStatus,
    fetchUser,
    logout,
    createGuestUser
  }
})

