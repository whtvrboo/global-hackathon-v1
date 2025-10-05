<template>
  <div id="app" class="min-h-screen bg-dark-950">
    <Navigation />
    <router-view />
    <GlobalToast />
    <PWAUpdatePrompt />

    <!-- Floating Action Button (FAB) for Quick Capture -->
    <Transition name="fab">
      <button v-if="showFAB" @click="openQuickCapture"
        class="fixed bottom-6 right-6 z-30 w-14 h-14 bg-accent-red hover:bg-accent-red/90 text-white rounded-full shadow-lg hover:shadow-xl transition-all duration-300 flex items-center justify-center group"
        :class="{ 'scale-110': fabPulse }" aria-label="Quick capture note">
        <svg class="w-6 h-6 transition-transform group-hover:rotate-90" fill="none" stroke="currentColor"
          viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
      </button>
    </Transition>

    <!-- Quick Capture Modal -->
    <QuickCaptureModal :is-open="isQuickCaptureOpen" @close="closeQuickCapture" @saved="handleAnnotationSaved" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'
import Navigation from './components/Navigation.vue'
import GlobalToast from './components/GlobalToast.vue'
import PWAUpdatePrompt from './components/PWAUpdatePrompt.vue'
import QuickCaptureModal from './components/QuickCaptureModal.vue'

const route = useRoute()
const authStore = useAuthStore()

const isQuickCaptureOpen = ref(false)
const fabPulse = ref(false)

// Show FAB only when user is authenticated and not on login page
const showFAB = computed(() => {
  return authStore.isAuthenticated && route.name !== 'login'
})

// Pulse animation on mount
onMounted(() => {
  setTimeout(() => {
    fabPulse.value = true
    setTimeout(() => {
      fabPulse.value = false
    }, 600)
  }, 1000)
})

const openQuickCapture = () => {
  isQuickCaptureOpen.value = true
}

const closeQuickCapture = () => {
  isQuickCaptureOpen.value = false
}

const handleAnnotationSaved = (annotation) => {
  // Could emit an event or update a global state if needed
  console.log('Annotation saved:', annotation)
}
</script>

<style scoped>
.fab-enter-active {
  animation: fabSlideIn 0.4s ease-out;
}

.fab-leave-active {
  animation: fabSlideOut 0.3s ease-in;
}

@keyframes fabSlideIn {
  from {
    transform: translateY(100px) scale(0);
    opacity: 0;
  }

  to {
    transform: translateY(0) scale(1);
    opacity: 1;
  }
}

@keyframes fabSlideOut {
  from {
    transform: translateY(0) scale(1);
    opacity: 1;
  }

  to {
    transform: translateY(100px) scale(0);
    opacity: 0;
  }
}
</style>
