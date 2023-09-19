package luceio

import "io"

// StringWriter is the interface that wraps the WriteString method. It is the
// same as io.StringWriter.
type StringWriter interface {
	WriteString(string) (int, error)
}

// StringsWriter is the interface that wraps the WriteStrings method. It writes
// each of the strings in order and returns the number of bytes written.
type StringsWriter interface {
	WriteStrings(...string) (int, error)
}

// TemplateExecutor is an interface representing the Execute and ExecuteTemplate
// methods on a template. It is fulfilled by the *Template types in text/template
// and html/template.
type TemplateExecutor interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
	Execute(io.Writer, interface{}) error
}
