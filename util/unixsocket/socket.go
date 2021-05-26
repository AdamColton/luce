package unixsocket

import (
	"io/fs"
	"net"
	"sync"
	"syscall"

	"github.com/adamcolton/luce/util/lfile"
)

type FileSystem interface {
	lfile.FSRemover
}

type Socket struct {
	Addr    string
	Handler func(conn net.Conn)
	stop    chan bool
	mux     sync.Mutex
	FS      FileSystem
}

func New(addr string, handler func(conn net.Conn)) *Socket {
	return &Socket{
		Addr:    addr,
		Handler: handler,
		FS:      lfile.OSRepository{},
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
	// TODO: move this to someplace (like lfile of los) as TryRemove
	if err := s.FS.Remove(addr); err != nil {
		if pathErr, isPathErr := err.(*fs.PathError); !isPathErr || pathErr.Err != syscall.ENOENT {
			return err
		}
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
		s.FS.Remove(addr)
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
