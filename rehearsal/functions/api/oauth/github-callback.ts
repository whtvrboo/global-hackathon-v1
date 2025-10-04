import { Env } from '../../../types/worker'

async function exchangeCodeForToken(env: Env, code: string, redirectUri: string): Promise<string> {
  const body = new URLSearchParams({
    client_id: env.GITHUB_CLIENT_ID,
    client_secret: env.GITHUB_CLIENT_SECRET,
    code,
    redirect_uri: redirectUri,
  })
  const resp = await fetch('https://github.com/login/oauth/access_token', {
    method: 'POST',
    headers: { Accept: 'application/json' },
    body,
  })
  if (!resp.ok) throw new Error('GitHub token exchange failed')
  const json = await resp.json<any>()
  if (!json.access_token) throw new Error('No access token in response')
  return json.access_token as string
}

async function fetchGithubUser(
  accessToken: string,
): Promise<{ id: number; login: string; avatar_url?: string }> {
  const resp = await fetch('https://api.github.com/user', {
    headers: { Authorization: `Bearer ${accessToken}`, 'User-Agent': 'rehearsal-app' },
  })
  if (!resp.ok) throw new Error('GitHub user fetch failed')
  const json = await resp.json<any>()
  return { id: json.id, login: json.login, avatar_url: json.avatar_url }
}

function base64UrlEncode(data: ArrayBuffer): string {
  const bytes = new Uint8Array(data)
  let binary = ''
  for (let i = 0; i < bytes.length; i++) binary += String.fromCharCode(bytes[i])
  const base64 = btoa(binary)
  return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '')
}

async function signJwtHS256(secret: string, payload: Record<string, unknown>): Promise<string> {
  const header = { alg: 'HS256', typ: 'JWT' }
  const encoder = new TextEncoder()
  const headerB64 = base64UrlEncode(encoder.encode(JSON.stringify(header)))
  const payloadB64 = base64UrlEncode(encoder.encode(JSON.stringify(payload)))
  const toSign = `${headerB64}.${payloadB64}`
  const key = await crypto.subtle.importKey(
    'raw',
    encoder.encode(secret),
    { name: 'HMAC', hash: 'SHA-256' },
    false,
    ['sign'],
  )
  const signature = await crypto.subtle.sign('HMAC', key, encoder.encode(toSign))
  const sigB64 = base64UrlEncode(signature)
  return `${toSign}.${sigB64}`
}

function parseCookies(cookieHeader: string | null): Record<string, string> {
  const out: Record<string, string> = {}
  if (!cookieHeader) return out
  cookieHeader.split(';').forEach((part) => {
    const [k, v] = part.split('=').map((s) => s.trim())
    if (k) out[k] = decodeURIComponent(v || '')
  })
  return out
}

export const onRequest: PagesFunction<Env> = async (context) => {
  const { request, env } = context

  const url = new URL(request.url)
  const origin = `${url.protocol}//${url.host}`
  const redirectUri = `${origin}/api/auth/github/callback`

  const code = url.searchParams.get('code')
  const state = url.searchParams.get('state')
  if (!code || !state) return new Response('Missing code/state', { status: 400 })

  const cookies = parseCookies(request.headers.get('Cookie'))
  if (!cookies.oauth_state || cookies.oauth_state !== state) {
    return new Response('Invalid OAuth state', { status: 400 })
  }

  try {
    const token = await exchangeCodeForToken(env, code, redirectUri)
    const ghUser = await fetchGithubUser(token)

    const nowSeconds = Math.floor(Date.now() / 1000)
    const expiresIn = 60 * 60 * 24 * 7 // 7 days
    const jwt = await signJwtHS256(env.JWT_SECRET, {
      sub: String(ghUser.id),
      username: ghUser.login,
      avatar: ghUser.avatar_url,
      iss: 'rehearsal',
      iat: nowSeconds,
      exp: nowSeconds + expiresIn,
      provider: 'github',
    })

    const appRedirect = `${env.APP_URL}/auth/callback?token=${encodeURIComponent(jwt)}`

    const headers = new Headers({ Location: appRedirect })
    headers.append('Set-Cookie', 'oauth_state=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Lax')
    return new Response(null, { status: 302, headers })
  } catch (err) {
    console.error('OAuth callback error', err)
    return new Response('OAuth failed', { status: 500 })
  }
}
