package luceio

import (
	"io"
	"unicode"
	"unicode/utf8"
)

// WrapWidth allows setting a default width to be used by new instances of
// LineWrapperContext. It is the width used by DefaulLineWrapperContext.
var WrapWidth = 80

// LineWrapperContext returns contextual information to guide the writing
// operation. It is assumed that these values will not change with respect to an
// instance.
//
// LineWrapperContext only handle Unix style line endings.
type LineWrapperContext interface {
	// WrapWidth is the maximum length of a line, including the padding.
	WrapWidth() int
	// Padding is written at the start of each new line.
	Padding() string
}

// DefaulLineWrapperContext is a LineWrapperContext that uses the package level
// WrapWidth and an empty string for the padding.
type DefaulLineWrapperContext struct{}

// WrapWidth returns the package level WrapWidth.
func (DefaulLineWrapperContext) WrapWidth() int { return WrapWidth }

// Padding returns an empty string
func (DefaulLineWrapperContext) Padding() string { return "" }

// LineWrappingWriter fulfills io.Writer. It breaks lines at whitespace so they
// are no longer than the WrapWidth (measured in runes), and writes the Padding
// after each new line. A word longer than the width is not broken. Line breaks
// already in the text are kept. Only Unix line endings are handled.
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
// underlying writer. If w is a *SumWriter it is used directly, and the context
// is looked for on its underlying Writer. If the writer fulfills
// LineWrapperContext it is used, otherwise the DefaulLineWrapperContext is. The
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

// Write b to the underlying io.Writer attempting to line wrap as it goes. The
// wrapping continues across calls to Write. The whitespace a line is broken at
// is replaced by the new line. It returns the number of bytes written, which
// can differ from len(b) because of the padding, the replaced whitespace and
// whitespace that is not written (see below). Whitespace that is not followed by
// a non-whitespace character before the next new line or the end of b is not
// written.
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
			i += size
			w.SumWriter.Write(b[start:i])
			w.SumWriter.Write(w.padding)
			w.onNewLine = true
			lineLen = w.lnPad
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

// WriteNewline writes a new line followed by the padding and starts a new line.
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

// WritePadding writes the padding set by the context. It is used to pad the
// first line, as the padding is otherwise written after each new line.
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
	// Width is returned by WrapWidth.
	Width int
	// Pad is returned by Padding.
	Pad string
}

// WrapWidth fulfills LineWrapperContext providing the Width
func (lwcw LineWrapperContextWriter) WrapWidth() int { return lwcw.Width }

// Padding fulfills LineWrapperContext providing the Padding
func (lwcw LineWrapperContextWriter) Padding() string { return lwcw.Pad }
