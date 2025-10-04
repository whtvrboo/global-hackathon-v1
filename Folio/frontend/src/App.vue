<template>
  <div id="app" class="min-h-screen bg-dark-950">
    <router-view />
    <GlobalToast />
  </div>
</template>

<script setup>
import GlobalToast from './components/GlobalToast.vue'
</script>
