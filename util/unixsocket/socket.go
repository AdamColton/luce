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

	mux     sync.Mutex
	stop    chan struct{} // closed by Close, nil unless running
	done    chan struct{} // closed when Run is finished with this run
	startup *startup      // what AwaitRunning waits for, made when needed
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

// startup is how a run tells AwaitRunning that it has started, or why it could
// not. err is set before done is closed.
type startup struct {
	done chan struct{}
	err  error
}

func (st *startup) finish(err error) {
	st.err = err
	close(st.done)
}

// current is the startup that AwaitRunning waits for. It is called with mux
// held.
func (s *Socket) current() *startup {
	if s.startup == nil {
		s.startup = &startup{done: make(chan struct{})}
	}
	return s.startup
}

// AwaitRunning waits until the Socket is listening, so that a client can
// connect. It can be called before Run is called. It returns nil once the
// Socket is running, and the error that Run returns if it could not start. A
// Run that did not start leaves its error for AwaitRunning until Run is called
// again. Once a Socket has been closed, AwaitRunning waits for the next Run.
func (s *Socket) AwaitRunning() error {
	s.mux.Lock()
	st := s.current()
	s.mux.Unlock()
	<-st.done
	return st.err
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

// start gets the Socket ready to accept connections and tells AwaitRunning. It
// is called with mux held.
func (s *Socket) start() (net.Listener, error) {
	if s.stop != nil {
		return nil, ErrRunning
	}
	st := s.current()
	select {
	case <-st.done:
		// a run that failed to start left this
		s.startup = nil
		st = s.current()
	default:
	}
	l, err := s.open()
	if err != nil {
		st.finish(err)
		return nil, err
	}
	s.stop = make(chan struct{})
	s.done = make(chan struct{})
	st.finish(nil)
	return l, nil
}

// open removes the stale file and listens.
func (s *Socket) open() (net.Listener, error) {
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
	return listen("unix", s.Addr)
}

// Run listens at Addr until Close is called, and tells AwaitRunning when it is
// listening. Before it listens it removes a file
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
		s.startup = nil
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
