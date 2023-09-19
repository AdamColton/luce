package luceio

import "io"

// WriterTos takes a slice of WriterTos and provides a single WriteTo method
// that will call each of them.
type WriterTos []io.WriterTo

// WriteTo invokes each of the WriterTos.
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

type WriterToSeperator struct {
	WriterTos
	Seperator []byte
}

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
