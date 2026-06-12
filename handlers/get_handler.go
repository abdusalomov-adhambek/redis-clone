package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"log"
	"net"
	"time"
)

// GetHandler handles the GET command.
// It retrieves the value associated with the given key from storage.
// Returns NULL if the key does not exist or has expired.
func GetHandler(conn net.Conn, args []string) {
	// Ensure at least one argument (key) is provided
	if len(args) < 1 {
		conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments")))
		return
	}

	key := args[0] // key to retrieve

	// Lock storage for safe concurrent access
	variables.Mu.Lock()
	log.Printf("GET %s\n", key)
	// Check if the key exists in storage
	value, exists := variables.Storage[key]
	if !exists {
		conn.Write([]byte(encode.EncodeNull()))
		variables.Mu.Unlock()
		return
	}

	log.Printf("GET %s = %s\n", key, value)
	// Check if the key has an expiration and whether it has expired
	expireAt, exists := variables.Expirations[key]
	if exists && time.Now().After(expireAt) {
		conn.Write([]byte(encode.EncodeNull()))

		delete(variables.Storage, key)
		delete(variables.Expirations, key)

		variables.Mu.Unlock()

		persistence.Save()
		return
	}

	conn.Write([]byte(encode.EncodeBulkString(value)))
	variables.Mu.Unlock()
}
