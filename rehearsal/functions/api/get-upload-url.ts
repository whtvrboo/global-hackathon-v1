import { Env } from '../../types/worker'

export const onRequest: PagesFunction<Env> = async (context) => {
  const { request, env } = context

  // Only allow POST requests
  if (request.method !== 'POST') {
    return new Response('Method not allowed', { status: 405 })
  }

  try {
    const { fileName, fileType } = await request.json()

    if (!fileName || !fileType) {
      return new Response('Missing fileName or fileType', { status: 400 })
    }

    // Validate file type (only allow audio files)
    const allowedTypes = ['audio/wav', 'audio/mp3', 'audio/mpeg', 'audio/aiff', 'audio/aif']
    if (!allowedTypes.includes(fileType)) {
      return new Response('Invalid file type. Only audio files are allowed.', { status: 400 })
    }

    // Generate a unique file key
    const fileKey = `stems/${Date.now()}-${Math.random().toString(36).substring(2)}-${fileName}`

    // Generate pre-signed upload URL for R2
    const uploadUrl = await env.R2_BUCKET.createPresignedUrl('PUT', fileKey, {
      expiresIn: 3600, // 1 hour
    })

    return new Response(
      JSON.stringify({
        success: true,
        uploadUrl: uploadUrl.url,
        fileKey,
        expiresAt: Date.now() + 3600000, // 1 hour from now
      }),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      },
    )
  } catch (error) {
    console.error('Error generating upload URL:', error)
    return new Response(
      JSON.stringify({
        success: false,
        error: 'Failed to generate upload URL',
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
