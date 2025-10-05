<template>
    <!-- Desktop Pill Sidebar -->
    <aside class="hidden lg:flex lg:w-20 lg:flex-col lg:fixed lg:inset-y-0 lg:z-50 lg:items-center">
        <div class="flex flex-col items-center mt-8 space-y-4">
            <!-- Logo -->
            <div
                class="w-12 h-12 bg-gradient-to-br from-accent-red to-accent-blue rounded-full flex items-center justify-center shadow-lg">
                <span class="text-white font-bold text-lg">F</span>
            </div>

            <!-- Navigation Pills -->
            <nav class="flex flex-col space-y-3">
                <router-link to="/"
                    class="group relative flex items-center justify-center w-12 h-12 rounded-full transition-all duration-200"
                    :class="route.path === '/' ? 'bg-accent-blue shadow-lg shadow-accent-blue/25' : 'bg-dark-800/50 hover:bg-dark-700/70 hover:scale-105'">
                    <svg class="h-5 w-5 transition-colors"
                        :class="route.path === '/' ? 'text-white' : 'text-dark-300 group-hover:text-white'" fill="none"
                        stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6">
                        </path>
                    </svg>
                    <!-- Tooltip -->
                    <div
                        class="absolute left-14 bg-dark-800 text-white text-sm px-3 py-2 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-50">
                        Home
                    </div>
                </router-link>

                <router-link v-if="authStore.isAuthenticated && authStore.user?.username"
                    :to="`/profile/${authStore.user.username}`"
                    class="group relative flex items-center justify-center w-12 h-12 rounded-full transition-all duration-200"
                    :class="route.path.includes('/profile/') ? 'bg-accent-red shadow-lg shadow-accent-red/25' : 'bg-dark-800/50 hover:bg-dark-700/70 hover:scale-105'">
                    <svg class="h-5 w-5 transition-colors"
                        :class="route.path.includes('/profile/') ? 'text-white' : 'text-dark-300 group-hover:text-white'"
                        fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z">
                        </path>
                    </svg>
                    <!-- Tooltip -->
                    <div
                        class="absolute left-14 bg-dark-800 text-white text-sm px-3 py-2 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-50">
                        Profile
                    </div>
                </router-link>
                <button v-else-if="authStore.isAuthenticated" @click="$router.push('/profile')"
                    class="group relative flex items-center justify-center w-12 h-12 rounded-full transition-all duration-200 bg-dark-800/50 hover:bg-dark-700/70 hover:scale-105">
                    <svg class="h-5 w-5 transition-colors text-dark-300 group-hover:text-white" fill="none"
                        stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z">
                        </path>
                    </svg>
                    <!-- Tooltip -->
                    <div
                        class="absolute left-14 bg-dark-800 text-white text-sm px-3 py-2 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-50">
                        Profile
                    </div>
                </button>

                <!-- Logout Button -->
                <button v-if="authStore.isAuthenticated" @click="authStore.logout()"
                    class="group relative flex items-center justify-center w-12 h-12 rounded-full bg-dark-800/50 hover:bg-red-600/20 hover:scale-105 transition-all duration-200">
                    <svg class="h-5 w-5 text-dark-300 group-hover:text-red-400 transition-colors" fill="none"
                        stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1">
                        </path>
                    </svg>
                    <!-- Tooltip -->
                    <div
                        class="absolute left-14 bg-dark-800 text-white text-sm px-3 py-2 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-50">
                        Logout
                    </div>
                </button>
            </nav>
        </div>
    </aside>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const authStore = useAuthStore()
</script>
