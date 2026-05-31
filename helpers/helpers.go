package helpers

import (
	"fmt"
	"goredisclone/handlers"
	"goredisclone/variables"
	"io"
	"net"
	"strings"
	"time"
)

// handleConnection reads commands from a client connection and dispatches them.
// It runs in its own goroutine for each connected client.
func HandleConnection(conn net.Conn) {
	fmt.Println("new client connection: ", conn.RemoteAddr())

	defer conn.Close()

	for {
		buffer := make([]byte, 1024)

		// Read incoming data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client disconnected")
			} else {
				fmt.Println("read error:", err)
			}
			conn.Close()
			break
		}

		// Parse the raw message into a command and its arguments
		message := string(buffer[:n])
		command, args := ParseCommand(message)
		Dispatch(conn, command, args)
	}
}

// cleanupWorker runs in the background every 5 seconds.
// It scans all keys with expirations and deletes those that have expired.
func CleanupWorker() {
	for {
		time.Sleep(5 * time.Second)

		variables.Mu.Lock()
		for key, timeEx := range variables.Expirations {
			// Delete the key from both storage and expiration map if expired
			if time.Now().After(timeEx) {
				delete(variables.Expirations, key)
				delete(variables.Storage, key)
			}
		}
		variables.Mu.Unlock()
	}
}

// parseCommand splits a raw message into a command name and its arguments.
// The command is normalized to uppercase.
func ParseCommand(message string) (string, []string) {
	tokens := strings.Fields(message)

	if len(tokens) == 0 {
		return "", nil
	}

	command := strings.ToUpper(tokens[0])
	args := tokens[1:]

	return command, args
}

// dispatch routes the command to the appropriate handler.
// Unknown commands return an error response to the client.
func Dispatch(conn net.Conn, command string, args []string) {
	switch command {
	case "PING":
		handlers.PingHandler(conn)
	case "SET":
		handlers.SetHanlder(conn, args)
	case "GET":
		handlers.GetHandler(conn, args)
	default:
		conn.Write([]byte("Unknown command\r\n"))
	}
}
