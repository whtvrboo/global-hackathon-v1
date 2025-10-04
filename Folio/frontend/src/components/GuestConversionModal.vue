<template>
  <transition
    enter-active-class="transition ease-out duration-200"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition ease-in duration-150"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="show"
      @click="$emit('close')"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
      <div
        @click.stop
        class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-8"
      >
        <div class="text-center mb-6">
          <div class="text-4xl mb-4">🎉</div>
          <h2 class="text-2xl font-bold text-gray-900 mb-2">Save Your Progress</h2>
          <p class="text-gray-600">
            You've been using Folio as a guest. Create a full account to save your reading logs permanently.
          </p>
        </div>

        <div class="space-y-4">
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h3 class="font-semibold text-blue-900 mb-2">What you'll get:</h3>
            <ul class="text-sm text-blue-800 space-y-1">
              <li>✅ Keep all your reading logs</li>
              <li>✅ Access from any device</li>
              <li>✅ Follow other readers</li>
              <li>✅ Get reading recommendations</li>
            </ul>
          </div>

          <div class="flex gap-3">
            <SecondaryButton @click="$emit('close')" class="flex-1">
              Maybe Later
            </SecondaryButton>
            <PrimaryButton @click="convertToFullUser" class="flex-1">
              Create Account
            </PrimaryButton>
          </div>

          <p class="text-xs text-gray-500 text-center">
            Your guest data will be preserved when you sign up
          </p>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'
import PrimaryButton from './ui/PrimaryButton.vue'
import SecondaryButton from './ui/SecondaryButton.vue'

const props = defineProps({
  show: Boolean
})

const emit = defineEmits(['close'])

const authStore = useAuthStore()

const convertToFullUser = () => {
  // Store the guest session ID for the conversion
  const guestSessionId = authStore.user?.guest_session_id
  if (guestSessionId) {
    localStorage.setItem('guest_session_id', guestSessionId)
  }
  
  // Redirect to Google OAuth with guest conversion
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const convertUrl = `${apiUrl}/api/auth/google/convert?guest_session_id=${guestSessionId}`
  window.location.href = convertUrl
}
</script>
