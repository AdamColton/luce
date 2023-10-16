package cli

import (
	"reflect"
	"strings"
	"time"

	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/reflector"
)

type Runner struct {
	*handler.Commands
	Exit, Close bool
	Context
	OnExit, OnClose func()
	RespHandler     *handler.Switch
	Timeout         int
	Prompt          string
	InputProc       func(string) []string
	StartMessage    string
}

func InputProc(str string) []string {
	return strings.Split(strings.TrimSpace(str), " ")
}

func (r *Runner) Run() {
	if r.InputProc == nil {
		r.InputProc = InputProc
	}
	r.WriteStrings(r.StartMessage)
	for !r.Exit {
		r.WriteStrings(r.Prompt)
		input := r.InputProc(r.ReadString())
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

func (r *Runner) ShowCommands() {
	r.Commands.WriteTo(r)
}

func (r *Runner) handleInput(input []string) {
	if len(input) == 0 {
		return
	}
	cmdName, input := input[0], input[1:]

	_, h := r.Commands.Get(cmdName)
	if h == nil {
		_, h = r.Commands.Get("")
		if h == nil {
			r.WriteStrings("unknown command: ", cmdName)
			return
		}
	}

	t := h.Type()
	var s reflect.Value
	var si any
	if t != nil {
		s = reflector.Make(h.Type())
		si = s.Interface()

		if len(input) > 0 {
			_, fields := parseCmd(input)
			for k, v := range fields {
				r.Parser().ParseValueFieldName(s, k, v)
			}
		} else {
			ok := r.PopulateStruct(cmdName, si)
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
