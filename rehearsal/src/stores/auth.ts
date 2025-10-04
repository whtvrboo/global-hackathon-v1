import { defineStore } from 'pinia'

interface UserInfo {
  id: string
  username: string
  avatar?: string
}

function decodeJwtPayload(token: string): any | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null

    const payload = parts[1]
    if (!payload) return null

    const json = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(json)
  } catch {
    return null
  }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: (typeof localStorage !== 'undefined' && localStorage.getItem('auth_token')) || '',
    user: null as UserInfo | null,
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
  },
  actions: {
    setToken(token: string) {
      this.token = token
      if (typeof localStorage !== 'undefined') localStorage.setItem('auth_token', token)
      const payload = decodeJwtPayload(token)
      if (payload) {
        this.user = { id: String(payload.sub), username: payload.username, avatar: payload.avatar }
      }
    },
    clear() {
      this.token = ''
      this.user = null
      if (typeof localStorage !== 'undefined') localStorage.removeItem('auth_token')
    },
    getAuthHeader(): HeadersInit {
      return this.token ? { Authorization: `Bearer ${this.token}` } : {}
    },
  },
})
