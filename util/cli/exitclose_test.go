package cli_test

import (
	"testing"

	"github.com/adamcolton/luce/util/cli"
	"github.com/stretchr/testify/assert"
)

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
