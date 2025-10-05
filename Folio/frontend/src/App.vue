<template>
  <div id="app" class="min-h-screen bg-dark-950">
    <Navigation />
    <router-view />
    <GlobalToast />
    <PWAUpdatePrompt />
  </div>
</template>

<script setup>
import Navigation from './components/Navigation.vue'
import GlobalToast from './components/GlobalToast.vue'
import PWAUpdatePrompt from './components/PWAUpdatePrompt.vue'
</script>
