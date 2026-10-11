package server

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	listener net.Listener
}

func (s *Server) New(addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: l,
	}, nil
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) Serve() error {
	srv, err := net.Listen("tcp", s.Addr().String())
	if err != nil {
		return err
	}
	defer srv.Close()
	for {
		cli, err := srv.Accept()
		if err != nil {
			return err
		}
		go handleConn(cli) // each conn gets a separate goroutine
	}
	// no return nil cause of infinite for loop and return wont be reachable, error handling is inside the loop
	// and if there is no err the conn is handled
}

func (s *Server) Close() error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(10 * time.Second)) // prevent client hanging
	buffer := make([]byte, 512)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("error reading %s", err)
			return
		}
		fmt.Printf("received buffer: %s", buffer[:n])
		conn.Write([]byte("Received message\n"))
	}
}
