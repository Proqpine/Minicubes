package main

import (
	"fmt"
	"io"
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

	// data := []byte("Hello World")
	data, err := sendRequest(buf)
	if err != nil {
		fmt.Println("Error forwarding to backend:", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return
	}

}

func sendRequest(data []byte) ([]byte, error) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	return response, nil
}
