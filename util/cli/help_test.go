package cli_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/stretchr/testify/assert"
)

func TestHelp(t *testing.T) {
	helper := cli.Helper("show the commands")
	assert.Equal(t, "show the commands", helper.HelpUsage().Usage)
	assert.Equal(t, &cli.HelpResp{Command: []string{"a", "b"}}, helper.HelpHandler(&cli.HelpReq{Command: []string{"a", "b"}}))

	// The words after the command are the path.
	req := &cli.HelpReq{}
	req.Init([]string{"foo", "bar"})
	assert.Equal(t, []string{"foo", "bar"}, req.Command)
}

func TestHelpCommand(t *testing.T) {
	cmds := []*handler.Command{
		{
			Name:   "help",
			Action: cli.Helper("").HelpHandler,
		}, {
			Name:   "parent",
			Usage:  "has children",
			Action: func() {},
			Subcmds: []*handler.Command{
				{Name: "child", Usage: "is one", Action: func() {}},
			},
		},
	}
	r, buf, _ := newRunner(t, cmds)
	r.RespHandler = lerr.Must(handler.Handlers(r.HelpRespHandler))

	r.Static([]string{"help"})
	assert.Equal(t, "   help\n   parent has children\n", buf.String())

	buf.Reset()
	r.Static([]string{"help", "parent"})
	assert.Equal(t, "   child is one\n", buf.String())

	// A command that isn't there is said so.
	buf.Reset()
	r.Static([]string{"help", "nosuch"})
	assert.Equal(t, "unknown command: nosuch\n", buf.String())
}
