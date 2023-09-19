package luceio

import (
	"io"
)

// TemplateWrapper allows a template to generate instances of TemplateTo.
type TemplateWrapper struct {
	TemplateExecutor
}

// TemplateTo combines the template, name and data into a TemplateTo that
// fulfills io.WriterTo.
func (t TemplateWrapper) TemplateTo(name string, data interface{}) *TemplateTo {
	return &TemplateTo{
		TemplateExecutor: t,
		Name:             name,
		Data:             data,
	}
}

// TemplateTo writes a template and fulfills io.WriterTo. If Name is blank, the
// base template is used, otherwise the named template is used.
type TemplateTo struct {
	TemplateExecutor
	Name string
	Data interface{}
}

// NewTemplateTo returns a TemplateTo which fulfills io.WriterTo.
func NewTemplateTo(template TemplateExecutor, name string, data interface{}) *TemplateTo {
	return &TemplateTo{
		TemplateExecutor: template,
		Name:             name,
		Data:             data,
	}
}

// WriteTo executes the template, writing to w. It fulfills io.WriterTo. If the
// template returns an error, the count returned is 0, even if some bytes were
// written.
func (t *TemplateTo) WriteTo(w io.Writer) (int64, error) {
	sw := NewSumWriter(w)

	var err error
	if t.Name == "" {
		err = t.Execute(sw, t.Data)
	} else {
		err = t.ExecuteTemplate(sw, t.Name, t.Data)
	}
	if err != nil {
		return 0, err
	}

	return sw.Sum, sw.Err
}
