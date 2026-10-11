package main

import (
	"fmt"
	"net"
	"time"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(10 * time.Second)) // this prevents hanging

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Received %s", buffer[:n])
		conn.Write([]byte("Received message\n"))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error listening: %s", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server running")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting: %s", err)
			continue
		}
		go handleConn(conn) // each conn is in a different goroutine
	}
}
