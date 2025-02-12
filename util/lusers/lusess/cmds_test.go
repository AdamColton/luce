package lusess_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/lusers/lusess"
	"github.com/stretchr/testify/assert"
)

type cmdr struct {
	cli.Helper
	*lusess.StoreCmds
	ec *cli.ExitClose
}

func (c *cmdr) Commands() *handler.Commands {
	cmds := handler.DefaultRegistrar.Commands(c)
	handler.AddAlias(cmds,
		"help", "h",
	)
	cs := cmds.Vals(nil).Sort(handler.CmdNameLT)

	return lerr.Must(handler.Cmds(cs))
}

func (c *cmdr) Handlers(rnr *cli.Runner) []any {
	hdlrs := []any{
		rnr.HelpRespHandler,
	}
	hdlrs = append(hdlrs, lusess.AllRespHandlers(rnr)...)

	return hdlrs
}

func (c *cmdr) EC() *cli.ExitClose {
	return c.ec
}

func TestCmds(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	poll := func() string {
		for {
			str := buf.String()
			if str != "" {
				buf.Reset()
				return str
			}
			time.Sleep(time.Millisecond)
		}
	}

	str := newStore()
	c := &cmdr{
		StoreCmds: &lusess.StoreCmds{
			Store: str,
		},
		ec: cli.NewExitClose(nil, nil),
	}
	in := make(chan []byte)
	ctx := cli.NewContext(buf, in, nil)
	r := cli.NewRunner(c, ctx)
	go r.Run()

	assert.Equal(t, "> ", poll())

	in <- []byte("help")
	help := poll()
	assert.Contains(t, help, "g, group")
	assert.Contains(t, help, "lg, listGroups")
	assert.Contains(t, help, "listUsers")
	assert.Contains(t, help, "user")
	assert.Contains(t, help, "ug, userGroup")

	in <- []byte("user Name:Adam Password:testing")
	createUser := poll()
	assert.Contains(t, createUser, "user created")
	in <- []byte("user Name:Lauren Password:testing")
	createUser = poll()
	assert.Contains(t, createUser, "user created")

	in <- []byte("listUsers")
	users := poll()
	assert.Contains(t, users, "Adam")
	assert.Contains(t, users, "Lauren")

	in <- []byte("g Name:admin")
	createGroup := poll()
	assert.Contains(t, createGroup, "group created")
	in <- []byte("g Name:editor")
	createGroup = poll()
	assert.Contains(t, createGroup, "group created")

	in <- []byte("listGroups")
	groups := poll()
	assert.Contains(t, groups, "admin")
	assert.Contains(t, groups, "editor")

	in <- []byte("ug User:Adam Group:admin")
	ug := poll()
	assert.Contains(t, ug, "user was added to group")

	in <- []byte("ug User:Bob Group:admin")
	ug = poll()
	assert.Contains(t, ug, "User not found")

	in <- []byte("ug User:Adam Group:dancers")
	ug = poll()
	assert.Contains(t, ug, "group not found")
}
