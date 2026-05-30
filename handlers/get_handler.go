package handlers

import (
	"goredisclone/variables"
	"net"
	"time"
)

func GetHandler(conn net.Conn, args []string) {
	if len(args) < 1 {
		conn.Write([]byte("ERR wrong number of arguments\n"))
		return
	}

	key := args[0]

	variables.Mu.Lock()
	defer variables.Mu.Unlock()

	value, exists := variables.Storage[key]
	if !exists {
		conn.Write([]byte("NULL\n"))
		return
	}

	expireAt, exists := variables.Expirations[key]
	if exists && time.Now().After(expireAt) {
		conn.Write([]byte("NULL\n"))
		delete(variables.Storage, key)
		delete(variables.Expirations, key)
		return
	}

	conn.Write([]byte(value + "\n"))
}
