package luce

import (
	"fmt"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
)

type Luce struct {
	*cli.ExitCloseHandler
	*UnixSocketClient
	*RandCommand
	*RandBase64Command
	cli.Helper
	StartMessage      string
	DisableUnixSocket bool
	args              []string
}

func New(args []string) *Luce {
	return &Luce{
		RandCommand:       &RandCommand{},
		RandBase64Command: &RandBase64Command{},
		Helper:            "List all commands",
		StartMessage:      "Welcome to the luce tool\nenter 'help' for more\n",
		args:              args,
	}

}

func (l *Luce) Cli(ctx cli.Context, onExit func()) {
	ec := cli.NewExitClose(onExit, nil)
	ec.CanExit = len(l.args) == 0
	l.ExitCloseHandler = ec.Commands()
	l.UnixSocketClient = &UnixSocketClient{
		Ctx: ctx,
	}

	r := cli.NewRunner(l, ctx)
	r.StartMessage = l.StartMessage

	if ec.CanExit {
		r.Run()
	} else {
		r.Static(l.args)
	}
}

func (l *Luce) Commands() *handler.Commands {
	cmds := handler.DefaultRegistrar.Commands(l)
	handler.AddAlias(cmds, "exit", "q", "help", "h")
	cs := cmds.Vals(nil).Sort(handler.CmdNameLT)
	return lerr.Must(handler.Cmds(cs))
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
