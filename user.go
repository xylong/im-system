package main

import (
	"fmt"
	"net"
	"strings"
)

// User 用户
type User struct {
	Name      string
	Addr      string
	C         chan string
	Heartbeat chan struct{}
	server    *Server
	conn      net.Conn
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:      userAddr,
		Addr:      userAddr,
		C:         make(chan string),
		Heartbeat: make(chan struct{}),
		conn:      conn,
		server:    server,
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
	if msg == "who" {
		for _, user := range u.server.OnlineMap {
			onlineMsg := fmt.Sprintf("[%s]%s\n", user.Addr, user.Name)
			u.Send(onlineMsg)
		}
	} else if strings.Contains(msg, "rename:") {
		name := strings.Split(msg, ":")[1]
		_, ok := u.server.OnlineMap[name]
		if !ok {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[name] = u
			u.server.mapLock.Unlock()
			u.Name = name

			u.Send("您已更新用户名:" + name + "\n")
		} else {
			u.Send("用户名已占用\n")
		}
	} else if strings.Contains(msg, "to:") { // to:张三:hello
		arr := strings.Split(msg, ":")
		to := arr[1]
		if to == "" {
			u.Send("请输入用户名")
			return
		}
		user, ok := u.server.OnlineMap[to]
		if !ok {
			u.Send("用户不存在")
			return
		}
		content := arr[2]
		if content != "" {
			user.Send(fmt.Sprintf("%s:%s", u.Name, content))
		}
		return
	} else {
		u.server.Broadcast(u, msg)
	}
}

// Send 发消息
func (u *User) Send(msg string) {
	u.conn.Write([]byte(msg))
}
