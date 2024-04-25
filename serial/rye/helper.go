package rye

func Reverse(b []byte) {
	ln := len(b)
	end := ln / 2
	ln--
	for i := 0; i < end; i++ {
		b[i], b[ln-i] = b[ln-i], b[i]
	}
}

func Inverse(bs []byte) {
	for i, b := range bs {
		bs[i] = ^b
	}
}
