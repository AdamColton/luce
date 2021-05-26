package main

import (
	"os"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/tools/luce"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
)

func main() {
	rdr := iobus.Config{
		Sleep: time.Millisecond,
	}.NewReader(os.Stdin)
	ctx := cli.NewContext(os.Stdout, rdr.Out, nil)

	args := os.Args[1:]
	//args = []string{"rand", "N:100"}
	ec := cli.NewExitClose(nil, nil)
	ec.CanExit = len(args) == 0
	l := &Modes{
		Luce:    luce.New(ec),
		Context: ctx,
	}
	r := cli.NewRunner(l)
	l.InitRunner(r)
	r.Context = ctx

	if r.CanExit {
		r.Run()
	} else {
		r.Static(args)
	}
}

type Modes struct {
	*luce.Luce
	cli.Context
}

func (m *Modes) Commands() *handler.Commands {
	cmds := handler.DefaultRegistrar.Commands(m)
	if ex, ok := cmds["exit"]; ok {
		ex.Alias = "q"
	}
	cs := cmds.Vals(nil).Sort(handler.CmdNameLT)
	return lerr.Must(handler.Cmds(cs))
}

func (m *Modes) Handlers(rnr *cli.Runner) []any {
	return m.Luce.Handlers(rnr)
}
