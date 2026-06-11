package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
	"strconv"
	"time"
)

func ExpireHandler(conn net.Conn, args []string) {
	if len(args) != 2 {
		conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments")))
		return
	}

	key := args[0]

	ttl, err := strconv.Atoi(args[1])
	if err != nil {
		conn.Write([]byte(encode.EncodeError("ERR wrong expire time")))
		return
	}

	if ttl <= 0 {
		conn.Write([]byte(encode.EncodeError("ERR expire time must be positive")))
		return
	}

	variables.Mu.Lock()

	_, exists := variables.Storage[key]
	if !exists {
		variables.Mu.Unlock()
		conn.Write([]byte(encode.EncodeInteger(0)))
		return
	}

	expireAt := time.Now().Add(time.Duration(ttl) * time.Second)
	variables.Expirations[key] = expireAt

	variables.Mu.Unlock()

	persistence.Save()
	conn.Write([]byte(encode.EncodeInteger(1)))
}
