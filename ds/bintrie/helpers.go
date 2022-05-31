package bintrie

import "unsafe"

func sizeOf[U Uint](u U) U {
	return U(unsafe.Sizeof(u)) * 8
}
