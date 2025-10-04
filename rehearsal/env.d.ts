/// <reference types="vite/client" />

// Minimal Cloudflare Pages function type to satisfy TypeScript
declare type PagesFunction<Env = unknown> = (context: {
  request: Request
  env: Env
  ctx: any
}) => Promise<Response>
