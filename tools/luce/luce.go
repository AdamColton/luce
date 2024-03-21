package luce

import (
	"fmt"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
)

type Luce struct {
	*cli.ExitCloseHandler
	Timeout int
	cli.Helper
	Prompt, StartMessage string
}

func New(ec *cli.ExitClose) *Luce {
	return &Luce{
		ExitCloseHandler: ec.Commands(),
		Timeout:          25,
		Helper:           "List all commands",
		Prompt:           "> ",
		StartMessage:     "Welcome to the luce tool\nenter 'help' for more\n",
	}

}

func (l *Luce) Commands() *handler.Commands {
	cmds := handler.DefaultRegistrar.Commands(l)
	if ex, ok := cmds["exit"]; ok {
		ex.Alias = "q"
	}
	cs := cmds.Vals(nil).Sort(handler.CmdNameLT)
	return lerr.Must(handler.Cmds(cs))
}

func (l *Luce) Runner() *cli.Runner {
	r := cli.NewRunner(l)
	l.InitRunner(r)
	return r
}

func (l *Luce) InitRunner(r *cli.Runner) {
	r.Prompt = l.Prompt
	r.StartMessage = l.StartMessage
}

func (l *Luce) EC() *cli.ExitClose {
	return l.ExitClose
}

func (l *Luce) Handlers(rnr *cli.Runner) []any {
	return []any{
		func(rr *RandResp) {
			fmt.Fprint(rnr, rr.R)
		},
		func(rr *Rand64Resp) {
			fmt.Fprint(rnr, rr.R)
		},
		rnr.CloseRespHandler,
		rnr.ExitRespHandler,
		rnr.HelpRespHandler,
	}
}
