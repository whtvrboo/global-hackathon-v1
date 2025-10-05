<template>
    <!-- Header Navigation -->
    <header v-if="shouldShowNavigation" class="glass-strong border-b border-dark-800 sticky top-0 z-50">
        <div class="container-mobile max-w-7xl mx-auto">
            <div class="flex items-center justify-between py-4">
                <!-- Logo -->
                <router-link to="/" class="flex items-center gap-3">
                    <div
                        class="w-8 h-8 bg-gradient-to-br from-accent-red to-accent-blue rounded-lg flex items-center justify-center">
                        <span class="text-white font-bold text-sm">F</span>
                    </div>
                    <h1 class="text-xl font-bold text-white">Folio</h1>
                </router-link>

                <!-- Desktop Navigation -->
                <nav class="hidden md:flex items-center gap-6">
                    <router-link to="/discover" class="btn-ghost">Discover</router-link>
                    <router-link v-if="authStore.isAuthenticated" to="/feed" class="btn-ghost">Feed</router-link>
                    <router-link v-if="authStore.isAuthenticated" :to="`/profile/${authStore.user?.username}`"
                        class="btn-ghost">Profile</router-link>
                    <router-link v-if="!authStore.isAuthenticated" to="/login" class="btn-primary">Login</router-link>
                    <button v-if="authStore.isAuthenticated" @click="authStore.logout()"
                        class="btn-ghost">Logout</button>
                </nav>

                <!-- Mobile Menu Button -->
                <button class="md:hidden p-2 text-dark-300 hover:text-white transition-colors">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M4 6h16M4 12h16M4 18h16"></path>
                    </svg>
                </button>
            </div>
        </div>
    </header>

    <!-- Bottom Navigation (Mobile) -->
    <nav v-if="shouldShowNavigation && authStore.isAuthenticated"
        class="md:hidden fixed bottom-0 left-0 right-0 glass-strong border-t border-dark-800 z-40">
        <div class="flex items-center justify-around py-2">
            <router-link to="/" class="flex flex-col items-center gap-1 p-3 transition-colors"
                :class="isActive('/') ? 'text-accent-red' : 'text-dark-400 hover:text-white'">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6">
                    </path>
                </svg>
                <span class="text-xs">Home</span>
            </router-link>

            <router-link to="/discover" class="flex flex-col items-center gap-1 p-3 transition-colors"
                :class="isActive('/discover') ? 'text-accent-blue' : 'text-dark-400 hover:text-white'">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                </svg>
                <span class="text-xs">Discover</span>
            </router-link>

            <router-link to="/feed" class="flex flex-col items-center gap-1 p-3 transition-colors"
                :class="isActive('/feed') ? 'text-accent-red' : 'text-dark-400 hover:text-white'">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M4 6h16M4 10h16M4 14h16M4 18h16"></path>
                </svg>
                <span class="text-xs">Feed</span>
            </router-link>

            <router-link :to="`/profile/${authStore.user?.username}`"
                class="flex flex-col items-center gap-1 p-3 transition-colors"
                :class="isActive(`/profile/${authStore.user?.username}`) ? 'text-accent-red' : 'text-dark-400 hover:text-white'">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                </svg>
                <span class="text-xs">Profile</span>
            </router-link>
        </div>
    </nav>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { computed } from 'vue'

const route = useRoute()
const authStore = useAuthStore()

// Routes where navigation should not be shown
const noNavigationRoutes = ['/login', '/auth/callback']

const shouldShowNavigation = computed(() => {
    return !noNavigationRoutes.includes(route.path)
})

const isActive = (path) => {
    return route.path === path
}
</script>
