package handlers

import (
	"net"
	"sync"
)

func GetHandler(conn net.Conn, args []string, storage map[string]string, mu *sync.Mutex) {
	if len(args) != 1 {
		conn.Write([]byte("ERR wrong number of arguments\n"))
		return
	}

	key := args[0]
	mu.Lock()
	value, exists := storage[key]
	defer mu.Unlock()

	if !exists {
		conn.Write([]byte("ERR key not found\n"))
		return
	}

	conn.Write([]byte(value + "\n"))
}
