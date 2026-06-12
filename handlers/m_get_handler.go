package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
	"time"
)

// MGetHandler handles the MGET command.
// Returns the values of all specified keys as an array; missing or expired keys
// are returned as null bulk strings in the corresponding positions.
func MGetHandler(conn net.Conn, args []string) {
	if len(args) == 0 {
		conn.Write([]byte(
			encode.EncodeError("ERR wrong number of arguments for 'mget' command"),
		))
		return
	}

	values := make([]any, 0, len(args)) // result slice; nil entries represent missing or expired keys

	variables.Mu.Lock()

	for _, key := range args {

		value, exists := variables.Storage[key] // current value stored at key
		if !exists {
			values = append(values, nil)
			continue
		}

		expireAt, hasExpire := variables.Expirations[key] // expiration timestamp for key, if set

		if hasExpire && time.Now().After(expireAt) {
			delete(variables.Storage, key)
			delete(variables.Expirations, key)
			values = append(values, nil)
			continue
		}

		values = append(values, value)
	}

	variables.Mu.Unlock()

	persistence.Save()

	conn.Write([]byte(encode.EncodeArrayWithNulls(values)))
}
