package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
	"time"
)

func MGetHandler(conn net.Conn, args []string) {
	if len(args) == 0 {
		conn.Write([]byte(
			encode.EncodeError("ERR wrong number of arguments for 'mget' command"),
		))
		return
	}

	values := make([]any, 0, len(args))

	variables.Mu.Lock()

	for _, key := range args {

		value, exists := variables.Storage[key]
		if !exists {
			values = append(values, nil)
			continue
		}

		expireAt, hasExpire := variables.Expirations[key]

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
