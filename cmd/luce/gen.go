package main

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/template"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/adamcolton/luce/lerr"
	"github.com/urfave/cli"
)

func randCmd(c *cli.Context) error {
	max := c.Int64("n")
	if b := c.Int("b"); b > 0 {
		fmt.Println(b)
		max = 1 << uint(b)
	}
	rand.Seed(time.Now().UnixMicro())
	fmt.Println(rand.Int63n(max))
	return nil
}

func randBase64(c *cli.Context) error {
	b := make([]byte, c.Int("b"))
	cryptorand.Read(b)

	fmt.Println(base64.URLEncoding.EncodeToString(b))
	return nil
}

var filterTmpl = template.Must(template.New("filter").Parse(`
// {{.FilterType}} provides tools to filter {{.Type}}s and compose filters
type {{.FilterType}} func({{.Type}}) bool

func ({{.Receiver}} {{.FilterType}}) Slice(vals []{{.Type}}) []{{.Type}} {
	var out []{{.Type}}
	for _, val := range vals {
		if {{.Receiver}}(val) {
			out = append(out, val)
		}
	}
	return out
}

// Chan runs a go routine listening on ch and any {{.Type}} that passes the  
// {{.FilterType}} is passed to the channel that is returned.
func ({{.Receiver}} {{.FilterType}}) Chan(ch <-chan {{.Type}}, buf int) <-chan {{.Type}} {
	out := make(chan {{.Type}}, buf)
	go func() {
		for in := range ch {
			if {{.Receiver}}(in) {
				out <- in
			}
		}
		close(out)
	}()
	return out
}

// Or builds a new {{.FilterType}} that will return true if either underlying
// {{.FilterType}} is true.
func ({{.Receiver}} {{.FilterType}}) Or({{.Receiver}}2 {{.FilterType}}) {{.FilterType}} {
	return func(val {{.Type}}) bool {
		return {{.Receiver}}(val) || {{.Receiver}}2(val)
	}
}

// And builds a new {{.FilterType}} that will return true if both underlying
// {{.FilterType}}s are true.
func ({{.Receiver}} {{.FilterType}}) And({{.Receiver}}2 {{.FilterType}}) {{.FilterType}} {
	return func(val {{.Type}}) bool {
		return {{.Receiver}}(val) && {{.Receiver}}2(val)
	}
}

// Not builds a new {{.FilterType}} that will return true if the underlying
// {{.FilterType}} is false.
func ({{.Receiver}} {{.FilterType}}) Not() {{.FilterType}} {
	return func(val {{.Type}}) bool {
		return !{{.Receiver}}(val)
	}
}
`))

type filterData struct {
	FilterType, Type, Receiver string
}

func (fd *filterData) update() {
	if fd.FilterType == "" {
		fd.FilterType = strings.Title(fd.Type)
	}
	if fd.Receiver == "" {
		r, _ := utf8.DecodeRuneInString(fd.FilterType)
		r = unicode.ToLower(r)
		fd.Receiver = string(r)
	}
}

func filter(c *cli.Context) error {
	t := c.Args().First()
	if t == "" {
		return lerr.Str("Must include type as argument")
	}
	fd := &filterData{
		Type:       t,
		FilterType: c.String("t"),
		Receiver:   c.String("r"),
	}
	fd.update()
	return filterTmpl.Execute(os.Stdout, fd)
}

func rand32(c *cli.Context) error {
	rand.Seed(time.Now().UnixMicro())
	fmt.Println(rand.Uint32())
	return nil
}
