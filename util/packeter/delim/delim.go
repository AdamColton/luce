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
