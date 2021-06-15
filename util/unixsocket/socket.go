package unixsocket

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
)

// Socket is used to setup and run the server side of a unix socket connection.
type Socket struct {
	Name         string
	Addr         string
	stop         chan bool
	Commands     []Command
	cmdMap       map[string]Command
	StartMessage string
	sync.Mutex
}

// Close a running socket.
func (s *Socket) Close() {
	s.Lock()
	if s.stop != nil {
		s.stop <- true
		<-s.stop
		s.stop = nil
	}
	s.Unlock()
}

// Run the socket
func (s *Socket) Run() error {
	s.populateCmdMap()
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

		go s.handleUnixClient(conn)
	}
}

func (s *Socket) handleUnixClient(conn net.Conn) {
	defer conn.Close()
	bus := iobus.NewReadWriter(conn)
	ctx := &Context{
		Socket: s,
		conn:   conn,
		in:     serialbus.String(bus.In),
	}
	dflt := s.cmdMap[""]
	if s.StartMessage != "" {
		ctx.WriteString(s.StartMessage)
	}
	ctx.WriteString(fmt.Sprintf("\n(%s) ", s.Name))
	for str := range ctx.in {
		ctx.rawStr = str
		args := strings.Fields(str)
		if len(args) == 0 {
			continue
		}
		cmd, ok := s.cmdMap[args[0]]
		if ok {
			ctx.Args = args[1:]
			cmd.Action(ctx)
		} else if dflt.Action != nil {
			ctx.Args = args
			dflt.Action(ctx)
		}
		if ctx.shouldClose {
			break
		}
		ctx.WriteString(fmt.Sprintf("\n(%s) ", s.Name))
	}
}

func (s *Socket) populateCmdMap() {
	s.cmdMap = make(map[string]Command, len(s.Commands))
	for _, cmd := range s.Commands {
		if cmd.Action != nil {
			s.cmdMap[cmd.Name] = cmd
		}
	}
}

// Help returns a description of all the commands registered with the socket as
// a string.
func (s *Socket) Help() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("  Commands:\n")
	for _, cmd := range s.Commands {
		name := cmd.Name
		if name == "" {
			name = "<default>"
		}
		fmt.Fprintf(buf, "    %s", name)
		if cmd.Usage != "" {
			fmt.Fprintf(buf, "\t %s", cmd.Usage)
		}
		fmt.Fprint(buf, "\n")
	}
	return buf.String()
}
