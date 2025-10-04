<script setup lang="ts">
import { ref, onMounted, onUnmounted, watchEffect, nextTick } from 'vue'
import WaveSurfer from 'wavesurfer.js'
import type { Stem } from '@/data/dummyData'

interface Props {
    stem: Stem
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
            responsive: true,
            height: 80,
            normalize: true,
            backend: 'WebAudio',
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
            const duration = wavesurfer.getDuration()
            if (duration > 0) {
                const timestamp = progress * duration
                emit('add-comment', timestamp, props.stem.id)
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
