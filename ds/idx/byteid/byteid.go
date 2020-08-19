package byteid

import (
	"crypto/rand"
)

// ID is a helper for managing byte slices as IDs
type ID []byte

// IDLen is used to generate byte slice IDs
type IDLen uint16

// Zero returns a byte slice of length IDLen that is all zeros
func (ln IDLen) Zero() ID {
	return make(ID, ln)
}

// Rand creates a random byte slice ID of length IDLen
func (ln IDLen) Rand() ID {
	id := ln.Zero()
	rand.Read(id)
	return id
}

// Inc returns a new byte slice ID that is incremented by one from the given
// ID.
func (id ID) Inc() ID {
	ln := len(id)
	inc := make(ID, ln)
	copy(inc, id)
	for ln--; ln >= 0; ln-- {
		inc[ln]++
		if inc[ln] != 0 {
			break
		}
	}
	return inc
}

// Dec returns a new byte slice ID that is decremented by one from the given
// ID.
func (id ID) Dec() ID {
	ln := len(id)
	inc := make(ID, ln)
	copy(inc, id)
	for ln--; ln >= 0; ln-- {
		inc[ln]--
		if inc[ln] != 255 {
			break
		}
	}
	return inc
}
