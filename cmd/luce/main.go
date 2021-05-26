package main

import (
	"os"

	"github.com/adamcolton/luce/tools/luce"
	"github.com/adamcolton/luce/util/cli"
)

func main() {
	l := luce.New(os.Args[1:])
	cli.StdIO(l)
}
