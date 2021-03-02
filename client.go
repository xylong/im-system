package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	// 服务器IP
	RemoteIP string
	// 服务器端口
	RemotePort int
)

func init() {
	flag.StringVar(&RemoteIP, "ip", "127.0.0.1", "server ip address")
	flag.IntVar(&RemotePort, "port", 8888, "server port")
}

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
	flag.Parse()

	client := NewClient(RemoteIP, RemotePort)
	if client == nil {
		fmt.Println("connect failed...")
		return
	}

	select {}
}
