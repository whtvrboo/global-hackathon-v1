import { ref, reactive, onMounted, onUnmounted } from 'vue'
import * as Y from 'yjs'
import { WebsocketProvider } from 'y-partykit/provider'
import type { Stem, Comment } from '@/data/dummyData'
import { useAuthStore } from '@/stores/auth'

export interface UserAwareness {
  name: string
  avatar?: string
}

export interface Commit {
  id: string // This will serve as the unique hash
  message: string
  stems: Stem[]
  comments: Comment[]
  createdAt: string
  author: {
    id: string
    name: string
  }
}

export interface TrackData {
  versions: Commit[]
  currentVersionId: string
  ownerId?: string
}

export function useCollaboration(roomId: string) {
  // Initialize Y.js document
  const yDoc = new Y.Doc()
  const auth = useAuthStore()

  // Initialize PartyKit provider
  const provider = new WebsocketProvider(
    'ws://localhost:1999', // TODO: Make this configurable
    `track-${roomId}`,
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

  // Add debugging for WebSocket connection (disabled for now)
  if (provider) {
    provider.on('status', (event: any) => {
      console.log('WebSocket status:', event)
    })

    provider.on('connection-error', (event: any) => {
      console.error('WebSocket connection error:', event)
    })

    provider.on('connection-close', (event: any) => {
      console.log('WebSocket connection closed:', event)
    })
  }

  // Get the trackData Y.Map
  const trackData = yDoc.getMap('trackData')

  // Reactive Vue state
  const versions = ref<Commit[]>([])
  const currentVersionId = ref<string>('')
  const stems = ref<Stem[]>([])
  const comments = ref<Comment[]>([])

  // Update state from Y.js document
  const updateState = () => {
    const trackDataObj = trackData.toJSON() as TrackData

    if (trackDataObj.versions) {
      versions.value = trackDataObj.versions
    }

    if (trackDataObj.currentVersionId) {
      currentVersionId.value = trackDataObj.currentVersionId

      // Update stems and comments for current version
      const currentVersion = versions.value.find((v) => v.id === trackDataObj.currentVersionId)
      if (currentVersion) {
        stems.value = currentVersion.stems || []
        comments.value = currentVersion.comments || []
      }
    }
  }

  // Initialize with default data if empty
  const initializeDefaultData = () => {
    if (versions.value.length === 0) {
      const initialCommit: Commit = {
        id: crypto.randomUUID().slice(0, 7),
        message: 'Initial commit',
        stems: [],
        comments: [],
        createdAt: new Date().toISOString(),
        author: {
          id: auth.user?.id || 'system',
          name: auth.user?.username || 'System',
        },
      }

      trackData.set('versions', [initialCommit])
      trackData.set('currentVersionId', initialCommit.id)
      if (auth.user?.id) {
        trackData.set('ownerId', auth.user.id)
      }
    }
  }

  // Methods to modify the shared state
  const addStem = (stem: Stem) => {
    const currentVersion = versions.value.find((v) => v.id === currentVersionId.value)
    if (currentVersion) {
      const enriched = { ...stem, authorId: auth.user?.id }
      const updatedStems = [...currentVersion.stems, enriched as any]
      const updatedVersion = { ...currentVersion, stems: updatedStems }

      const updatedVersions = versions.value.map((v) =>
        v.id === currentVersionId.value ? updatedVersion : v,
      )

      trackData.set('versions', updatedVersions)
    }
  }

  const addComment = (comment: Comment) => {
    const currentVersion = versions.value.find((v) => v.id === currentVersionId.value)
    if (currentVersion) {
      const enriched = { ...comment, authorId: auth.user?.id }
      const updatedComments = [...currentVersion.comments, enriched as any]
      const updatedVersion = { ...currentVersion, comments: updatedComments }

      const updatedVersions = versions.value.map((v) =>
        v.id === currentVersionId.value ? updatedVersion : v,
      )

      trackData.set('versions', updatedVersions)
    }
  }

  const createNewVersion = (message: string) => {
    const newId = crypto.randomUUID().slice(0, 7)
    const currentVersion = versions.value.find((v) => v.id === currentVersionId.value)

    if (currentVersion && auth.user) {
      const newCommit: Commit = {
        id: newId,
        message: message,
        stems: JSON.parse(JSON.stringify(currentVersion.stems)), // Deep copy
        comments: JSON.parse(JSON.stringify(currentVersion.comments)), // Deep copy
        createdAt: new Date().toISOString(),
        author: {
          id: auth.user.id,
          name: auth.user.username || 'Anonymous',
        },
      }

      const updatedVersions = [...versions.value, newCommit]
      trackData.set('versions', updatedVersions)
      trackData.set('currentVersionId', newId)
    }
  }

  const switchVersion = (versionId: string) => {
    trackData.set('currentVersionId', versionId)
  }

  const updateStem = (stemId: string, updates: Partial<Stem>) => {
    const currentVersion = versions.value.find((v) => v.id === currentVersionId.value)
    if (currentVersion) {
      const updatedStems = currentVersion.stems.map((stem) =>
        stem.id === stemId ? { ...stem, ...updates } : stem,
      )
      const updatedVersion = { ...currentVersion, stems: updatedStems }

      const updatedVersions = versions.value.map((v) =>
        v.id === currentVersionId.value ? updatedVersion : v,
      )

      trackData.set('versions', updatedVersions)
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
    versions,
    currentVersionId,
    stems,
    comments,
    connectedUsers,

    // Methods
    addStem,
    addComment,
    createNewVersion,
    switchVersion,
    updateStem,

    // Internal state (for advanced usage)
    trackData,
  }
}
