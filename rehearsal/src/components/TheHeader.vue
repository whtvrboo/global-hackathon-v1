<script setup lang="ts">
// TheHeader component for the main navigation/header
</script>

<template>
    <header class="bg-base-100 border-b border-base-300 p-4">
        <div class="flex items-center justify-between">
            <h1 class="text-xl font-bold">Rehearsal</h1>
            <div class="flex items-center gap-4">
                <span class="text-sm text-base-content/70">Track ID: {{ $route.params.id }}</span>
            </div>
        </div>
    </header>
</template>

<style scoped>
/* Header styles */
</style>
