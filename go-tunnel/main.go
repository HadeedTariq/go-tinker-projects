package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

var (
	clients     = make(map[string]net.Conn)
	clientsLock sync.Mutex
)

func main() {
	// Listen for tunnel client connections
	go startTunnelServer(":8080")

	// Listen for external requests
	startPublicServer(":9090")
}

// Tunnel Server: Handles connections from tunnel clients

func startTunnelServer(addr string) {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		fmt.Println("Error starting tunnel server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Tunnel Server listening on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting tunnel client:", err)
			continue
		}

		buf := make([]byte, 256)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			conn.Close()
			continue
		}

		clientID := string(buf[:n])
		clientsLock.Lock()
		clients[clientID] = conn // Store client
		clientsLock.Unlock()

		fmt.Println("New tunnel client registered:", clientID)
	}
}

// Public Server: Handles external requests and forwards them
func startPublicServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error starting public server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Public Server listening on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting public request:", err)
			continue
		}

		go handlePublicRequest(conn)
	}
}

// Handles external requests and forwards to the tunnel client
func handlePublicRequest(conn net.Conn) {
	defer conn.Close()

	clientsLock.Lock()
	var client_conn net.Conn

	for _, c := range clients {
		client_conn = c
		break
	}
	clientsLock.Unlock()

	if client_conn == nil {
		fmt.Println("No tunnel clients available")
		conn.Write([]byte("No tunnel clients available"))
		return
	}

	go io.Copy(client_conn, conn)

	io.Copy(conn, client_conn)
}
