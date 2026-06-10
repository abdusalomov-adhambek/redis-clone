package handlers

import (
	"fmt"
	"goredisclone/variables"
	"net"
)

func ExistsHandler(conn net.Conn, args []string) {
	if len(args) != 1 {
		_, _ = conn.Write([]byte("ERR wrong number of arguments for 'exists' command\r\n"))
		return
	}
	key := args[0]
	variables.Mu.Lock()
	_, exists := variables.Storage[key]
	if !exists {
		fmt.Fprintf(conn, ":0\r\n")
		variables.Mu.Unlock()
		return
	}
	fmt.Fprintf(conn, ":1\r\n")
	variables.Mu.Unlock()
}
