package logic

import (
	"strconv"

	"github.com/adamcolton/luce/util/cli"
)

func New(ec *cli.ExitClose) *HandlerObject {
	return &HandlerObject{
		ExitCloseHandler: ec.Commands(),
		Timeout:          25,
		Helper:           "List all commands",
	}
}

func (ho *HandlerObject) Runner() *cli.Runner {
	r := cli.NewRunner(ho)
	r.Prompt = "> "
	r.StartMessage = "Welcome to the commands demo\nenter 'help' for more\n"
	return r
}

func (ho *HandlerObject) Handlers(rnr *cli.Runner) []any {
	return []any{
		func(r *PersonResp) {
			rnr.WriteStrings("Created Person: ", r.Name)
		},
		func(r string) {
			rnr.WriteString(r)
		},
		func(r []string) {
			for i, c := range r {
				if i > 0 {
					rnr.WriteString("\n")
				}
				rnr.WriteStrings(" * ", c)
			}
		},
		func(r int) {
			rnr.WriteStrings(strconv.Itoa(r), " empty requests")
		},
		func(r *SetTimeoutResp) {
			rnr.WriteStrings("Timeout set to ", strconv.Itoa(r.Timeout))
		},
		rnr.ExitRespHandler,
		rnr.CloseRespHandler,
		rnr.HelpRespHandler,
	}
}
