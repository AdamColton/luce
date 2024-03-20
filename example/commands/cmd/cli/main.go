package main

import (
	"os"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/example/commands/logic"
	"github.com/adamcolton/luce/util/cli"
)

func main() {
	in, _ := iobus.Config{
		Sleep: time.Millisecond,
	}.NewReader(os.Stdin)

	enableExit := func() {}
	r := logic.RWHandler(enableExit, nil)
	r.Context = cli.NewContext(os.Stdout, in, nil)

	args := os.Args[1:]
	if len(args) > 0 {
		r.Static(args)
	} else {
		r.Run()
	}
}
