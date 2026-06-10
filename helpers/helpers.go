package helpers

import (
	"bufio"
	"fmt"
	"goredisclone/handlers"
	"goredisclone/persistence"
	"goredisclone/variables"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// handleConnection reads commands from a client connection and dispatches them.
// It runs in its own goroutine for each connected client.
func HandleConnection(conn net.Conn) {
	fmt.Println("new client connection: ", conn.RemoteAddr())

	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		args, err := ParseRESP(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client disconnected")
				return
			} else {
				fmt.Println("parse error:", err)
			}
			return
		}

		if len(args) == 0 {
			continue
		}

		command := strings.ToUpper(args[0])
		Dispatch(conn, command, args[1:])
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
		persistence.Save()

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
	log.Printf("Dispatching command: %s %#v\n", command, args)
	switch command {
	case "PING":
		handlers.PingHandler(conn)
	case "SET":
		handlers.SetHanlder(conn, args)
	case "GET":
		handlers.GetHandler(conn, args)
	case "DEL":
		handlers.DelHandler(conn, args)
	case "EXPIRE":
		handlers.ExpireHandler(conn, args)
	case "TTL":
		handlers.TTLHandler(conn, args)
	case "EXISTS":
		handlers.ExistsHandler(conn, args)
	default:
		conn.Write([]byte("-ERR unknown command\r\n"))
	}
}

// ParseRESP reads one RESP (Redis Serialization Protocol) array command from
// the buffered reader and returns its elements as a string slice.
//
// RESP array format:
//
//	*<count>\r\n          — number of elements
//	$<length>\r\n         — byte length of the next element
//	<data>\r\n            — the element itself
//
// The function returns an error if the stream is malformed or if the
// connection is closed mid-message (io.EOF).
func ParseRESP(reader *bufio.Reader) ([]string, error) {
	// Read the first line: must start with '*' followed by the element count.
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)

	if len(line) == 0 || line[0] != '*' {
		return nil, fmt.Errorf("expected array")
	}

	// Parse the number of bulk-string elements that follow.
	count, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	if count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}

	args := make([]string, 0, count)
	for len(args) < count {
		// Each element begins with a bulk-string header: '$<length>\r\n'.
		bulkHeader, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		bulkHeader = strings.TrimSpace(bulkHeader)

		if len(bulkHeader) == 0 || bulkHeader[0] != '$' {
			return nil, fmt.Errorf("expected bulk header")
		}

		// Parse the byte length declared in the bulk-string header.
		length, err := strconv.Atoi(bulkHeader[1:])
		if err != nil {
			return nil, err
		}

		if length <= 0 {
			return nil, fmt.Errorf("length must be positive")
		}

		// Read exactly <length> data bytes plus the trailing '\r\n'.
		data := make([]byte, length+2)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			return nil, err
		}

		// Append only the actual data, stripping the trailing '\r\n'.
		args = append(args, string(data[:length]))
	}

	return args, nil
}
