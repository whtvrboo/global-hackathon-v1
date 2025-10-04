import { Env } from '../../types/worker'

export const onRequest = async (context: { request: Request; env: Env; ctx: any }) => {
  const { request, env } = context

  // CORS preflight
  if (request.method === 'OPTIONS') {
    return new Response(null, {
      status: 200,
      headers: {
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'POST, OPTIONS',
        'Access-Control-Allow-Headers': 'Content-Type, Authorization',
      },
    })
  }
  // Only allow POST requests
  if (request.method !== 'POST') {
    return new Response('Method not allowed', {
      status: 405,
      headers: {
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Headers': 'Content-Type, Authorization',
      },
    })
  }

  try {
    // Require Authorization
    const { getBearerToken, verifyJwt } = await import('./_auth')
    const token = getBearerToken(request)
    const payload = token ? await verifyJwt(env as any, token) : null
    if (!payload)
      return new Response('Unauthorized', {
        status: 401,
        headers: { 'Access-Control-Allow-Origin': '*' },
      })

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
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Headers': 'Content-Type, Authorization',
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
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Headers': 'Content-Type, Authorization',
        },
      },
    )
  }
}
