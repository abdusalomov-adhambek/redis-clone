package handlers

import (
	"goredisclone/persistence"
	"goredisclone/variables"
	"log"
	"net"
	"strconv"
	"time"
)

func ExpireHandler(conn net.Conn, args []string) {
	variables.Mu.Lock()
	log.Println("ExpireHandler args:", args)
	if len(args) != 2 {
		variables.Mu.Unlock()
		conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}

	key := args[0]
	ttl, err := strconv.Atoi(args[1])
	if err != nil {
		variables.Mu.Unlock()
		conn.Write([]byte("-ERR wrong expire time\r\n"))
		return
	}

	_, exists := variables.Storage[key]
	if !exists {
		variables.Mu.Unlock()
		conn.Write([]byte(":0\r\n"))
		return
	}
	expireAt := time.Now().Add(time.Duration(ttl) * time.Second)

	variables.Expirations[key] = expireAt
	variables.Mu.Unlock()

	persistence.Save()
	conn.Write([]byte(":1\r\n"))
}
