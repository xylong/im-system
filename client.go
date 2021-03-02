package main

import (
	"fmt"
	"net"
)

type Client struct {
	IP   string
	Port int
	Name string
	conn net.Conn
}

func NewClient(IP string, port int) *Client {
	client := &Client{IP: IP, Port: port}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.IP, client.Port))
	if err != nil {
		fmt.Println("client dial error:", err)
		return nil
	}

	client.conn = conn

	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("connect failed...")
		return
	}

	select {}
}
