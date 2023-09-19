package luceio

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
