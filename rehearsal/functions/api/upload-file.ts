import { Env } from '../../types/worker'

export const onRequest: PagesFunction<Env> = async (context) => {
  const { request, env } = context

  // Only allow POST requests
  if (request.method !== 'POST') {
    return new Response('Method not allowed', { status: 405 })
  }

  try {
    const formData = await request.formData()
    const file = formData.get('file') as File

    if (!file) {
      return new Response('No file provided', { status: 400 })
    }

    // Validate file type
    const allowedTypes = ['audio/wav', 'audio/mp3', 'audio/mpeg', 'audio/aiff', 'audio/aif']
    if (!allowedTypes.includes(file.type)) {
      return new Response('Invalid file type. Only audio files are allowed.', { status: 400 })
    }

    // For development: Create a mock file key and return a simple URL
    // In production, this would upload to R2
    const fileKey = `stems/${Date.now()}-${Math.random().toString(36).substring(2)}-${file.name}`

    // For development, we'll return a mock URL
    // In a real implementation, this would be the R2 URL
    const mockUrl = `https://rehearsal-stems.example.com/${fileKey}`

    return new Response(
      JSON.stringify({
        success: true,
        fileKey,
        url: mockUrl, // Mock URL for development
        size: file.size,
        type: file.type,
      }),
      {
        headers: {
          'Content-Type': 'application/json',
          'Cache-Control': 'no-cache, no-store, must-revalidate',
          Pragma: 'no-cache',
          Expires: '0',
        },
      },
    )
  } catch (error) {
    console.error('Error uploading file:', error)
    return new Response(
      JSON.stringify({
        success: false,
        error: 'Failed to upload file',
      }),
      {
        status: 500,
        headers: {
          'Content-Type': 'application/json',
        },
      },
    )
  }
}
