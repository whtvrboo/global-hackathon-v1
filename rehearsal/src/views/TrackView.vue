<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import TheHeader from '@/components/TheHeader.vue'
import TakesExplorer from '@/components/TakesExplorer.vue'
import WaveformStem from '@/components/WaveformStem.vue'
import CommentSidebar from '@/components/CommentSidebar.vue'
import { useCollaboration } from '@/composables/useCollaboration'
import { useAuthStore } from '@/stores/auth'
import type { Stem, Comment } from '@/data/dummyData'
import type { Take } from '@/composables/useCollaboration'

const route = useRoute()
const trackId = route.params.id as string

// Use the collaboration composable
const { takes, currentTakeId, stems, comments, connectedUsers, isConnected, createNewTake, switchTake, addStem, addComment, updateStem, uploadAndAddStem, trackData } = useCollaboration(trackId)
const auth = useAuthStore()

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

// Take state
const showTakeDialog = ref(false)
const takeName = ref('')
const takeDescription = ref('')

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

// Take handling
const handleCreateNewTake = () => {
    takeName.value = ''
    takeDescription.value = ''
    showTakeDialog.value = true
}

const confirmNewTake = () => {
    if (takeName.value.trim() && takeDescription.value.trim()) {
        createNewTake(takeName.value.trim(), takeDescription.value.trim())
        showTakeDialog.value = false
        takeName.value = ''
        takeDescription.value = ''
    }
}

const cancelNewTake = () => {
    showTakeDialog.value = false
    takeName.value = ''
    takeDescription.value = ''
}

// Handle file upload
const handleFileUpload = async (event: Event) => {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]

    if (!file) return

    isUploading.value = true
    uploadProgress.value = 0

    const result = await uploadAndAddStem(file)

    if (result.success) {
        // Reset file input
        target.value = ''
    } else {
        alert(result.error)
    }

    isUploading.value = false
    uploadProgress.value = 0
}

</script>

<template>
    <div class="track-view min-h-screen bg-base-200">
        <!-- Header -->
        <TheHeader :users="connectedUsers" />

        <!-- Connection Status -->
        <div v-if="!isConnected" class="container mx-auto p-6">
            <div class="flex items-center justify-center min-h-[400px]">
                <div class="text-center">
                    <div class="loading loading-spinner loading-lg mb-4"></div>
                    <h2 class="text-xl font-semibold mb-2">Connecting to track...</h2>
                    <p class="text-base-content/70">Setting up collaboration room</p>
                </div>
            </div>
        </div>

        <!-- Main Layout -->
        <div v-else class="container mx-auto p-6">
            <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
                <!-- Central Column for Track -->
                <div class="lg:col-span-3">
                    <!-- Takes Explorer -->
                    <TakesExplorer :takes="takes" :current-take-id="currentTakeId" @switch-take="switchTake"
                        @create-new-take="handleCreateNewTake" />

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
                        <WaveformStem v-for="stem in stems" :key="stem.id" :stem="stem" :comments="comments"
                            @mounted="registerStem" @unmounted="unregisterStem" @add-comment="handleAddComment" />
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

        <!-- Take Dialog -->
        <div v-if="showTakeDialog"
            class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 backdrop-blur-sm">
            <div class="bg-base-100 rounded-lg p-6 w-full max-w-md mx-4 shadow-2xl border border-base-300">
                <div class="flex items-center gap-3 mb-4">
                    <div class="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4">
                            </path>
                        </svg>
                    </div>
                    <div>
                        <h3 class="font-semibold text-lg">Save as New Take</h3>
                        <p class="text-sm text-base-content/70">Capture this moment in your creative journey</p>
                    </div>
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium mb-2">Take Name:</label>
                    <input v-model="takeName" class="input input-bordered w-full focus:input-primary"
                        placeholder="e.g., Take 2, Chorus Version, Final Mix" autofocus @keydown.enter="confirmNewTake"
                        @keydown.escape="cancelNewTake">
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium mb-2">Describe this Take:</label>
                    <textarea v-model="takeDescription" class="textarea textarea-bordered w-full focus:textarea-primary"
                        placeholder="What makes this take special? What did you change or try?"
                        @keydown.enter.ctrl="confirmNewTake" @keydown.escape="cancelNewTake" rows="3">
                    </textarea>
                </div>

                <!-- Current state preview -->
                <div class="mb-4 p-3 bg-base-200 rounded-lg">
                    <div class="text-xs text-base-content/60 mb-2">Current state:</div>
                    <div class="flex items-center gap-2 text-sm">
                        <span class="badge badge-sm badge-outline">{{ stems.length }} stem{{ stems.length !== 1 ? 's' :
                            '' }}</span>
                        <span class="badge badge-sm badge-outline">{{ comments.length }} comment{{ comments.length !== 1
                            ? 's' : '' }}</span>
                    </div>
                </div>

                <div class="flex gap-2 justify-end">
                    <button @click="cancelNewTake" class="btn btn-outline btn-sm">Cancel</button>
                    <button @click="confirmNewTake" class="btn btn-primary btn-sm"
                        :disabled="!takeName.trim() || !takeDescription.trim()">
                        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7">
                            </path>
                        </svg>
                        Save as New Take
                    </button>
                </div>

                <div class="text-xs text-base-content/60 mt-3 text-center">
                    Press Ctrl+Enter to save, or Escape to cancel
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
