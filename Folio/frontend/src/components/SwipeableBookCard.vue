<template>
    <div ref="cardRef" @touchstart="handleTouchStart" @touchmove="handleTouchMove" @touchend="handleTouchEnd"
        @mousedown="handleMouseDown" :style="{
            transform: `translate(${translateX}px, ${translateY}px) rotate(${rotation}deg)`,
            opacity: opacity,
            transition: isAnimating ? 'all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275)' : 'none'
        }" class="absolute inset-0 cursor-grab active:cursor-grabbing select-none">
        <div class="h-full bg-white rounded-3xl shadow-2xl overflow-hidden flex flex-col">
            <!-- Book Cover - Takes up most of the card -->
            <div class="relative flex-1 bg-gradient-to-br from-gray-100 to-gray-200">
                <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title" class="w-full h-full object-contain"
                    @load="imageLoaded = true" />
                <div v-else class="w-full h-full flex items-center justify-center text-8xl">
                    📚
                </div>

                <!-- Swipe Indicators -->
                <div class="absolute top-8 left-8 transform -rotate-12 transition-opacity duration-200"
                    :style="{ opacity: likeOpacity }">
                    <div
                        class="px-6 py-3 border-4 border-green-500 text-green-500 font-bold text-2xl rounded-xl bg-white/90">
                        LIKE
                    </div>
                </div>

                <div class="absolute top-8 right-8 transform rotate-12 transition-opacity duration-200"
                    :style="{ opacity: passOpacity }">
                    <div
                        class="px-6 py-3 border-4 border-red-500 text-red-500 font-bold text-2xl rounded-xl bg-white/90">
                        PASS
                    </div>
                </div>
            </div>

            <!-- Book Info -->
            <div class="p-6 bg-white">
                <h2 class="text-2xl font-bold text-gray-900 mb-2 line-clamp-2">
                    {{ book.title }}
                </h2>
                <p v-if="book.authors?.length" class="text-lg text-gray-600 mb-3 line-clamp-1">
                    {{ book.authors.join(', ') }}
                </p>

                <!-- Quick Stats -->
                <div class="flex items-center gap-4 text-sm text-gray-500 mb-4">
                    <div v-if="book.rating" class="flex items-center gap-1">
                        <span class="text-yellow-500">★</span>
                        <span class="font-medium">{{ book.rating }}</span>
                        <span v-if="book.ratings_count">({{ formatNumber(book.ratings_count) }})</span>
                    </div>
                    <div v-if="book.page_count" class="flex items-center gap-1">
                        📖 {{ book.page_count }} pages
                    </div>
                </div>

                <!-- Categories -->
                <div v-if="book.categories?.length" class="flex flex-wrap gap-2 mb-4">
                    <span v-for="(category, index) in book.categories.slice(0, 3)" :key="index"
                        class="px-3 py-1 text-xs font-medium bg-primary/10 text-primary rounded-full">
                        {{ category }}
                    </span>
                </div>

                <!-- Description Preview -->
                <p v-if="book.description" class="text-sm text-gray-700 line-clamp-3 leading-relaxed">
                    {{ book.description }}
                </p>
            </div>
        </div>

        <!-- Action Buttons (fallback for non-swipe) -->
        <div class="absolute bottom-24 left-0 right-0 flex justify-center gap-6 pointer-events-none">
            <button @click.stop="handlePass"
                class="pointer-events-auto w-16 h-16 rounded-full bg-white shadow-xl flex items-center justify-center text-2xl hover:scale-110 transition-transform active:scale-95">
                ✕
            </button>
            <button @click.stop="handleLike"
                class="pointer-events-auto w-16 h-16 rounded-full bg-primary shadow-xl flex items-center justify-center text-2xl hover:scale-110 transition-transform active:scale-95">
                ❤️
            </button>
            <button @click.stop="$emit('info')"
                class="pointer-events-auto w-16 h-16 rounded-full bg-white shadow-xl flex items-center justify-center text-2xl hover:scale-110 transition-transform active:scale-95">
                ℹ️
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
    book: {
        type: Object,
        required: true
    }
})

const emit = defineEmits(['swipe-left', 'swipe-right', 'info'])

const cardRef = ref(null)
const translateX = ref(0)
const translateY = ref(0)
const startX = ref(0)
const startY = ref(0)
const isDragging = ref(false)
const isAnimating = ref(false)
const imageLoaded = ref(false)

const rotation = computed(() => {
    return translateX.value / 20 // Subtle rotation based on drag
})

const opacity = computed(() => {
    const distance = Math.abs(translateX.value)
    return Math.max(1 - distance / 400, 0.5)
})

const likeOpacity = computed(() => {
    return translateX.value > 0 ? Math.min(translateX.value / 100, 1) : 0
})

const passOpacity = computed(() => {
    return translateX.value < 0 ? Math.min(Math.abs(translateX.value) / 100, 1) : 0
})

const handleTouchStart = (e) => {
    isDragging.value = true
    startX.value = e.touches[0].clientX
    startY.value = e.touches[0].clientY
}

const handleMouseDown = (e) => {
    isDragging.value = true
    startX.value = e.clientX
    startY.value = e.clientY

    const handleMouseMove = (e) => {
        if (!isDragging.value) return
        translateX.value = e.clientX - startX.value
        translateY.value = (e.clientY - startY.value) * 0.5 // Less vertical movement
    }

    const handleMouseUp = () => {
        document.removeEventListener('mousemove', handleMouseMove)
        document.removeEventListener('mouseup', handleMouseUp)
        handleDragEnd()
    }

    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
}

const handleTouchMove = (e) => {
    if (!isDragging.value) return
    translateX.value = e.touches[0].clientX - startX.value
    translateY.value = (e.touches[0].clientY - startY.value) * 0.5
}

const handleTouchEnd = () => {
    handleDragEnd()
}

const handleDragEnd = () => {
    isDragging.value = false

    const threshold = 100 // Minimum swipe distance

    if (translateX.value > threshold) {
        // Swiped right - Like
        animateSwipe('right')
    } else if (translateX.value < -threshold) {
        // Swiped left - Pass
        animateSwipe('left')
    } else {
        // Return to center
        isAnimating.value = true
        translateX.value = 0
        translateY.value = 0
        setTimeout(() => {
            isAnimating.value = false
        }, 300)
    }
}

const animateSwipe = (direction) => {
    isAnimating.value = true
    translateX.value = direction === 'right' ? 1000 : -1000
    translateY.value = translateY.value + 100

    setTimeout(() => {
        if (direction === 'right') {
            emit('swipe-right', props.book)
        } else {
            emit('swipe-left', props.book)
        }
    }, 300)
}

const handleLike = () => {
    animateSwipe('right')
}

const handlePass = () => {
    animateSwipe('left')
}

const formatNumber = (num) => {
    if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'k'
    }
    return num.toString()
}
</script>
