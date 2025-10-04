import { ref, reactive, onMounted, onUnmounted } from 'vue'
import * as Y from 'yjs'
import { WebsocketProvider } from 'y-partykit/provider'
import type { Stem, Comment } from '@/data/dummyData'

export interface Version {
  id: string
  name: string
  stems: Stem[]
  comments: Comment[]
  createdAt: string
}

export interface TrackData {
  versions: Version[]
  currentVersionId: string
}

export function useCollaboration(roomId: string) {
  // Initialize Y.js document
  const yDoc = new Y.Doc()

  // Initialize PartyKit provider
  // For now, let's use a simple approach without real-time sync
  // We'll implement a basic state management without WebSocket
  const provider = null // Disable WebSocket for now

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
  const versions = ref<Version[]>([])
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
      const defaultVersion: Version = {
        id: 'v1',
        name: 'Version 1',
        stems: [],
        comments: [],
        createdAt: new Date().toISOString(),
      }

      trackData.set('versions', [defaultVersion])
      trackData.set('currentVersionId', 'v1')
    }
  }

  // Methods to modify the shared state
  const addStem = (stem: Stem) => {
    const currentVersion = versions.value.find((v) => v.id === currentVersionId.value)
    if (currentVersion) {
      const updatedStems = [...currentVersion.stems, stem]
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
      const updatedComments = [...currentVersion.comments, comment]
      const updatedVersion = { ...currentVersion, comments: updatedComments }

      const updatedVersions = versions.value.map((v) =>
        v.id === currentVersionId.value ? updatedVersion : v,
      )

      trackData.set('versions', updatedVersions)
    }
  }

  const createNewVersion = () => {
    const newVersionId = `v${versions.value.length + 1}`
    const currentVersion = versions.value.find((v) => v.id === currentVersionId.value)

    if (currentVersion) {
      const newVersion: Version = {
        id: newVersionId,
        name: `Version ${versions.value.length + 1}`,
        stems: [...currentVersion.stems], // Deep copy stems
        comments: [...currentVersion.comments], // Deep copy comments
        createdAt: new Date().toISOString(),
      }

      const updatedVersions = [...versions.value, newVersion]
      trackData.set('versions', updatedVersions)
      trackData.set('currentVersionId', newVersionId)
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

    // Initialize with default data if needed
    updateState()
    initializeDefaultData()
  })

  onUnmounted(() => {
    trackData.unobserve(updateState)
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
