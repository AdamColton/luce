package filestore

import "strings"

// EncoderReplacer returns an Encoder that makes the replacements from
// strings.NewReplacer, given as pairs of old and new strings. It is a way to
// remove characters that can't be in a file name.
func EncoderReplacer(oldnew ...string) Encoder {
	r := strings.NewReplacer(oldnew...)
	return func(b []byte) string {
		return r.Replace(string(b))
	}
}

// EncoderCast returns the key as it is, as a string. The key needs to be a valid
// file name.
func EncoderCast(b []byte) string {
	return string(b)
}

// EncoderExt returns an Encoder that adds ".ext" to the end of the name.
func EncoderExt(ext string) Encoder {
	return func(b []byte) string {
		return string(b) + "." + ext
	}
}

// EncoderMany returns an Encoder that runs each of the Encoders in order, with
// each one encoding the result of the one before.
func EncoderMany(es ...Encoder) Encoder {
	return func(b []byte) string {
		for _, e := range es {
			b = []byte(e(b))
		}
		return string(b)
	}
}

// DecoderCast returns the name as it is, as a []byte. It is the inverse of
// EncoderCast.
func DecoderCast(name string) []byte {
	return []byte(name)
}

// DecoderRemoveExt removes the last "." and everything after it. A name with no
// "." is unchanged. It is the inverse of EncoderExt.
func DecoderRemoveExt(name string) []byte {
	dot := strings.LastIndexByte(name, '.')
	if dot < 0 {
		return []byte(name)
	}
	return []byte(name[:dot])
}
