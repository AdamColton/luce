package rye

// Reverse reverses the order of the bytes in b, in place.
func Reverse(b []byte) {
	ln := len(b)
	end := ln / 2
	ln--
	for i := 0; i < end; i++ {
		b[i], b[ln-i] = b[ln-i], b[i]
	}
}

// Inverse flips every bit of every byte in bs, in place.
func Inverse(bs []byte) {
	for i, b := range bs {
		bs[i] = ^b
	}
}
