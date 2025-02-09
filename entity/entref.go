package entity

// EntPtr enforces that the Entity is a pointer of type *T.
type EntPtr[T any] interface {
	*T
	Entity
}

type entPtr[P any, E EntPtr[P]] struct {
	p *P
}

func (ep entPtr[P, E]) EntKey() Key {
	return E(ep.p).EntKey()
}

func (ep entPtr[P, E]) EntVal(buf []byte) ([]byte, error) {
	return E(ep.p).EntVal(buf)
}

// Ref holds a reference to *T
type Ref[T any, E EntPtr[T]] struct {
	key Key
	ent entPtr[T, E]
	idx int
}

func (er *Ref[T, E]) EntKey() Key {
	if er == nil {
		return nil
	}
	return er.key
}

// WeakGet will return the current value of the reference, but will not load the
// entity if the reference is nil.
func (er *Ref[T, E]) WeakGet() (e E, ok bool) {
	if er.ent.p != nil {
		e = E(er.ent.p)
		ok = true
	}
	return
}
