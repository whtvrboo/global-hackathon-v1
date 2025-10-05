<template>
  <div id="app" class="min-h-screen bg-dark-950">
    <Navigation />
    <Sidebar />
    <router-view />
    <GlobalToast />
    <PWAUpdatePrompt />
  </div>
</template>

<script setup>
import Navigation from './components/Navigation.vue'
import Sidebar from './components/Sidebar.vue'
import GlobalToast from './components/GlobalToast.vue'
import PWAUpdatePrompt from './components/PWAUpdatePrompt.vue'
</script>
