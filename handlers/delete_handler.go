package handlers

import (
	"fmt"
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

	fmt.Fprintf(conn, ":%d\r\n", deleted)
	variables.Mu.Unlock()

	persistence.Save()
}
