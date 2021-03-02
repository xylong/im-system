package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

func NewServer(IP string, port int) *Server {
	return &Server{IP: IP, Port: port}
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

}
