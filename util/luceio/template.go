package luceio

import (
	"io"
)

// TODO: move TemplateExecutor to interfaces

// TemplateExecutor is an interface representing the ExecuteTemplate method on
// a template.
type TemplateExecutor interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
	Execute(io.Writer, interface{}) error
}

// TemplateWrapper allows a template to generate insances of TemplateTo
type TemplateWrapper struct {
	TemplateExecutor
}

// TemplateTo combines the template, name and data and fulfills io.WriterTo
func (t TemplateWrapper) TemplateTo(name string, data interface{}) *TemplateTo {
	return &TemplateTo{
		TemplateExecutor: t,
		Name:             name,
		Data:             data,
	}
}

// TemplateTo writes a template and fulfils WriterTo. If Name is blank, the base
// template is used, otherwise the named template is used.
type TemplateTo struct {
	TemplateExecutor
	Name string
	Data interface{}
}

// NewTemplateTo returns a TemplateTo which fulfils WriterTo
func NewTemplateTo(template TemplateExecutor, name string, data interface{}) *TemplateTo {
	return &TemplateTo{
		TemplateExecutor: template,
		Name:             name,
		Data:             data,
	}
}

// WriteTo writes a template and fulfils WriterTo.
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
