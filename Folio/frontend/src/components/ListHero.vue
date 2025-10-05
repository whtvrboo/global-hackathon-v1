<template>
    <div class="relative w-full h-[60vh] text-white flex items-center justify-center text-center p-8"
        :style="{ backgroundColor: themeColor }">
        <div v-if="headerImageUrl" class="absolute inset-0">
            <img :src="headerImageUrl" :alt="title" class="w-full h-full object-cover" />
            <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/50 to-transparent"></div>
        </div>
        <div class="relative z-10">
            <h1 class="text-4xl md:text-6xl font-bold leading-tight mb-4">{{ title }}</h1>
            <p v-if="description" class="max-w-3xl mx-auto text-lg md:text-xl text-white/90 mb-6">{{ description }}</p>
            <div v-if="creator" class="relative flex items-center justify-center gap-3">
                <div class="flex items-center gap-3 cursor-pointer hover:bg-white/10 rounded-lg p-2 transition-all duration-200"
                    @mouseenter="showUserPopupImmediate" @mouseleave="hideUserPopup">
                    <img v-if="creator.picture" :src="creator.picture" :alt="creator.name"
                        class="w-12 h-12 rounded-full border-2 border-white/50" />
                    <div>
                        <p class="font-semibold">{{ creator.name }}</p>
                        <p class="text-sm text-white/70">@{{ creator.username }}</p>
                    </div>
                </div>

                <!-- User Hover Popup -->
                <UserHoverPopup :show="showUserPopup" :user="creator" position="bottom-right"
                    @close="showUserPopup = false" />
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import UserHoverPopup from './UserHoverPopup.vue'

defineProps({
    title: String,
    description: String,
    headerImageUrl: String,
    themeColor: {
        type: String,
        default: '#1a1a1a'
    },
    creator: Object
})

const showUserPopup = ref(false)
let hoverTimeout = null

const hideUserPopup = () => {
    // Add a small delay to prevent flickering when moving between elements
    hoverTimeout = setTimeout(() => {
        showUserPopup.value = false
    }, 100)
}

const showUserPopupImmediate = () => {
    // Clear any pending close timeout when showing popup
    if (hoverTimeout) {
        clearTimeout(hoverTimeout)
        hoverTimeout = null
    }
    showUserPopup.value = true
}
</script>
