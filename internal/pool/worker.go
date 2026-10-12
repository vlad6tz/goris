package pool

import (
	"bufio"
	"net"
	"sync"
)

type Session struct {
	Conn net.Conn
	R    *bufio.Reader
}

type Pool struct {
	ch      chan *Session
	handler func(*Session) bool
	workers int
	quit    chan struct{}
	once    sync.Once
}

func (p *Pool) worker() {
	for {
		select {
		case <-p.quit:
			return
		case s, ok := <-p.ch:
			if !ok {
				return
			}
			if !p.handler(s) {
				s.Conn.Close()
				continue
			}
			select {
			case p.ch <- s:
			case <-p.quit:
				s.Conn.Close()
				return
			}
		}
	}
}
func New(handler func(*Session) bool, workers, queue int) *Pool {
	return &Pool{
		ch:      make(chan *Session, queue),
		handler: handler,
		workers: workers,
		quit:    make(chan struct{}),
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.workers; i++ {
		go p.worker()
	}
}

func (p *Pool) Submit(s *Session) {
	select {
	case p.ch <- s:
	case <-p.quit:
		s.Conn.Close()
	default:
		s.Conn.Close()
	}
}

func (p *Pool) Close() {
	p.once.Do(func() {
		close(p.quit)
	})
}
