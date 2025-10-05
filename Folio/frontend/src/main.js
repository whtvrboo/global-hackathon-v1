import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './style.css'

// PWA Service Worker Registration
import { registerSW } from 'virtual:pwa-register'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')

// Register service worker for PWA
if ('serviceWorker' in navigator) {
  registerSW({
    onNeedRefresh(updateSW) {
      // Show a notification to the user that a new version is available
      console.log('New content available. Please refresh.')
      // Trigger the update prompt component
      if (window.showPWAUpdatePrompt) {
        window.showPWAUpdatePrompt(updateSW)
      }
    },
    onOfflineReady() {
      // Show a notification that the app is ready to work offline
      console.log('App ready to work offline')
    },
  })
}

