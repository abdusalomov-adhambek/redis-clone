package handlers

import (
	"fmt"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
	"time"
)

// GetHandler handles the GET command.
// It retrieves the value associated with the given key from storage.
// Returns NULL if the key does not exist or has expired.
func GetHandler(conn net.Conn, args []string) {
	// Ensure at least one argument (key) is provided
	if len(args) < 1 {
		conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}

	key := args[0]

	// Lock storage for safe concurrent access
	variables.Mu.Lock()
	defer variables.Mu.Unlock()

	// Check if the key exists in storage
	value, exists := variables.Storage[key]
	if !exists {
		conn.Write([]byte("$-1\r\n"))
		return
	}

	// Check if the key has an expiration and whether it has expired
	expireAt, exists := variables.Expirations[key]
	if exists && time.Now().After(expireAt) {
		conn.Write([]byte("$-1\r\n"))
		delete(variables.Storage, key)
		delete(variables.Expirations, key)
		persistence.Save()
		return
	}

	fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)

}
