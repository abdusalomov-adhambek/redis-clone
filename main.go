// Package main is the entry point of the Redis clone TCP server.
package main

import (
	"fmt"
	"goredisclone/helpers"
	"goredisclone/persistence"
	"log"
	"net"
)

// main initializes the server by loading persisted data from disk, starting
// the TCP listener on the configured port, launching a background goroutine
// that periodically removes expired keys, and then entering an accept loop
// that spawns a new goroutine for every incoming client connection.
func main() {
	port := ":8001" // TCP address the server binds to

	// Start persistence load in the background
	persistence.Load()

	// Start TCP listener on the given port
	listener, err := net.Listen("tcp", port) // TCP listener that accepts incoming connections
	if err != nil {
		panic(err)
	}

	log.Println("Redis clone connectect port", port)

	// Start background worker to clean up expired keys
	go helpers.CleanupWorker()

	// Accept incoming client connections in a loop
	for {
		conn, err := listener.Accept() // newly accepted client connection
		if err != nil {
			fmt.Println("connectoin error:", err)
			continue
		}

		fmt.Println("new client connected:", conn.LocalAddr())

		// Handle each client connection in a separate goroutine
		go helpers.HandleConnection(conn)
	}
}
