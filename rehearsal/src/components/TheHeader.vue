<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import type { UserAwareness } from '@/composables/useCollaboration'

interface Props {
    users?: UserAwareness[]
}

defineProps<Props>()

const auth = useAuthStore()

const handleLogout = () => {
    auth.clear()
    // Redirect to home or refresh the page
    window.location.href = '/'
}

const handleGitHubLogin = () => {
    window.location.href = '/api/auth/github'
}
</script>

<template>
    <header class="bg-base-100 border-b border-base-300">
        <div class="container mx-auto px-6 py-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center gap-4">
                    <h1 class="text-xl font-bold">Rehearsal</h1>
                </div>

                <div class="flex items-center gap-4">
                    <!-- Connected Users Avatars -->
                    <div v-if="users && users.length > 0" class="flex items-center -space-x-2">
                        <div class="text-xs text-base-content/60 mr-2">
                            {{ users.length }} collaborator{{ users.length !== 1 ? 's' : '' }}
                        </div>
                        <div v-for="user in users" :key="user.name" class="avatar tooltip group" :data-tip="user.name">
                            <div
                                class="w-8 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2 transition-all duration-300 group-hover:scale-110 group-hover:ring-2">
                                <img v-if="user.avatar" :src="user.avatar" :alt="user.name" class="rounded-full" />
                                <div v-else
                                    class="w-full h-full bg-gradient-to-br from-primary/20 to-primary/40 rounded-full flex items-center justify-center">
                                    <span class="text-xs font-medium text-primary">{{ user.name.charAt(0).toUpperCase()
                                    }}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Login button when not authenticated -->
                    <div v-if="!auth.isAuthenticated">
                        <button @click="handleGitHubLogin" class="btn btn-primary btn-sm">
                            <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd"
                                    d="M10 0C4.477 0 0 4.484 0 10.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0110 4.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.203 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.942.359.31.678.921.678 1.856 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0020 10.017C20 4.484 15.522 0 10 0z"
                                    clip-rule="evenodd" />
                            </svg>
                            Sign in with GitHub
                        </button>
                    </div>

                    <!-- User dropdown when authenticated -->
                    <div v-if="auth.isAuthenticated" class="dropdown dropdown-end">
                        <label tabindex="0" class="btn btn-ghost btn-circle avatar">
                            <div class="w-10 rounded-full">
                                <img v-if="auth.user?.avatar" :src="auth.user.avatar" alt="User Avatar" />
                                <span v-else>{{ auth.user?.username?.charAt(0).toUpperCase() }}</span>
                            </div>
                        </label>
                        <ul tabindex="0"
                            class="mt-3 z-[1] p-2 shadow menu menu-sm dropdown-content bg-base-100 rounded-box w-52">
                            <li>
                                <a class="justify-between">
                                    Profile
                                    <span class="badge">New</span>
                                </a>
                            </li>
                            <li><a>Settings</a></li>
                            <li><a @click="handleLogout">Logout</a></li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </header>
</template>

<style scoped>
/* Header styles */
</style>
