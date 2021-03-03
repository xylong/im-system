package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
	flag int
}

func NewClient(IP string, port int) *Client {
	client := &Client{
		IP:   IP,
		Port: port,
		flag: -1,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.IP, client.Port))
	if err != nil {
		fmt.Println("client dial error:", err)
		return nil
	}

	client.conn = conn

	return client
}

func (c *Client) menu() bool {
	var flag int

	fmt.Println("0 退出")
	fmt.Println("1 公聊")
	fmt.Println("2 私聊")
	fmt.Println("3 更新用户名")
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("请输入0-3选择模式")
		return false
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}

		switch c.flag {
		case 1:
			fmt.Println("公聊")
		case 2:
			fmt.Println("私聊")
		case 3:
			c.UpdateName()
		}
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println(">>>请输入用户名")
	fmt.Scanln(&c.Name)
	msg := fmt.Sprintf("rename:%s\n", c.Name)
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn write error:", err)
		return false
	}

	return true
}

func (c *Client) Deal() {
	io.Copy(os.Stdout, c.conn)
}

func main() {
	flag.Parse()

	client := NewClient(RemoteIP, RemotePort)
	if client == nil {
		fmt.Println("connect failed...")
		return
	}

	go client.Deal()
	client.Run()

	select {}
}
