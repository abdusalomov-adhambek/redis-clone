package handlers

import (
	"fmt"
	"goredisclone/variables"
	"net"
	"time"
)

func TTLHandler(conn net.Conn, args []string) {
	variables.Mu.Lock()
	defer variables.Mu.Unlock()

	if len(args) != 1 {
		conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}

	key := args[0]
	if _, ok := variables.Storage[key]; !ok {
		fmt.Fprint(conn, ":-2\r\n")
		return
	}

	expireAt, exists := variables.Expirations[key]
	if !exists {
		fmt.Fprint(conn, ":-1\r\n")
		return
	}

	ttl := int(time.Until(expireAt).Seconds())
	fmt.Fprintf(conn, ":%d\r\n", ttl)
}
