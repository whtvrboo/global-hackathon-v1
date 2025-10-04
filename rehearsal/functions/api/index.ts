import { Env } from '../../types/worker'

export default {
  async fetch(request: Request, env: Env, ctx: any): Promise<Response> {
    const url = new URL(request.url)
    const pathname = url.pathname

    // Route to appropriate handler based on path
    if (pathname === '/api/create-track') {
      // Import and call create-track handler
      const { default: createTrackHandler } = await import('./create-track')
      return createTrackHandler.fetch(request, env, ctx)
    } else if (pathname === '/api/upload-file') {
      // Import and call upload-file handler
      const { onRequest: uploadFileHandler } = await import('./upload-file')
      return uploadFileHandler({ request, env, ctx })
    } else if (pathname === '/api/get-upload-url') {
      // Import and call get-upload-url handler
      const { onRequest: getUploadUrlHandler } = await import('./get-upload-url')
      return getUploadUrlHandler({ request, env, ctx })
    } else {
      return new Response('Not found', { status: 404 })
    }
  },
}
