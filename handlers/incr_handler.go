package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
	"strconv"
)

func INCRHandler(conn net.Conn, args []string) {
	if len(args) != 1 {
		conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments")))
		return
	}

	key := args[0]

	variables.Mu.Lock()

	value, exists := variables.Storage[key]
	if !exists {
		conn.Write([]byte(encode.EncodeError("ERR key not found")))
		variables.Mu.Unlock()
		return
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		conn.Write([]byte(encode.EncodeError("ERR value is not an integer")))
		variables.Mu.Unlock()
		return
	}

	valueInt++
	variables.Storage[key] = strconv.Itoa(valueInt)

	conn.Write([]byte(encode.EncodeInteger(valueInt)))

	variables.Mu.Unlock()

	persistence.Save()
}
