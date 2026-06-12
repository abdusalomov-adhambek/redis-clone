package handlers

import (
	"goredisclone/encode"
	"goredisclone/variables"
	"net"
	"time"
)

// ExistsHandler handles the EXISTS command.
// Returns 1 if the key exists and has not expired, 0 otherwise.
func ExistsHandler(conn net.Conn, args []string) {
	if len(args) != 1 {
		_, _ = conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments for 'exists' command")))
		return
	}
	key := args[0] // key to check for existence

	variables.Mu.Lock()
	defer variables.Mu.Unlock()

	_, exists := variables.Storage[key]
	if !exists {
		_, _ = conn.Write([]byte(encode.EncodeInteger(0)))
		return
	}

	expireAt, exists := variables.Expirations[key] // registered expiration timestamp, if any
	if exists && time.Now().After(expireAt) {
		delete(variables.Storage, key)
		delete(variables.Expirations, key)

		_, _ = conn.Write([]byte(encode.EncodeInteger(0)))
		return
	}

	conn.Write([]byte(encode.EncodeInteger(1)))
}
