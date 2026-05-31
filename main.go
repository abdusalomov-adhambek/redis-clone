package main

import (
	"fmt"
	"goredisclone/helpers"
	"goredisclone/persistence"
	"log"
	"net"
)

func main() {
	port := ":8001"

	// Start persistence load in the background
	persistence.Load()

	// Start TCP listener on the given port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	log.Println("Redis clone connectect port", port)

	// Start background worker to clean up expired keys
	go helpers.CleanupWorker()

	// Accept incoming client connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("connectoin error:", err)
			continue
		}

		fmt.Println("new client connected:", conn.LocalAddr())

		// Handle each client connection in a separate goroutine
		go helpers.HandleConnection(conn)
	}
}
