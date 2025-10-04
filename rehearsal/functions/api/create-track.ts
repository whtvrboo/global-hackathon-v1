export interface Env {
  // Environment variables can be added here later
}

export default {
  async fetch(request: Request, env: Env, ctx: any): Promise<Response> {
    // Handle CORS preflight requests
    if (request.method === 'OPTIONS') {
      return new Response(null, {
        status: 200,
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': 'POST, OPTIONS',
          'Access-Control-Allow-Headers': 'Content-Type',
        },
      })
    }

    // Only handle POST requests
    if (request.method !== 'POST') {
      return new Response(JSON.stringify({ error: 'Method not allowed' }), {
        status: 405,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Headers': 'Content-Type, Authorization',
        },
      })
    }

    try {
      // Validate Authorization header with JWT
      const { getBearerToken, verifyJwt } = await import('./_auth')
      const token = getBearerToken(request)
      if (!token) {
        return new Response(JSON.stringify({ error: 'Unauthorized' }), {
          status: 401,
          headers: { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' },
        })
      }
      const payloadJson = await verifyJwt(env as any, token)
      if (!payloadJson) {
        return new Response(JSON.stringify({ error: 'Unauthorized' }), {
          status: 401,
          headers: { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' },
        })
      }
      const userId = payloadJson.sub
      const username = payloadJson.username
      if (!userId || !username) {
        return new Response(JSON.stringify({ error: 'Invalid token payload' }), {
          status: 401,
          headers: { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' },
        })
      }

      // Generate a random trackId
      const trackId = crypto.randomUUID()

      // Log the owner and trackId
      console.log('Creating track:', { ownerId: userId, username, trackId })

      // Return success response
      return new Response(
        JSON.stringify({
          success: true,
          trackId,
          ownerId: userId,
        }),
        {
          status: 200,
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Headers': 'Content-Type, Authorization',
          },
        },
      )
    } catch (error) {
      console.error('Error creating track:', error)

      return new Response(
        JSON.stringify({
          error: 'Internal server error',
        }),
        {
          status: 500,
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*',
          },
        },
      )
    }
  },
}
