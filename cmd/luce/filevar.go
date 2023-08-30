package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
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
		str := string(f)
		if strings.Contains(str, "`") {
			buf := bytes.NewBuffer(nil)
			fmt.Fprintf(buf, "var %s = string([]byte{%d", varName, f[0])
			for _, b := range f[1:] {
				fmt.Fprintf(buf, ", %d", b)
			}
			fmt.Fprint(buf, "})")
			luceio.NewLineWrappingWriter(out).Write(buf.Bytes())
		} else {
			fmt.Fprintf(out, "var %s = `%s`", varName, str)
		}
		fmt.Fprint(out, "\n\n")
	}

	return nil
}
