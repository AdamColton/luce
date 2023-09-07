package merkle

func uint32ToSlice(u uint32) []byte {
	out := make([]byte, 4)
	for i := 0; u > 0; i++ {
		out[i] = byte(u)
		u >>= 8
	}
	return out
}

type intish interface {
	int | int32
}

func divUp[T intish](a, b T) T {
	// TODO: move this to a shared math lib
	out := a / b
	if a%b != 0 {
		out++
	}
	return out
}
