package handlers

import (
	"fmt"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
	"strconv"
)

func INCRHandler(conn net.Conn, args []string) {
	if len(args) != 1 {
		conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}
	key := args[0]

	variables.Mu.Lock()
	value, exists := variables.Storage[key]
	if !exists {
		conn.Write([]byte("-ERR key not found\r\n"))
		variables.Mu.Unlock()
		return
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		conn.Write([]byte("-ERR value is not an integer\r\n"))
		variables.Mu.Unlock()
		return
	}

	valueInt++

	variables.Storage[key] = strconv.Itoa(valueInt)
	fmt.Fprintf(conn, ":%d\r\n", valueInt)

	variables.Mu.Unlock()

	persistence.Save()
}
