package main

import (
	"fmt"
	"net"
)

// User 用户
type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

// ListenMessage 监听user的channel
// 一旦有消息则发送给客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}

// Online 上线
func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	u.server.Broadcast(u, "已上线")
}

// Offline 下线
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	u.server.Broadcast(u, "下线")
}

// Handle 消息处理
func (u *User) Handle(msg string) {
	switch msg {
	case "who":
		for _, user := range u.server.OnlineMap {
			onlineMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)
			u.Send(onlineMsg)
		}
	default:
		u.server.Broadcast(u, msg)
	}
}

// Send 发消息
func (u *User) Send(msg string) {
	u.conn.Write([]byte(msg))
}
