package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP   string
	Port int
	// 在线用户
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// 消息广播
	Message chan string
}

func NewServer(IP string, port int) *Server {
	return &Server{
		IP:        IP,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// Start 启动服务
func (s *Server) Start() {
	// listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net listen err:", err)
		return
	}

	// close listen
	defer listener.Close()

	// 启动监听Message的goroutine
	go s.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listening err: ", err)
			continue
		}

		go s.Handle(conn)
	}
}

func (s *Server) Handle(conn net.Conn) {
	// 上线
	user := NewUser(conn, s)
	user.Online()

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 1024*5)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err: ", err)
				return
			}

			// 提取用户消息（去除\n）
			msg := string(buf[:n-1])
			user.Handle(msg)

			// 维持心跳
			user.Heartbeat <- struct{}{}
		}
	}()

	// 阻塞
	for {
		select {
		case <-user.Heartbeat:

		case <-time.After(time.Minute * 1):
			user.Send("Off line")
			// 销毁用户的资源
			close(user.C)
			// 关闭连接
			conn.Close()

			return
		}
	}
}

// Broadcast 广播消息
func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)
	s.Message <- sendMsg
}

// ListenMessage 监听广播消息
// 一旦有消息就发送给所有在线用户
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, client := range s.OnlineMap {
			client.C <- msg
		}
		s.mapLock.Unlock()
	}
}
