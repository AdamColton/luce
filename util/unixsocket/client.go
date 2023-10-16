package unixsocket

import (
	"net"
	"strconv"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/lfile"
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

	// I'm breaking stuff, I'm going to need a way to "cancel"
	pipe := ConnPipe(conn)

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
				pipe.Snd <- []byte(str)
			}
		}
	}()

	for m := range pipe.Rcv {
		ctx.Write(m)
	}
	done = true
	cancel <- true
	return nil
}

var CoreFS lfile.CoreFS = lfile.OSRepository{}

const (
	ErrNilContext = lerr.Str("got nil Context")
)

func getSock(ctx cli.Context) (string, error) {
	if ctx == nil {
		return "", ErrNilContext
	}
	m := lerr.Must(lfile.RegexMatch(`.+\.sock`, "", ".*"))

	mr := m.Root("")
	mr.CoreFS = CoreFS
	local := slice.FromIterFactory(mr.Factory, nil)

	mr.Root = "/tmp/"
	tmp := slice.FromIterFactory(mr.Factory, nil)

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
	if idx < len(all) {
		return all[idx], nil
	}
	return "", nil
}
