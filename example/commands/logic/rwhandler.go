package logic

import (
	"strconv"
	"strings"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
)

func RWHandler(onExit, onClose func()) *cli.Runner {
	ho := &HandlerObject{
		Closer:  onClose != nil,
		Exiter:  onExit != nil,
		Timeout: 25,
	}
	rnr := &cli.Runner{
		Commands:     ho.Commands(),
		OnExit:       onExit,
		OnClose:      onClose,
		Timeout:      ho.Timeout,
		Prompt:       "> ",
		StartMessage: "Welcome to the commands demo\nenter 'help' for more\n",
		InputProc: func(s string) []string {
			return strings.Split(strings.TrimSpace(s), " ")
		},
	}

	rnr.RespHandler = lerr.Must(handler.Handlers(
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
		func(r *ExitResp) {
			rnr.Exit = true
		},
		func(r *CloseResp) {
			rnr.Close = true
			rnr.Exit = true
		},
		func(r *HelpResp) {
			rnr.ShowCommands()
		},
		func(r *SetTimeoutResp) {
			rnr.WriteStrings("Timeout set to ", strconv.Itoa(r.Timeout))
		},
	))

	return rnr
}
