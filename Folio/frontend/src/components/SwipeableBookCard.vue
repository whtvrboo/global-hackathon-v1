<template>
    <div class="relative w-full max-w-sm mx-auto">
        <!-- Swipeable Card Container -->
        <div ref="cardRef" class="card card-hover group relative overflow-hidden"
            :style="{ transform: `translateX(${translateX}px) rotate(${rotation}deg)` }" @touchstart="handleTouchStart"
            @touchmove="handleTouchMove" @touchend="handleTouchEnd" @mousedown="handleMouseDown"
            @mousemove="handleMouseMove" @mouseup="handleMouseEnd" @mouseleave="handleMouseEnd">
            <!-- Book Cover -->
            <div class="aspect-[2/3] relative overflow-hidden bg-dark-800 rounded-xl mb-4">
                <img v-if="book.cover_url" :src="book.cover_url" :alt="book.title"
                    class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                <div v-else class="w-full h-full flex items-center justify-center text-4xl text-dark-400">
                    📚
                </div>

                <!-- Status Badge -->
                <div v-if="book.status"
                    class="absolute top-3 right-3 px-3 py-1 text-xs font-semibold rounded-full glass-strong"
                    :class="statusColor(book.status)">
                    {{ statusLabel(book.status) }}
                </div>

                <!-- Rating Overlay -->
                <div v-if="book.rating"
                    class="absolute bottom-3 left-3 flex items-center gap-1 px-2 py-1 glass-strong rounded-full">
                    <span class="text-accent-orange text-sm">★</span>
                    <span class="text-white text-sm font-semibold">{{ book.rating }}</span>
                </div>
            </div>

            <!-- Book Info -->
            <div class="p-4">
                <h3 class="font-semibold text-white line-clamp-2 mb-2 text-lg leading-tight">
                    {{ book.title }}
                </h3>
                <p v-if="book.authors?.length" class="text-sm text-dark-400 line-clamp-1 mb-3">
                    by {{ book.authors.join(', ') }}
                </p>

                <!-- Description -->
                <p v-if="book.description" class="text-caption text-dark-300 line-clamp-3 mb-4">
                    {{ book.description }}
                </p>

                <!-- Action Buttons -->
                <div class="flex gap-2">
                    <button @click="$emit('like', book)" class="flex-1 btn-primary text-sm py-2">
                        ❤️ Like
                    </button>
                    <button @click="$emit('add-to-list', book)" class="flex-1 btn-secondary text-sm py-2">
                        📚 Add to List
                    </button>
                </div>
            </div>

            <!-- Swipe Indicators -->
            <div class="absolute inset-0 pointer-events-none">
                <!-- Like Indicator -->
                <div v-if="swipeDirection === 'right'"
                    class="absolute top-1/2 left-8 transform -translate-y-1/2 text-accent-green text-6xl font-bold opacity-50">
                    LIKE
                </div>

                <!-- Pass Indicator -->
                <div v-if="swipeDirection === 'left'"
                    class="absolute top-1/2 right-8 transform -translate-y-1/2 text-accent-red text-6xl font-bold opacity-50">
                    PASS
                </div>
            </div>
        </div>

        <!-- Swipe Actions -->
        <div class="flex justify-center gap-4 mt-4">
            <button @click="swipeLeft"
                class="w-12 h-12 bg-dark-800 hover:bg-accent-red/20 rounded-full flex items-center justify-center transition-colors">
                <svg class="w-6 h-6 text-accent-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12">
                    </path>
                </svg>
            </button>

            <button @click="swipeRight"
                class="w-12 h-12 bg-dark-800 hover:bg-accent-green/20 rounded-full flex items-center justify-center transition-colors">
                <svg class="w-6 h-6 text-accent-green" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                    </path>
                </svg>
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

const emit = defineEmits(['like', 'add-to-list', 'pass'])

const cardRef = ref(null)
const translateX = ref(0)
const rotation = ref(0)
const swipeDirection = ref(null)
const isDragging = ref(false)
const startX = ref(0)
const startY = ref(0)

const statusLabel = (status) => {
    const labels = {
        'want_to_read': 'Want to Read',
        'reading': 'Reading',
        'read': 'Read',
        'dnf': 'DNF'
    }
    return labels[status] || status
}

const statusColor = (status) => {
    const colors = {
        'want_to_read': 'text-accent-blue bg-accent-blue/20',
        'reading': 'text-accent-green bg-accent-green/20',
        'read': 'text-accent-purple bg-accent-purple/20',
        'dnf': 'text-dark-400 bg-dark-600/20'
    }
    return colors[status] || 'text-dark-400 bg-dark-600/20'
}

const handleTouchStart = (e) => {
    isDragging.value = true
    startX.value = e.touches[0].clientX
    startY.value = e.touches[0].clientY
}

const handleTouchMove = (e) => {
    if (!isDragging.value) return

    const currentX = e.touches[0].clientX
    const deltaX = currentX - startX.value

    translateX.value = deltaX
    rotation.value = deltaX * 0.1

    // Determine swipe direction
    if (Math.abs(deltaX) > 50) {
        swipeDirection.value = deltaX > 0 ? 'right' : 'left'
    } else {
        swipeDirection.value = null
    }
}

const handleTouchEnd = () => {
    if (!isDragging.value) return

    isDragging.value = false

    if (Math.abs(translateX.value) > 100) {
        if (translateX.value > 0) {
            swipeRight()
        } else {
            swipeLeft()
        }
    } else {
        resetPosition()
    }
}

const handleMouseDown = (e) => {
    isDragging.value = true
    startX.value = e.clientX
    startY.value = e.clientY
    e.preventDefault()
}

const handleMouseMove = (e) => {
    if (!isDragging.value) return

    const deltaX = e.clientX - startX.value

    translateX.value = deltaX
    rotation.value = deltaX * 0.1

    if (Math.abs(deltaX) > 50) {
        swipeDirection.value = deltaX > 0 ? 'right' : 'left'
    } else {
        swipeDirection.value = null
    }
}

const handleMouseEnd = () => {
    if (!isDragging.value) return

    isDragging.value = false

    if (Math.abs(translateX.value) > 100) {
        if (translateX.value > 0) {
            swipeRight()
        } else {
            swipeLeft()
        }
    } else {
        resetPosition()
    }
}

const swipeLeft = () => {
    translateX.value = -300
    rotation.value = -30
    swipeDirection.value = 'left'

    setTimeout(() => {
        emit('pass', props.book)
    }, 300)
}

const swipeRight = () => {
    translateX.value = 300
    rotation.value = 30
    swipeDirection.value = 'right'

    setTimeout(() => {
        emit('like', props.book)
    }, 300)
}

const resetPosition = () => {
    translateX.value = 0
    rotation.value = 0
    swipeDirection.value = null
}
</script>