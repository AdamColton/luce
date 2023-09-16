package filestore

import "strings"

func EncoderReplacer(oldnew ...string) Encoder {
	r := strings.NewReplacer(oldnew...)
	return func(b []byte) string {
		return r.Replace(string(b))
	}
}

func EncoderCast(b []byte) string {
	return string(b)
}

func EncoderExt(ext string) Encoder {
	return func(b []byte) string {
		return string(b) + "." + ext
	}
}

func EncoderMany(es ...Encoder) Encoder {
	return func(b []byte) string {
		for _, e := range es {
			b = []byte(e(b))
		}
		return string(b)
	}
}

func DecoderCast(name string) []byte {
	return []byte(name)
}

func DecoderRemoveExt(name string) []byte {
	dot := strings.LastIndexByte(name, '.')
	if dot < 0 {
		return []byte(name)
	}
	return []byte(name[:dot])
}
