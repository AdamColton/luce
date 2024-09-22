package cli_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// ctxRunner is a CLIRunner that hands the Context it is given to a channel.
type ctxRunner struct {
	ctx    chan cli.Context
	onExit chan bool
}

func (cr ctxRunner) Cli(ctx cli.Context, onExit func()) {
	cr.ctx <- ctx
	cr.onExit <- onExit == nil
}

func TestStdIO(t *testing.T) {
	out := bytes.NewBuffer(nil)
	cli.StdIn = bytes.NewBufferString("hello\n")
	cli.StdOut = out
	defer func() {
		cli.StdIn = os.Stdin
		cli.StdOut = os.Stdout
	}()

	cr := ctxRunner{ctx: make(chan cli.Context), onExit: make(chan bool)}
	go cli.StdIO(cr)

	// The Context reads from StdIn and writes to StdOut.
	var ctx cli.Context
	assert.NoError(t, timeout.After(1000, func() { ctx = <-cr.ctx }))
	assert.NoError(t, timeout.After(1000, func() {
		assert.Equal(t, "hello", ctx.ReadString(nil))
		ctx.WriteString("out")
		assert.True(t, <-cr.onExit)
	}))
	assert.Equal(t, "out", out.String())
}
