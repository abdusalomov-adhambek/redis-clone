package handlers

import (
	"goredisclone/encode"
	"goredisclone/persistence"
	"goredisclone/variables"
	"net"
	"strconv"
)

// INCRHandler handles the INCR command.
// Increments the integer value of an existing key by 1 and returns the new value.
// Returns an error if the key does not exist or its value is not a valid integer.
func INCRHandler(conn net.Conn, args []string) {
	if len(args) != 1 {
		conn.Write([]byte(encode.EncodeError("ERR wrong number of arguments")))
		return
	}

	key := args[0] // key whose integer value will be incremented

	variables.Mu.Lock()

	value, exists := variables.Storage[key] // current string value stored at key
	if !exists {
		conn.Write([]byte(encode.EncodeError("ERR key not found")))
		variables.Mu.Unlock()
		return
	}

	valueInt, err := strconv.Atoi(value) // current value parsed as an integer
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
