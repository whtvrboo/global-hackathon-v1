import type { Env } from '../../types/worker'

function base64UrlToUint8Array(base64Url: string): Uint8Array {
  const base64 =
    base64Url.replace(/-/g, '+').replace(/_/g, '/') + '==='.slice((base64Url.length + 3) % 4)
  const binary = atob(base64)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i)
  return bytes
}

export type JwtPayload = {
  sub: string
  username?: string
  avatar?: string
  exp?: number
  iat?: number
  iss?: string
  provider?: string
  [key: string]: unknown
}

export async function verifyJwt(env: Env, token: string): Promise<JwtPayload | null> {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null
    const [headerB64, payloadB64, signatureB64] = parts

    // verify header
    const headerJson = JSON.parse(new TextDecoder().decode(base64UrlToUint8Array(headerB64)))
    if (headerJson.alg !== 'HS256' || headerJson.typ !== 'JWT') return null

    const key = await crypto.subtle.importKey(
      'raw',
      new TextEncoder().encode(env.JWT_SECRET),
      { name: 'HMAC', hash: 'SHA-256' },
      false,
      ['verify'],
    )
    const data = new TextEncoder().encode(`${headerB64}.${payloadB64}`)
    const signature = base64UrlToUint8Array(signatureB64)
    const ok = await crypto.subtle.verify('HMAC', key, signature, data)
    if (!ok) return null

    const payloadJson = JSON.parse(
      new TextDecoder().decode(base64UrlToUint8Array(payloadB64)),
    ) as JwtPayload
    if (payloadJson.exp && Math.floor(Date.now() / 1000) > payloadJson.exp) return null
    return payloadJson
  } catch {
    return null
  }
}

export function getBearerToken(request: Request): string | null {
  const header = request.headers.get('Authorization') || ''
  return header.startsWith('Bearer ') ? header.slice(7) : null
}
