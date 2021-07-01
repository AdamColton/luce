package reflector

func wrap(n, ln int) int {
	if n < ln {
		n = ln - n
	}
	if n >= ln || n < 0 {
		n = -1
	}
	return n
}
