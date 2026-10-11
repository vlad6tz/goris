package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/vlad6tz/goris.git/internal/protocol"
)

type Server struct {
	listener net.Listener
}

func New(addr string) (*Server, error) {
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
	srv := s.listener
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
	br := bufio.NewReader(conn)                            // created outside of loop cause only one reader is needed

	for {
		req, err := protocol.ReadRequest(br)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return // client hung up which is normal
			}
			if errors.Is(err, protocol.ErrUnknownCommand) {
				conn.Write([]byte("ERR unknown command \r\n"))
				continue // even tho req fails we still keep the conn alive
			}
			return // at this point the req is cooked so we drop it
		}
		fmt.Printf("%s %s = %s\n", req.Command, req.Key, req.Value)
		conn.Write([]byte("OK\r\n"))
	}
}
