package handlers

import "net"

// PingHandler handles the PING command.
// It responds with "PONG" to confirm the server is alive.
func PingHandler(conn net.Conn) {
	conn.Write([]byte("PONG\r\n"))
}
