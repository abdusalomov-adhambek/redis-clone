package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
)

// DelHandler handles the DEL command.
// It deletes one or more keys and returns the count of keys that were actually removed.
func DelHandler(conn net.Conn, args []string) {
	variables.Mu.Lock()

	deleted := 0 // counts keys that were actually present and removed

	for _, key := range args {
		if _, exists := variables.Storage[key]; exists {
			delete(variables.Storage, key)
			delete(variables.Expirations, key)
			deleted++
		}
	}

	conn.Write([]byte(encode.EncodeInteger(deleted)))
	variables.Mu.Unlock()

	if deleted > 0 {
		persistence.Save()
	}
}
