export interface Stem {
  id: string
  name: string
  url: string
  duration: number
  isMuted: boolean
  isSolo: boolean
  authorId?: string
  isRecording?: boolean
}

export interface Comment {
  id: string
  timestamp: number
  text: string
  author: string
  createdAt: string
  authorId?: string
}

export const hardcodedStems: Stem[] = [
  {
    id: '1',
    name: 'Vocal Track',
    url: '/audio/vocal.wav',
    duration: 180,
    isMuted: false,
    isSolo: false,
  },
  {
    id: '2',
    name: 'Guitar',
    url: '/audio/guitar.wav',
    duration: 180,
    isMuted: false,
    isSolo: false,
  },
  {
    id: '3',
    name: 'Bass',
    url: '/audio/bass.wav',
    duration: 180,
    isMuted: false,
    isSolo: false,
  },
  {
    id: '4',
    name: 'Drums',
    url: '/audio/drums.wav',
    duration: 180,
    isMuted: false,
    isSolo: false,
  },
]

export const hardcodedComments: Comment[] = [
  {
    id: '1',
    timestamp: 15.5,
    text: 'Love this intro! Maybe add some reverb?',
    author: 'Producer',
    createdAt: '2024-01-15T10:30:00Z',
  },
  {
    id: '2',
    timestamp: 45.2,
    text: 'The chorus needs more energy here',
    author: 'Songwriter',
    createdAt: '2024-01-15T10:32:00Z',
  },
  {
    id: '3',
    timestamp: 120.8,
    text: 'Perfect transition into the bridge',
    author: 'Producer',
    createdAt: '2024-01-15T10:35:00Z',
  },
]
