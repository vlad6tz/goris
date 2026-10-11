package server

import "net"

type Server struct {
	listener net.Listener
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) Serve() {

}
