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
                    <div v-if="users" class="flex items-center -space-x-2">
                        <div v-for="user in users" :key="user.name" class="avatar tooltip" :data-tip="user.name">
                            <div class="w-8 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                                <img v-if="user.avatar" :src="user.avatar" :alt="user.name" />
                                <span v-else class="text-xs">{{ user.name.charAt(0).toUpperCase() }}</span>
                            </div>
                        </div>
                    </div>

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
