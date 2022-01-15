package ltmpl

import (
	"html/template"
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/luceio"
)

// Trimmer takes a filepath and trims it to a useful portion. The
// lfile.PathLength type is a useful instance.
type Trimmer interface {
	Trim(string) string
}

// HTMLLoader defines a list of Globs to match and the path length to preserve.
// This allows for the file structure to be retained for useful template names.
type HTMLLoader struct {
	lfile.IteratorSource
	Trimmer
}

// Load the templates
func (l *HTMLLoader) Load() (*template.Template, error) {
	var t *template.Template
	var addTemplate func(string) *template.Template

	addTemplate = func(name string) *template.Template {
		t = template.New(name)
		addTemplate = t.New
		return t
	}

	fn := func(name string, i lfile.Iterator) error {
		_, err := addTemplate(name).Parse(string(i.Data()))
		return err
	}

	err := l.Each(fn)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (l *HTMLLoader) WriteTo(w io.Writer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sep := ""
	fn := func(name string, i lfile.Iterator) error {
		sw.Fprint("%s{{define \"%s\" -}}\n%s\n{{- end}}", sep, name, string(i.Data()))
		sep = "\n"
		return sw.Err
	}

	return sw.Sum, l.Each(fn)
}

type Each func(name string, i lfile.Iterator) error

func (l *HTMLLoader) Each(fn Each) error {
	i, done := l.Iterator()
	var err error
	for ; err == nil && !done; done = i.Next() {
		tmplname := i.Path()
		if l.Trimmer != nil {
			tmplname = l.Trim(tmplname)
		}
		err = fn(tmplname, i)
	}
	return lerr.Any(err, i.Err())
}
