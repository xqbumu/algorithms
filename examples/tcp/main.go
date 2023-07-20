package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Starting TCP listener...")

	// Listen for incoming TCP connections on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("TCP listener started")

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleSocks5Connection(conn)
		continue

		tconn := NewTeeConn(conn)

		// Read the first few bytes of the incoming data
		header := make([]byte, 1024)
		reader := bufio.NewReader(tconn)
		_, err = reader.Read(header)
		if err != nil {
			fmt.Println("Error reading data:", err)
			continue
		}

		// Match the data against the protocol headers and commands
		headerStr := string(header)
		if strings.HasPrefix(headerStr, "SSH-") {
			fmt.Println("SSH connection detected")
			// Handle SSH connection
		} else if strings.HasPrefix(headerStr, "\x05") {
			fmt.Println("SOCKS5 connection detected")
			// Handle SOCKS5 connection
			go handleSocks5Connection(tconn)
		} else if strings.HasPrefix(headerStr, "GET ") ||
			strings.HasPrefix(headerStr, "POST ") ||
			strings.HasPrefix(headerStr, "HEAD ") ||
			strings.HasPrefix(headerStr, "PUT ") ||
			strings.HasPrefix(headerStr, "DELETE ") {
			fmt.Println("HTTP connection detected")
			// Handle HTTP connection
		} else if strings.HasPrefix(headerStr, "\x16\x03") {
			fmt.Println("TLS connection detected")
			// Handle TLS connection
		} else if strings.HasPrefix(headerStr, "CONNECT") {
			fmt.Println("HTTPS connection detected")
			// Handle HTTPS connection
		} else {
			fmt.Println("Unknown connection detected")
			// Handle unknown connection
		}

		// Close the connection
		// conn.Close()
	}
}
