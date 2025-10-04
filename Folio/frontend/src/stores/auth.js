import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

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
    } else {
      localStorage.removeItem('token')
    }
  }

  function setGuestStatus(guestStatus) {
    isGuest.value = guestStatus
  }

  async function fetchUser() {
    if (!token.value) return

    loading.value = true
    try {
      // Try regular user endpoint first
      const response = await axios.get('/api/me')
      user.value = response.data
      isGuest.value = false
    } catch (error) {
      // If regular user fails, try guest endpoint
      try {
        const guestResponse = await axios.get('/api/guest/me')
        user.value = guestResponse.data
        isGuest.value = true
      } catch (guestError) {
        console.error('Failed to fetch user:', guestError)
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

