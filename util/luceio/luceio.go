package luceio

import (
	"io"
)

// StringWriter is the interface that wraps the WriteString method
type StringWriter interface {
	WriteString(string) (int, error)
}

// MultiWrite takes a Writer and a slice of WriterTos and passes the Writer into
// each of them and writes the seperator between each.
func MultiWrite(w io.Writer, tos []io.WriterTo, seperator string) (int64, error) {
	var sbs []byte
	if seperator != "" {
		sbs = []byte(seperator)
	}
	var s int64
	for i, t := range tos {
		if sbs != nil && i != 0 {
			n, err := w.Write(sbs)
			if err != nil {
				return s, err
			}
			s += int64(n)
		}
		n, err := t.WriteTo(w)
		if err != nil {
			return s, err
		}
		s += int64(n)
	}
	return s, nil
}

// WriterToMerge merges multiple WriterTo into a single one. It tries to cast
// the first to an instance of WriterTos, making it efficient to merge multiple
// with successive calls.
func WriterToMerge(wts ...io.WriterTo) io.WriterTo {
	var w WriterTos
	if len(wts) == 0 {
		return w
	}
	if c, ok := wts[0].(WriterTos); ok {
		w = c
		wts = wts[1:]
	}
	for _, wt := range wts {
		if wt == nil {
			continue
		}
		w = append(w, wt)
	}
	return w
}
