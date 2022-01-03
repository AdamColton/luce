package ltmpl

import (
	"html/template"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
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

	var err error
	i, done := l.Iterator()
	for ; !done && err == nil; _, done = i.Next() {
		tmplname := i.Path()
		if l.Trimmer != nil {
			tmplname = l.Trim(tmplname)
		}
		_, err = addTemplate(tmplname).Parse(string(i.Data()))
	}
	err = lerr.NewMany(i.Err(), err).First()
	if err != nil {
		return nil, err
	}

	return t, nil
}
