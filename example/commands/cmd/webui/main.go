package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/adamcolton/luce/example/commands/logic"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp/formdecoder"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/gorilla/mux"
)

var (
	decoder = formdecoder.New()

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
	}
	cmds := ho.Commands()
	hm, err := cmds.Switch()
	lerr.Panic(err)

	names := cmds.Names()
	menuSlc := make([]string, 0, len(names))
	for _, c := range names {
		menuSlc = append(menuSlc, fmt.Sprintf(`<a href="/%s">%s</a>`, c, c))
	}
	menu := fmt.Sprintf("<div>%s</div>", strings.Join(menuSlc, " | "))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sw := luceio.NewSumWriter(w)
		sw.WriteStrings(headerStart, "Index", headerEnd, menu, "Test Server", bodyEnd)
	})

	for _, name := range names {
		_, h := cmds.Get(name)
		t := h.Type()
		frm := structToForm(t)
		n := name

		if t.Elem().NumField() == 0 {
			r.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
				sw := luceio.NewSumWriter(w)
				sw.WriteStrings(headerStart, name, headerEnd, menu)

				i := reflector.Make(t).Interface()
				i, err = hm.Handle(i)
				if err != nil {
					sw.WriteStrings("Error: ", err.Error(), bodyEnd)
					return
				}

				respWriter(sw, i)
				sw.WriteStrings(bodyEnd)
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

				respWriter(sw, i)
				sw.WriteStrings(bodyEnd)
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
	case *logic.ExitResp:
		exit = true
		sw.WriteString("exit has no meaning in this context")
	case *logic.CloseResp:
		exit = true
		cls = true
		sw.WriteString("close server not implemented")
	case *logic.HelpResp:
		sw.WriteString(`I have no help to offer. See <a href="/">index</a>.`)
	case *logic.SetTimeoutResp:
		sw.WriteStrings("Timeout set to ", strconv.Itoa(r.Timeout))
	}

	return
}
