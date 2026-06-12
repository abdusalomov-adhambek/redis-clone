package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
)

func MSetHandler(conn net.Conn, args []string) {
	if len(args) == 0 || len(args)%2 != 0 {
		conn.Write([]byte(
			encode.EncodeError("ERR wrong number of arguments for 'mset' command"),
		))
		return
	}
	variables.Mu.Lock()

	for i := 0; i < len(args); i += 2 {
		key := args[i]
		value := args[i+1]
		variables.Storage[key] = value
	}

	variables.Mu.Unlock()

	persistence.Save()
	conn.Write([]byte(encode.EncodeSimpleString("OK")))

}
