package logic_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/adamcolton/luce/example/commands/logic"
	"github.com/adamcolton/luce/util/cli"
	"github.com/stretchr/testify/assert"
)

// safeBuffer is written by the Runner and read by the test.
type safeBuffer struct {
	mux sync.Mutex
	buf bytes.Buffer
}

func (sb *safeBuffer) Write(p []byte) (int, error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	return sb.buf.Write(p)
}

func (sb *safeBuffer) take() string {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	s := sb.buf.String()
	sb.buf.Reset()
	return s
}

func TestCommands(t *testing.T) {
	ho := logic.New(cli.NewExitClose(func() {}, func() {}))
	r := ho.Runner()
	out := &safeBuffer{}
	r.Context = cli.NewContext(out, nil, nil)
	run := func(words ...string) string {
		r.Static(words)
		return out.take()
	}

	assert.Equal(t, "Created Person: Ada\n", run("person", "Name:Ada", "Age:36"))
	assert.Equal(t, "Created Person: Bob\n", run("p", "Name:Bob"), "p is the alias of person")
	assert.Equal(t, []string{"Ada", "Bob"}, ho.People)

	assert.Equal(t, " * Ada\n * Bob\n", run("person", "list"), "list is a sub command of person")
	assert.Equal(t, "1 empty requests\n", run("empty"))
	assert.Equal(t, "2 empty requests\n", run("empty"))

	assert.Equal(t, "Timeout set to 5\n", run("setTimeout", "Timeout:5"))
	assert.Equal(t, 5, ho.Timeout)
	assert.Equal(t, "\n", run("timeout"), "a command with no response writes nothing")

	help := run("help")
	assert.Contains(t, help, "q, exit")
	assert.Equal(t, help, run(""), "an empty command shows the help")
	assert.Equal(t, 25, logic.New(cli.NewExitClose(nil, nil)).Timeout, "the default Timeout")
}

func TestCommandsWithoutExit(t *testing.T) {
	// The cli command cannot exit when it runs a single command line.
	ho := logic.New(cli.NewExitClose(nil, nil))
	var r *cli.Runner
	assert.NotPanics(t, func() { r = ho.Runner() })
	out := &safeBuffer{}
	r.Context = cli.NewContext(out, nil, nil)
	r.Static([]string{"help"})
	assert.NotContains(t, out.take(), "exit", "exit is disabled, so it is not listed")
}

func TestHandlers(t *testing.T) {
	ho := logic.New(cli.NewExitClose(func() {}, func() {}))
	r := ho.Runner()
	out := &safeBuffer{}
	r.Context = cli.NewContext(out, nil, nil)

	hs := ho.Handlers(r)
	hs[1].(func(string))("a string")
	assert.Equal(t, "a string", out.take())
	hs[2].(func([]string))([]string{"one", "two"})
	assert.Equal(t, " * one\n * two", out.take())
	assert.Same(t, ho.EC(), ho.ExitClose)
}
