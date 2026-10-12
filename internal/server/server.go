package server

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"

	"github.com/vlad6tz/goris.git/internal/pool"
	"github.com/vlad6tz/goris.git/internal/protocol"
	"github.com/vlad6tz/goris.git/internal/store"
)

type Server struct {
	listener net.Listener
	store    *store.Store
	pool     *pool.Pool
}

func New(addr string, s *store.Store, worker, queue int) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		listener: l,
		store:    s,
	}
	srv.pool = pool.New(srv.handleOne, worker, queue)
	return srv, nil
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) Serve() error {
	s.pool.Run() // spawns construction workers
	srv := s.listener
	for {
		cli, err := srv.Accept()
		if err != nil {
			return err
		}
		sess := &pool.Session{
			Conn: cli,
			R:    bufio.NewReader(cli),
		}
		s.pool.Submit(sess)

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

func (s *Server) handleOne(ses *pool.Session) bool {
	ses.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	req, err := protocol.ReadRequest(ses.R)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false // client disconnected
		}
		if errors.Is(err, protocol.ErrUnknownCommand) {
			ses.Conn.Write([]byte("ERR unknown command \r\n"))
			return true // conn is still alive
		}
		return false
	}
	switch req.Command {
	case protocol.GET:
		v, ok := s.store.Get(req.Key)
		if !ok {
			ses.Conn.Write([]byte("ERR key not found\r\n"))
			return true
		}
		ses.Conn.Write([]byte("OK\r\n"))
		ses.Conn.Write(v)
	case protocol.SET:
		s.store.Set(req.Key, req.Value)
		ses.Conn.Write([]byte("OK\r\n"))
	}
	return true
}
