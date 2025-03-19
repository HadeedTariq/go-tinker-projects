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
		buff := make([]byte, 256)

		n, err := conn.Read(buff)

		if err != nil {
			fmt.Println("Error reading from client:", err)
			conn.Close()
			continue
		}

		clientId := string(buff[:n])

		clientsLock.Lock()
		clients[clientId] = conn
		clientsLock.Unlock()

		fmt.Println("New tunnel client registered:", clientId)

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

	// Assume we only have one client for simplicity
	clientsLock.Lock()
	var clientConn net.Conn
	for _, c := range clients {
		clientConn = c
		break
	}
	clientsLock.Unlock()

	if clientConn == nil {
		fmt.Println("No tunnel clients available")
		conn.Write([]byte("No tunnel clients available"))
		return
	}

	// Forward request to the tunnel client
	go io.Copy(clientConn, conn)
	io.Copy(conn, clientConn)
}
