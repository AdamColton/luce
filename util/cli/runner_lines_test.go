package cli_test

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

type greetReq struct {
	Name  string
	Times int
}

type wordsReq struct {
	words []string
}

// Init fulfills cli.Initer, so it takes the words that follow the command.
func (w *wordsReq) Init(input []string) {
	w.words = input
}

type greetResp struct {
	Msg string
}

// newRunner creates a Runner with the commands, that reads the lines, which are
// waiting on a channel that is not closed. The responses of type *greetResp and
// string are written to 'resp'.
func newRunner(t *testing.T, cmds []*handler.Command, lines ...string) (*cli.Runner, *bytes.Buffer, *[]string) {
	ctx, buf := newTestContext(lines...)
	resp := &[]string{}
	r := &cli.Runner{
		Context:  ctx,
		Commands: lerr.Must(handler.Cmds(cmds)),
		RespHandler: lerr.Must(handler.Handlers(
			func(g *greetResp) { *resp = append(*resp, g.Msg) },
			func(s string) { *resp = append(*resp, s) },
		)),
	}
	return r, buf, resp
}

func greetCommands() []*handler.Command {
	return []*handler.Command{
		{
			Name:  "greet",
			Usage: "say hello",
			Action: func(g *greetReq) *greetResp {
				return &greetResp{Msg: fmt.Sprintf("%s x%d", g.Name, g.Times)}
			},
		}, {
			Name:   "words",
			Action: func(w *wordsReq) string { return fmt.Sprint(w.words) },
		}, {
			Name:   "ping",
			Action: func() string { return "pong" },
		}, {
			Name:   "fail",
			Action: func() error { return errors.New("boom") },
		}, {
			Name:   "unhandled",
			Action: func() int { return 5 },
		},
	}
}

func TestStaticArguments(t *testing.T) {
	r, buf, resp := newRunner(t, greetCommands())

	// The arguments after the command are name:value pairs.
	r.Static([]string{"greet", "Name:Ada", "Times:3"})
	assert.Equal(t, []string{"Ada x3"}, *resp)
	assert.Equal(t, "\n", buf.String())

	// Other words are not used.
	r.Static([]string{"greet", "ignored", "Times:2"})
	assert.Equal(t, " x2", (*resp)[1])

	// An argument type that takes the words itself gets all of them.
	r.Static([]string{"words", "a", "b:c", "d"})
	assert.Equal(t, "[a b:c d]", (*resp)[2])

	// A command with no argument.
	r.Static([]string{"ping"})
	assert.Equal(t, "pong", (*resp)[3])

	// Nothing to run.
	buf.Reset()
	r.Static(nil)
	assert.Equal(t, "\n", buf.String())
}

func TestStaticErrors(t *testing.T) {
	r, buf, resp := newRunner(t, greetCommands())

	// A value that can't be parsed stops the command.
	r.Static([]string{"greet", "Name:Ada", "Times:many"})
	assert.Contains(t, buf.String(), "could not parse Times:many: ")
	assert.Empty(t, *resp)

	buf.Reset()
	r.Static([]string{"greet", "Nope:1"})
	assert.Equal(t, "could not parse Nope:1: field 'Nope' not found\n", buf.String())

	// An error from the handler is written as it is.
	buf.Reset()
	r.Static([]string{"fail"})
	assert.Equal(t, "boom\n", buf.String())

	// A command that isn't there.
	buf.Reset()
	r.Static([]string{"nope", "with", "words"})
	assert.Equal(t, "unknown command: nope with words\n", buf.String())
}

func TestStaticFallback(t *testing.T) {
	cmds := append(greetCommands(), &handler.Command{
		Name:   "",
		Action: func() string { return "did not match" },
	})
	r, _, resp := newRunner(t, cmds)

	r.Static([]string{"nope"})
	assert.Equal(t, []string{"did not match"}, *resp)

	// A command that matches is not the fallback.
	r.Static([]string{"ping"})
	assert.Equal(t, "pong", (*resp)[1])
}

func TestStaticResponses(t *testing.T) {
	r, buf, resp := newRunner(t, greetCommands())

	// A response that has no handler is reported.
	r.Static([]string{"unhandled"})
	assert.Equal(t, "no handler found: int\n", buf.String())

	// So is an error from the handler of the response.
	buf.Reset()
	r.RespHandler = lerr.Must(handler.Handlers(func(s string) error { return errors.New("bad response") }))
	r.Static([]string{"ping"})
	assert.Equal(t, "bad response: string\n", buf.String())

	// Without a RespHandler the responses are dropped.
	buf.Reset()
	r.RespHandler = nil
	r.Static([]string{"ping"})
	assert.Equal(t, "\n", buf.String())
	assert.Empty(t, *resp)
}

func TestStaticPrompts(t *testing.T) {
	r, buf, resp := newRunner(t, greetCommands(), "Ada", "2")
	r.Static([]string{"greet"})
	assert.Equal(t, "(greet:Name) (greet:Times) \n", buf.String())
	assert.Equal(t, []string{"Ada x2"}, *resp)

	// Cancelling leaves the command unrun.
	r, buf, resp = newRunner(t, greetCommands(), "Ada", "\x18")
	r.Static([]string{"greet"})
	assert.Empty(t, *resp)
}

func TestShowCommands(t *testing.T) {
	r, buf, _ := newRunner(t, []*handler.Command{
		{
			Name:   "parent",
			Usage:  "has children",
			Action: func() {},
			Subcmds: []*handler.Command{
				{Name: "child", Usage: "is one", Action: func() {}},
			},
		},
	})
	r.ShowCommands(nil)
	assert.Equal(t, "   parent has children", buf.String())

	buf.Reset()
	r.ShowCommands([]string{"parent"})
	assert.Equal(t, "   child is one", buf.String())

	buf.Reset()
	r.ShowCommands([]string{"no", "such"})
	assert.Equal(t, "unknown command: no such", buf.String())
}

func TestRun(t *testing.T) {
	var exited, closed int
	r, buf, resp := newRunner(t, append(greetCommands(),
		&handler.Command{Name: "quit", Action: func() *cli.ExitResp { return &cli.ExitResp{} }},
	), "greet Name:Ada  Times:2", "", "quit", "ping")
	r.StartMessage = "welcome\n"
	r.Prompt = "> "
	r.ExitClose = cli.NewExitClose(func() { exited++ }, func() { closed++ })
	r.RespHandler = lerr.Must(handler.Handlers(
		func(g *greetResp) { *resp = append(*resp, g.Msg) },
		r.ExitRespHandler,
	))

	assert.NoError(t, timeout.After(1000, r.Run))

	// The empty line is the empty word, and nothing matches it. Exiting ends the
	// loop, so the last line is not read.
	assert.Equal(t, "welcome\n> \n> unknown command: \n> \n", buf.String())
	assert.Equal(t, []string{"Ada x2"}, *resp)
	assert.Equal(t, 1, exited)
	assert.Equal(t, 0, closed)
}

func TestRunClose(t *testing.T) {
	var exited, closed int
	r, _, _ := newRunner(t, []*handler.Command{
		{Name: "close", Action: func() *cli.CloseResp { return &cli.CloseResp{} }},
	}, "close")
	r.ExitClose = cli.NewExitClose(func() { exited++ }, func() { closed++ })
	r.RespHandler = lerr.Must(handler.Handlers(r.CloseRespHandler))

	assert.NoError(t, timeout.After(1000, r.Run))
	assert.Equal(t, 1, exited)
	assert.Equal(t, 1, closed)
}

func TestRunDefaults(t *testing.T) {
	// The ExitClose and InputProc are filled in, and Run stops when the input is
	// closed.
	in := make(chan []byte, 1)
	in <- []byte("ping")
	close(in)
	buf := bytes.NewBuffer(nil)
	r := &cli.Runner{
		Context: cli.NewContext(buf, in, nil),
		Commands: lerr.Must(handler.Cmds([]*handler.Command{
			{Name: "ping", Action: func() string { return "pong" }},
		})),
		Prompt: "$ ",
	}
	assert.NoError(t, timeout.After(1000, r.Run))
	assert.Equal(t, "$ \n$ ", buf.String())
	assert.NotNil(t, r.ExitClose)
	assert.NotNil(t, r.InputProc)
}

func TestInputProc(t *testing.T) {
	tt := map[string][]string{
		"foo":              {"foo"},
		"  foo   bar\tbaz": {"foo", "bar", "baz"},
		"a b:c":            {"a", "b:c"},
		"":                 {""},
		"   ":              {""},
	}
	for in, expected := range tt {
		t.Run(strconv.Quote(in), func(t *testing.T) {
			assert.Equal(t, expected, cli.InputProc(in))
		})
	}
}

// stubCommander is a Commander, and if it has handlers, a HandlerLister.
type stubCommander struct {
	cmds     []*handler.Command
	handlers []any
}

func (sc *stubCommander) Commands() *handler.Commands { return lerr.Must(handler.Cmds(sc.cmds)) }
func (sc *stubCommander) EC() *cli.ExitClose          { return &cli.ExitClose{} }

type stubHandlerCommander struct {
	stubCommander
}

func (sc *stubHandlerCommander) Handlers(r *cli.Runner) []any {
	return sc.handlers
}

func TestNewRunnerDefaults(t *testing.T) {
	ctx, _ := newTestContext()
	cmds := []*handler.Command{{Name: "ping", Action: func() string { return "pong" }}}

	r := cli.NewRunner(&stubCommander{cmds: cmds}, ctx)
	assert.Equal(t, "> ", r.Prompt)
	assert.Equal(t, 25, r.Timeout)
	assert.Equal(t, []string{"a", "b"}, r.InputProc("a  b"))
	assert.Nil(t, r.RespHandler)
	assert.NotNil(t, r.ExitClose)

	// A Commander that handles its own responses has a RespHandler.
	hc := &stubHandlerCommander{stubCommander{cmds: cmds, handlers: []any{func(s string) {}}}}
	r = cli.NewRunner(hc, ctx)
	assert.NotNil(t, r.RespHandler)
	_, err := r.RespHandler.Handle("a string")
	assert.NoError(t, err)

	// A handler that isn't one is a mistake in the program.
	hc.handlers = []any{"not a handler"}
	assert.Panics(t, func() { cli.NewRunner(hc, ctx) })
}
