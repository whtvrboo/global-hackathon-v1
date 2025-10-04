import type * as Party from 'partykit/server'
import { onConnect } from 'y-partykit'
import { verifyJwt } from './auth'

export default class YjsServer implements Party.Server {
  constructor(public room: Party.Room) {}

  async onConnect(conn: Party.Connection) {
    const token = new URL(conn.uri).searchParams.get('token')
    if (!token) {
      console.log('Connection rejected: Missing token')
      return conn.close(1002, 'Authentication failed')
    }

    try {
      const payload = await verifyJwt(this.room.env.JWT_SECRET as string, token)
      if (!payload) {
        console.log('Connection rejected: Invalid token')
        return conn.close(1002, 'Authentication failed')
      }
      // You can attach the user info to the connection for later use
      // conn.setState({ user: { id: payload.sub, username: payload.username } });
    } catch (err) {
      console.error('Token verification error:', err)
      return conn.close(1002, 'Authentication failed')
    }

    console.log('Authenticated connection established for room:', this.room.id)
    return onConnect(conn, this.room, {
      // Persistence can be enabled by setting the `persist` option
      // persist: true,
    })
  }
}
