package types

import "net"

type Client struct {
	Conn net.Conn
}

// Send message to client
func (client Client) Send(input string) {
	client.Conn.Write([]byte(input + "\n"))
}
