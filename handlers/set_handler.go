package handlers

import (
	"goredisclone/variables"
	"net"
	"strconv"
	"strings"
	"time"
)

func SetHanlder(conn net.Conn, args []string) {
	if len(args) < 2 {
		conn.Write([]byte("ERR wrong number of arguments \n"))
		return
	}

	key := args[0]
	value := args[1]
	variables.Mu.Lock()
	defer variables.Mu.Unlock()

	if len(args) == 4 {
		if strings.ToUpper(args[2]) == "EX" {
			ttl, err := strconv.Atoi(args[3])
			if err != nil {
				conn.Write([]byte("ERR invalid TTL \n"))
				return
			}
			expireAt := time.Now().Add(time.Duration(ttl) * time.Second)
			variables.Expirations[key] = expireAt
		}
	}

	variables.Storage[key] = value
	conn.Write([]byte("OK\n"))

}
