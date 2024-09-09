package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Order int
	Addr  string
}

var (
	mu              sync.Mutex
	activeServers   []Server
	inactiveServers []Server
)

func (s *Server) getServer(num int) string {
	if s.Order == num {
		fmt.Println(s.Addr)
		return s.Addr
	}
	return ""
}

func main() {
	activeServers = []Server{
		{
			Order: 1,
			Addr:  ":8080",
		},
		{
			Order: 2,
			Addr:  ":8081",
		},
		{
			Order: 3,
			Addr:  ":8082",
		},
	}

	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}

	defer ln.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on port: %s\n", port)

	go startHealthChecks()

	numOfReq := 0

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		mu.Lock()
		serverIndex := numOfReq % len(activeServers)
		sddress := activeServers[serverIndex].Addr
		numOfReq++
		mu.Unlock()

		go handleConnections(conn, sddress)
	}
}

func startHealthChecks() {
	for {
		time.Sleep(10 * time.Second)
		mu.Lock()

		for i := 0; i < len(activeServers); i++ {
			server := activeServers[i]
			if !handleHealthCheck("http://localhost" + server.Addr) {
				fmt.Printf("Server %s is down, moving to inactive\n", server.Addr)
				inactiveServers = append(inactiveServers, server)
				activeServers = append(activeServers[:i], activeServers[i+1:]...)
				i--
			}
		}

		for i := 0; i < len(inactiveServers); i++ {
			server := inactiveServers[i]
			if handleHealthCheck("http://localhost" + server.Addr) {
				fmt.Printf("Server %s is back online, restoring to active\n", server.Addr)
				activeServers = append(activeServers, server)
				inactiveServers = append(inactiveServers[:i], inactiveServers[i+1:]...)
				i-- // Adjust the index to account for the restored server
			}
		}
		mu.Unlock()
	}
}

func handleHealthCheck(addr string) bool {
	resp, err := http.Get(addr)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
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
