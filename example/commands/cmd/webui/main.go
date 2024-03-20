package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/example/commands/logic"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp/valuedecoder"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/gorilla/mux"
)

// TODO: a lot of this logic should be moved to a util

var (
	decoder = valuedecoder.Form()

	headerStart = `<!doctype html><html><head><title>`
	headerEnd   = `</title></head><body>`
	bodyEnd     = `</body></html>`
)

func main() {
	r := mux.NewRouter()
	s := &http.Server{
		Addr:    ":6161",
		Handler: r,
	}

	ho := &logic.HandlerObject{
		Timeout: 25,
		ExitCloseHandler: cli.NewExitClose(
			nil,
			func() { s.Shutdown(context.Background()) },
		).Commands(),
	}
	cmds := ho.Commands()
	hm := lerr.Must(cmds.Switch())

	names := cmds.Names()
	menuSlc := slice.TransformSlice(names, func(name string, idx int) (string, bool) {
		return fmt.Sprintf(`<a href="/%s">%s</a>`, name, name), true
	})
	menu := fmt.Sprintf("<div>%s</div>", strings.Join(menuSlc, " | "))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sw := luceio.NewSumWriter(w)
		sw.WriteStrings(headerStart, "Index", headerEnd, menu, "Test Server", bodyEnd)
	})

	finish := func(sw *luceio.SumWriter, i any) {
		ho.Exit, ho.Close = respWriter(sw, i)
		sw.WriteStrings(bodyEnd)
		if ho.Close {
			ho.OnClose()
		}
	}

	for _, name := range names {
		if name == "" {
			continue
		}
		_, h := cmds.Get([]string{name})
		t := h.Type()
		frm := structToForm(t)
		n := name

		if t.Elem().NumField() == 0 {
			r.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
				sw := luceio.NewSumWriter(w)
				sw.WriteStrings(headerStart, name, headerEnd, menu)

				i := reflector.Make(t).Interface()
				i, err := hm.Handle(i)
				if err != nil {
					sw.WriteStrings("Error: ", err.Error(), bodyEnd)
					return
				}

				finish(sw, i)
			}).Methods("GET")
		} else {
			r.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
				sw := luceio.NewSumWriter(w)
				sw.WriteStrings(headerStart, name, headerEnd, menu, "Form:", n, frm, bodyEnd)
			}).Methods("GET")
			r.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
				sw := luceio.NewSumWriter(w)
				sw.WriteStrings(headerStart, name, headerEnd, menu)

				i := reflector.Make(t).Interface()
				err := decoder.Decode(i, r)
				if err != nil {
					sw.WriteStrings("ERROR:", err.Error(), bodyEnd)
					return
				}

				i, err = hm.Handle(i)
				if err != nil {
					sw.WriteStrings("Error: ", err.Error(), bodyEnd)
					return
				}

				finish(sw, i)
			}).Methods("POST")
		}
	}

	s.ListenAndServe()
}

func structToForm(t reflect.Type) string {
	t = t.Elem()
	ln := t.NumField()
	out := make([]string, 0, ln+1)
	out = append(out, `<form method=POST>`)
	for i := 0; i < ln; i++ {
		f := t.Field(i)
		prompt := f.Tag.Get("prompt")

		if prompt == "" {
			prompt = f.Name
		} else {
			prompt = fmt.Sprintf("%s %s", f.Name, prompt)
		}
		out = append(out, fmt.Sprintf(`<div>%s <input type=text name="%s" /></div>`, prompt, f.Name))
	}
	out = append(out, `<div><input type=submit /></div></form>`)
	return strings.Join(out, "\n")
}

func respWriter(sw *luceio.SumWriter, i any) (exit, cls bool) {
	switch r := i.(type) {
	case *logic.PersonResp:
		sw.WriteStrings("Created Person: ", r.Name)
	case string:
		r = strings.ReplaceAll(r, "\n", "<br \\>\n")
		sw.WriteString(r)
	case []string:
		sw.WriteString(`<ul>`)
		for _, s := range r {
			sw.WriteStrings(`<li>`, s, `</li>`)
		}
		sw.WriteString(`</ul>`)
	case int:
		sw.WriteStrings(strconv.Itoa(r), " empty requests")
	case *cli.ExitResp:
		exit = true
		sw.WriteString("exit has no meaning in this context")
	case *cli.CloseResp:
		exit = true
		cls = true
		sw.WriteString("close server not implemented")
	case *cli.HelpResp:
		sw.WriteString(`I have no help to offer. See <a href="/">index</a>.`)
	case *logic.SetTimeoutResp:
		sw.WriteStrings("Timeout set to ", strconv.Itoa(r.Timeout))
	}

	return
}
