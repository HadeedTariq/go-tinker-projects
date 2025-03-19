package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	// Connect to the Tunnel Server
	serverConn, err := net.Dial("tcp", "192.168.10.4:8080") // Change to your relay server IP
	if err != nil {
		fmt.Println("Error connecting to tunnel server:", err)
		return
	}
	defer serverConn.Close()

	// Send a unique ID to register with the Tunnel Server
	clientID := "my-local-tunnel"
	serverConn.Write([]byte(clientID))

	// Listen on a local port
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		fmt.Println("Error starting local listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Tunnel Client listening on localhost:3000")

	for {
		localConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting local request:", err)
			continue
		}

		// Forward request to the tunnel server
		go handleConnection(localConn, serverConn)
	}
}

func handleConnection(localConn, serverConn net.Conn) {
	defer localConn.Close()

	// Forward data in both directions
	go io.Copy(serverConn, localConn)
	io.Copy(localConn, serverConn)
}
