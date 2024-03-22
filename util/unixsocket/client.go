package unixsocket

import (
	"net"
	"path/filepath"
	"strconv"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/util/cli"
)

func Client(ctx cli.Context) error {
	addr, err := getSock(ctx)
	if err != nil {
		return err
	}
	if addr == "" {
		return nil
	}
	ctx.WriteStrings("  Connecting to", addr, "\n\n")
	conn, err := net.Dial("unix", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	//TODO: this ends up consuming one input on close
	//need a better way
	cancel := make(chan bool)
	done := false
	go func() {
		for {
			str := ctx.ReadString(cancel)
			if done {
				break
			} else {
				conn.Write([]byte(str))
			}
		}
	}()

	connRdr := iobus.Config{
		CloseOnEOF: true,
	}.NewReader(conn)
	for m := range connRdr.Out {
		ctx.Write(m)
	}
	done = true
	cancel <- true
	return nil
}

func getSock(ctx cli.Context) (string, error) {
	local, err := filepath.Glob("*.sock")
	if err != nil {
		return "", err
	}

	tmp, err := filepath.Glob("/tmp/*.sock")
	if err != nil {
		return "", err
	}

	all := append(local, tmp...)
	if len(all) == 0 {
		ctx.WriteString("No sockets found")
		return "", nil
	}

	ctx.WriteString("  Sockets:\n")
	for i, s := range all {
		is := strconv.Itoa(i)
		ctx.WriteStrings("    ", is, "\t", s, "\n")
	}
	var idx int
	ctx.Input("(socket) ", &idx)
	if err == nil && idx < len(all) {
		return all[idx], nil
	}
	return "", nil
}
