package handlers

import (
	"goredisclone/encode"
	"goredisclone/variables"
	"net"
)

// KeysHandler handles the KEYS command.
// Currently only the "*" pattern is supported; it returns all keys in storage.
func KeysHandler(conn net.Conn, args []string) {
	if len(args) != 1 {
		conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments")))
		return
	}

	if args[0] != "*" {
		conn.Write([]byte(encode.EncodeError("ERR only '*' pattern supported")))
		return
	}

	variables.Mu.Lock()

	keys := make([]string, 0, len(variables.Storage)) // result slice of all key names in storage
	for key := range variables.Storage {
		keys = append(keys, key)
	}

	conn.Write([]byte(encode.EncodeArray(keys)))

	variables.Mu.Unlock()
}
