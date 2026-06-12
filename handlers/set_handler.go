package handlers

import (
	"goredisclone/encode"
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
		conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments")))
		return
	}

	key := args[0]   // key to store
	value := args[1] // value to associate with the key

	// Lock storage for safe concurrent access
	variables.Mu.Lock()

	// Handle optional EX flag for TTL: SET key value EX <seconds>
	if len(args) == 4 {
		if strings.ToUpper(args[2]) == "EX" {

			ttl, err := strconv.Atoi(args[3]) // TTL in seconds from the EX option
			if err != nil {
				conn.Write([]byte(encode.EncodeError("ERR invalid TTL")))
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

	conn.Write([]byte(encode.EncodeSimpleString("OK")))
}
