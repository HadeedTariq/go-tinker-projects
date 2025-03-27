package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// Connect to the Tunnel Server
	serverConn, err := net.Dial("tcp", "192.168.10.4:8080") // Change to your relay server IP
	if err != nil {
		fmt.Println("Error connecting to tunnel server:", err)
		return
	}

	defer serverConn.Close()
	client_id := "my-tunnel-server"

	serverConn.Write([]byte(client_id))
	localPort := os.Args[1]
	serverAddr := os.Args[2]
	listener, err := net.Listen("tcp", ":"+localPort)
	if err != nil {
		fmt.Println("Error starting local listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Tunnel client listening on", localPort, "forwarding to", serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting local request:", err)
			continue
		}

		go handleConnection(conn, serverConn)
	}
}

func handleConnection(localConn, serverConn net.Conn) {
	defer localConn.Close()

	go io.Copy(serverConn, localConn)
	io.Copy(localConn, serverConn)
}
