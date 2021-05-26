package unixsocket

import (
	"net"
	"os"
	"sync"
)

type Socket struct {
	Addr    string
	Handler func(conn net.Conn)
	stop    chan bool
	mux     sync.Mutex
}

func New(addr string, handler func(conn net.Conn)) *Socket {
	return &Socket{
		Addr:    addr,
		Handler: handler,
	}
}

// Close a running socket.
func (s *Socket) Close() {
	s.mux.Lock()
	if s.stop != nil {
		s.stop <- true
		<-s.stop
		s.stop = nil
	}
	s.mux.Unlock()
}

// Run the socket
func (s *Socket) Run() error {
	addr := s.Addr
	if err := os.RemoveAll(addr); err != nil {
		return err
	}

	l, err := net.Listen("unix", addr)
	if err != nil {
		return err
	}

	s.stop = make(chan bool)
	closed := false

	go func() {
		<-s.stop
		closed = true
		l.Close()
		os.RemoveAll(addr)
		close(s.stop)
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			if closed {
				return nil
			}
			return err
		}

		go s.Handler(conn)
	}
}
