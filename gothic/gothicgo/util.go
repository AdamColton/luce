package gothicgo

import (
	"bytes"

	"github.com/adamcolton/luce/lerr"
)

// PrefixWriteToString takes a PrefixWriterTo and Prefixer and uses them to
// generate a string. It is intended for testing rather than production use.
func PrefixWriteToString(w PrefixWriterTo, p Prefixer) string {
	buf := bytes.NewBuffer(nil)
	_, err := w.PrefixWriteTo(buf, p)
	lerr.Panic(err)
	return buf.String()
}
