package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
	"strconv"
	"time"
)

// ExpireHandler handles the EXPIRE command.
// Sets a TTL (in seconds) on an existing key; returns 1 on success, 0 if the key does not exist.
func ExpireHandler(conn net.Conn, args []string) {
	if len(args) != 2 {
		conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments")))
		return
	}

	key := args[0] // key to set expiration on

	ttl, err := strconv.Atoi(args[1]) // TTL in seconds
	if err != nil {
		conn.Write([]byte(encode.EncodeError("ERR wrong expire time")))
		return
	}

	if ttl <= 0 {
		conn.Write([]byte(encode.EncodeError("ERR expire time must be positive")))
		return
	}

	variables.Mu.Lock()

	_, exists := variables.Storage[key] // verify the key exists before applying an expiration
	if !exists {
		variables.Mu.Unlock()
		conn.Write([]byte(encode.EncodeInteger(0)))
		return
	}

	expireAt := time.Now().Add(time.Duration(ttl) * time.Second) // absolute timestamp when the key should expire
	variables.Expirations[key] = expireAt

	variables.Mu.Unlock()

	persistence.Save()
	conn.Write([]byte(encode.EncodeInteger(1)))
}
