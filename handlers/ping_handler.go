package handlers

import "net"

func PingHandler(conn net.Conn) {
	conn.Write([]byte("PONG\r\n"))
}
