import { ref, reactive, onMounted, onUnmounted } from 'vue'
import * as Y from 'yjs'
import { WebsocketProvider } from 'y-partykit/provider'
import type { Stem, Comment } from '@/data/dummyData'
import { useAuthStore } from '@/stores/auth'

// R2 Configuration
const R2_PUBLIC_URL = 'https://pub-a3525dca50814e51ab25628446c7acff.r2.dev' // Replace with your actual R2 public URL

export interface UserAwareness {
  name: string
  avatar?: string
}

export interface Take {
  id: string
  name: string
  description: string
  stems: Stem[]
  comments: Comment[]
  createdAt: string
  author: {
    id: string
    name: string
  }
}

export interface TrackData {
  takes: Take[]
  currentTakeId: string
  ownerId?: string
}

export function useCollaboration(roomId: string) {
  // Initialize Y.js document
  const yDoc = new Y.Doc()
  const auth = useAuthStore()

  // Initialize PartyKit provider
  const provider = new WebsocketProvider(
    'ws://127.0.0.1:1999/parties/rehearsal_partykit',
    roomId,
    yDoc,
    {
      params: {
        token: auth.token,
      },
    },
  )

  const awareness = provider?.awareness
  const connectedUsers = ref<UserAwareness[]>([])

  // Set local user's awareness state
  if (auth.user && awareness) {
    awareness.setLocalStateField('user', {
      name: auth.user.username,
      avatar: auth.user.avatar,
    })
  }

  // Listen for awareness changes
  const onAwarenessChange = () => {
    if (!awareness) return
    const users = []
    for (const [, state] of awareness.getStates()) {
      if (state.user) {
        users.push(state.user)
      }
    }
    connectedUsers.value = users
  }

  // Add debugging for WebSocket connection
  if (provider) {
    provider.on('connection-error', (event: any) => {
      console.error('WebSocket connection error:', event)
      isConnected.value = false
      // Retry connection after a short delay
      setTimeout(() => {
        if (provider) {
          provider.connect()
        }
      }, 1000)
    })

    provider.on('connection-close', (event: any) => {
      console.log('WebSocket connection closed:', event)
      isConnected.value = false
    })

    provider.on('connection-open', () => {
      console.log('WebSocket connection opened')
      isConnected.value = true
    })

    provider.on('status', (event: any) => {
      console.log('WebSocket status:', event)
      if (event.status === 'connected') {
        isConnected.value = true
      } else if (event.status === 'disconnected') {
        isConnected.value = false
      }
    })
  }

  // Get the trackData Y.Map
  const trackData = yDoc.getMap('trackData')

  // Reactive Vue state
  const takes = ref<Take[]>([])
  const currentTakeId = ref<string>('')
  const stems = ref<Stem[]>([])
  const comments = ref<Comment[]>([])
  const isConnected = ref(false)

  // Update state from Y.js document
  const updateState = () => {
    const trackDataObj = trackData.toJSON() as TrackData

    if (trackDataObj.takes) {
      takes.value = trackDataObj.takes
    }

    if (trackDataObj.currentTakeId) {
      currentTakeId.value = trackDataObj.currentTakeId

      // Update stems and comments for current take
      const currentTake = takes.value.find((t) => t.id === trackDataObj.currentTakeId)
      if (currentTake) {
        stems.value = currentTake.stems || []
        comments.value = currentTake.comments || []
      }
    }
  }

  // Initialize with default data if empty
  const initializeDefaultData = () => {
    if (takes.value.length === 0) {
      const initialTake: Take = {
        id: crypto.randomUUID().slice(0, 7),
        name: 'Take 1',
        description: 'Initial take',
        stems: [],
        comments: [],
        createdAt: new Date().toISOString(),
        author: {
          id: auth.user?.id || 'system',
          name: auth.user?.username || 'System',
        },
      }

      trackData.set('takes', [initialTake])
      trackData.set('currentTakeId', initialTake.id)
      if (auth.user?.id) {
        trackData.set('ownerId', auth.user.id)
      }
    }
  }

  // Methods to modify the shared state
  const addStem = (stem: Stem) => {
    const currentTake = takes.value.find((t) => t.id === currentTakeId.value)
    if (currentTake) {
      const enriched = {
        ...stem,
        authorId: auth.user?.id,
        authorName: auth.user?.username,
        authorAvatar: auth.user?.avatar,
      }
      const updatedStems = [...currentTake.stems, enriched as any]
      const updatedTake = { ...currentTake, stems: updatedStems }

      const updatedTakes = takes.value.map((t) => (t.id === currentTakeId.value ? updatedTake : t))

      trackData.set('takes', updatedTakes)
    }
  }

  const addComment = (comment: Comment) => {
    const currentTake = takes.value.find((t) => t.id === currentTakeId.value)
    if (currentTake) {
      const enriched = { ...comment, authorId: auth.user?.id }
      const updatedComments = [...currentTake.comments, enriched as any]
      const updatedTake = { ...currentTake, comments: updatedComments }

      const updatedTakes = takes.value.map((t) => (t.id === currentTakeId.value ? updatedTake : t))

      trackData.set('takes', updatedTakes)
    }
  }

  const createNewTake = (name: string, description: string) => {
    const newId = crypto.randomUUID().slice(0, 7)
    const currentTake = takes.value.find((t) => t.id === currentTakeId.value)

    if (currentTake && auth.user) {
      const newTake: Take = {
        id: newId,
        name: name,
        description: description,
        stems: JSON.parse(JSON.stringify(currentTake.stems)), // Deep copy
        comments: JSON.parse(JSON.stringify(currentTake.comments)), // Deep copy
        createdAt: new Date().toISOString(),
        author: {
          id: auth.user.id,
          name: auth.user.username || 'Anonymous',
        },
      }

      const updatedTakes = [...takes.value, newTake]
      trackData.set('takes', updatedTakes)
      trackData.set('currentTakeId', newId)
    }
  }

  const switchTake = (takeId: string) => {
    trackData.set('currentTakeId', takeId)
  }

  const updateStem = (stemId: string, updates: Partial<Stem>) => {
    const currentTake = takes.value.find((t) => t.id === currentTakeId.value)
    if (currentTake) {
      const updatedStems = currentTake.stems.map((stem) =>
        stem.id === stemId ? { ...stem, ...updates } : stem,
      )
      const updatedTake = { ...currentTake, stems: updatedStems }

      const updatedTakes = takes.value.map((t) => (t.id === currentTakeId.value ? updatedTake : t))

      trackData.set('takes', updatedTakes)
    }
  }

  // Upload and add stem function
  const uploadAndAddStem = async (
    file: File,
  ): Promise<{ success: true } | { success: false; error: string }> => {
    try {
      // Validate file type
      const allowedTypes = [
        'audio/wav',
        'audio/mp3',
        'audio/mpeg',
        'audio/aiff',
        'audio/aif',
        'audio/webm',
      ]
      if (!allowedTypes.includes(file.type)) {
        return {
          success: false,
          error: 'Please select a valid audio file (WAV, MP3, AIFF, or WebM)',
        }
      }

      // Validate file size (max 100MB)
      const maxSize = 100 * 1024 * 1024 // 100MB
      if (file.size > maxSize) {
        return { success: false, error: 'File size must be less than 100MB' }
      }

      // 1. Get presigned upload URL from our API
      const response = await fetch('/api/get-upload-url', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...auth.getAuthHeader(),
        },
        body: JSON.stringify({
          fileName: file.name,
          fileType: file.type,
        }),
      })

      if (!response.ok) {
        throw new Error(`Failed to get upload URL: ${response.status} ${response.statusText}`)
      }

      const { uploadUrl, fileKey } = await response.json()

      // 2. Upload the file to R2 using the presigned URL
      let publicUrl: string

      if (uploadUrl.includes('mock-upload.r2.dev')) {
        // Development: Mock upload - use object URL for local file
        console.log('Using mock upload for development')
        publicUrl = URL.createObjectURL(file)
      } else {
        // Production: Upload to R2
        const uploadResponse = await fetch(uploadUrl, {
          method: 'PUT',
          body: file,
          headers: {
            'Content-Type': file.type,
          },
        })

        if (!uploadResponse.ok) {
          throw new Error(
            `Upload to R2 failed: ${uploadResponse.status} ${uploadResponse.statusText}`,
          )
        }

        // Construct the final public URL
        publicUrl = `${R2_PUBLIC_URL}/${fileKey}`
      }

      // 4. Create stem object with proper attribution
      const stem: Stem = {
        id: `stem-${Date.now()}-${Math.random().toString(36).substring(2)}`,
        name: file.name,
        url: publicUrl,
        duration: 0, // Will be updated when waveform is loaded
        isMuted: false,
        isSolo: false,
        authorId: auth.user?.id,
        authorName: auth.user?.username || 'Anonymous',
        authorAvatar: auth.user?.avatar,
      }

      // 5. Add stem to the current take with full author object
      const currentTake = takes.value.find((t) => t.id === currentTakeId.value)
      if (currentTake) {
        const enriched = {
          ...stem,
          authorId: auth.user?.id,
          authorName: auth.user?.username,
          authorAvatar: auth.user?.avatar,
        }
        const updatedStems = [...currentTake.stems, enriched as any]
        const updatedTake = { ...currentTake, stems: updatedStems }

        const updatedTakes = takes.value.map((t) =>
          t.id === currentTakeId.value ? updatedTake : t,
        )

        trackData.set('takes', updatedTakes)
      }

      console.log('File uploaded successfully to R2:', stem)
      return { success: true }
    } catch (error) {
      console.error('Upload failed:', error)
      const errorMessage =
        error instanceof Error ? error.message : 'Upload failed. Please try again.'
      return { success: false, error: errorMessage }
    }
  }

  // Lifecycle hooks
  onMounted(() => {
    if (provider) {
      provider.connect()
    }

    // Observe changes to the Y.js document
    trackData.observe(updateState)
    if (awareness) {
      awareness.on('change', onAwarenessChange)
    }

    // Initialize with default data if needed
    updateState()
    initializeDefaultData()

    // Fallback: set connected state after a short delay if not already set
    setTimeout(() => {
      if (!isConnected.value && provider) {
        console.log('Setting connection state to true (fallback)')
        isConnected.value = true
      }
    }, 2000)
  })

  onUnmounted(() => {
    trackData.unobserve(updateState)
    if (awareness) {
      awareness.off('change', onAwarenessChange)
    }
    if (provider) {
      provider.disconnect()
    }
  })

  return {
    // Reactive state
    takes,
    currentTakeId,
    stems,
    comments,
    connectedUsers,
    isConnected,

    // Methods
    addStem,
    addComment,
    createNewTake,
    switchTake,
    updateStem,
    uploadAndAddStem,

    // Internal state (for advanced usage)
    trackData,
  }
}
