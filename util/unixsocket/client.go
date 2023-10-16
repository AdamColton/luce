package unixsocket

import (
	"net"
	"slices"
	"strconv"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/lfile"
)

const (
	// ErrNilContext is returned by Client if it is given a nil Context.
	ErrNilContext = lerr.Str("got nil Context")
	// ErrNoSuchSocket is returned by Client if the number that was chosen is not
	// in the list of sockets.
	ErrNoSuchSocket = lerr.Str("no such socket")
)

// Client lets the user pick one of the sockets found in the current directory
// and in /tmp/ and connects ctx to it. The sockets are the files that end in
// ".sock". Lines that are read from ctx are sent to the socket and what the
// socket sends is written to ctx, until the socket closes. It only works with a
// socket that uses ConnPipe, such as a CLISocket. Nothing is connected if no
// socket is found or if the user cancels the choice.
func Client(ctx cli.Context) error {
	return ClientFS(ctx, lfile.OSRepository{})
}

// ClientFS is Client with the sockets looked up in fs.
func ClientFS(ctx cli.Context, fs lfile.CoreFS) error {
	addr, err := getSock(ctx, fs)
	if err != nil || addr == "" {
		return err
	}
	return connect(ctx, addr)
}

// connect runs ctx against the socket at addr until the socket is closed.
func connect(ctx cli.Context, addr string) error {
	ctx.WriteStrings("  Connecting to ", addr, "\n\n")
	conn, err := net.Dial("unix", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	pipe := ConnPipe(conn)

	// == projects.Code.luce.unixsocket ==
	// [ ] Client loses a line when the socket closes
	//  If the user's next line is read at the same moment the socket closes,
	//  ReadString returns it instead of seeing the cancel, and the line is lost.
	//  It needs a way to give a line back or to cancel a read that has started.
	cancel := make(chan bool)
	exited := make(chan struct{})
	go func() {
		defer close(exited)
		for {
			str := ctx.ReadString(cancel)
			select {
			case <-cancel:
				return
			default:
			}
			if ctx.Closed() {
				return
			}
			select {
			case pipe.Snd <- []byte(str):
			case <-cancel:
				return
			}
		}
	}()

	for m := range pipe.Rcv {
		ctx.Write(m)
	}
	close(cancel)
	<-exited
	return nil
}

// getSock lists the sockets and asks which one to use. It returns "" if there
// are none or the user cancels.
func getSock(ctx cli.Context, fs lfile.CoreFS) (string, error) {
	if ctx == nil {
		return "", ErrNilContext
	}
	m := lerr.Must(lfile.RegexMatch(`.+\.sock`, "", ".*"))

	mr := m.Root("")
	mr.CoreFS = fs
	local := slice.FromIterFactory(mr.Factory, nil)

	mr.Root = "/tmp/"
	tmp := slice.FromIterFactory(mr.Factory, nil)

	all := slices.Concat(local, tmp)
	if len(all) == 0 {
		ctx.WriteString("No sockets found\n")
		return "", nil
	}

	ctx.WriteString("  Sockets:\n")
	for i, s := range all {
		is := strconv.Itoa(i)
		ctx.WriteStrings("    ", is, "\t", s, "\n")
	}
	var idx int
	if !ctx.Input("(socket) ", &idx) {
		return "", nil
	}
	if idx < 0 || idx >= len(all) {
		return "", ErrNoSuchSocket
	}
	return all[idx], nil
}
