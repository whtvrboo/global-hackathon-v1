<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import TheHeader from '@/components/TheHeader.vue'
import CommitHistory from '@/components/CommitHistory.vue'
import WaveformStem from '@/components/WaveformStem.vue'
import CommentSidebar from '@/components/CommentSidebar.vue'
import { useCollaboration } from '@/composables/useCollaboration'
import { useAuthStore } from '@/stores/auth'
import { useRecorder } from '@/composables/useRecorder'
import type { Stem, Comment } from '@/data/dummyData'
import type { Commit as Version } from '@/composables/useCollaboration' // Alias for compatibility

const route = useRoute()
const trackId = route.params.id as string

// Use the collaboration composable
const { versions, currentVersionId, stems, comments, connectedUsers, createNewVersion, switchVersion, addStem, addComment, updateStem, trackData } = useCollaboration(trackId)
const auth = useAuthStore()
const recorder = useRecorder()

// File upload state
const isUploading = ref(false)
const uploadProgress = ref(0)

// Comment state
const showCommentDialog = ref(false)
const commentText = ref('')
const commentTimestamp = ref(0)
const commentStemId = ref('')
const commentAuthor = ref('Anonymous')
if (auth.user?.username) {
    commentAuthor.value = auth.user.username
}

// Version state
const showCommitDialog = ref(false)
const commitMessage = ref('')

// Master playback controls
const isPlaying = ref(false)
const stemRefs = ref<InstanceType<typeof WaveformStem>[]>([])

// Master play/pause functionality
const toggleMasterPlayback = async () => {
    if (isPlaying.value) {
        // Pause all stems
        stemRefs.value.forEach(stem => {
            stem.pause()
        })
        isPlaying.value = false
    } else {
        // Play all stems
        stemRefs.value.forEach(stem => {
            stem.play()
        })
        isPlaying.value = true
    }
}

const stopAll = () => {
    stemRefs.value.forEach(stem => {
        stem.stop()
    })
    isPlaying.value = false
}

const handleRecord = async () => {
    if (recorder.isRecording.value) {
        recorder.stop()
        // The rest of the logic (upload, etc.) will be handled by a watcher on recorder.audioBlob
    } else {
        await recorder.start()
        if (!recorder.error.value) {
            // Add a placeholder stem
            const placeholderStem: Stem = {
                id: `recording-${Date.now()}`,
                name: 'New Recording...',
                url: '', // No URL yet
                duration: 0,
                isMuted: false,
                isSolo: false,
                authorId: auth.user?.id,
                isRecording: true, // Custom flag
            }
            addStem(placeholderStem)
        } else {
            alert(recorder.error.value)
        }
    }
}

// Watch for the recording to finish
watch(() => recorder.audioBlob.value, (newBlob) => {
    if (newBlob) {
        uploadRecording(newBlob)
    }
})

// Register stem refs
const registerStem = (stemInstance: any) => {
    if (stemInstance && !stemRefs.value.includes(stemInstance)) {
        stemRefs.value.push(stemInstance)
    }
}

const unregisterStem = (stemInstance: any) => {
    const index = stemRefs.value.indexOf(stemInstance)
    if (index > -1) {
        stemRefs.value.splice(index, 1)
    }
}

// Comment handling
const handleAddComment = (timestamp: number, stemId: string) => {
    commentTimestamp.value = timestamp
    commentStemId.value = stemId
    commentText.value = ''
    showCommentDialog.value = true
}

const submitComment = () => {
    if (commentText.value.trim()) {
        const comment: Comment = {
            id: `comment-${Date.now()}-${Math.random().toString(36).substring(2)}`,
            timestamp: commentTimestamp.value,
            text: commentText.value.trim(),
            author: commentAuthor.value,
            createdAt: new Date().toISOString(),
        }

        addComment(comment)
        showCommentDialog.value = false
        commentText.value = ''
    }
}

const cancelComment = () => {
    showCommentDialog.value = false
    commentText.value = ''
}

// Version handling
const handleCreateNewCommit = () => {
    commitMessage.value = ''
    showCommitDialog.value = true
}

const confirmNewCommit = () => {
    if (commitMessage.value.trim()) {
        createNewVersion(commitMessage.value.trim())
        showCommitDialog.value = false
        commitMessage.value = ''
    }
}

const cancelNewCommit = () => {
    showCommitDialog.value = false
    commitMessage.value = ''
}

// Handle file upload
const handleFileUpload = async (event: Event) => {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]

    if (!file) return

    // Validate file type
    const allowedTypes = ['audio/wav', 'audio/mp3', 'audio/mpeg', 'audio/aiff', 'audio/aif']
    if (!allowedTypes.includes(file.type)) {
        alert('Please select a valid audio file (WAV, MP3, or AIFF)')
        return
    }

    // Validate file size (max 100MB)
    const maxSize = 100 * 1024 * 1024 // 100MB
    if (file.size > maxSize) {
        alert('File size must be less than 100MB')
        return
    }

    isUploading.value = true
    uploadProgress.value = 0

    try {
        // For development: Create a data URL from the file
        // This allows wavesurfer.js to load the audio directly
        const arrayBuffer = await file.arrayBuffer()
        const blob = new Blob([arrayBuffer], { type: file.type })
        const dataUrl = URL.createObjectURL(blob)

        console.log('Creating data URL for file:', file.name, 'Type:', file.type, 'Size:', file.size)

        // Create stem object and add to collaboration
        const stem: Stem = {
            id: `stem-${Date.now()}-${Math.random().toString(36).substring(2)}`,
            name: file.name,
            url: dataUrl, // Use the data URL for development
            duration: 0, // Will be updated when waveform is loaded
            isMuted: false,
            isSolo: false,
            authorId: auth.user?.id,
        }

        // Add stem to the current version
        addStem(stem)

        // Reset file input
        target.value = ''

        console.log('File uploaded successfully:', stem)
    } catch (error) {
        console.error('Upload failed:', error)
        alert('Upload failed. Please try again.')
    } finally {
        isUploading.value = false
        uploadProgress.value = 0
    }
}

const uploadRecording = async (blob: Blob) => {
    isUploading.value = true
    const recordingStem = stems.value.find(s => s.isRecording)
    if (!recordingStem) {
        console.error("Couldn't find placeholder stem for recording.")
        isUploading.value = false
        return
    }

    try {
        const fileName = `recording-${Date.now()}.webm`
        const fileType = 'audio/webm'

        // 1. Get presigned upload URL from our API
        const response = await fetch('/api/get-upload-url', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...auth.getAuthHeader(),
            },
            body: JSON.stringify({ fileName, fileType }),
        })

        if (!response.ok) throw new Error('Failed to get upload URL')
        const { uploadUrl, fileKey } = await response.json()

        // 2. Upload the audio blob to R2 using the presigned URL
        const uploadResponse = await fetch(uploadUrl, {
            method: 'PUT',
            body: blob,
            headers: {
                'Content-Type': fileType,
            },
        })

        if (!uploadResponse.ok) throw new Error('Upload to R2 failed')

        // 3. Construct the final public URL (adjust if you have a custom domain)
        const finalUrl = `https://rehearsal-stems.example.com/${fileKey}`

        // 4. Update the placeholder stem with the final URL
        updateStem(recordingStem.id, {
            url: finalUrl,
            name: `Rec ${new Date().toLocaleTimeString()}`,
            isRecording: false,
        })

    } catch (error) {
        console.error('Recording upload failed:', error)
        alert('Failed to save recording. Please try again.')
        // Here you might want to remove the placeholder stem
    } finally {
        isUploading.value = false
    }
}
</script>

<template>
    <div class="track-view min-h-screen bg-base-200">
        <!-- Header -->
        <TheHeader :users="connectedUsers" />

        <!-- Main Layout -->
        <div class="container mx-auto p-6">
            <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
                <!-- Central Column for Track -->
                <div class="lg:col-span-3">
                    <!-- Version Selector -->
                    <CommitHistory :versions="versions" :current-version-id="currentVersionId"
                        @switch-version="switchVersion" @create-new-version="handleCreateNewCommit" />

                    <!-- Master Controls -->
                    <div class="master-controls bg-base-100 border border-base-300 rounded-lg p-4 mb-6">
                        <div class="flex items-center gap-4">
                            <button @click="toggleMasterPlayback" class="btn btn-primary"
                                :class="{ 'btn-outline': !isPlaying }">
                                <svg v-if="!isPlaying" class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                    <path d="M8 5v10l8-5-8-5z" />
                                </svg>
                                <svg v-else class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                    <path d="M5 4h3v12H5V4zm7 0h3v12h-3V4z" />
                                </svg>
                                {{ isPlaying ? 'Pause All' : 'Play All' }}
                            </button>
                            <button @click="stopAll" class="btn btn-outline btn-sm">
                                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                    <path
                                        d="M10 18a8 8 0 100-16 8 8 0 000 16zM8 7a1 1 0 00-1 1v4a1 1 0 102 0V8a1 1 0 00-1-1zm4 0a1 1 0 00-1 1v4a1 1 0 102 0V8a1 1 0 00-1-1z" />
                                </svg>
                                Stop All
                            </button>
                            <div class="divider divider-horizontal"></div>
                            <button @click="handleRecord" class="btn btn-error"
                                :class="{ 'btn-outline': !recorder.isRecording.value }">
                                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                    <circle cx="10" cy="10" r="6" />
                                </svg>
                                {{ recorder.isRecording.value ? 'Stop' : 'Record' }}
                            </button>
                            <div class="text-sm text-base-content/70">
                                {{ stems.length }} stem{{ stems.length !== 1 ? 's' : '' }}
                            </div>
                        </div>
                    </div>

                    <!-- File Upload -->
                    <div class="upload-section bg-base-100 border border-base-300 rounded-lg p-4 mb-6">
                        <div class="flex items-center gap-4">
                            <input type="file" accept="audio/*" @change="handleFileUpload" :disabled="isUploading"
                                class="file-input file-input-bordered file-input-sm w-full max-w-xs" />
                            <div v-if="isUploading" class="flex items-center gap-2">
                                <span class="loading loading-spinner loading-sm"></span>
                                <span class="text-sm">Uploading...</span>
                            </div>
                        </div>
                        <p class="text-xs text-base-content/70 mt-2">
                            Supported formats: WAV, MP3, AIFF (max 100MB)
                        </p>
                    </div>

                    <!-- Stems -->
                    <div class="stems-container">
                        <WaveformStem v-for="stem in stems" :key="stem.id" :stem="stem" @mounted="registerStem"
                            @unmounted="unregisterStem" @add-comment="handleAddComment" />
                    </div>
                </div>

                <!-- Right Sidebar -->
                <div class="lg:col-span-1">
                    <CommentSidebar :comments="comments" />
                </div>
            </div>
        </div>

        <!-- Comment Dialog -->
        <div v-if="showCommentDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div class="bg-base-100 rounded-lg p-6 w-full max-w-md mx-4">
                <h3 class="font-semibold text-lg mb-4">Add Comment</h3>

                <div class="mb-4">
                    <label class="block text-sm font-medium mb-2">Timestamp: {{ commentTimestamp.toFixed(2) }}s</label>
                    <label class="block text-sm font-medium mb-2">Stem: {{stems.find(s => s.id === commentStemId)?.name
                        || 'Unknown'}}</label>
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium mb-2">Comment:</label>
                    <textarea v-model="commentText" class="textarea textarea-bordered w-full"
                        placeholder="Add your comment here..." rows="3" @keydown.enter.ctrl="submitComment"
                        @keydown.escape="cancelComment">
                    </textarea>
                </div>

                <div class="flex gap-2 justify-end">
                    <button @click="cancelComment" class="btn btn-outline btn-sm">Cancel</button>
                    <button @click="submitComment" class="btn btn-primary btn-sm" :disabled="!commentText.trim()">
                        Add Comment
                    </button>
                </div>

                <div class="text-xs text-base-content/60 mt-2">
                    Press Ctrl+Enter to submit, or Escape to cancel
                </div>
            </div>
        </div>

        <!-- Version Dialog -->
        <div v-if="showCommitDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div class="bg-base-100 rounded-lg p-6 w-full max-w-md mx-4">
                <h3 class="font-semibold text-lg mb-4">Create New Commit</h3>

                <div class="mb-4">
                    <p class="text-sm text-base-content/70 mb-3">
                        This will save a snapshot of the current track state.
                    </p>
                    <label class="block text-sm font-medium mb-2">Commit Message:</label>
                    <input v-model="commitMessage" class="input input-bordered w-full"
                        placeholder="e.g., Added bassline harmonies" @keydown.enter="confirmNewCommit"
                        @keydown.escape="cancelNewCommit">
                </div>

                <div class="flex gap-2 justify-end">
                    <button @click="cancelNewCommit" class="btn btn-outline btn-sm">Cancel</button>
                    <button @click="confirmNewCommit" class="btn btn-primary btn-sm" :disabled="!commitMessage.trim()">
                        Create Commit
                    </button>
                </div>

                <div class="text-xs text-base-content/60 mt-2">
                    Press Enter to create, or Escape to cancel
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.track-view {
    min-height: 100vh;
}

.stems-container {
    max-height: calc(100vh - 200px);
    overflow-y: auto;
}
</style>
