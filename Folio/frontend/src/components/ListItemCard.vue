<template>
    <div class="flex flex-col md:flex-row gap-6 md:gap-8 p-6 bg-dark-900 rounded-xl">
        <div class="flex-shrink-0">
            <img v-if="item.book.cover_url" :src="item.book.cover_url" :alt="item.book.title"
                class="w-32 md:w-40 h-auto object-cover rounded-lg shadow-lg mx-auto" />
            <div v-else class="w-32 md:w-40 h-48 md:h-56 bg-dark-800 rounded-lg flex items-center justify-center">
                <span class="text-3xl text-dark-400"></span>
            </div>
        </div>
        <div class="flex-1">
            <h3 class="text-2xl font-bold text-white mb-1">{{ item.book.title }}</h3>
            <p v-if="item.book.authors" class="text-lg text-dark-300 mb-4">by {{ item.book.authors.join(', ') }}</p>

            <blockquote v-if="item.notes" class="relative p-4 bg-dark-800 border-l-4 rounded-r-lg"
                :style="{ borderColor: themeColor }">
                <p class="text-white/90 italic">{{ item.notes }}</p>
            </blockquote>

            <p v-if="item.book.description" class="mt-4 text-sm text-dark-400 line-clamp-3">
                {{ item.book.description }}
            </p>
        </div>
    </div>
</template>

<script setup>
defineProps({
    item: Object,
    themeColor: {
        type: String,
        default: '#6366f1'
    }
})
</script>
