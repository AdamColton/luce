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

func (er *Ref[T, E]) isEntRef() {}

func (er *Ref[T, E]) setIdx(idx int) {
	er.idx = idx
}

func (er *Ref[T, E]) allRefsRm() {
	if er.idx > 0 {
		last := len(allRefs) - 1
		swap := allRefs[last]
		swap.setIdx(er.idx)
		allRefs[er.idx] = swap
		allRefs = allRefs[:last]
		er.idx = 0
	}
}

func (er *Ref[T, E]) addToAllRefs() {
	if er.idx == 0 {
		er.idx = len(allRefs)
		allRefs = append(allRefs, er)
	}
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

// Clear sets the underlying pointer to nil allowing the Entity to be
// garbage collected. This should only be temorary for testing.
func (er *Ref[T, E]) Clear(cacheRm bool) {
	er.ent.p = nil
	er.allRefsRm()
	if cacheRm {
		cache.Delete(er.key.Hash64())
	}
}
