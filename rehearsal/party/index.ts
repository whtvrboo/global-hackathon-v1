import type * as Party from 'partykit/server'
import { onConnect } from 'y-partykit'

export default class YjsServer implements Party.Server {
  constructor(public room: Party.Room) {
    console.log('YjsServer created for room:', room.id)
  }

  async onRequest(req: Party.Request) {
    // Handle CORS for WebSocket connections
    if (req.method === 'OPTIONS') {
      return new Response(null, {
        status: 200,
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
          'Access-Control-Allow-Headers': 'Content-Type',
        },
      })
    }

    // For GET requests, return a simple response
    if (req.method === 'GET') {
      return new Response('PartyKit WebSocket Server', {
        status: 200,
        headers: {
          'Access-Control-Allow-Origin': '*',
        },
      })
    }

    return new Response('Method not allowed', { status: 405 })
  }

  onConnect(conn: Party.Connection) {
    console.log('WebSocket connection established for room:', this.room.id)
    return onConnect(conn, this.room, {
      // Optional: Add configuration options here
      // You can configure persistence, awareness, etc.
    })
  }

  onClose(conn: Party.Connection) {
    console.log('WebSocket connection closed for room:', this.room.id)
  }

  onError(conn: Party.Connection, error: Error) {
    console.error('WebSocket error for room:', this.room.id, error)
  }
}
