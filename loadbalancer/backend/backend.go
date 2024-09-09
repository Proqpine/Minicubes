package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Backend listening on port: %s\n", port)

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

	response := "HTTP/1.1 200 OK\r\nContent-Length: 25\r\n\r\nHello From Backend Server"
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return
	}
	fmt.Println("Replied with a hello message")
}
