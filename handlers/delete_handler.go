package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
)

func DelHandler(conn net.Conn, args []string) {
	variables.Mu.Lock()

	deleted := 0

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
