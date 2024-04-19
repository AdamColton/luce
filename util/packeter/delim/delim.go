// Package delim frames messages by separating them with a delimiter.
package delim

import "bytes"

// Delimiter frames messages by ending each one with the delimiter bytes. It
// fulfills packeter.Packer. A Delimiter must not be empty.
type Delimiter []byte

// Pack splits data on the delimiter and returns each piece followed by the
// delimiter. If data does not end with the delimiter, one is added, so data
// with no delimiter in it returns a single message. The messages share memory
// with data. When the delimiter is added to the last message it is appended in
// place, which overwrites any spare capacity of data.
func (d Delimiter) Pack(data []byte) [][]byte {
	lnDelim := len(d)
	lnData := len(data)
	n := bytes.Count(data, d)
	endsWithD := lnData > lnDelim && bytes.Equal(d, data[lnData-lnDelim:])
	var out [][]byte
	if endsWithD {
		out = make([][]byte, n)
	} else {
		out = make([][]byte, n, n+1)
	}

	for i := range out {
		idx := bytes.Index(data, d) + lnDelim
		out[i], data = data[:idx], data[idx:]
	}
	if !endsWithD {
		out = append(out, data)
		out[n] = append(out[n], d...)
	}
	return out
}
