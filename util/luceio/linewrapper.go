package luceio

import (
	"io"
	"unicode"
	"unicode/utf8"
)

// WrapWidth allows setting a default width to be used by new instances of
// LineWrapperContext
var WrapWidth = 80

// LineWrapperContext returns contextual information to guide the writing
// operation. It is assumed that these values will not change with respect to an
// instance.
//
// LineWrapperContext only handle Unix style line endings.
type LineWrapperContext interface {
	WrapWidth() int
	Padding() string
}

// DefaulLineWrapperContext will return the package level WrapWidth and an empty
// string for the padding
type DefaulLineWrapperContext struct{}

// WrapWidth returns package level WrapWidth
func (DefaulLineWrapperContext) WrapWidth() int { return WrapWidth }

// Padding returns an empty string
func (DefaulLineWrapperContext) Padding() string { return "" }

// LineWrappingWriter fulfills io.Writer.
type LineWrappingWriter struct {
	LineWrapperContext
	*SumWriter
	sw         StringWriter
	padding    []byte
	onNewLine  bool
	start      int
	lineLength int
	lnPad      int
}

// NewLineWrappingWriter returns a LineWrappingWriter that will write to the
// underlying writer. It will try to upgrade the writer to LineWrapperContext
// but if that fails, it will use the DefaulLineWrapperContext. The
// LineWrapperContextWriter can be used to wrap the Writer and set Width and
// Padding.
func NewLineWrappingWriter(w io.Writer) *LineWrappingWriter {
	sw, ok := w.(*SumWriter)
	if !ok {
		sw = &SumWriter{
			Writer: w,
		}
	} else {
		w = sw.Writer
	}

	lwc, ok := w.(LineWrapperContext)
	if !ok {
		lwc = DefaulLineWrapperContext{}
	}

	return &LineWrappingWriter{
		LineWrapperContext: lwc,
		SumWriter:          sw,
	}
}

func (w *LineWrappingWriter) setPadding() {
	w.padding = []byte(w.Padding())
	w.lnPad = utf8.RuneCount(w.padding)
}

func (w *LineWrappingWriter) Write(b []byte) (int, error) {
	if w.Err != nil {
		return 0, w.Err
	}
	if w.padding == nil {
		w.setPadding()
	}

	ww := w.WrapWidth()
	s0 := w.Sum

	start := 0
	lineLen := w.lineLength
	lastWS := -1
	llAtLastWS := -1
	done := true
	i := 0
	for i < len(b) {
		r, size := utf8.DecodeRune(b[i:])
		if r == '\n' {
			w.SumWriter.Write(b[start:i])
			w.SumWriter.Write(w.padding)
			w.onNewLine = true
			lineLen = w.lnPad
			i += size
			start = i // skip \n
			done = true
			continue
		}

		// 0xA0 is non-breaking space
		if unicode.IsSpace(r) && r != 0xA0 {
			lastWS = i
			llAtLastWS = lineLen
		} else {
			done = false
		}
		lineLen++
		if lineLen > ww && lastWS > 0 {
			w.SumWriter.Write(b[start:lastWS])
			start = lastWS + 1
			lastWS = -1
			w.WriteNewline()
			lineLen += -llAtLastWS + w.lnPad
		}
		i += size
	}
	if !done {
		rest := b[start:]
		w.lineLength = len(rest)
		w.SumWriter.Write([]byte(rest))
	}
	return int(w.Sum - s0), w.Err
}

var nl = []byte("\n")

// WriteNewline writes a new line followed by the necessary padding
func (w *LineWrappingWriter) WriteNewline() (int, error) {
	if w.padding == nil {
		w.setPadding()
	}
	s0 := w.Sum
	w.SumWriter.Write(nl)
	w.SumWriter.Write(w.padding)
	w.lineLength = 0
	return int(w.Sum - s0), w.Err
}

// WritePadding writes the padding set by the context.
func (w *LineWrappingWriter) WritePadding() (int, error) {
	if w.padding == nil {
		w.setPadding()
	}
	n, err := w.SumWriter.Write(w.padding)
	w.lineLength += w.lnPad
	return n, err
}

// LineWrapperContextWriter provides a method to add Width and Padding Context
// to a Writer.
type LineWrapperContextWriter struct {
	io.Writer
	Width int
	Pad   string
}

// WrapWidth fulfills LineWrapperContext providing the Width
func (lwcw LineWrapperContextWriter) WrapWidth() int { return lwcw.Width }

// Padding fulfills LineWrapperContext providing the Padding
func (lwcw LineWrapperContextWriter) Padding() string { return lwcw.Pad }
