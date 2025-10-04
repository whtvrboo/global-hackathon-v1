// Simple service worker to prevent caching of POST requests
// This helps avoid the "attempt-to-cache-non-get-request" error

self.addEventListener('fetch', (event) => {
  // Only handle GET requests for caching
  if (event.request.method !== 'GET') {
    // For non-GET requests, just pass through without caching
    return
  }

  // For GET requests, you can add caching logic here if needed
  // For now, we'll just pass through all requests
  event.respondWith(fetch(event.request))
})

// Prevent the service worker from trying to cache POST requests
self.addEventListener('install', (event) => {
  self.skipWaiting()
})

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim())
})
