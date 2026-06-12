package handlers

import (
	"goredisclone/encode"
	"goredisclone/variables"
	"net"
	"time"
)

// TTLHandler handles the TTL command.
// Returns the remaining time-to-live in seconds, -2 if the key does not exist,
// or null if the key exists but has no expiration set.
func TTLHandler(conn net.Conn, args []string) {
	variables.Mu.Lock()
	defer variables.Mu.Unlock()

	if len(args) != 1 {
		conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments")))
		return
	}

	key := args[0] // key to query TTL for
	if _, ok := variables.Storage[key]; !ok {
		conn.Write([]byte(encode.EncodeInteger(-2)))
		return
	}

	expireAt, exists := variables.Expirations[key] // registered expiration timestamp for the key
	if !exists {
		conn.Write([]byte(encode.EncodeNull()))
		return
	}

	ttl := int(time.Until(expireAt).Seconds()) // remaining seconds until the key expires

	conn.Write([]byte(encode.EncodeInteger(ttl)))
}
