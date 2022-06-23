package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adamcolton/luce/lerr"
	"github.com/urfave/cli"
)

func filevars(c *cli.Context) error {
	n := c.NArg()
	if n%2 != 0 {
		return lerr.Str("Expect (varname filename)+ ; got odd number of arguments")
	}
	out, err := os.Create(c.String("o") + ".go")
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "package %s\n\n", c.String("p"))
	for i := 0; i < n; i++ {
		varName := c.Args().Get(i)
		f, err := ioutil.ReadFile(c.Args().Get(i + 1))
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "var %s = `%s`\n\n", varName, string(f))
	}

	return nil
}
