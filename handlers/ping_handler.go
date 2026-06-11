// Package handlers contains one handler function per Redis command.
// Each handler receives the client connection and the parsed argument list,
// writes the appropriate RESP response back to the client, and returns.
package handlers

import (
	"goredisclone/encode"
	"net"
)

// PingHandler handles the PING command.
// Responds with the RESP simple-string "+PONG" to confirm the server is alive.
func PingHandler(conn net.Conn) {
	conn.Write([]byte(encode.EncodeSimpleString("PONG")))
}
