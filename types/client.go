package types

import "net"

type Client struct {
	Conn net.Conn
}

func (client Client) Send(input string) {
	client.Conn.Write([]byte(input))
}
