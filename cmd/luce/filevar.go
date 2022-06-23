package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

func filevars(c *cli.Context) error {
	out, err := os.Create(c.String("o") + ".go")
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "package %s\n\n", c.String("p"))
	n := c.NArg()
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
