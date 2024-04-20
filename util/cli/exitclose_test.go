package cli_test

import (
	"testing"

	"github.com/adamcolton/luce/util/cli"
	"github.com/stretchr/testify/assert"
)

func TestExitCloseCallbacks(t *testing.T) {
	var exited, closed int
	ec := cli.NewExitClose(func() { exited++ }, func() { closed++ })
	h := ec.Commands()

	// The handlers only call them when asked to.
	h.CloseHandler(&cli.CloseReq{})
	h.ExitHandler(&cli.ExitReq{})
	assert.Equal(t, 0, exited)
	assert.Equal(t, 0, closed)

	ec.RunClose, ec.RunExit = true, true
	h.CloseHandler(&cli.CloseReq{})
	h.ExitHandler(&cli.ExitReq{})
	assert.Equal(t, 1, exited)
	assert.Equal(t, 1, closed)

	// There does not have to be something to call.
	ec = cli.NewExitClose(nil, nil)
	ec.RunClose, ec.RunExit = true, true
	assert.False(t, ec.CanExit)
	assert.False(t, ec.CanClose)
	assert.NotPanics(t, func() {
		ec.Commands().CloseHandler(&cli.CloseReq{})
		ec.Commands().ExitHandler(&cli.ExitReq{})
	})
}

func TestExitClose(t *testing.T) {
	eFn := func() {}
	cFn := func() {}
	ec := cli.NewExitClose(eFn, cFn)
	assert.Equal(t, ec, ec.EC())

	expected := &cli.ExitCloseHandler{
		ExitClose: ec,
		CloseDesc: "Close the server",
		ExitDesc:  "Exit the client",
	}
	cmds := ec.Commands()
	assert.Equal(t, expected, cmds)

	assert.Equal(t, &cli.CloseResp{}, cmds.CloseHandler(&cli.CloseReq{}))
	assert.Equal(t, &cli.ExitResp{}, cmds.ExitHandler(&cli.ExitReq{}))

	dets := cmds.CloseUsage()
	assert.Equal(t, cmds.CloseDesc, dets.Usage)
	assert.Equal(t, !cmds.CanClose, dets.Disabled)

	dets = cmds.ExitUsage()
	assert.Equal(t, cmds.ExitDesc, dets.Usage)
	assert.Equal(t, !cmds.CanExit, dets.Disabled)
}
