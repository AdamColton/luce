package luceio

import "io"

// StringWriter is the interface that wraps the WriteString method
type StringWriter interface {
	WriteString(string) (int, error)
}

// StringWriter is the interface that wraps the WriteString method
type StringsWriter interface {
	WriteStrings(...string) (int, error)
}

// TemplateExecutor is an interface representing the ExecuteTemplate method on
// a template.
type TemplateExecutor interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
	Execute(io.Writer, interface{}) error
}
