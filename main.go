package main

import (
	"fmt"
	"goredisclone/handlers"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var storage = make(map[string]string)
var mu sync.Mutex

func main() {
	port := ":8001"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	log.Println("Redis clone connectect port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("connectoin error:", err)
			continue
		}

		fmt.Println("new client connected:", conn.LocalAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("new client connection: ", conn.RemoteAddr())

	defer conn.Close()

	for {
		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client disconnected")
			} else {
				fmt.Println("read error:", err)
			}
			conn.Close()
			break
		}

		message := string(buffer[:n])
		command, args := parseCommand(message)
		dispatch(conn, command, args)
	}
}

func parseCommand(message string) (string, []string) {
	tokens := strings.Fields(message)

	if len(tokens) == 0 {
		return "", nil
	}

	command := strings.ToUpper(tokens[0])
	args := tokens[1:]

	return command, args
}

func dispatch(conn net.Conn, command string, args []string) {
	switch command {
	case "PING":
		handlers.PingHandler(conn)
	case "SET":
		handlers.SetHanlder(conn, args, storage, &mu)
	case "GET":
		handlers.GetHandler(conn, args, storage, &mu)
	default:
		conn.Write([]byte("Unknown command\r\n"))
	}

}
