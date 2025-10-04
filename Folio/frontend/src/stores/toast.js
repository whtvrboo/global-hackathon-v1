import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useToastStore = defineStore('toast', () => {
  const toasts = ref([])
  let toastId = 0

  const addToast = (message, type = 'info', duration = 4000) => {
    const id = ++toastId
    const toast = {
      id,
      message,
      type, // 'success', 'error', 'info', 'warning'
      duration,
      timestamp: Date.now()
    }
    
    toasts.value.push(toast)
    
    // Auto remove after duration
    if (duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, duration)
    }
    
    return id
  }

  const removeToast = (id) => {
    const index = toasts.value.findIndex(toast => toast.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }

  const clearAll = () => {
    toasts.value = []
  }

  // Convenience methods
  const success = (message, duration = 4000) => addToast(message, 'success', duration)
  const error = (message, duration = 6000) => addToast(message, 'error', duration)
  const info = (message, duration = 4000) => addToast(message, 'info', duration)
  const warning = (message, duration = 5000) => addToast(message, 'warning', duration)

  return {
    toasts,
    addToast,
    removeToast,
    clearAll,
    success,
    error,
    info,
    warning
  }
})
