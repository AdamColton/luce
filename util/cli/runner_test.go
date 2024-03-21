package cli_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	buf := bytes.NewBuffer(nil)
	in := make(chan []byte)
	wg := sync.WaitGroup{}

	var r *cli.Runner
	r = &cli.Runner{
		Context:      cli.NewContext(buf, in, nil),
		StartMessage: "Start Message\n",
		Prompt:       ">",
		Commands: lerr.Must(handler.Cmds([]*handler.Command{
			{
				Name:   "test",
				Usage:  "this is a test",
				Action: func() {},
			}, {
				Name: "exit",
				Action: func() {
					r.Exit = true
				},
			}, {
				Name: "person",
				Action: func(p *Person) {
					assert.Equal(t, p.Name, "Adam")
					assert.Equal(t, p.Age, 39)
					wg.Done()
				},
			},
		})),
	}

	wg.Add(1)
	go func() {
		r := reader{buf}

		assert.Equal(t, "Start Message\n>", r.read())
		wg.Add(1)
		in <- []byte("person")
		assert.Equal(t, "(person:Name) ", r.read())
		in <- []byte("Adam")
		assert.Equal(t, "(person:Age) ", r.read())
		in <- []byte("39")

		in <- []byte("exit")
		wg.Done()
	}()

	err := timeout.After(25, func() {
		r.Run()
		wg.Wait()

		wg.Add(1)
		r.Static([]string{"person", "Name:Adam", "Age:39"})
		wg.Wait()
	})
	assert.NoError(t, err)
}

type cmdr struct {
	*cli.ExitCloseHandler
	cli.Helper
	out chan<- string
}

type SayHiReq struct {
	Name string
}
type SayHiResp struct {
	Msg string
}

func (c *cmdr) SayHiHandler(req *SayHiReq) *SayHiResp {
	return &SayHiResp{
		Msg: "Hi " + req.Name,
	}
}

func (c *cmdr) SayHiUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "say hi",
	}
}

func (c *cmdr) Handlers(rnr *cli.Runner) []any {
	// TODO: could I do an auto-response to pull handlers from a source?
	// because if this isn't updated the handlers don't work as expected.
	return []any{
		rnr.ExitRespHandler,
		rnr.CloseRespHandler,
		rnr.HelpRespHandler,
		func(r *SayHiResp) {
			c.out <- r.Msg
		},
	}
}
func (c *cmdr) Commands() *handler.Commands {
	cmds := handler.DefaultRegistrar.Commands(c)
	// TODO: AddAlias needs to take lmap.`Wrapper
	handler.AddAlias(cmds,
		"exit", "q",
		"close", "cls",
		"foo", "f",
	)
	cs := cmds.Vals(nil).Sort(handler.CmdNameLT)

	return lerr.Must(handler.Cmds(cs))
}

type reader struct {
	buf *bytes.Buffer
}

func (r *reader) read() string {
	str := ""
	for ; str == ""; time.Sleep(time.Millisecond) {
		str = r.buf.String()
	}
	r.buf.Reset()
	return str
}

type domainObject struct {
	out chan string
	wg  sync.WaitGroup
}

func (do *domainObject) Cli(ctx cli.Context, onExit func()) {
	c := &cmdr{
		ExitCloseHandler: cli.NewExitClose(onExit, nil).Commands(),
		out:              do.out,
	}
	rnr := cli.NewRunner(c, ctx)
	rnr.Run()
	do.wg.Done()
}

func TestNewRunner(t *testing.T) {
	r := reader{bytes.NewBuffer(nil)}
	in := make(chan []byte)
	ctx := cli.NewContext(r.buf, in, nil)

	do := &domainObject{
		out: make(chan string),
	}
	do.wg.Add(1)
	go do.Cli(ctx, func() {})

	assert.Equal(t, "> ", r.read())
	in <- []byte("sayHi")
	assert.Equal(t, "(sayHi:Name) ", r.read())
	in <- []byte("Adam")
	assert.Equal(t, "Hi Adam", <-do.out)

	assert.Equal(t, "\n> ", r.read())
	in <- []byte("help")
	help := []string{
		"q, exit Exit the client",
		"help",
		"sayHi   say hi",
	}
	tfn := slice.ForAll(func(cmd string) string { return "   " + cmd })
	help = tfn.Slice(help, nil)
	expected := strings.Join(help, "\n") + "\n> "
	assert.Equal(t, expected, r.read())

	in <- []byte("exit")

	assert.NoError(t, timeout.After(25, &(do.wg)))
}
