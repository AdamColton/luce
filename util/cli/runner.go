package cli

import (
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/upgrade"
)

type Runner struct {
	*handler.Commands
	*ExitClose
	Context
	RespHandler  *handler.Switch
	Timeout      int
	Prompt       string
	InputProc    func(string) []string
	StartMessage string
}

func InputProc(str string) []string {
	return strings.Split(strings.TrimSpace(str), " ")
}

func (r *Runner) Run() {
	if r.ExitClose == nil {
		r.ExitClose = &ExitClose{}
	}
	if r.InputProc == nil {
		r.InputProc = InputProc
	}
	r.WriteStrings(r.StartMessage)
	for !r.Exit {
		r.WriteStrings(r.Prompt)
		input := r.InputProc(r.ReadString(nil))
		r.handleInput(input)
		r.WriteStrings("\n")
	}
	if r.OnExit != nil {
		r.OnExit()
	}
	if r.Close && r.OnClose != nil {
		r.OnClose()
	}
}

func (r *Runner) Static(input []string) {
	r.handleInput(input)
	r.WriteString("\n")
	time.Sleep(time.Millisecond)
}

func (r *Runner) ShowCommands(path []string) {
	r.Commands.Writer(path).WriteTo(r)
}

type Initer interface {
	Init(input []string)
}

func (r *Runner) handleInput(input []string) {
	if len(input) == 0 {
		return
	}
	_, h, idx := r.Commands.Seek(input)
	cmds, input := input[:idx], input[idx:]
	if h == nil {
		_, h = r.Commands.Get([]string{""})
		if h == nil {
			r.WriteString("unknown command: ")
			r.WriteStrings(input...)
			return
		}
	}

	t := h.Type()
	var s reflect.Value
	var si any
	if t != nil {
		s = reflector.Make(h.Type())
		si = s.Interface()

		if initer, ok := si.(Initer); ok {
			initer.Init(input)
		} else if len(input) > 0 {
			_, fields := parseCmd(input)
			for k, v := range fields {
				r.Parser().ParseValueFieldName(s, k, v)
			}
		} else {
			ok := r.PopulateStruct(cmds[idx-1], si)
			if !ok {
				return
			}
		}
	}

	i, err := h.Handle(si)

	if err != nil {
		r.WriteStrings("unknown command: ", err.Error())
	}
	if i != nil && r.RespHandler != nil {
		r.RespHandler.Handle(i)
	}
}

func parseCmd(input []string) ([]string, map[string]string) {
	var out []string
	named := make(map[string]string)
	for _, s := range input {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) == 2 {
			named[parts[0]] = parts[1]
		} else {
			out = append(out, s)
		}
	}
	return out, named
}

type Commander interface {
	Commands() *handler.Commands
	EC() *ExitClose
}

type HandlerLister interface {
	Handlers(*Runner) []any
}

func NewRunner(c Commander, ctx Context) *Runner {
	rnr := &Runner{
		Commands:  c.Commands(),
		ExitClose: c.EC(),
		Timeout:   25,
		Prompt:    "> ",
		InputProc: func(s string) []string {
			return strings.Split(strings.TrimSpace(s), " ")
		},
		Context: ctx,
	}

	if hl, ok := upgrade.To[HandlerLister](c); ok {
		rnr.RespHandler = lerr.Must(handler.Handlers(hl.Handlers(rnr)...))
	}

	return rnr
}

type CLIRunner interface {
	Cli(ctx Context, onExit func())
}

var StdIn io.Reader = os.Stdin
var StdOut io.Writer = os.Stdout

func StdIO(rnr CLIRunner) {
	rdr := iobus.Config{
		Sleep: time.Millisecond,
	}.NewReader(StdIn)
	ctx := NewContext(StdOut, rdr.Out, nil)
	rnr.Cli(ctx, nil)
}
