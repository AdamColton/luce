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

// Commands is used to setup and run the server side of a unix socket connection.
type Commands struct {
	Name         string
	Addr         string
	stop         chan bool
	Commands     []Command
	cmdMap       map[string]Command
	StartMessage string
	sync.Mutex
}

// Close a running socket.
func (c *Commands) Close() {
	c.Lock()
	if c.stop != nil {
		c.stop <- true
		<-c.stop
		c.stop = nil
	}
	c.Unlock()
}

// Run the socket
func (c *Commands) Run() error {
	c.populateCmdMap()
	addr := c.Addr
	if err := os.RemoveAll(addr); err != nil {
		return err
	}

	l, err := net.Listen("unix", addr)
	if err != nil {
		return err
	}

	c.stop = make(chan bool)
	closed := false

	go func() {
		<-c.stop
		closed = true
		l.Close()
		os.RemoveAll(addr)
		close(c.stop)
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			if closed {
				return nil
			}
			return err
		}

		go c.handleUnixClient(conn)
	}
}

func (c *Commands) handleUnixClient(conn net.Conn) {
	defer conn.Close()
	bus := iobus.NewReadWriter(conn)
	ctx := &Context{
		Socket: c,
		conn:   conn,
		in:     serialbus.String(bus.In),
	}
	dflt := c.cmdMap[""]
	if c.StartMessage != "" {
		ctx.WriteString(c.StartMessage)
	}
	ctx.WriteString(fmt.Sprintf("\n(%s) ", c.Name))
	for str := range ctx.in {
		ctx.rawStr = str
		args := strings.Fields(str)
		if len(args) == 0 {
			continue
		}
		cmd, ok := c.cmdMap[args[0]]
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
		ctx.WriteString(fmt.Sprintf("\n(%s) ", c.Name))
	}
}

func (c *Commands) populateCmdMap() {
	c.cmdMap = make(map[string]Command, len(c.Commands))
	for _, cmd := range c.Commands {
		if cmd.Action != nil {
			c.cmdMap[cmd.Name] = cmd
		}
	}
}

// Help returns a description of all the commands registered with the socket as
// a string.
func (c *Commands) Help() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("  Commands:\n")
	for _, cmd := range c.Commands {
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
