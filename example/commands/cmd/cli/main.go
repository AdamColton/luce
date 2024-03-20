package main

import (
	"os"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/example/commands/logic"
	"github.com/adamcolton/luce/util/cli"
)

func main() {
	rdr := iobus.Config{
		Sleep: time.Millisecond,
	}.NewReader(os.Stdin)

	args := os.Args[1:]
	ec := cli.NewExitClose(nil, nil)
	ec.CanExit = len(args) == 0
	ho := logic.New(ec)
	r := ho.Runner()
	r.Context = cli.NewContext(os.Stdout, rdr.Out, nil)

	if r.CanExit {
		r.Run()
	} else {
		r.Static(args)
	}
}
