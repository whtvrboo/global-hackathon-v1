import { Env } from '../../../types/worker'

function generateRandomState(): string {
  const randomBytes = new Uint8Array(32)
  crypto.getRandomValues(randomBytes)
  // base64url encode
  const base64 = btoa(String.fromCharCode(...randomBytes))
  return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '')
}

export const onRequest: PagesFunction<Env> = async (context) => {
  const { request, env } = context

  const url = new URL(request.url)
  const origin = `${url.protocol}//${url.host}`
  const callbackUrl = `${origin}/api/auth/github/callback`

  const state = generateRandomState()

  const authorizeUrl = new URL('https://github.com/login/oauth/authorize')
  authorizeUrl.searchParams.set('client_id', env.GITHUB_CLIENT_ID)
  authorizeUrl.searchParams.set('redirect_uri', callbackUrl)
  authorizeUrl.searchParams.set('state', state)
  authorizeUrl.searchParams.set('scope', 'read:user user:email')

  const headers = new Headers({ Location: authorizeUrl.toString() })
  const cookie = `oauth_state=${state}; Path=/; HttpOnly; Secure; SameSite=Lax; Max-Age=600`
  headers.append('Set-Cookie', cookie)

  return new Response(null, { status: 302, headers })
}
