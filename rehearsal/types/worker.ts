// Minimal R2 types used by this project (to avoid bringing in workers-types)
export interface R2PresignedUrl {
  url: string
}

export interface R2Bucket {
  createPresignedUrl(
    method: 'PUT' | 'GET' | 'DELETE',
    key: string,
    options: { expiresIn: number },
  ): Promise<R2PresignedUrl>
}

export interface Env {
  R2_BUCKET: R2Bucket
  // GitHub OAuth
  GITHUB_CLIENT_ID: string
  GITHUB_CLIENT_SECRET: string
  // JWT secret for HS256 signing
  JWT_SECRET: string
  // Public app URL to redirect back to after OAuth (e.g., https://your-app.com)
  APP_URL: string
}
