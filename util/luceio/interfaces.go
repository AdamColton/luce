package luceio

// StringWriter is the interface that wraps the WriteString method
type StringWriter interface {
	WriteString(string) (int, error)
}

// StringWriter is the interface that wraps the WriteString method
type StringsWriter interface {
	WriteStrings(...string) (int, error)
}
