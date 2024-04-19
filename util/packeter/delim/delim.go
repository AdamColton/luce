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

// Unpacker returns an Unpacker that splits messages on d.
func (d Delimiter) Unpacker() *Unpacker {
	return &Unpacker{
		Delimiter: d,
	}
}

// Unpacker fulfills packeter.Unpacker for messages that end in a delimiter. It
// keeps the start of an unfinished message between calls to Unpack, so a
// message can arrive split across several calls. It is not safe for concurrent
// use. Because it embeds the Delimiter it also has the Delimiter's Pack method.
type Unpacker struct {
	Delimiter
	buf []byte
}

// Unpack returns every complete message found in data, without the delimiter.
// Data after the last delimiter is kept and becomes the start of the next
// message. A delimiter that is split across two calls to Unpack is not
// recognized. A message that is complete within data shares memory with it, so
// data must not be reused while the message is still needed.
func (u *Unpacker) Unpack(data []byte) [][]byte {
	var out [][]byte
	for {
		idx := bytes.Index(data, u.Delimiter)
		if idx == -1 {
			u.buf = append(u.buf, data...)
			return out
		}
		ln := len(u.Delimiter)
		if len(u.buf) > 0 {
			u.buf = append(u.buf, data[:idx]...)
			out = append(out, u.buf)
			u.buf = nil
		} else {
			out = append(out, data[:idx])
		}
		data = data[idx+ln:]
	}
}
