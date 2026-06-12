package auth

import (
	"goredisclone/encode"
	"goredisclone/variables"
	"net"
)

// Auth handles the AUTH command. It validates the supplied password against
// the configured server password and marks the connection as authenticated on success.
func Auth(conn net.Conn, args []string) {
	if len(args) != 1 {
		conn.Write([]byte(encode.EncodeError("ERR no password provided")))
		return
	}

	if args[0] != variables.Password {
		conn.Write([]byte(encode.EncodeError("ERR wrong password")))
		return
	}

	variables.AuthenticatedClients[conn] = true
	conn.Write([]byte(encode.EncodeSimpleString("OK")))
}

// IsAuthenticated reports whether conn has previously passed AUTH.
func IsAuthenticated(conn net.Conn) bool {
	_, ok := variables.AuthenticatedClients[conn]
	return ok
}

// AuthRemove removes conn from the authenticated-clients map when a connection closes.
func AuthRemove(conn net.Conn) {
	delete(variables.AuthenticatedClients, conn)
}
