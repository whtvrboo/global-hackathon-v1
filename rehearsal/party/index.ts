import type * as Party from 'partykit/server'
import { onConnect } from 'y-partykit'
import { verifyJwt } from './auth'

export default class YjsServer implements Party.Server {
  constructor(public room: Party.Room) {}

  async onConnect(conn: Party.Connection, ctx: Party.ConnectionContext) {
    const token = new URL(conn.uri).searchParams.get('token')
    if (!token) {
      conn.close(1002, 'Authentication failed')
      return
    }

    try {
      const payload = await verifyJwt(this.room.env.JWT_SECRET as string, token)
      if (!payload) {
        conn.close(1002, 'Authentication failed')
        return
      }

      // Set user state on the connection
      conn.setState({ user: { id: payload.sub, username: payload.username } })
    } catch (err) {
      console.error('Token verification error:', err)
      conn.close(1002, 'Authentication failed')
      return
    }

    return onConnect(conn, this.room, {
      // Persistence can be enabled by setting the `persist` option
      // persist: true,
    })
  }
}
