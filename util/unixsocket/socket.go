package unixsocket

import (
	"net"
	"sync"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
)

const (
	// ErrNilHandler is returned by Run if the Socket has no Handler.
	ErrNilHandler = lerr.Str("Socket has no Handler")
	// ErrRunning is returned by Run if the Socket is already running.
	ErrRunning = lerr.Str("Socket is already running")
)

// FileSystem is what a Socket needs from a file system: to remove the file at
// its address.
type FileSystem = lfile.FSRemover

// Socket listens on a unix domain socket and calls Handler, in its own Go
// routine, for every connection. Set Addr, Handler and FS before Run and leave
// them alone until Run returns. Use New to create a Socket.
type Socket struct {
	// Addr is the path of the socket file.
	Addr string
	// Handler is called with each connection. It is responsible for closing it.
	Handler func(conn net.Conn)
	// FS is used to remove the socket file: a stale one before listening, and
	// the Socket's own when it stops. It is lfile.OSRepository by default.
	FS FileSystem

	mux  sync.Mutex
	stop chan struct{} // closed by Close, nil unless running
	done chan struct{} // closed when Run is finished with this run
	// listen can be replaced to test Run.
	listen func(network, addr string) (net.Listener, error)
}

// New creates a Socket that will listen at addr and call handler for each
// connection. The Socket removes the file at addr with lfile.OSRepository.
func New(addr string, handler func(conn net.Conn)) *Socket {
	return &Socket{
		Addr:    addr,
		Handler: handler,
		FS:      lfile.OSRepository{},
	}
}

// Close stops a running Socket and waits for Run to finish: the listener is
// closed and the socket file is removed. Connections that were accepted are
// not closed. It does nothing if the Socket is not running, and it is safe to
// call more than once, at the same time, and from a Handler.
func (s *Socket) Close() {
	s.mux.Lock()
	stop, done := s.stop, s.done
	if stop != nil {
		select {
		case <-stop:
		default:
			close(stop)
		}
	}
	s.mux.Unlock()
	if done != nil {
		<-done
	}
}

// start gets the Socket ready to accept connections. It is called with mux
// held.
func (s *Socket) start() (net.Listener, error) {
	if s.stop != nil {
		return nil, ErrRunning
	}
	if s.Handler == nil {
		return nil, ErrNilHandler
	}
	if err := lfile.TryRemove(s.FS, s.Addr); err != nil {
		return nil, err
	}
	listen := s.listen
	if listen == nil {
		listen = net.Listen
	}
	l, err := listen("unix", s.Addr)
	if err != nil {
		return nil, err
	}
	s.stop = make(chan struct{})
	s.done = make(chan struct{})
	return l, nil
}

// Run listens at Addr until Close is called. Before it listens it removes a file
// that an earlier run left at Addr. It returns nil if it was stopped by Close,
// ErrRunning if the Socket is already running, ErrNilHandler if there is no
// Handler, and otherwise the error that stopped it. However it stops, the
// listener is closed and the socket file is removed. A Socket can be run again
// after Run returns.
func (s *Socket) Run() error {
	s.mux.Lock()
	l, err := s.start()
	if err != nil {
		s.mux.Unlock()
		return err
	}
	stop, done := s.stop, s.done
	addr, handler, fs := s.Addr, s.Handler, s.FS
	s.mux.Unlock()

	go func() {
		select {
		case <-stop:
			l.Close()
		case <-done:
		}
	}()
	defer func() {
		l.Close()
		fs.Remove(addr)
		s.mux.Lock()
		s.stop = nil
		s.mux.Unlock()
		close(done)
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-stop:
				return nil
			default:
				return err
			}
		}
		go handler(conn)
	}
}
