<script setup lang="ts">
import { ref, onMounted, onUnmounted, watchEffect, nextTick, computed } from 'vue'
import WaveSurfer from 'wavesurfer.js'
import type { Stem, Comment } from '@/data/dummyData'

interface Props {
    stem: Stem
    comments?: Comment[]
}

const props = defineProps<Props>()

// Template refs
const waveformEl = ref<HTMLDivElement>()
let wavesurfer: WaveSurfer | null = null

// Loading and error states
const isLoading = ref(false)
const hasError = ref(false)
const errorMessage = ref('')

// Define emits
const emit = defineEmits<{
    mounted: [instance: any]
    unmounted: [instance: any]
    'add-comment': [timestamp: number, stemId: string]
}>()

// Computed properties for comment markers
const commentMarkers = computed(() => {
    if (!props.comments || !wavesurfer) return []

    const duration = wavesurfer.getDuration()
    if (duration <= 0) return []

    return props.comments
        .filter(comment => comment.timestamp <= duration)
        .map(comment => ({
            ...comment,
            position: (comment.timestamp / duration) * 100
        }))
})

// Initialize wavesurfer in onMounted
onMounted(async () => {
    if (!waveformEl.value) return

    try {
        wavesurfer = WaveSurfer.create({
            container: waveformEl.value,
            waveColor: '#4f46e5',
            progressColor: '#7c3aed',
            cursorColor: '#ef4444',
            barWidth: 2,
            barRadius: 3,
            height: 80,
            normalize: true,
            mediaControls: false,
        })

        // Don't load audio if the stem is a recording placeholder
        if (props.stem.isRecording) {
            isLoading.value = true
            errorMessage.value = 'Recording in progress...'
            return
        }

        // Load audio if URL is available and valid
        if (props.stem.url && props.stem.url.trim() !== '') {
            try {
                // Check if it's a valid URL or data URL
                const isValidUrl = props.stem.url.startsWith('http') ||
                    props.stem.url.startsWith('blob:') ||
                    props.stem.url.startsWith('data:')

                if (isValidUrl) {
                    isLoading.value = true
                    hasError.value = false
                    await wavesurfer.load(props.stem.url)
                    isLoading.value = false
                } else {
                    console.warn('Invalid URL format for initial audio:', props.stem.url)
                    hasError.value = true
                    errorMessage.value = 'Invalid audio URL format'
                }
            } catch (error) {
                console.error('Failed to load initial audio:', error)
                isLoading.value = false
                hasError.value = true
                errorMessage.value = 'Failed to load audio file'
            }
        }

        // Handle mute state
        if (props.stem.isMuted) {
            wavesurfer.setMuted(true)
        }

        // Add click event listener for comments
        wavesurfer.on('click', (progress: number) => {
            if (wavesurfer) {
                const duration = wavesurfer.getDuration()
                if (duration > 0) {
                    const timestamp = progress * duration
                    emit('add-comment', timestamp, props.stem.id)
                }
            }
        })

        // Emit mounted event with the component instance
        emit('mounted', {
            play: () => wavesurfer?.play(),
            pause: () => wavesurfer?.pause(),
            stop: () => wavesurfer?.stop(),
            isPlaying: () => wavesurfer?.isPlaying() || false,
            setMuted: (muted: boolean) => wavesurfer?.setMuted(muted),
            setSolo: (solo: boolean) => {
                if (wavesurfer) {
                    wavesurfer.setMuted(!solo)
                }
            }
        })

    } catch (error) {
        console.error('Failed to initialize wavesurfer:', error)
    }
})

// Watch for URL changes and load new audio
watchEffect(async () => {
    if (wavesurfer && props.stem.url) {
        try {
            // Only try to load if the URL is valid and not empty
            if (props.stem.url && props.stem.url.trim() !== '') {
                // Check if it's a valid URL or data URL
                const isValidUrl = props.stem.url.startsWith('http') ||
                    props.stem.url.startsWith('blob:') ||
                    props.stem.url.startsWith('data:')

                if (isValidUrl) {
                    isLoading.value = true
                    hasError.value = false
                    await wavesurfer.load(props.stem.url)
                    isLoading.value = false
                } else {
                    console.warn('Invalid URL format for audio:', props.stem.url)
                    hasError.value = true
                    errorMessage.value = 'Invalid audio URL format'
                }
            }
        } catch (error) {
            console.error('Failed to load audio:', error)
            isLoading.value = false
            hasError.value = true
            errorMessage.value = 'Failed to load audio file'
        }
    }
})

// Watch for mute state changes
watchEffect(() => {
    if (wavesurfer) {
        wavesurfer.setMuted(props.stem.isMuted)
    }
})

// Cleanup on unmount
onUnmounted(() => {
    // Emit unmounted event
    emit('unmounted', {
        play: () => wavesurfer?.play(),
        pause: () => wavesurfer?.pause(),
        stop: () => wavesurfer?.stop(),
        isPlaying: () => wavesurfer?.isPlaying() || false,
        setMuted: (muted: boolean) => wavesurfer?.setMuted(muted),
        setSolo: (solo: boolean) => {
            if (wavesurfer) {
                wavesurfer.setMuted(!solo)
            }
        }
    })

    if (wavesurfer) {
        wavesurfer.destroy()
        wavesurfer = null
    }
})

// Expose wavesurfer instance for parent control
defineExpose({
    wavesurfer: () => wavesurfer,
    play: () => wavesurfer?.play(),
    pause: () => wavesurfer?.pause(),
    stop: () => wavesurfer?.stop(),
    isPlaying: () => wavesurfer?.isPlaying() || false,
    setMuted: (muted: boolean) => wavesurfer?.setMuted(muted),
    setSolo: (solo: boolean) => {
        // Solo logic would be handled by parent component
        // This is just for the individual stem
        if (wavesurfer) {
            wavesurfer.setMuted(!solo)
        }
    }
})

const toggleMute = () => {
    // TODO: Implement mute toggle logic
    console.log('Toggle mute for stem:', props.stem.id)
}

const toggleSolo = () => {
    // TODO: Implement solo toggle logic
    console.log('Toggle solo for stem:', props.stem.id)
}
</script>

<template>
    <div class="waveform-stem bg-base-100 border border-base-300 rounded-lg p-4 mb-4">
        <div class="flex items-center justify-between mb-3">
            <h3 class="font-semibold text-lg">{{ stem.name }}</h3>
            <div class="flex gap-2">
                <span v-if="stem.isRecording" class="badge badge-error animate-pulse">REC</span>
                <button @click="toggleMute" class="btn btn-sm" :class="stem.isMuted ? 'btn-error' : 'btn-outline'">
                    {{ stem.isMuted ? 'Unmute' : 'Mute' }}
                </button>
                <button @click="toggleSolo" class="btn btn-sm" :class="stem.isSolo ? 'btn-warning' : 'btn-outline'">
                    {{ stem.isSolo ? 'Unsolo' : 'Solo' }}
                </button>
            </div>
        </div>

        <!-- Waveform container -->
        <div ref="waveformEl" class="w-full h-24 bg-base-200 rounded border relative">
            <!-- WaveSurfer.js will be rendered here -->

            <!-- Comment markers -->
            <div v-for="marker in commentMarkers" :key="marker.id"
                class="absolute top-0 bottom-0 w-0.5 bg-primary/60 hover:bg-primary transition-colors cursor-pointer group"
                :style="{ left: `${marker.position}%` }" :title="`${marker.author}: ${marker.text}`">
                <!-- Comment marker dot -->
                <div
                    class="absolute -top-1 -left-1 w-3 h-3 bg-primary rounded-full border-2 border-base-100 shadow-sm group-hover:scale-110 transition-transform">
                </div>
                <!-- Comment preview on hover -->
                <div
                    class="absolute top-0 left-2 bg-base-100 border border-base-300 rounded-lg p-2 shadow-lg opacity-0 group-hover:opacity-100 transition-opacity z-10 min-w-48 max-w-64">
                    <div class="flex items-center gap-2 mb-1">
                        <div class="w-4 h-4 bg-primary/20 rounded-full flex items-center justify-center">
                            <span class="text-xs font-medium text-primary">{{ marker.author.charAt(0).toUpperCase()
                                }}</span>
                        </div>
                        <span class="text-xs font-medium">{{ marker.author }}</span>
                        <span class="text-xs text-base-content/60">{{ marker.timestamp.toFixed(1) }}s</span>
                    </div>
                    <p class="text-xs text-base-content/80">{{ marker.text }}</p>
                </div>
            </div>

            <!-- Loading state -->
            <div v-if="isLoading" class="absolute inset-0 flex items-center justify-center bg-base-200/80">
                <div class="flex items-center gap-2">
                    <span class="loading loading-spinner loading-sm"></span>
                    <span class="text-sm">Loading audio...</span>
                </div>
            </div>

            <!-- Error state -->
            <div v-if="hasError" class="absolute inset-0 flex items-center justify-center bg-error/10">
                <div class="text-center">
                    <div class="text-error text-sm font-medium">{{ errorMessage }}</div>
                    <div class="text-xs text-base-content/70 mt-1">Check console for details</div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.waveform-stem {
    min-height: 120px;
}
</style>
