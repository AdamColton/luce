package cli

import (
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/upgrade"
)

// Runner reads command lines and runs the commands. Run writes a prompt, reads a
// line and splits it with InputProc. The words that match Commands are the
// command, and the Handler for it is called with an argument, which is a pointer
// to a struct. If more words follow the command, they are name:value pairs that
// set the fields of the struct, unless the struct fulfills Initer and takes the
// words itself. If nothing follows the command the Runner asks for each field.
// What the handler returns goes to RespHandler.
type Runner struct {
	*handler.Commands
	*ExitClose
	Context
	// RespHandler handles what the commands return, by type.
	RespHandler *handler.Switch
	// Timeout is for the handlers to use. The Runner does not use it.
	Timeout int
	// Prompt is written before every line is read.
	Prompt string
	// InputProc splits a line into words. Run uses InputProc if it is nil.
	InputProc func(string) []string
	// StartMessage is written when Run starts.
	StartMessage string
}

// InputProc splits a line at white space into words. An empty line is one empty
// word, which is the name of the command that runs when nothing matches. It is
// the default for Runner.InputProc.
func InputProc(str string) []string {
	words := strings.Fields(str)
	if len(words) == 0 {
		return []string{""}
	}
	return words
}

// Run writes the StartMessage. Then, until Exit is set, it writes the Prompt,
// reads a line, runs it and writes a new line. When it is done it calls OnExit,
// and OnClose if Close was set. It stops early if the input is closed. It sets
// ExitClose and InputProc if they are nil.
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
		line := r.ReadString(nil)
		if r.Closed() {
			break
		}
		r.handleInput(r.InputProc(line))
		r.WriteStrings("\n")
	}
	if r.OnExit != nil {
		r.OnExit()
	}
	if r.Close && r.OnClose != nil {
		r.OnClose()
	}
}

// Static runs one command line that is already split into words, writes a new
// line and waits a millisecond.
func (r *Runner) Static(input []string) {
	r.handleInput(input)
	r.WriteString("\n")
	time.Sleep(time.Millisecond)
}

// ShowCommands writes the commands below the one at the path, or all of them if
// the path is empty. If there is no command at the path it says so.
func (r *Runner) ShowCommands(path []string) {
	w := r.Commands.Writer(path)
	if w == nil {
		r.WriteStrings("unknown command: ", strings.Join(path, " "))
		return
	}
	w.WriteTo(r)
}

// Initer can be fulfilled by the argument of a handler. The Runner calls Init
// with the words after the command, instead of parsing them as name:value pairs.
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
			r.WriteString(strings.Join(input, " "))
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
			keys := make([]string, 0, len(fields))
			for k := range fields {
				keys = append(keys, k)
			}
			slices.Sort(keys)
			for _, k := range keys {
				err := r.Parser().ParseValueFieldName(s, k, fields[k])
				if err != nil {
					r.WriteStrings("could not parse ", k, ":", fields[k], ": ", err.Error())
					return
				}
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
		r.WriteString(err.Error())
	}
	if i != nil && r.RespHandler != nil {
		if _, err = r.RespHandler.Handle(i); err != nil {
			r.WriteStrings(err.Error(), ": ", reflect.TypeOf(i).String())
		}
	}
}

// parseCmd separates the words of the form name:value from the others.
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

// Commander is what NewRunner needs: the commands to run and the ExitClose.
type Commander interface {
	Commands() *handler.Commands
	EC() *ExitClose
}

// HandlerLister can be fulfilled by a Commander that handles the responses to
// its own commands. The handlers are given the Runner, so that they can act on
// it, as ExitRespHandler does.
type HandlerLister interface {
	Handlers(*Runner) []any
}

// NewRunner creates a Runner for the Commander that reads and writes with ctx.
// The prompt is "> " and Timeout is 25. If c fulfills HandlerLister, its Handlers
// are registered as the RespHandler, and NewRunner panics if one is not valid.
func NewRunner(c Commander, ctx Context) *Runner {
	rnr := &Runner{
		Commands:  c.Commands(),
		ExitClose: c.EC(),
		Timeout:   25,
		Prompt:    "> ",
		InputProc: InputProc,
		Context:   ctx,
	}

	if hl, ok := upgrade.To[HandlerLister](c); ok {
		rnr.RespHandler = lerr.Must(handler.Handlers(hl.Handlers(rnr)...))
	}

	return rnr
}

// CLIRunner can run a command line on a Context. The onExit func is called when
// it exits, if it is not nil.
type CLIRunner interface {
	Cli(ctx Context, onExit func())
}

// StdIn is what StdIO reads. It can be replaced, for instance to test a command
// line.
var StdIn io.Reader = os.Stdin

// StdOut is what StdIO writes. It can be replaced, for instance to capture the
// output of a command line.
var StdOut io.Writer = os.Stdout

// StdIO runs rnr on StdIn and StdOut. The input is read by an iobus.Reader that
// polls every millisecond, and the input is not closed at the end of it.
func StdIO(rnr CLIRunner) {
	rdr := iobus.Config{
		Sleep: time.Millisecond,
	}.NewReader(StdIn)
	ctx := NewContext(StdOut, rdr.Out, nil)
	rnr.Cli(ctx, nil)
}
