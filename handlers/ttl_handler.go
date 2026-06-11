package handlers

import (
	"goredisclone/encode"
	"goredisclone/variables"
	"net"
	"time"
)

func TTLHandler(conn net.Conn, args []string) {
	variables.Mu.Lock()
	defer variables.Mu.Unlock()

	if len(args) != 1 {
		conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments")))
		return
	}

	key := args[0]
	if _, ok := variables.Storage[key]; !ok {
		conn.Write([]byte(encode.EncodeInteger(-2)))
		return
	}

	expireAt, exists := variables.Expirations[key]
	if !exists {
		conn.Write([]byte(encode.EncodeNull()))
		return
	}

	ttl := int(time.Until(expireAt).Seconds())

	conn.Write([]byte(encode.EncodeInteger(ttl)))
}
