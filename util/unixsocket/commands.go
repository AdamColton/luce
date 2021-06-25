package unixsocket

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
)

// Commands is used to setup and run the server side of a unix socket connection.
type Commands struct {
	Addr         string
	Name         string
	Commands     []Command
	cmdMap       map[string]Command
	StartMessage string
	cls          func()
}

// Close the underlying socket
func (c *Commands) Close() {
	c.cls()
}

// Run the underlying socket
func (c *Commands) Run() error {
	c.populateCmdMap()
	s := New(c.Addr, c.handleUnixClient)
	c.cls = s.Close
	return s.Run()
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

// Help returns a description of all the registered commands as a string.
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
