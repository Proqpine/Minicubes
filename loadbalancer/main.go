package main

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	Order int
	Addr  string
}

func (s *Server) getServer(num int) string {
	if s.Order == num {
		fmt.Println(s.Addr)
		return s.Addr
	}
	return ""
}

func main() {
	servers := []Server{
		{
			Order: 1,
			Addr:  ":8080",
		},
		{
			Order: 2,
			Addr:  ":8081",
		},
	}

	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening port: %s\n", port)

	numOfReq := 0

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		serverIndex := numOfReq % len(servers)
		sddress := servers[serverIndex].Addr

		numOfReq++
		fmt.Println("Adrress", sddress)
		go handleConnections(conn, sddress)
	}
}

func handleConnections(conn net.Conn, addr string) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Received: %s", buf)

	data, err := handleRequestToBackend(buf, addr)
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

func handleRequestToBackend(data []byte, addr string) ([]byte, error) {
	conn, err := net.Dial("tcp", addr)
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
