package luceio

import "io"

// WriterTos takes a slice of WriterTos and provides a single WriteTo method
// that will call each of them.
type WriterTos []io.WriterTo

// WriteTo fulfills io.WriterTo. Each element in the list writes to w.
func (wts WriterTos) WriteTo(w io.Writer) (int64, error) {
	var sum int64
	for _, wt := range wts {
		n, err := wt.WriteTo(w)
		sum += n
		if err != nil {
			return sum, err
		}
	}
	return sum, nil
}

// Merge flattens the slice of WriterTos.
func (wts WriterTos) Merge() WriterTos {
	var out WriterTos
	for _, wt := range wts {
		out = out.merge(wt)
	}
	return out
}

func (wts WriterTos) merge(wt io.WriterTo) WriterTos {
	if sub, ok := wt.(WriterTos); ok {
		for _, wt := range sub {
			wts = wts.merge(wt)
		}
	} else {
		wts = append(wts, wt)
	}
	return wts
}

// Seperator creates a WriterToSeperator that will write a seperator
// between each WriterTo.
func (wts WriterTos) Seperator(seperator any) WriterToSeperator {
	s := WriterToSeperator{
		WriterTos: wts,
	}
	switch st := seperator.(type) {
	case string:
		s.Seperator = []byte(st)
	case []byte:
		s.Seperator = st
	}
	return s
}

// WriterToSeperator will write a seperator between each WriterTo.
type WriterToSeperator struct {
	WriterTos
	Seperator []byte
}

// WriteTo will call WriteTo on each element and write the Seperator between
// each WriterTo call.
func (wts WriterToSeperator) WriteTo(w io.Writer) (int64, error) {
	var s int64
	for i, wt := range wts.WriterTos {
		if i != 0 {
			n, err := w.Write(wts.Seperator)
			if err != nil {
				return s, err
			}
			s += int64(n)
		}
		n, err := wt.WriteTo(w)
		if err != nil {
			return s, err
		}
		s += int64(n)
	}
	return s, nil
}
