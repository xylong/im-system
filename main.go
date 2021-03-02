package main

import "flag"

var (
	IP   string
	Port int
)

func init() {
	flag.StringVar(&IP, "ip", "127.0.0.1", "server ip address")
	flag.IntVar(&Port, "port", 8888, "server port")
}

func main() {
	flag.Parse()
	server := NewServer(IP, Port)
	server.Start()
}
