<template>
  <div id="app" class="min-h-screen bg-dark-950">
    <Navigation />
    <router-view />
    <GlobalToast />
  </div>
</template>

<script setup>
import Navigation from './components/Navigation.vue'
import GlobalToast from './components/GlobalToast.vue'
</script>
