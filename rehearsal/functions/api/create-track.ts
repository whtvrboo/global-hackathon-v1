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
        },
      })
    }

    try {
      // Parse the request body
      const body = await request.json()
      const { email } = body

      // Validate email
      if (!email || typeof email !== 'string') {
        return new Response(JSON.stringify({ error: 'Email is required' }), {
          status: 400,
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*',
          },
        })
      }

      // Generate a random trackId
      const trackId = crypto.randomUUID()

      // Log the email and trackId (as specified in requirements)
      console.log('Creating track:', { email, trackId })

      // Return success response
      return new Response(
        JSON.stringify({
          success: true,
          trackId,
        }),
        {
          status: 200,
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*',
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
