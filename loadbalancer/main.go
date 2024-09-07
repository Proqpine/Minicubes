package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening port: %s\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnections(conn)
	}
}

func handleConnections(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Received: %s", buf)
}
