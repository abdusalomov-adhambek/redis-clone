package handlers

import (
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
	"strconv"
	"strings"
	"time"
)

// SetHanlder handles the SET command.
// It stores a key-value pair in storage.
// Optionally accepts EX <seconds> to set a TTL (expiration time).
func SetHanlder(conn net.Conn, args []string) {
	// Ensure at least two arguments (key and value) are provided
	if len(args) < 2 {
		conn.Write([]byte("ERR wrong number of arguments \n"))
		return
	}

	key := args[0]
	value := args[1]

	// Lock storage for safe concurrent access
	variables.Mu.Lock()

	// Handle optional EX flag for TTL: SET key value EX <seconds>
	if len(args) == 4 {
		if strings.ToUpper(args[2]) == "EX" {
			ttl, err := strconv.Atoi(args[3])
			if err != nil {
				conn.Write([]byte("ERR invalid TTL \n"))
				variables.Mu.Unlock()
				return
			}
			// Calculate and store the expiration timestamp
			expireAt := time.Now().Add(time.Duration(ttl) * time.Second)
			variables.Expirations[key] = expireAt
		}
	}

	// Store the key-value pair
	variables.Storage[key] = value
	variables.Mu.Unlock()

	persistence.Save()

	conn.Write([]byte("OK\n"))
}
