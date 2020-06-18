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
		if err != nil {
			return int64(sum + n), err
		}
		sum += n
	}
	return int64(sum), nil
}
