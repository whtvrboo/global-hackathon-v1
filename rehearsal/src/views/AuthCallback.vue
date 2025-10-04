<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

onMounted(() => {
    const token = String(route.query.token || '')
    if (token) {
        auth.setToken(token)
        router.replace('/')
    } else {
        router.replace('/')
    }
})
</script>

<template>
    <div class="min-h-screen flex items-center justify-center">
        <span class="loading loading-spinner loading-lg"></span>
    </div>
</template>
