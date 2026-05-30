package handlers

import (
	"net"
	"sync"
)

func SetHanlder(conn net.Conn, args []string, storage map[string]string, mu *sync.Mutex) {
	if len(args) != 2 {
		conn.Write([]byte("ERR wrong number of arguments \n"))
		return
	}

	key := args[0]
	value := args[1]
	mu.Lock()
	storage[key] = value
	defer mu.Unlock()
	conn.Write([]byte("OK\n"))

}
