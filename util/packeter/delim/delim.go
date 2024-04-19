package delim

import "bytes"

type Delimiter []byte

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

func (d Delimiter) Unpacker() *Unpacker {
	return &Unpacker{
		Delimiter: d,
	}
}

type Unpacker struct {
	Delimiter
	buf []byte
}

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
