package handlers

import (
	"fmt"
	"goredisclone/variables"
	"net"
)

func KeysHandler(conn net.Conn, args []string) {
	if len(args) != 1 {
		conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}

	if args[0] != "*" {
		conn.Write([]byte("-ERR only '*' pattern supported\r\n"))
		return
	}

	variables.Mu.Lock()

	keys := make([]string, 0, len(variables.Storage))
	for key := range variables.Storage {
		keys = append(keys, key)
	}

	variables.Mu.Unlock()

	fmt.Fprintf(conn, "*%d\r\n", len(keys))

	for _, key := range keys {
		fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(key), key)
	}
}
